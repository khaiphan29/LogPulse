package setup

import (
   "github.com/khaiphan29/logpulse/internal/config/redis"
   "github.com/khaiphan29/logpulse/internal/redis"
)

func setupRedisClient() *redis.Client {
   redisConfig := redisconfig.GetRedisConfig()
   cacheClient := redis.NewClient(*redisConfig)
   return cacheClient
}
