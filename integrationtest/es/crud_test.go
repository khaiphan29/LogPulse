package es_test

import (
	"testing"

	"github.com/khaiphan29/logpulse/internal/es/doc"
	"github.com/stretchr/testify/assert"
)

func TestCreateDocument(t *testing.T) {
   esdocService := esdoc.New(esClient)
	if esdocService == nil {
		t.Fatal("Elasticsearch client is not initialized")
	}

	document := map[string]any{
		"logId":       "12345",
		"timestamp":   "2023-10-01T12:00:00Z",
		"logLevel":    "INFO",
		"message":     "This is a test log message",
		"metadata":    map[string]string{"key": "value"},
		"source":      "unit_test",
		"environment": "test",
		"type":        "application_log",
	}

	err := esdocService.CreateDocument(TEST_INDEX, document)
	assert.Equal(t, nil, err, "Sending document to index should not return an error")
}

func DeleteDocumentByQuery(t *testing.T) {
   esdocService := esdoc.New(esClient)
	if esdocService == nil {
		t.Fatal("Elasticsearch client is not initialized")
	}

   query := map[string]any{
      "query": map[string]any{
         "match": map[string]any{
            "logId": "12345",
         },
      },
   }

   err := esdocService.DeleteDocumentByQuery(TEST_INDEX, query)
   assert.Equal(t, nil, err, "Deleting document by query should not return an error")
}

