package setup

import (
	"context"

	kafkaconfig "github.com/khaiphan29/logpulse/internal/config/kafka"
	"github.com/khaiphan29/logpulse/internal/constants"
	"github.com/khaiphan29/logpulse/internal/kafka"
	"github.com/khaiphan29/logpulse/pkg/logger"
)

func SetupKafkaProducer() (*kafka.Producer, func ()) {
   // Kafka Producer
   producerCfg := kafkaconfig.LoadProducerConfig("KAFKA_PRODUCER_LOG")
   producer, err := kafka.NewProducer(producerCfg)
   if err != nil {
      logger.Fatal("Failed to initialize Kafka producer", map[string]any{
         "error": err,
      })
   }
   producer.StartDeliveryReport()
   return producer, producer.Shutdown
}

func SetupConsumerGroups(msgProcessor kafka.MessageProcessor) (*kafka.ConsumerGroup, func ()) {
   consumerGroupCfgs := kafkaconfig.LoadConsumerGroupConfig("KAFKA_CONSUMER_LOG")
   consumers, err := kafka.CreateConsumerGroup(consumerGroupCfgs, msgProcessor)
   if err != nil {
      logger.Fatal("Failed to initialize Kafka consumer groups", map[string]any{
         "error": err,
      })
   }
   ctx, cancel := context.WithCancel(context.Background())
   consumers.ListenForMessages(-1, ctx)
   return consumers, cancel
}

func GetKafkaTopicsCfg() (map[string]kafkaconfig.TopicConfig) {
   return kafkaconfig.LoadTopicsConfig(constants.KAFKA_TOPICS_CONFIG_PATH).Topics
}
