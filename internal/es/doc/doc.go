package esdoc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/khaiphan29/logpulse/pkg/logger"
)

type Client struct {
   *elasticsearch.Client
}

func New(esClient *elasticsearch.Client) *Client {
   return &Client{
      Client: esClient,
   }
}

// Send a document to the specified index in ElasticSearch
func (client *Client) CreateDocument(index string, document []byte) error {
	// Send the document to the ElasticSearch index
	res, err := client.Index(
		index,
		bytes.NewReader(document),
		client.Index.WithRefresh("true"), // Immediately make it available for search
	)
	if err != nil {
      return err
	} else {
      logger.Debug("Document indexed successfully", map[string]any{
         "index": index,
         "response": res,
         "document": string(document),
      })
   }
	defer res.Body.Close()

	// Check for response errors
	if res.IsError() {
      return fmt.Errorf("error response from ElasticSearch: %s", res.String())
   }
   return nil
}

func (client *Client) DeleteDocumentByID(index string, documentID string) error {
   res, err := client.Delete(
      index,
      documentID,
      client.Delete.WithRefresh("true"), // The change is visible immediately for search
   )
   if err != nil {
      return err
   }
   defer res.Body.Close()

   // Check for response errors
   if res.IsError() {
      return fmt.Errorf("error response from ElasticSearch: %s", res.String())
   }
   return nil
}

func (client *Client) DeleteDocumentByQuery(index string, query any) error {
   // Marshal the query into JSON
   queryJSON, err := json.Marshal(query)
   if err != nil {
      return err
   }

   res, err := client.DeleteByQuery(
      []string{index},
      bytes.NewReader(queryJSON),
      client.DeleteByQuery.WithRefresh(true), // The change is visible immediately for search
   )
   if err != nil {
      return err
   }
   defer res.Body.Close()

   // Check for response errors
   if res.IsError() {
      return fmt.Errorf("error response from ElasticSearch: %s", res.String())
   }
   return nil
}
