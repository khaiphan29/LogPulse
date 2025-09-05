package main

import (
   "os"
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

   "github.com/khaiphan29/logpulse/internal/setup"
   "github.com/khaiphan29/logpulse/pkg/logger"
)


func buildTopicSpecifications() []kafka.TopicSpecification {
   topicsCfg := setup.GetKafkaTopicsCfg()
   topicSpecs := make([]kafka.TopicSpecification, 0, len(topicsCfg))
   for _, cfg := range topicsCfg {
      topicSpec := kafka.TopicSpecification{
         Topic:         cfg.Name,
         NumPartitions: cfg.Partitions,
         Config:        cfg.GetConfigMap(),
      }
      topicSpecs = append(topicSpecs, topicSpec)
   }

   return topicSpecs
}

func main() {
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

   // Create topics
   topicConfigs := buildTopicSpecifications()
	results, err := admin.CreateTopics(ctx, topicConfigs, kafka.SetAdminOperationTimeout(30*time.Second))

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
