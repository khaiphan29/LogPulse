package handlers

import (
   "errors"
   "encoding/json"
   "net/http"
   "github.com/gin-gonic/gin"

   "github.com/khaiphan29/logpulse/internal/api/parsing"
   "github.com/khaiphan29/logpulse/pkg/logger"
   "github.com/khaiphan29/logpulse/internal/constants"
)

type Producer interface {
   SendMessage(topic *string, key, value []byte) error
   Shutdown()
}

type Handler struct {
   producer Producer
}

func NewHandler(producer Producer) *Handler {
   return &Handler{
      producer: producer,
   }
}

func (h *Handler) GETLog(c *gin.Context) {
   // Log the request
   logger.Info("Received GET request", nil)

   // Respond with a simple message
   c.JSON(200, gin.H{"message": "This is LOG"})
}

func (h *Handler) POSTLog(c *gin.Context) {
   // Validate the incoming JSON payload
   var logData parser.LogPayload
   if err := validateLog(c, &logData); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
   }

   // Log the data
   logger.Info("Received log data", map[string]any{
      "data": logData,
   })

   // Respond with a success message
   c.JSON(http.StatusOK, gin.H{"message": "Log received"})

   // Send the log data to Kafka
   topic := constants.KafkaTopicLogs
   key := []byte(logData.Source)
   value, err := json.Marshal(logData)
   if err != nil {
      logger.Error("Failed to marshal log data", map[string]any{
         "error": err,
      })
      return
   }
   h.producer.SendMessage(&topic, key, value)
}


func validateLog(c *gin.Context, logData *parser.LogPayload) error {
   // Validate the incoming JSON payload
   if err := c.ShouldBindJSON(logData); err != nil {
      logger.Error("Failed to bind JSON", map[string]any{
         "error": err,
      })
      return errors.New("Invalid JSON")
   }

   // Validate LogLevel
   if _, ok := parser.AllowedLogLevels[logData.LogLevel]; !ok {
      logger.Error("Invalid log level", map[string]any{
         "logLevel": logData.LogLevel,
      })
      return errors.New("Invalid log level")
   }

   // Validate Source
   if _, ok := parser.AllowedSources[logData.Source]; !ok {
      logger.Error("Invalid source", map[string]any{
         "source": logData.Source,
      })
      return errors.New("Invalid source")
   }

   return nil
}
