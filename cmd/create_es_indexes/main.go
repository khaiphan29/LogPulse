package main

import (
   "os"

   "github.com/khaiphan29/logpulse/pkg/logger"
   es "github.com/khaiphan29/logpulse/internal/elasticsearch"
   "github.com/khaiphan29/logpulse/internal/constants"
)

func main() {
   if len(os.Args) < 2 {
      logger.Error("Arg: ES_PORT missing", nil)
      return
   }
   // Create the index
   // Define index name and mapping
	indexName := constants.ESIndexLogs
	mapping := `{
		"mappings": {
			"properties": {
				"logId": { "type": "keyword" },
				"timestamp": { "type": "date" },
				"logLevel": { "type": "keyword" },
				"message": { "type": "text", "analyzer": "standard" },
				"metadata": { "type": "object" },
				"source": { "type": "keyword" },
				"environment": { "type": "keyword" },
				"type": { "type": "keyword" }
			}
		}
	}`

   client := es.InitClient(os.Args[1])
   err := client.CreateIndex(indexName, []byte(mapping))
   if err != nil {
      logger.Error("Failed to create index", map[string]any{
         "error": err,
      })
   }
}

