package helpers

import (
   "github.com/khaiphan29/logpulse/pkg/logger"
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
