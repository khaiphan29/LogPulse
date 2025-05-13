package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v9"
   "github.com/elastic/go-elasticsearch/v9/esapi"

	"github.com/khaiphan29/logpulse/pkg/logger"
)

var es *elasticsearch.Client

func init() {
   if es != nil {
      return
   }
	var err error
	es, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"http://localhost:9200", // Replace with your custom address and port
		},
	})
	if err != nil {
      logger.Fatal("Error creating Elasticsearch client", map[string]any{
         "error": err,
      })
	} else {
      logger.Info("Elasticsearch client created successfully", nil)
   }
}

// CreateIndex creates an ElasticSearch index with the specified mapping
func CreateIndex(indexName string, mapping []byte) {
	res, err := es.Indices.Create(indexName, es.Indices.Create.WithBody(bytes.NewReader(mapping)))
	if err != nil {
      logger.Fatal("Error creating index", map[string]any{
         "index": indexName,
         "error": err,
      })
	}
	defer res.Body.Close()
	if res.IsError() {
      logger.Fatal("Error response from ElasticSearch", map[string]any{
         "error": res.String(),
      })
	} else {
      logger.Info("Index created successfully", map[string]any{
         "index": indexName,
      })
	}
}

// SendToIndex sends a document to the specified index in ElasticSearch
func SendToIndex(indexName string, document interface{}) error {
	// Marshal the document into JSON
	payloadJSON, err := json.Marshal(document)
	if err != nil {
      return err
	}

	// Send the document to the ElasticSearch index
	res, err := es.Index(
		indexName,
		bytes.NewReader(payloadJSON),
		es.Index.WithRefresh("true"), // Immediately make it available for search
	)
	if err != nil {
      return err
	}
	defer res.Body.Close()

	// Check for response errors
	if res.IsError() {
      return fmt.Errorf("error response from ElasticSearch: %s", res.String())
	} else {
      logger.Info("Document indexed successfully", map[string]any{
         "status": res.Status(),
      })
   }
   return nil
}

// CountLogsByLevelAndSource counts logs grouped by log level and source
// ExecuteQuery takes an ElasticSearch query as input and executes it
func ExecuteQuery(indexName string, query string) (*esapi.Response, error) {
	// Execute the search query
	res, err := es.Search(
		es.Search.WithIndex(indexName),
		es.Search.WithBody(bytes.NewReader([]byte(query))),
	)
	if err != nil {
      return nil, err
	}

   return res, nil
}
