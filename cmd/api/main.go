package main

import (
	"os"
   "os/signal"
   "syscall"
   "github.com/khaiphan29/logpulse/pkg/logger"
	"github.com/khaiphan29/logpulse/internal/setup"
)

func main() {
   shutdown := setup.InitService()
   defer shutdown()
   // Set up a graceful shutdown
   quit := make(chan os.Signal, 1)
   // Wait for a signal to shut down
   // the process will not terminate immediately when the signal is received.
   signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
   <-quit
   logger.Info("Shutting down server...", nil)
   os.Exit(0) // Ensure the program exits with a zero status code
}
