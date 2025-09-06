package kafkaconfig

import (
   "os"
   "strings"
   "github.com/khaiphan29/logpulse/internal/config/utils"
)

type ConsumerConfig struct {
   BootstrapServers string
   GroupID string
   Topics []string
   AutoOffsetReset string
   EnableAutoCommit bool
   EnableAutoOffsetStore bool
}

type ConsumerGroupConfig struct {
   NumOfConsumers int
   Config ConsumerConfig
}

type ProducerConfig struct {
   BootstrapServers string // the initial contact point so that the client (consumer/producer) can connect to the cluster.
}

func LoadConsumerGroupConfig(prefix string) *ConsumerGroupConfig {
   numberOfConsumers := utils.LoadIntEnv(prefix + "_NUM_OF_CONSUMERS")
   autoCommit := utils.LoadBoolEnv(prefix + "_AUTO_COMMIT")
   autoOffsetStore := utils.LoadBoolEnv(prefix + "_AUTO_OFFSET_STORE")

   return &ConsumerGroupConfig{
      NumOfConsumers: numberOfConsumers,
      Config: ConsumerConfig{
         BootstrapServers: os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
         GroupID: os.Getenv(prefix + "_GROUP_ID"),
         Topics: strings.Split(os.Getenv(prefix + "_TOPICS"), ","),
         AutoOffsetReset: os.Getenv(prefix + "_AUTO_OFFSET_RESET"),
         EnableAutoCommit: autoCommit,
         EnableAutoOffsetStore: autoOffsetStore,
      },
   }
}

func LoadProducerConfig(predix string) *ProducerConfig {
   return &ProducerConfig{
      BootstrapServers: os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
   }
}
