package setup

import (
   "os"
   "os/signal"
   "syscall"
   "context"
   "net/http"
   "time"

   "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/khaiphan29/logpulse/internal/api/router"
   "github.com/khaiphan29/logpulse/internal/api/handlers"
   mykafka "github.com/khaiphan29/logpulse/internal/kafka"
	"github.com/khaiphan29/logpulse/pkg/logger"
   "github.com/khaiphan29/logpulse/internal/alert"
)

type ConsumerGroupConfig struct {
   Count int
   Processor mykafka.MessageProcessor
   Topics []string
   Config *kafka.ConfigMap
}

type ServiceConfig struct {
   Port string
   KafkaBrokers string
   ProducerConfig *kafka.ConfigMap
   ConsumerGroupConfig []ConsumerGroupConfig
}

func InitService(cfg *ServiceConfig) {
   // Initialize Kafka producer
   producer, err := mykafka.InitProducer(cfg.ProducerConfig)
   if err != nil {
      logger.Fatal("Failed to initialize Kafka producer", map[string]any{
         "error": err,
      })
   }
   defer producer.Shutdown()

   // Initialize Kafka consumer groups
   consumerCtx, cancel := context.WithCancel(context.Background())
   defer cancel()
   for _, consumerGroup := range cfg.ConsumerGroupConfig {
      for i := 0; i < consumerGroup.Count; i++ {
         consumer, err := mykafka.NewConsumer(consumerGroup.Config, consumerGroup.Processor)
         if err != nil {
            logger.Fatal("Failed to initialize Kafka consumer", map[string]any{
               "error": err,
            })
         }

         // Subscribe to the topics
         if err := consumer.SubscribeTopics(consumerGroup.Topics, nil); err != nil {
            logger.Fatal("Failed to subscribe to topics", map[string]any{
               "error": err,
            })
         } else {
            topics, _ := consumer.Subscription()
            logger.Info("Subscribed to topics", map[string]any{
               "topics": topics,
            })
         }

         go consumer.ListenForMessages(-1, consumerCtx)
      }
   }

   // Initialize handlers
   logHandler := handlers.NewHandler(producer)

   // Initialize HTTP server
   r := router.NewRouter("release", logHandler)
   s := &http.Server{
      Addr: cfg.Port,
      Handler: r,
   }

   go func() {
      logger.Info("Server started on port", map[string]any{
         "port": cfg.Port,
      })
      if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
         logger.Fatal("Listen: %s\n", map[string]any{
            "error": err,
         })
      }
   }()
   defer shutdownHTTPServer(s)

   // Initialize alerting system
   alertCtx, alertCancel := context.WithCancel(context.Background())
   go func() {
      logger.Info("Starting alerting system", nil)
      for {
         select {
         case <-alertCtx.Done():
            logger.Info("Stopping alerting system", nil)
            return
         default:
            alert.ErrorAlert()
         }
      time.Sleep(5 * time.Second)
      }
   }()
   defer alertCancel()

   // Set up a graceful shutdown
   quit := make(chan os.Signal, 1)
   // Wait for a signal to shut down
   // the process will not terminate immediately when the signal is received.
   signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
   <-quit
   logger.Info("Shutting down server...", nil)
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
