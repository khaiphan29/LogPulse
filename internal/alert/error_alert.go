package alert

import (
	"time"

	"github.com/khaiphan29/logpulse/internal/constants"
	"github.com/khaiphan29/logpulse/pkg/logger"
)

const (
   PERIOD = "7d"
   THRESHOLD = 2
   CACHE_KEY = "error_alert_cache"
   CACHE_TTL = 3600 // in seconds
)

type CacheClient interface {
   SetRecord(key string, value any, exp int) error
   GetRecord(key string) (string, error)
   DeleteRecord(key string) error
   ErrKeyNotExist() bool
   IncreaseValue(key string, increment int64) (int64, error)
   DecreaseValue(key string, decrement int64) (int64, error)
}

type Aggregator interface {
   CountTotalLogsByLevel(index, level string, timegte, timelte time.Time) (int, error)
}

type ErrorAlert struct {
   Cache CacheClient
   Agg   Aggregator
}

var ea *ErrorAlert

// trigger by Kafka consumer?
func EvaluateError(cb func ()) {
   // Define the time interval and threshold
   gte := time.Now().Add(-7 * 24 * time.Hour) // 7 days ago
   lte := time.Now()

   // Call the function to get error logs by sources
   results, err := ea.Agg.CountTotalLogsByLevel(constants.ES_INDEX_LOG, constants.LOG_LEVEL_ERROR, gte, lte)
   if err != nil {
      logger.Error("Error getting error logs by sources", map[string]any{
         "error": err,
      })
      return
   }

   if results >= THRESHOLD {
      cb()
   }
}

