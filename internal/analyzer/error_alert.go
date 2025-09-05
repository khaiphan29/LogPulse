package analyzer

import (
	"time"

	"github.com/khaiphan29/logpulse/internal/constants"
	"github.com/khaiphan29/logpulse/pkg/logger"
)

// const (
//    PERIOD = "7d"
//    THRESHOLD = 2
//    CACHE_KEY = "error_alert_cache"
//    CACHE_TTL = 3600 // in seconds
// )

type CacheClient interface {
   SetRecord(key string, value any, exp time.Duration) error
   GetRecord(key string) (string, error)
   DeleteRecord(key string) error
   ErrKeyNotExist(err error) bool
   IncreaseValue(key string, increment int64) (int64, error)
   DecreaseValue(key string, decrement int64) (int64, error)
}

type LogErrAggregator interface {
   CountTotalLogsByLevel(index, level string, timegte, timelte time.Time) (int, error)
}

type LogErrAnalyzer struct {
   Agg   LogErrAggregator
   AggIndex string

   Threshold int
   TimeWindow time.Duration

   Cache CacheClient
   CacheErrKey string
   CacheTTL time.Duration
}

func NewLogErrorAnalyzer(agg LogErrAggregator, aggIndex string, threshold int, timeWindow time.Duration, cache CacheClient, cacheErrKey string, cacheTTL time.Duration) *LogErrAnalyzer {
   return &LogErrAnalyzer{
      Agg: agg,
      AggIndex: aggIndex,
      Threshold: threshold,
      TimeWindow: timeWindow,
      Cache: cache,
      CacheErrKey: cacheErrKey,
      CacheTTL: cacheTTL,
   }
}

func (a *LogErrAnalyzer) AnalyzeError() {
   // Define the time interval and threshold
   gte := time.Now().Add(-a.TimeWindow)
   lte := time.Now()

   // Call the function to get error logs by sources
   results, err := a.Agg.CountTotalLogsByLevel(a.AggIndex, constants.LOG_LEVEL_ERROR, gte, lte)
   if err != nil {
      logger.Error("Error getting error logs by sources", map[string]any{
         "error": err,
      })
      return
   }

   if results >= a.Threshold {
      a.errorAlert()
   }
}

func (a *LogErrAnalyzer) errorAlert() {
   // send alert
   logger.Warn("High number of error logs detected", map[string]any{
      "threshold": a.Threshold,
      "Window": a.TimeWindow,
   })
}
