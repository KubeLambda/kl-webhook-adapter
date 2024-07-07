package broker

import (
	"fmt"
	"serverless-service-webhook-adapter/internal/core/app"
	"sync"

	"go.uber.org/zap"

	"github.com/nats-io/nats.go"
)

var (
	responseChans  map[string]chan *nats.Msg // Map to store channels for responses
	mu             sync.Mutex
	NATSConnection *nats.Conn
	JetStreamCtx   *nats.JetStreamContext
)

func JetStreamInit(cfg *app.BrokerConfig) (nats.JetStreamContext, error) {
	addrURL := fmt.Sprintf("nats://%s:%d", cfg.Addr, cfg.Port)
	zap.S().Infof("Initialize NATS connection: %s", addrURL)

	NATSConnection, err := nats.Connect(addrURL)
	if err != nil {
		return nil, err
	}
	zap.S().Infof("Connected to '%s'", addrURL)

	jetStreamCtx, err := NATSConnection.JetStream() //nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}
	JetStreamCtx = &jetStreamCtx

	createStream(&jetStreamCtx, cfg.Stream, []string{"request.*", "response.*"})

	if responseChans == nil {
		responseChans = make(map[string]chan *nats.Msg)
	}

	zap.S().Info("Create consumer for request.*")
	jetStreamCtx.Subscribe("response.*", handleResponse)
	return jetStreamCtx, nil
}

func createStream(js *nats.JetStreamContext, streamName string, subjects []string) error {
	zap.S().Infof("Get stream '%s'", streamName)
	jsctx := *js
	stream, err := jsctx.StreamInfo(streamName)

	if stream == nil {
		zap.S().Infof("No '%s' stream found, create one", streamName)
		_, err = jsctx.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: subjects,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func handleResponse(msg *nats.Msg) {
	zap.S().Debugf("Received message: %s", msg.Data)

	// Get the correlation ID from the message
	correlationID := msg.Subject

	mu.Lock()
	if ch, exists := responseChans[correlationID]; exists {
		err := msg.Ack()
		if err != nil {
			zap.S().Error("Unable to Ack %s", msg.Subject)
			return
		}
		ch <- msg

		close(ch)
		delete(responseChans, correlationID)
	}
	mu.Unlock()
}

func GetResponseChan(correlationID string) chan *nats.Msg {
	zap.S().Debugf("Create channel '%s'", correlationID)

	mu.Lock()
	defer mu.Unlock()
	ch := make(chan *nats.Msg, 1)

	responseChans[correlationID] = ch
	return ch
}
