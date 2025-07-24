package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"

   "github.com/elastic/go-elasticsearch/v9/esapi"

	"github.com/khaiphan29/logpulse/pkg/logger"
)

// CreateIndex creates an ElasticSearch index with the specified mapping
func (client *Client) CreateIndex(indexName string, mapping []byte) {
   esClient := client.esClient
	res, err := esClient.Indices.Create(indexName, esClient.Indices.Create.WithBody(bytes.NewReader(mapping)))
	if err != nil {
      logger.Fatal("Error creating index", map[string]any{
         "index": indexName,
         "error": err,
      })
	}
	defer res.Body.Close()
	if res.IsError() {
      logger.Error("Error response from ElasticSearch", map[string]any{
         "error": res.String(),
      })
	} else {
      logger.Info("Index created successfully", map[string]any{
         "index": indexName,
      })
	}
}

// SendToIndex sends a document to the specified index in ElasticSearch
func (client *Client) SendToIndex(indexName string, document interface{}) error {
   esClient := client.esClient
	// Marshal the document into JSON
	payloadJSON, err := json.Marshal(document)
	if err != nil {
      return err
	}

	// Send the document to the ElasticSearch index
	res, err := esClient.Index(
		indexName,
		bytes.NewReader(payloadJSON),
		esClient.Index.WithRefresh("true"), // Immediately make it available for search
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
func (client *Client) ExecuteQuery(indexName string, query string) (*esapi.Response, error) {
   esClient := client.esClient
	// Execute the search query
	res, err := esClient.Search(
		esClient.Search.WithIndex(indexName),
		esClient.Search.WithBody(bytes.NewReader([]byte(query))),
	)
	if err != nil {
      return nil, err
	}

   return res, nil
}
