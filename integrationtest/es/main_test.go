package es_test

import (
	"os"
	"testing"

	"github.com/khaiphan29/logpulse/internal/config/es"
	"github.com/khaiphan29/logpulse/internal/setup"
	"github.com/khaiphan29/logpulse/pkg/logger"
)

const (
   TEST_INDEX = "test_logs"
)

var (
   esService *setup.ESServices
)

func TestMain(m *testing.M) {
   // Initialize the Elasticsearch client
   esService = setup.SetUpESSerive()
   esIndex := esService.Index

   indexConfig := esconfig.LoadIndexConfig("test_log_mapping.json")
   mapping := indexConfig.Mapping

   err := esIndex.CreateIndex(TEST_INDEX, []byte(mapping))
   if err != nil {
      logger.Fatal("Failed to create index", map[string]any{
         "error": err,
      })
   }

   // Run tests
   m.Run()

   // Cleanup: Delete the index after tests
   err = esIndex.DeleteIndex(TEST_INDEX)
   if err != nil {
      logger.Error("Failed to delete index", map[string]any{
         "error": err,
      })
   }

   os.Exit(0)
}
