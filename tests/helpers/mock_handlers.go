package helpers

import (
   "github.com/khaiphan29/logpulse/internal/api/handlers"
)

func NewMockHandler() *handlers.Handler {
   mockProducer := &MockProducer{}
   logHandler := handlers.NewHandler(mockProducer)
   return logHandler
}
