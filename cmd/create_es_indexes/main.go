package main

import (
	"github.com/khaiphan29/logpulse/internal/config/es"
	"github.com/khaiphan29/logpulse/internal/constants"
	"github.com/khaiphan29/logpulse/internal/setup"
	"github.com/khaiphan29/logpulse/pkg/logger"
)

func main() {
   // Load Elasticsearch configuration
   esService := setup.SetUpESSerive()
   indexCfg := esconfig.LoadIndexConfig(constants.ES_LOGS_INDEX_CFG_PATH)

   // Create the index
   err := esService.Index.CreateIndex(indexCfg.Name, indexCfg.Mapping)
   if err != nil {
      logger.Error("Failed to create index", map[string]any{
         "index": indexCfg.Name,
         "mapping": string(indexCfg.Mapping),
         "error": err,
      })
   } else {
      logger.Info("Index created successfully", map[string]any{
         "index": indexCfg.Name,
      })
   }
}

