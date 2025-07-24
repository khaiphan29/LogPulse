package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v9"

	"github.com/khaiphan29/logpulse/pkg/logger"
)

type Client struct {
   *elasticsearch.Client
   port string
}

var client *Client

func InitClient(port string) *Client {
   if client != nil {
      logger.Error("Elasticsearch client is already initialized", nil)
      return client
   }

   esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"http://localhost:" + port,
		},
	})
	if err != nil {
      logger.Fatal("Error creating Elasticsearch client", map[string]any{
         "error": err,
      })
	} else {
      logger.Info("Elasticsearch client created successfully", nil)
   }

   client = &Client{
      Client: esClient,
      port:   port,
   }

   return client
}

func GetClient() *Client {
   if client == nil {
      logger.Error("Elasticsearch client is not initialized", nil)
      return nil
   }
   return client
}

func (client *Client) ChangePort(port string) {
   if client == nil {
      logger.Error("Elasticsearch client is not initialized", nil)
      return
   }
   client.port = port
   logger.Info("Elasticsearch client port changed", map[string]any{
      "port": port,
   })
}

func (client *Client) GetPort() string {
   if client == nil {
      logger.Error("Elasticsearch client is not initialized", nil)
      return ""
   }
   return client.port
}
