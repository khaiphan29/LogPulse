package esclient

import (
	"github.com/elastic/go-elasticsearch/v9"
)

type ESCLient = elasticsearch.Client

func NewClient(port string) (*ESCLient, error) {
   esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"http://localhost:" + port,
		},
	})

   if err != nil {
      return nil, err
   }
   return esClient, nil
}

