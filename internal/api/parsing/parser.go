package parser

import (
	"time"
)

// Use struct tag to validate using gin ShouldBindJSON
type LogPayload struct {
   LogID       string                 `json:"logId" binding:"required"`
   Timestamp   time.Time              `json:"timestamp" binding:"required"`
   LogLevel    string                 `json:"logLevel" binding:"required"`
   Message     string                 `json:"message" binding:"required"`
   Metadata    map[string]any         `json:"metadata,omitempty"`
   Source      string                 `json:"source" binding:"required"`
   Environment string                 `json:"environment,omitempty"`
   Type        string                 `json:"type,omitempty"`
}

var AllowedLogLevels = map[string]bool{
	"DEBUG": true,
	"INFO":  true,
	"WARN":  true,
	"ERROR": true,
	"FATAL": true,
}
