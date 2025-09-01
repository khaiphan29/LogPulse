package esagg_test

import (
   "time"
   "testing"
   "encoding/json"
   "github.com/stretchr/testify/assert"

   "github.com/khaiphan29/logpulse/internal/es/agg"
)

func TestGetTotalError(t *testing.T) {
   gte := time.Now().Add(-7 * 24 * time.Hour) // 7 days ago
   lte := time.Now()

   opts := esagg.DefaultOptions().
      WithSize(0).
      WithTerm("logLevel", "error").
      WithTimeRange(gte, lte)

   total_errors := map[string]any{
      "value_count": map[string]any{
         "field": "logLevel",
      },
   }
   errors_by_source := map[string]any{
      "terms": map[string]any{
         "field": "source",
         "size": 20,
      },
   }

   opts.WithAggs("total_errors", total_errors).
      WithAggs("errors_by_source", errors_by_source)
   jsonStr, err := opts.ToJSON()
   assert.NoError(t, err, "Converting to JSON should not return an error")

   expected := `{
      "size": 0,
      "query": {
         "bool": {
            "must": [
               { "term": { "logLevel": "error" } },
               { "range": { "timestamp": { "gte": "` + gte.Format(time.RFC3339) + `", "lte": "` + lte.Format(time.RFC3339) + `" } } }
            ]
         }
      },
      "aggs": {
         "total_errors": {
            "value_count": { "field": "logLevel" }
         },
         "errors_by_source": {
            "terms": { "field": "source", "size": 20 }
         }
      }
   }`

   var expectedMap, actualMap map[string]any
   err = json.Unmarshal([]byte(expected), &expectedMap)
   assert.NoError(t, err, "Unmarshal of expected JSON failed")

   err = json.Unmarshal([]byte(jsonStr), &actualMap)
   assert.NoError(t, err, "Unmarshal of actual JSON failed")

   assert.Equal(t, expectedMap, actualMap, "Generated JSON does not match expected JSON")
}
