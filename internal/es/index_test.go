package esindex_test

import (
	"testing"

	"github.com/khaiphan29/logpulse/internal/constants"
	"github.com/khaiphan29/logpulse/internal/es/client"
	"github.com/khaiphan29/logpulse/internal/es/index"
	"github.com/khaiphan29/logpulse/pkg/logger"
	"github.com/stretchr/testify/assert"
)

var c *esclient.ESCLient

func init() {
   var err error
   c, err = esclient.NewClient(constants.ES_PORT)
   if err != nil {
      logger.Fatal("Failed to initialize Elasticsearch client", map[string]any{
         "error": err,
      })
   }
}

func TestCreateIndex(t *testing.T) {
   indexService := esindex.New(c)
	if indexService == nil {
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

	err := indexService.CreateIndex(indexName, []byte(mapping))

	assert.Equal(t, nil, err, "Index creation should not return an error")
}

func TestDeleteIndex(t *testing.T) {
   indexService := esindex.New(c)
	if indexService == nil {
		t.Fatal("Elasticsearch client is not initialized")
	}

	indexName := "test_index"

	err := indexService.DeleteIndex(indexName)

	assert.Equal(t, err, nil)
}

