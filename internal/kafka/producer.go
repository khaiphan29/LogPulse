package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

   "github.com/khaiphan29/logpulse/pkg/logger"
   "github.com/khaiphan29/logpulse/internal/config/kafka"
)

type Producer struct {
   *kafka.Producer
}

// NewProducer initializes and returns a new Kafka producer
// brokers exmaple: localhost:9092,localhost:9093
func NewProducer(cfg *kafkaconfig.ProducerConfig) (*Producer, error) {
   config := &kafka.ConfigMap{
      "bootstrap.servers":        cfg.BootstrapServers,
   }
   producer, err := kafka.NewProducer(config)
   if err != nil {
      return nil, err
   }

   return &Producer{ producer }, nil
}

// SendMessage sends a message to the Kafka topic
func (p *Producer) SendMessage(topic string, key, value []byte) error {
	// Produce the message asynchronously
   err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
         Topic: &topic,
      },
		// Currenly do not have key since we just use Kafka for broadcasting logs
		Key:        key,
		Value:      value,
	}, nil)

   return err
}

func (p *Producer) StartDeliveryReport() {
   go func() {
      logger.Info("Starting Kafka Producer delivery report", nil)
      // Event() return a channel
      for e := range p.Events() {
         switch ev := e.(type) {
         case *kafka.Message:
            if ev.TopicPartition.Error != nil {
               logger.Error("Failed to deliver message", map[string]any{
                  "topic": *ev.TopicPartition.Topic,
                  "error": ev.TopicPartition.Error,
               })
            } else {
               logger.Info("PRODUCER: Message delivered", map[string]any{
                  "topic": *ev.TopicPartition.Topic,
                  "partition": ev.TopicPartition.Partition,
                  "key":   string(ev.Key),
                  "value": string(ev.Value),
               })
            }
         case *kafka.Error:
            logger.Error("Kafka error", map[string]any{
               "error": ev,
            })
         }
      }
   }()
}

// Close closes the Kafka producer
func (p *Producer) Shutdown() {
	// Wait for all messages to be delivered
	p.Flush(15 * 1000)
	p.Close()
}
