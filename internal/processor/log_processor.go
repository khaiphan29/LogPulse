package processor

import (
	"encoding/json"
	"fmt"

	"github.com/khaiphan29/logpulse/internal/api/parsing"
	"github.com/khaiphan29/logpulse/internal/constants"
)

type LogProcessor struct {
   IndexName string
   RetryTopic  string
   Storage LogStorage
   Retry LogRetry
   Analyzer LogAnalyzer
}

type LogStorage interface {
   CreateDocument(index string, document []byte) error
}

type LogRetry interface {
   SendMessage(topic string, key, value []byte) error
}

type LogAnalyzer interface {
   AnalyzeError()
}

func NewLogProcessor(indexName, retryTopic string, storage LogStorage, retry LogRetry, analyzer LogAnalyzer) *LogProcessor {
   return &LogProcessor{
      IndexName: indexName,
      RetryTopic: retryTopic,
      Storage: storage,
      Retry: retry,
      Analyzer: analyzer,
   }
}

func (lp *LogProcessor) Process(msg []byte) error {
   // Check if the message is valid JSON
   var log parser.LogPayload
   if err := json.Unmarshal(msg, &log); err != nil {
      return err
   }

   if err := lp.Storage.CreateDocument(lp.IndexName, msg); err != nil {
      // If there's an error storing the log, send it to the retry topic
      key := []byte(log.Source)
      if sendErr := lp.Retry.SendMessage(lp.RetryTopic, key, msg); sendErr != nil {
         return fmt.Errorf("failed to store log: %v, also failed to send to retry: %v", err, sendErr)
      }
   }

   if log.LogLevel == constants.LOG_LEVEL_ERROR {
      lp.Analyzer.AnalyzeError()
   }

   return nil
}

