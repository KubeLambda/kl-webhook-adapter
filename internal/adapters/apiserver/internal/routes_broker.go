package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"kl-webhook-adapter/internal/adapters/broker"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const AdapterRoutePath = "/api/adapter"

func SendToBroker() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read the request body
		requestBody, err := io.ReadAll(r.Body)
		if err != nil {
			zap.S().Errorf("Failed to read request body: %v", err)
			http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Generate a unique correlation ID
		correlationID := uuid.New().String()

		// Create a response channel and store it in the map
		responseChan := broker.GetResponseChan(correlationID)

		id := fmt.Sprintf("request.%s", correlationID)
		zap.S().Debugf("New subject id: %s", id)

		js := *broker.JetStreamCtx
		// Send the message
		_, err = js.Publish(id, requestBody)
		if err != nil {
			zap.S().Errorf("Failed to publish message: %v", err)
			http.Error(w, fmt.Sprintf("Failed to send message: %v", err), http.StatusInternalServerError)
			return
		}

		// Wait for the response
		select {
		case responseMsg := <-responseChan:
			// Process and send the response back to the client
			var responseData map[string]interface{}
			json.Unmarshal(responseMsg.Data, &responseData)
			responseData["subject"] = "response"

			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(responseData)
		case <-time.After(10 * time.Second): // Timeout after 10 seconds
			zap.S().Errorf("Timeout waiting for response")
			http.Error(w, "Timeout waiting for response", http.StatusGatewayTimeout)
		}
	}
}
