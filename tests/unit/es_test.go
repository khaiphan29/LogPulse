package unit_test

import (
	"testing"

	es "github.com/khaiphan29/logpulse/internal/elasticsearch"
	"github.com/khaiphan29/logpulse/pkg/logger"
   "github.com/stretchr/testify/assert"
   "github.com/khaiphan29/logpulse/internal/constants"
)

func init() {
   // Initialize the Elasticsearch client
   esClient := es.InitClient(constants.ES_PORT)
   if esClient == nil {
      logger.Fatal("Failed to initialize Elasticsearch client", nil)
   }
}

func TestCreateIndex(t *testing.T) {
   esClient := es.GetClient()
   if esClient == nil {
      t.Fatal("Elasticsearch client is not initialized")
   }

   indexName := "test_index"
   mapping := `{
      "mappings": {
         "properties": {
            "logId": { "type": "keyword" },
            "timestamp": { "type": "date" },
            "logLevel": { "type": "keyword" },
            "message": { "type": "text", "analyzer": "standard" },
            "metadata": { "type": "object" },
            "source": { "type": "keyword" },
            "environment": { "type": "keyword" },
            "type": { "type": "keyword" }
         }
      }
   }`

   err := esClient.CreateIndex(indexName, []byte(mapping))

   assert.Equal(t, nil, err, "Index creation should not return an error")
}

func TestSendToIndex(t *testing.T) {
   esClient := es.GetClient()
   if esClient == nil {
      t.Fatal("Elasticsearch client is not initialized")
   }

   indexName := "test_index"
   document := map[string]any{
      "logId":      "12345",
      "timestamp":  "2023-10-01T12:00:00Z",
      "logLevel":   "INFO",
      "message":    "This is a test log message",
      "metadata":   map[string]string{"key": "value"},
      "source":     "unit_test",
      "environment": "test",
      "type":       "application_log",
   }

   err := esClient.SendToIndex(indexName, document)

   assert.Equal(t, nil, err, "Sending document to index should not return an error")
}

func TestExecuteQuery(t *testing.T) {
   esClient := es.GetClient()
   if esClient == nil {
      t.Fatal("Elasticsearch client is not initialized")
   }

   indexName := "test_index"
   query := `{
      "query": {
         "match_all": {}
      }
   }`

   res, err := esClient.ExecuteQuery(indexName, query)
   if err != nil {
      t.Fatalf("Failed to execute query: %v", err)
   }
   defer res.Body.Close()

   assert.NotNil(t, res, "Response should not be nil")
}

func TestDeleteIndex(t *testing.T) {
   esClient := es.GetClient()
   if esClient == nil {
      t.Fatal("Elasticsearch client is not initialized")
   }

   indexName := "test_index"

   err := esClient.DeleteIndex(indexName)

   assert.Equal(t, err, nil)
}
