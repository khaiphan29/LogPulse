package setup

import (
   "github.com/khaiphan29/logpulse/pkg/logger"

   "github.com/khaiphan29/logpulse/internal/es/client"
   "github.com/khaiphan29/logpulse/internal/es/index"
   "github.com/khaiphan29/logpulse/internal/es/doc"
   "github.com/khaiphan29/logpulse/internal/es/agg"

   "github.com/khaiphan29/logpulse/internal/config/es"
   "github.com/khaiphan29/logpulse/internal/constants"
)

type ESServices struct {
   Client *esclient.ESCLient
   Index  *esindex.Client
   Doc    *esdoc.Client
   Agg    *esagg.Client
}

func SetUpESSerive() (*ESServices) {
   adminCfg := esconfig.LoadClientConfig(constants.ES1_PREFIX)
   esClient, err := esclient.NewClient(adminCfg.Host, adminCfg.Port)
   if err != nil {
      logger.Fatal("Failed to connect to ES", map[string]any{
         "error": err,
      })
   }
   esIndex := esindex.New(esClient.Client)
   esDoc := esdoc.New(esClient.Client)
   esAgg := esagg.New(esClient.Client)
   return &ESServices{
      Client: esClient,
      Index:  esIndex,
      Doc:    esDoc,
      Agg:    esAgg,
   }
}

func GetESIndexCfg(path string) (*esconfig.IndexConfig) {
   return esconfig.LoadIndexConfig(path)
}
