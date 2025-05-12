package mykafka

import (
	"fmt"
	"time"
   "context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
   "github.com/khaiphan29/logpulse/pkg/logger"
)

type MessageProcessor interface {
	Process(message *kafka.Message) error
}

type Consumer struct {
   *kafka.Consumer
   messageProcessor MessageProcessor
}

func NewConsumer(config *kafka.ConfigMap, processor MessageProcessor) (*Consumer, error) {
   // Create a new consumer instance
   c, err := kafka.NewConsumer(config)
   if err != nil {
      return nil, fmt.Errorf("failed to create consumer: %w", err)
   }

   return &Consumer {
      Consumer:         c,
      messageProcessor: processor,
   }, nil
}

func (c *Consumer) ListenForMessages(timeoutMs time.Duration, ctx context.Context) {
   // Poll for messages
   for {
      select {
      case <-ctx.Done():
         logger.Info("Stopping consumer polling", nil)
         return
      default:
         topics, _ := c.Subscription()
         fmt.Printf("Consumer: %s on topic: %v is waiting for message\n", c, topics)
         msg, err := c.ReadMessage(timeoutMs)
         if err != nil {
            logger.Error("Failed to read message", map[string]any{
               "error": err,
            })
            continue
         }

         fmt.Printf("Consumer: %s on topic: %v start processing msg...\n", c, topics)
         // Process the message
         if err := c.messageProcessor.Process(msg); err != nil {
            logger.Error("Failed to process message", map[string]any{
               "error": err,
            })
         } else {
            c.Commit()
         }
      }
   }
}
