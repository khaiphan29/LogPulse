package constants
// Package constants defines constants and types used throughout the application.

const (
   LOG_LEVEL_ERROR = "ERROR"

   KAFKA_LOGS_TOPIC_KEY = "logs"
   KAFKA_LOG_RETRY_TOPIC_KEY = "logs_dlq"
   KAFKA_TOPICS_CONFIG_PATH = "./internal/config/kafka/topics.yaml"

   ES1_PREFIX = "ES"
   ES_LOGS_INDEX_CFG_PATH = "./internal/config/es/log_mapping.json"
)
