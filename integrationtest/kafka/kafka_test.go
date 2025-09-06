package kafka_test

import (
   "os"
   "testing"
   "time"
   "context"

   "github.com/khaiphan29/logpulse/internal/setup"
   mykafka "github.com/khaiphan29/logpulse/internal/kafka"
   "github.com/khaiphan29/logpulse/internal/config/kafka"
   "github.com/confluentinc/confluent-kafka-go/v2/kafka"
   "github.com/khaiphan29/logpulse/pkg/logger"
)

type MockMsgProcessor struct{
   Received chan []byte
}

func (m *MockMsgProcessor) Process(msg []byte) error {
   m.Received <- msg
   logger.Info("MockMsgProcessor received message", map[string]any{
      "message": string(msg),
   })
   return nil
}

var (
   producer *mykafka.Producer
   consumers *mykafka.ConsumerGroup
   mockProcessor *MockMsgProcessor
)

func setupTopics() {
   topics := kafkaconfig.LoadTopicsConfig("test_topics.yaml").Topics
   topicSpecs := make([]kafka.TopicSpecification, 0, len(topics))
   for _, cfg := range topics {
      topicSpec := kafka.TopicSpecification{
         Topic:         cfg.Name,
         NumPartitions: cfg.Partitions,
         Config:        cfg.GetConfigMap(),
      }
      topicSpecs = append(topicSpecs, topicSpec)
   }

   // Create admin client
   admin, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS")})
   if err != nil {
      logger.Error("Failed to create Kafka Admin client", map[string]any{
         "error": err,
      })
      return
   }
   defer admin.Close()

   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
   defer cancel()

   results, err := admin.CreateTopics(ctx, topicSpecs, kafka.SetAdminOperationTimeout(30*time.Second))

   // API Level failure
   if err != nil {
      logger.Error("Failed to create topics", map[string]any{
         "error": err,
      })
      return
   }

   // API Level success, but topic creation may still fail
   for _, result := range results {
      if result.Error.Code() != kafka.ErrNoError {
         logger.Error("Failed to create topic", map[string]any{
            "topic": result.Topic,
            "error": result.Error.String(),
         })
      } else {
         logger.Info("Topic created successfully", map[string]any{
            "topic": result.Topic,
         })
      }
   }

}

func TestMain(m *testing.M) {
   // Setup environment variables for testing
   setupTopics()
   os.Setenv("KAFKA_CONSUMER_LOG_TOPICS", "test_logs")
   p, shutdownProducer := setup.SetupKafkaProducer()
   p.StartDeliveryReport()
   mockProcessor = &MockMsgProcessor{
      Received: make(chan []byte, 1),
   }
   c, shutdownConsumers := setup.SetupConsumerGroups(mockProcessor)
   producer = p
   consumers = c

   m.Run()

   shutdownProducer()
   shutdownConsumers()
}

func TestKafkaIntegration(t *testing.T) {
   // Produce a test message
   topic := "test_logs"
   message := "Hello, Kafka!"

   err := producer.SendMessage(topic, nil, []byte(message))
   if err != nil {
      t.Fatalf("Failed to produce message: %v", err)
   }

   // Allow some time for the consumer to process the message
   time.Sleep(5 * time.Second)

   // Wait for the consumer to process it
    select {
    case msg := <-mockProcessor.Received:
        if string(msg) != message {
            t.Errorf("expected 'hello world', got %s", msg)
        }
    case <-time.After(5 * time.Second):
        t.Error("did not receive message in time")
    }
}

