package constants
// Package constants defines constants and types used throughout the application.

const (
   //HTTP Server
   HTTP_SERVER_PORT = ":8080"

   //Kafka Topic
   KAFKA_BROKER_URL = "localhost:9094"
   KAFKA_TOPIC_LOGS = "logs"
   KAFKA_TOPIC_LOGS_DLQ = "logs-dlq"
   KAFKA_TOPIC_LOGS_DLQ_PERMANENT = "logs-dlq-permanent"

   //ES
   ES_INDEX_LOG = "logs"
   ES_PORT = "9200"

   // Redis
   REDIS_PORT = "6379"
)
