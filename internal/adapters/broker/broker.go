package broker

import (
	"context"
	"fmt"
	"serverless-service-webhook-adapter/internal/core/app"
	"sync"

	rocketmq "github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var (
	QProducer rocketmq.Producer
	QConsumer rocketmq.PushConsumer

	responseChans map[string]chan *primitive.MessageExt // Map to store channels for responses
	mu            sync.Mutex
)

func InitRocketMQ(cfg *app.BrokerConfig) error {
	var err error

  if responseChans == nil {
    responseChans = make(map[string]chan *primitive.MessageExt)
  }

	listenAddr := fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port)

	// Initialize Producer
	QProducer, err = rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{listenAddr})),
    producer.WithRetry(2),
		producer.WithGroupName("webhook-adapter"),
	)
	if err != nil {
		return fmt.Errorf("failed to create producer: %v", err)
	}

	if err = QProducer.Start(); err != nil {
		return fmt.Errorf("failed to start producer: %v", err)
	}

  // create topics
  var wg sync.WaitGroup
  wg.Add(1)
  err = QProducer.SendAsync(context.Background(),
    func(ctx context.Context, result *primitive.SendResult, e error) {
      if e != nil {
        fmt.Printf("receive message error: %s\n", err)
      } else {
        fmt.Printf("send message success: result=%s\n", result.String())
      }
      wg.Done()
    }, primitive.NewMessage(cfg.Topics.Request, []byte("Create topic")))

  if err != nil {
    fmt.Printf("send message error: %s\n", err)
  }

	// Initialize Consumer
	QConsumer, err = rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{listenAddr}),
    consumer.WithRetry(2),
    consumer.WithAutoCommit(true),
		consumer.WithGroupName("webhook-adapter"),
	)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %v", err)
	}

	if err = QConsumer.Subscribe(cfg.Topics.Response, consumer.MessageSelector{}, handleResponse); err != nil {
		return fmt.Errorf("failed to subscribe to topic: %v", err)
	}

	if err = QConsumer.Start(); err != nil {
		return fmt.Errorf("failed to start consumer: %v", err)
	}

	return nil
}

func handleResponse(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgs {
		fmt.Printf("Received message: %s\n", msg.Body)

		// Get the correlatio ID from the message
		correlationID := msg.TransactionId

		mu.Lock()
		if ch, exists := responseChans[correlationID]; exists {
			ch <- msg
			close(ch)
			delete(responseChans, correlationID)
		}
		mu.Unlock()
	}
	return consumer.ConsumeSuccess, nil
}

func GetResponseChan(correlationID string) chan *primitive.MessageExt {
	mu.Lock()
	defer mu.Unlock()
	ch := make(chan *primitive.MessageExt, 1)

	responseChans[correlationID] = ch
	return ch
}
