package setup

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/khaiphan29/logpulse/internal/api/handlers"
	"github.com/khaiphan29/logpulse/internal/api/router"
	"github.com/khaiphan29/logpulse/internal/processor"
	"github.com/khaiphan29/logpulse/pkg/logger"

	"github.com/khaiphan29/logpulse/internal/constants"
)


func InitService() func() {
   // Load configuration
   esLogIndex := GetESIndexCfg(constants.ES_LOGS_INDEX_CFG_PATH)
   kafkaTopics := GetKafkaTopicsCfg()

   // Initialize Elasticsearch services
   esServices := SetUpESSerive()

   // Initialize Kafka producer
   producer, shutdownProducer := setUpKafkaProducer()

   // Cache Instance
   redisClient := setupRedisClient()

   // Initialize analyzer
   logErrAnalyzer := setupLogErrorAnalyzer(esServices.Agg, redisClient, esLogIndex.Name)

   // Initialize message processor
   logProcessor := processor.NewLogProcessor(esLogIndex.Name, kafkaTopics[constants.KAFKA_LOG_RETRY_TOPIC_KEY].Name, esServices.Doc, producer, logErrAnalyzer)

   // Initialize Kafka consumers
   _, cancelConsumers := setUpConsumerGroups(logProcessor)

   // Initialize handlers
   logHandler := handlers.NewLogHandler(kafkaTopics[constants.KAFKA_LOGS_TOPIC_KEY].Name, producer)
   handlers := []router.Registrar{
      logHandler,
   }

   // Initialize router
   r := router.New(handlers)

   // Start HTTP server
   s := &http.Server{
      Addr: os.Getenv("SERVER_PORT"),
      Handler: r,
   }

   go func() {
      logger.Info("Server started on", map[string]any{
         "Addr": s.Addr,
      })
      if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
         logger.Fatal("Listen: %s\n", map[string]any{
            "error": err,
         })
      }
   }()

   return func() {
      // Shutdown Kafka consumers
      cancelConsumers()
      // Shutdown Kafka producer
      shutdownProducer()
      // Shutdown HTTP server
      shutdownHTTPServer(s)
   }
}

func shutdownHTTPServer(s *http.Server) {
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()

   if err := s.Shutdown(ctx); err != nil {
      logger.Fatal("Server forced to shutdown:", map[string]any{
         "error": err,
      })
   } else {
      logger.Info("HTTP Server shutdown successfully", nil)
   }
}
