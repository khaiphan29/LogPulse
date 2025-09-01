package esindex

import (
	"bytes"
	"fmt"

   "github.com/elastic/go-elasticsearch/v9"
)

type Client struct {
   *elasticsearch.Client
}

func New(client *elasticsearch.Client) *Client {
   return &Client{
      Client: client,
   }
}

// CreateIndex creates an ElasticSearch index with the specified mapping
func (client *Client) CreateIndex(indexName string, mapping []byte) error {
	res, err := client.Indices.Create(indexName, client.Indices.Create.WithBody(bytes.NewReader(mapping)))
	if err != nil {
      return err
	}
	defer res.Body.Close()

	if res.IsError() {
      return fmt.Errorf("Error response from ElasticSearch: %s", res.String())
	}
   return nil
}

func (client *Client) DeleteIndex(indexName string) error {
   res, err := client.Indices.Delete([]string{indexName})
   if err != nil {
      return err
   }
   defer res.Body.Close()

   if res.IsError() {
      return fmt.Errorf("error deleting index %s: %s", indexName, res.String())
   }
   return nil
}

