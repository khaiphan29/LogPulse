package main

import (
   "os"
   "fmt"
	"context"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

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
   if len(os.Args) < 2 {
      logger.Error("Usage: create_kafka_topics <broker1_port>", nil)
      return
   }

   // Create admin client
   port := os.Args[1]
   server := fmt.Sprintf("localhost:%s", port)
	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": server})
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
	results, err := admin.CreateTopics(ctx, topics)

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
