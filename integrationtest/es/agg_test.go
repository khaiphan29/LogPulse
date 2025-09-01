package es_test

import (
   "fmt"
   "time"
   "testing"
   "github.com/stretchr/testify/assert"

   "github.com/khaiphan29/logpulse/internal/es/agg"
   "github.com/khaiphan29/logpulse/internal/es/doc"
   "github.com/khaiphan29/logpulse/pkg/logger"

)

func setup() {
   esDoc := esdoc.New(esClient)

   // Create records
   for i := 1; i <= 10; i++ {
      document := map[string]any{
         "logId":       fmt.Sprintf("log_%d", i),
         "timestamp":   time.Now().Format(time.RFC3339),
         "logLevel":    "INFO",
         "message":     "This is a test log message",
         "metadata":    map[string]string{"key": "value"},
         "source":      "unit_test",
         "environment": "test",
         "type":        "application_log",
      }
      err := esDoc.CreateDocument(TEST_INDEX, document)
      if err != nil {
         logger.Fatal("Failed to create document", map[string]any{
            "error": err,
         })
      }
   }
}

func TestCountTotalLogsByLevel(t *testing.T) {
   setup()
   lte := time.Now()
   gte := lte.Add(-24 * time.Hour)

   esAgg := esagg.New(esClient)
   count, err := esAgg.CountTotalLogsByLevel(TEST_INDEX, "INFO", gte, lte)
   assert.Equal(t, nil, err, "CountTotalLogsByLevel should not return an error")
   assert.Equal(t, 10, count, "There should be 10 logs with level INFO")
}
