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

