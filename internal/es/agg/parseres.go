package esagg

import (
   "fmt"
   "encoding/json"
   "github.com/elastic/go-elasticsearch/v9/esapi"
)

type CountAgg struct {
   Value int `json:"value"`
}

type CountAggResponse struct {
   Aggregations map[string]CountAgg `json:"aggregations"`
}

func GetTotalCount(res *esapi.Response) (int, error) {
   if res.IsError() {
      return 0, fmt.Errorf("error response from ElasticSearch: %s", res.String())
   }

   var out CountAggResponse
   if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
      return 0, fmt.Errorf("error parsing the response body: %w", err)
   }
   total := int(out.Aggregations[TOTAL_COUNT_QNAME].Value)

   return total, nil
}
