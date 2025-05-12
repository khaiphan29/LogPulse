package constants
// Package constants defines constants and types used throughout the application.

const (
   //HTTP Server
   HTTPServerPort = ":8080"

   //Kafka Topic
   KafkaBrokers = "localhost:9094"
   KafkaTopicLogs = "logs"
   KafkaTopicLogsDLQ = "logs-dlq"
   KafkaTopicLogsDLQPermanent = "logs-dlq-permanent"
)
