package redisconfig

import (
   "os"
   "strconv"
)

type RedisConfig struct {
   Host     string
   Port     string
   Password string
   DB       int
   PrefixKey   string
}

func GetRedisConfig() *RedisConfig {
   db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
   if err != nil {
      db = 0 // Default DB
   }
   return &RedisConfig{
      Host:     os.Getenv("REDIS_HOST"),
      Port:     os.Getenv("REDIS_PORT"),
      Password: os.Getenv("REDIS_PASSWORD"),
      DB:       db, // Default DB
      PrefixKey:   os.Getenv("REDIS_PREFIX_KEY"),
   }
}
