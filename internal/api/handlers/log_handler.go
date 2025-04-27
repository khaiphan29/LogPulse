package handlers

import (
   "net/http"
   "github.com/gin-gonic/gin"

   "github.com/khaiphan29/logpulse/internal/api/parsing"
   "github.com/khaiphan29/logpulse/pkg/logger"
)

func GetLogHandler(c *gin.Context) {
   // Log the request
   logger.Info("Received GET request", nil)

   // Respond with a simple message
   c.JSON(200, gin.H{"message": "This is LOG"})
}

func PostLogHandler(c *gin.Context) {
   // Validate the incoming JSON payload
   var logData parser.LogPayload
   if err := c.ShouldBindJSON(&logData); err != nil {
      logger.Error("Failed to bind JSON", map[string]any{
         "error": err,
      })
      c.JSON(http.StatusBadRequest, gin.H{
         "error": "Invalid JSON",
      })
      return
   }

   // Validate LogLevel
   if _, ok := parser.AllowedLogLevels[logData.LogLevel]; !ok {
      logger.Error("Invalid log level", map[string]any{
         "logLevel": logData.LogLevel,
      })
      c.JSON(http.StatusBadRequest, gin.H{
         "error": "Invalid log level",
      })
      return
   }

   // Log the data
   logger.Info("Received log data", map[string]any{
      "data": logData,
   })

   // Respond with a success message
   c.JSON(http.StatusOK, gin.H{"message": "Log received"})
}
