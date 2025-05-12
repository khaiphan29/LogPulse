package main

import (
   "os"
   "fmt"
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
   "github.com/joho/godotenv"

   "github.com/khaiphan29/logpulse/pkg/logger"
)

// Topic specification with compression
var topics = []kafka.TopicSpecification{
	{
		Topic:             "logs",
		NumPartitions:     3,
		Config: map[string]string{
			"compression.type": "lz4",        // Set compression type at the topic level
			"retention.ms":     "604800000", // 7 days retention
		},
	},
	{
		Topic:             "logs-dlq",
		NumPartitions:     3,
		Config: map[string]string{
			"compression.type": "lz4",        // Set compression type at the topic level
			"retention.ms":     "604800000", // 7 days retention
		},
	},
	{
		Topic:             "logs-dlq-permanent",
		NumPartitions:     3,
		Config: map[string]string{
			"compression.type": "lz4",        // Set compression type at the topic level
			"retention.ms":     "-1",         // keep it forever
		},
	},
   {
      Topic:             "__consumer_offsets",
      NumPartitions:     5,
      Config: map[string]string{
         "retention.ms":     "-1",         // keep it forever
         "cleanup.policy":   "compact",    // Compact the topic
      },
   },
}

func main() {
   err := godotenv.Load()
   if err != nil {
      logger.Fatal("Error loading .env file", map[string]any{
         "error": err,
      })
   }

   brokerPort := os.Getenv("KAFKA_BROKER1_PORT")
   server := fmt.Sprintf("localhost:%s", brokerPort)
	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": server})
	if err != nil {
      logger.Fatal("Failed to create Kafka Admin client", map[string]any{
         "error": err,
      })
	}
	defer admin.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results, err := admin.CreateTopics(ctx, topics)
	if err != nil {
      logger.Fatal("Failed to create topics", map[string]any{
         "error": err,
      })
	}

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
