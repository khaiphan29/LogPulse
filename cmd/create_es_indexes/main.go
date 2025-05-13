package main

import (
   es "github.com/khaiphan29/logpulse/internal/elasticsearch"
   "github.com/khaiphan29/logpulse/internal/constants"
)

func main() {
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

   es.CreateIndex(indexName, []byte(mapping))
}

