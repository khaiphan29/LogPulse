package alert

import (
   "fmt"
   "encoding/json"

   "github.com/khaiphan29/logpulse/pkg/logger"
   "github.com/khaiphan29/logpulse/internal/constants"
   es "github.com/khaiphan29/logpulse/internal/elasticsearch"
)

func ErrorAlert() {
   // Define the time interval and threshold
   timeInterval := "7d"
   threshold := 2

   // Call the function to get error logs by sources
   results, err := GetErrorBySources(timeInterval, threshold)
   if err != nil {
      logger.Error("Error getting error logs by sources", map[string]any{
         "error": err,
      })
      return
   }

   // Process the results
   for _, result := range results {
      for source, count := range result {
         logger.Info("Error alert", map[string]any{
            "source": source,
            "count":  count,
         })
      }
   }
}

func GetErrorBySources(timeInterval string, threshold int) ([]map[string]int, error) {
	query := fmt.Sprintf(`{
		"size": 0,
		"query": {
			"range": {
				"timestamp": {
					"gte": "now-%s",
					"lte": "now"
				}
			}
		},
		"aggs": {
			"by_source": {
				"terms": {
					"field": "source",
					"size": 100
				},
				"aggs": {
					"by_log_level": {
						"terms": {
							"field": "logLevel",
							"include": "ERROR"
						}
					}
				}
			}
		}
	}`, timeInterval)

	indexName := constants.ESIndexLogs

	// Execute the query
	res, err := es.ExecuteQuery(indexName, query)
	if err != nil {
		logger.Error("Error executing query", map[string]any{
			"index": indexName,
			"error": err,
		})
		return nil, err
	}
	defer res.Body.Close()

	// Check if the response contains errors
	if res.IsError() {
		logger.Error("Error response from ElasticSearch", map[string]any{
			"index": indexName,
			"error": res.String(),
		})
		return nil, fmt.Errorf("error response from ElasticSearch: %s", res.String())
	} else {
      logger.Info("Query executed successfully", map[string]any{
         "index": indexName,
         "status": res.Status(),
         "body": res.String(),
      })
   }

	// Decode the response
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		logger.Error("Error decoding response body", map[string]any{
			"index": indexName,
			"error": err,
		})
		return nil, err
	} else {
      logger.Info("Response decoded successfully", map[string]any{
         "result": result,
      })
   }

   // Extract the "aggregations" field
   aggregations := result["aggregations"].(map[string]interface{})
   logger.Info("Aggregations", map[string]any{
      "aggregations": aggregations,
   })

   // Extract the "by_source" aggregation
   bySourceAgg := aggregations["by_source"].(map[string]interface{})
   logger.Info("bySourceAgg", map[string]any{
      "aggregations": bySourceAgg,
   })


   // Extract the "buckets" array from the "by_source" aggregation
   buckets := bySourceAgg["buckets"].([]interface{})
   logger.Info("buckets of source", map[string]any{
      "aggregations": buckets,
   })

	var filteredResults []map[string]int

	for _, bucket := range buckets {
		sourceBucket := bucket.(map[string]interface{})
		source := sourceBucket["key"].(string)
		errorLogs := 0

		// Check if there are error logs for this source
		if logLevelAgg, exists := sourceBucket["by_log_level"].(map[string]interface{}); exists {
			for _, logLevelBucket := range logLevelAgg["buckets"].([]interface{}) {
				logLevelData := logLevelBucket.(map[string]interface{})
				if logLevelData["key"].(string) == "ERROR" {
					errorLogs = int(logLevelData["doc_count"].(float64))
				}
			}
		}

      logger.Info("Error logs count", map[string]any{
         "source": source,
         "count":  errorLogs,
      })

		// Add to results if error count exceeds threshold
		if errorLogs > threshold {
			filteredResults = append(filteredResults, map[string]int{
				source: errorLogs,
			})
		}
	}

	logger.Info("Query executed successfully", map[string]any{
		"index": indexName,
		"results": filteredResults,
	})

	return filteredResults, nil
}
