package setup

import (
   "os"
   "time"
   "strconv"

   "github.com/khaiphan29/logpulse/pkg/logger"
   "github.com/khaiphan29/logpulse/internal/analyzer"
)

func setupLogErrorAnalyzer(agg analyzer.LogErrAggregator, cacheClient analyzer.CacheClient, indexName string) *analyzer.LogErrAnalyzer {
   cacheKey := os.Getenv("CACHE_LOG_COUNT_KEY")
   cacheTTL, err := time.ParseDuration(os.Getenv("CACHE_LOG_COUNT_TTL"))
   if err != nil {
      logger.Fatal("Failed to parse cache CACHE_LOG_COUNT_TTL", map[string]any{
         "error": err,
      })
   }
   threshold, err := strconv.Atoi(os.Getenv("ANALYZER_LOG_ERROR_THRESHOLD"))
   if err != nil {
      logger.Fatal("Failed to parse ANALYZER_LOG_ERROR_THRESHOLD", map[string]any{
         "error": err,
      })
   }
   timeWindow, err := time.ParseDuration(os.Getenv("ANALYZER_LOG_ERROR_TIME_WINDOW"))
   if err != nil {
      logger.Fatal("Failed to parse ANALYZER_LOG_ERROR_TIME_WINDOW", map[string]any{
         "error": err,
      })
   }

   la := analyzer.NewLogErrorAnalyzer(agg, indexName, threshold, timeWindow, cacheClient, cacheKey, cacheTTL)
   return la
}
