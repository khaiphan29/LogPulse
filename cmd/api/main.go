package main

import (
   "os"

   "github.com/khaiphan29/logpulse/internal/api/router"
   "github.com/khaiphan29/logpulse/pkg/logger"
)

func main() {
   // Initialize the logger
   logger.Initialize()

   // Create a new Gin router
   router := router.SetupRouter()

   // Start the server
   port := ":8080"
   if (len(os.Args) > 1) {
      port = ":" + os.Args[1]
   }

   if err := router.Run(port); err != nil {
      logger.Fatal("Failed to start server", map[string]any{
         "error": err,
      })
   } else {
      logger.Info("Server started on port 8080", map[string]any{
         "port": 8080,
      })
   }
}
