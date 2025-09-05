package utils

import (
   "os"
   "strconv"
   "github.com/khaiphan29/logpulse/pkg/logger"
)

func LoadBoolEnv(key string) bool {
   valStr := os.Getenv(key)
   val, err := strconv.ParseBool(valStr)
   if err != nil {
      logger.Fatal("Invalid boolean value for " + key, map[string]any{
         "error": err,
      })
   }
   return val
}

func LoadIntEnv(key string) int {
   valStr := os.Getenv(key)
   val, err := strconv.Atoi(valStr)
   if err != nil {
      logger.Fatal("Invalid integer value for " + key, map[string]any{
         "error": err,
      })
   }
   return val
}
