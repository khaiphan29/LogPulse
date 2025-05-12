package mykafka

import (
   "errors"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

   "github.com/khaiphan29/logpulse/pkg/logger"
)

type KafkaProducer struct {
   *kafka.Producer
}

// One Kafka producer for my app
var (
   kafkaProducer *KafkaProducer
   mutex     sync.Mutex
)

// NewKafkaProducer initializes and returns a new Kafka producer
// brokers exmaple: localhost:9092,localhost:9093
func InitProducer(config *kafka.ConfigMap) (*KafkaProducer, error) {
   mutex.Lock()
   defer mutex.Unlock()

   if kafkaProducer != nil {
      err := errors.New("Kafka producer already exists")
      return kafkaProducer, err
   }

   producer, err := kafka.NewProducer(config)
   if err == nil {
      kafkaProducer = &KafkaProducer{producer}
      startDeliveryReport()
   }

   return kafkaProducer, err
}

func GetProducer() (*KafkaProducer, error) {
   var err error
   if kafkaProducer == nil {
      err = errors.New("Kafka producer not initialized")
   }
   return kafkaProducer, err
}

// SendMessage sends a message to the Kafka topic
func (p *KafkaProducer) SendMessage(topic *string, key, value []byte) error {
	// Produce the message asynchronously
   err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: topic},
		// Currenly do not have key since we just use Kafka for broadcasting logs
		Key:        key,
		Value:      value,
	}, nil)

   return err
}

func startDeliveryReport() {
   go func() {
      logger.Info("Starting Kafka Producer delivery report", nil)
      // Event() return a channel
      for e := range kafkaProducer.Events() {
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
func (p *KafkaProducer) Shutdown() {
	// Wait for all messages to be delivered
	kafkaProducer.Flush(15 * 1000)
	kafkaProducer.Close()
}
