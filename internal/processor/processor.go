package processor

import (
   "fmt"

   "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type LogProcessor struct {}

func (lp *LogProcessor) Process(message *kafka.Message) error {
   // Process the message
   fmt.Printf("Received message: %s\n", string(message.Value))
   return nil
}

type LogDLQProcessor struct {}

func (lp *LogDLQProcessor) Process(message *kafka.Message) error {
   // Process the message
   fmt.Printf("Received message: %s\n", string(message.Value))
   return nil
}

type LogDLQPermanentProcessor struct {}

func (lp *LogDLQPermanentProcessor) Process(message *kafka.Message) error {
   // Process the message
   fmt.Printf("Received message: %s\n", string(message.Value))
   return nil
}

