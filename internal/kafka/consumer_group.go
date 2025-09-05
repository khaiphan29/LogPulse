package kafka

import (
	"fmt"
   "context"
   "time"

   "github.com/khaiphan29/logpulse/internal/config/kafka"
)

type ConsumerGroup struct {
   Consumers []*Consumer
}

func CreateConsumerGroup(cfg *kafkaconfig.ConsumerGroupConfig, msgProcessor MessageProcessor) (*ConsumerGroup, error) {
   consumers := make([]*Consumer, 0, cfg.NumOfConsumers)
   for range cfg.NumOfConsumers {
      consumer, err := NewConsumer(cfg.Config, msgProcessor)
      if err != nil {
         return nil, fmt.Errorf("failed to create consumer: %w", err)
      }

      consumers = append(consumers, consumer)
   }
   return &ConsumerGroup{ consumers }, nil
}

func (cg *ConsumerGroup) ListenForMessages(timeoutMs time.Duration, ctx context.Context) {
   for _, c := range cg.Consumers {
      go c.ListenForMessages(timeoutMs, ctx)
   }
}
