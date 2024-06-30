package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"serverless-service-webhook-adapter/internal/adapters/broker"
	"time"

	"go.uber.org/zap"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/google/uuid"
)

const AdapterRoutePath = "/api/adapter"

func SendToBroker() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read the request body
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Generate a unique correlation ID
		correlationID := uuid.New().String()

		// correlationID is used to communicate id between producer and consumer
		msg := &primitive.Message{
			Topic:         "requestTopic",
			Body:          requestBody,
			TransactionId: correlationID,
		}

		// Create a response channel and store it in the map
		responseChan := broker.GetResponseChan(correlationID)

		// Send the message
		res, err := broker.QProducer.SendSync(context.Background(), msg)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
			return
		}

		zap.S().Info("Message sent successfully: result=%s\n", res.String())

		// Wait for the response
		select {
		case responseMsg := <-responseChan:
			// Process and send the response back to the client
			var responseData map[string]interface{}
			json.Unmarshal(responseMsg.Body, &responseData)

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(responseData)
		case <-time.After(10 * time.Second): // Timeout after 10 seconds
			http.Error(w, "Timeout waiting for response", http.StatusGatewayTimeout)
		}
	}
}
