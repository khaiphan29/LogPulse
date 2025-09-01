package processor

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/khaiphan29/logpulse/internal/api/parsing"
   "github.com/khaiphan29/logpulse/internal/constants"
)

type LogProcessor struct {
   LogConsumer LogConsumer
}

type LogConsumer interface {
   CreateDocument(index string, document any) error
}

func (lp *LogProcessor) Process(message *kafka.Message) error {
   // Process the message
   fmt.Printf("Received message: %s\n", string(message.Value))

   var log parser.LogPayload
   if err := json.Unmarshal(message.Value, &log); err != nil {
      return err
   }

   if err := lp.LogConsumer.CreateDocument(constants.ES_INDEX_LOG, log); err != nil {
      return err
   }

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

