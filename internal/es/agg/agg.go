package esagg

import (
	"bytes"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
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

// CountLogsByLevelAndSource counts logs grouped by log level and source
// ExecuteQuery takes an ElasticSearch query as input and executes it
func (client *Client) executeQuery(indexName string, query string) (*esapi.Response, error) {
	// Execute the search query
	res, err := client.Search(
		client.Search.WithIndex(indexName),
		client.Search.WithBody(bytes.NewReader([]byte(query))),
	)
	if err != nil {
      return nil, err
	}

   return res, nil
}

func (client *Client) CountTotalLogsByLevel(index, level string, timegte, timelte time.Time) (int, error) {
   opts, err := DefaultOptions().
      WithTimeRange(timegte, timelte).
      WithTerm("logLevel", level).
      WithAggs(TOTAL_COUNT_QNAME, map[string]any{
         "value_count": map[string]any{
            "field": "logLevel",
         },
      }).
      ToJSON()
   if err != nil {
      return 0, err
   }

   logger.Info("Counting total logs by level", map[string]any{
      "options": opts,
   })

   res, err := client.executeQuery(index, opts)
   if err != nil {
      return 0, err
   }
   defer res.Body.Close()
   return GetTotalCount(res)
}
