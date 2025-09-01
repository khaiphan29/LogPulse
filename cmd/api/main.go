package main

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/khaiphan29/logpulse/internal/constants"
	"github.com/khaiphan29/logpulse/internal/setup"
   "github.com/khaiphan29/logpulse/internal/processor"
)

func main() {
   // Start the server
   port := constants.HTTP_SERVER_PORT
   if (len(os.Args) > 1) {
      port = ":" + os.Args[1]
   }

   // Initialize Processors
   logProcessor := &processor.LogProcessor{}
   logDLQProcessor := &processor.LogDLQProcessor{}
   logDLQPermanentProcessor := &processor.LogDLQPermanentProcessor{}

   cfg := setup.ServiceConfig{
      Port: port,
      KafkaBrokers: constants.KAFKA_BROKER_URL,
      ProducerConfig: &kafka.ConfigMap{
         "bootstrap.servers": "localhost:9094",
      },
      ConsumerGroupConfig: []setup.ConsumerGroupConfig{
         {
            Count: 3,
            Topics: []string{constants.KAFKA_TOPIC_LOGS},
            Config: &kafka.ConfigMap{
               "bootstrap.servers": constants.KAFKA_BROKER_URL,
               "group.id":          constants.KAFKA_TOPIC_LOGS,
               "auto.offset.reset": "earliest",
               "enable.auto.commit": false,
               "enable.auto.offset.store": true, // auto in-mem offset update
            },
            Processor: logProcessor,
         },
                  {
            Count: 3,
            Topics: []string{constants.KAFKA_TOPIC_LOGS_DLQ},
            Config: &kafka.ConfigMap{
               "bootstrap.servers": constants.KAFKA_BROKER_URL,
               "group.id":          constants.KAFKA_TOPIC_LOGS_DLQ,
               "auto.offset.reset": "earliest",
               "enable.auto.commit": false,
               "enable.auto.offset.store": true, // auto in-mem offset update
            },
            Processor: logDLQProcessor,
         },
         {
            Count: 3,
            Topics: []string{constants.KAFKA_TOPIC_LOGS_DLQ_PERMANENT},
            Config: &kafka.ConfigMap{
               "bootstrap.servers": constants.KAFKA_BROKER_URL,
               "group.id":          constants.KAFKA_TOPIC_LOGS_DLQ_PERMANENT,
               "auto.offset.reset": "earliest",
               "enable.auto.commit": false,
               "enable.auto.offset.store": true, // auto in-mem offset update
            },
            Processor: logDLQPermanentProcessor,
         },

      },
   }

   setup.InitService(&cfg)
	os.Exit(1) // Ensure the program exits with a zero status code
}
