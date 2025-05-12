package router

import (
   "github.com/gin-gonic/gin"

   "github.com/khaiphan29/logpulse/internal/api/handlers"
)

func NewRouter(mode string, logHandler *handlers.Handler) *gin.Engine {
   var engine *gin.Engine
   switch mode {
      case "test":
         // Set Gin to test mode
         gin.SetMode(gin.TestMode)
         engine = gin.New()
      default:
         // Set Gin to release mode
         gin.SetMode(gin.ReleaseMode)
         engine = gin.Default()
      }

   setupRouter(engine, logHandler)
   return engine
}

func setupRouter(router *gin.Engine, logHandler *handlers.Handler) error {
   // Define a simple GET endpoint
   router.GET("/logs", logHandler.GETLog)

   // Define a POST endpoint for log data
   router.POST("/logs", logHandler.POSTLog)

   return nil
}
