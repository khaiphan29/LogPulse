package esagg

import (
	"encoding/json"
	"time"

	"github.com/khaiphan29/logpulse/internal/es/helper"
)

const (
   TOTAL_COUNT_QNAME = "total_count"
)

type AggOptions map[string]any

func DefaultOptions() AggOptions {
   var defaultOptions = map[string]any{
      "size": 0,
      "query": map[string]any{
         "bool": map[string]any{},
      },
      "aggs": map[string]any{},
   }

   return defaultOptions
}

func (a AggOptions) WithSize(size int) AggOptions {
   a["size"] = size
   return a
}

func (a AggOptions) WithMust(clause AggOptions) AggOptions {
   query := eshelper.EnsureMap(a, "query")
   boolQuery := query["bool"].(map[string]any)

   if existingMust, ok := boolQuery["must"]; ok {
      boolQuery["must"] = append(existingMust.([]map[string]any), clause)
   } else {
      boolQuery["must"] = []map[string]any{clause}
   }

   return a
}

func toESTime(t time.Time) string {
   return t.Format(time.RFC3339)
}

func (a AggOptions) WithTimeRange(gte, lte time.Time) AggOptions {
   rangeQuery := map[string]any{
      "range": map[string]any{
         "timestamp": map[string]any{
            "gte": toESTime(gte),
            "lte": toESTime(lte),
         },
      },
   }

   return a.WithMust(rangeQuery)
}

func (a AggOptions) WithTerm(field string, value any) AggOptions {
   termQuery := map[string]any{
      "term": map[string]any{
         field: value,
      },
   }

   return a.WithMust(termQuery)
}

func (a AggOptions) WithAggs (name string, opts map[string]any) AggOptions {
   aggs := eshelper.EnsureMap(a, "aggs")
   aggs[name] = opts
   return a
}

func (a AggOptions) ToJSON() (string, error) {
   jsonData, err := json.Marshal(a)
   if err != nil {
      return "", err
   }
   return string(jsonData), nil
}
