package router

import (
   "github.com/khaiphan29/logpulse/pkg/logger"
   "github.com/khaiphan29/logpulse/internal/api/handlers"
)

type MockProducer struct{}

func (m *MockProducer) SendMessage(topic *string, key, value []byte) error {
   logger.Info("Mock SendMessage called", map[string]any{
      "topic": topic,
      "key":   string(key),
      "value": string(value),
   })
   return nil
}

func (m *MockProducer) Shutdown() {
   // Mock shutdown logic
}

func NewMockHandler() *handlers.Handler {
   mockProducer := &MockProducer{}
   logHandler := handlers.NewHandler(mockProducer)
   return logHandler
}
