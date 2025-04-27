package router

import (
   "github.com/gin-gonic/gin"

   "github.com/khaiphan29/logpulse/internal/api/handlers"
)

func SetupRouter(router *gin.Engine) error {
   // Define a simple GET endpoint
   router.GET("/logs", handlers.GetLogHandler)

   // Define a POST endpoint for log data
   router.POST("/logs", handlers.PostLogHandler)

   return nil
}
