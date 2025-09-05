package es_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndex(t *testing.T) {
   indexService := esService.Index
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
   indexService := esService.Index
	if indexService == nil {
		t.Fatal("Elasticsearch client is not initialized")
	}

	indexName := "test_index"

	err := indexService.DeleteIndex(indexName)

	assert.Equal(t, err, nil)
}

