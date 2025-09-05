package main

import (
   "github.com/khaiphan29/logpulse/internal/config/es"
   "github.com/khaiphan29/logpulse/pkg/logger"
   "github.com/khaiphan29/logpulse/internal/setup"
)

func main() {
   // Load Elasticsearch configuration
   esService := setup.SetUpESSerive()
   indexCfg := esconfig.LoadIndexConfig("./internal/config/es/log_mapping.json")

   // Create the index
   err := esService.Index.CreateIndex(indexCfg.Name, []byte(indexCfg.Mapping))
   if err != nil {
      logger.Error("Failed to create index", map[string]any{
         "error": err,
      })
   } else {
      logger.Info("Index created successfully", map[string]any{
         "index": indexCfg.Name,
      })
   }
}

