package esclient

import (
	"github.com/elastic/go-elasticsearch/v9"
)

type ESCLient struct {
   *elasticsearch.Client
}

func NewClient(host, port string) (*ESCLient, error) {
   esClient, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{
			"http://" + host + ":" + port,
		},
	})

   if err != nil {
      return nil, err
   }
   return &ESCLient { esClient }, nil
}
