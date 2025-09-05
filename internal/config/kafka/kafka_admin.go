package kafkaconfig

import (
   "os"
   "gopkg.in/yaml.v3"
)

type BrokerConfig struct {
   Host string
   Port string
}

type TopicConfig struct {
   Name              string   `yaml:"name"`
   Partitions        int      `yaml:"partitions"`
   CompressionType   string   `yaml:"compression_type"`
   RetentionMs       string   `yaml:"retention_ms"`
   CleanupPolicy     string   `yaml:"cleanup_policy"`
}

type TopicConfigs struct {
   Topics map[string]TopicConfig `yaml:"topics"`
}

func LoadBrokerConfig(prefix string) *BrokerConfig {
   return &BrokerConfig{
      Host: os.Getenv(prefix + "_HOST"),
      Port: os.Getenv(prefix + "_PORT"),
   }
}

func LoadTopicsConfig(path string) TopicConfigs {
   data, err := os.ReadFile(path)
   if err != nil {
      panic(err)
   }

   var config TopicConfigs
   err = yaml.Unmarshal(data, &config)
   if err != nil {
      panic(err)
   }
   return config
}

func (cfg *TopicConfig) GetConfigMap() map[string]string {
   config := make(map[string]string)
   setField(config, "cleanup.policy", cfg.CleanupPolicy)
   setField(config, "compression.type", cfg.CompressionType)
   setField(config, "retention.ms", cfg.RetentionMs)
   return config
}

func setField(m map[string]string, key, value string) {
   if value != "" {
      m[key] = value
   }
}

//
// var cfg = setup.ServiceConfig{
//    Port: port,
//    KafkaBrokers: constants.KAFKA_BROKER_URL,
//    ProducerConfig: &kafka.ConfigMap{
//       "bootstrap.servers": "localhost:9094",
//    },
//    ConsumerGroupConfig: []setup.ConsumerGroupConfig{
//       {
//          Count: 3,
//          Topics: []string{constants.KAFKA_TOPIC_LOGS},
//          Config: &kafka.ConfigMap{
//             "bootstrap.servers": constants.KAFKA_BROKER_URL,
//             "group.id":          constants.KAFKA_TOPIC_LOGS,
//             "auto.offset.reset": "earliest",
//             "enable.auto.commit": false,
//             "enable.auto.offset.store": true, // auto in-mem offset update
//          },
//          Processor: logProcessor,
//       },
//       {
//          Count: 3,
//          Topics: []string{constants.KAFKA_TOPIC_LOGS_DLQ},
//          Config: &kafka.ConfigMap{
//             "bootstrap.servers": constants.KAFKA_BROKER_URL,
//             "group.id":          constants.KAFKA_TOPIC_LOGS_DLQ,
//             "auto.offset.reset": "earliest",
//             "enable.auto.commit": false,
//             "enable.auto.offset.store": true, // auto in-mem offset update
//          },
//          Processor: logDLQProcessor,
//       },
//       {
//          Count: 3,
//          Topics: []string{constants.KAFKA_TOPIC_LOGS_DLQ_PERMANENT},
//          Config: &kafka.ConfigMap{
//             "bootstrap.servers": constants.KAFKA_BROKER_URL,
//             "group.id":          constants.KAFKA_TOPIC_LOGS_DLQ_PERMANENT,
//             "auto.offset.reset": "earliest",
//             "enable.auto.commit": false,
//             "enable.auto.offset.store": true, // auto in-mem offset update
//          },
//          Processor: logDLQPermanentProcessor,
//       },
//
//    },
// }

