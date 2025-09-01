package main

import (
   "os"

   "github.com/khaiphan29/logpulse/pkg/logger"
   "github.com/khaiphan29/logpulse/internal/es/client"
   "github.com/khaiphan29/logpulse/internal/es/index"
   "github.com/khaiphan29/logpulse/internal/constants"
)

func main() {
   // Define index name and mapping
	indexName := constants.ES_INDEX_LOG
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
   // Initialize the Elasticsearch client
   client, err := esclient.NewClient(constants.ES_PORT)
   if err != nil {
      logger.Error("Failed to create index", map[string]any{
         "error": err,
      })
   }

   esIndex := esindex.New(client)
   // Create the index
   err = esIndex.CreateIndex(indexName, []byte(mapping))
   if err != nil {
      logger.Error("Failed to create index", map[string]any{
         "error": err,
      })
      os.Exit(1)
   }
}

