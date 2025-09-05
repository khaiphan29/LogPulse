package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
   redisconfig "github.com/khaiphan29/logpulse/internal/config/redis"
)

const (
	TIMEOUT = 200 * time.Millisecond // Default timeout
   DB      = 0                     // Default DB
   DEFAULT_EXPIRATION = 24 * time.Hour // Default expiration time for cache records
   DEFAULT_CONTEXT_TIMEOUT = 100 * time.Millisecond // Default context timeout
)

type Client struct {
   *redis.Client
   appPrefix string
}

func NewClient(cfg redisconfig.RedisConfig) *Client {
   reidsClient := redis.NewClient(&redis.Options{
      Addr:     cfg.Host + ":" + cfg.Port,
      Password: cfg.Password,
      DB:       cfg.DB, // Use default DB
   })
   return &Client{
      Client: reidsClient,
      appPrefix: cfg.PrefixKey,
   }
}

func (client *Client) SetRecord(key string, value any, exp time.Duration) error {
   fullKey := client.appPrefix + key
   ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
   defer cancel()
   err := client.Set(ctx, fullKey, value, exp).Err()
   if err != nil {
      return err
   }
   return nil
}

// Normal GetRecord function that returns a string value
func (client *Client) GetRecord(key string) (string, error) {
   fullKey := client.appPrefix + key
   ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
   defer cancel()
   value, err := client.Get(ctx, fullKey).Result()
   if err != nil {
      return "", err
   }
   return value, nil
}

func (client *Client) ErrKeyNotExist(err error) bool {
   if errors.Is(err, redis.Nil) {
      return true
   }
   return false
}

func (client *Client) DeleteRecord(key string) error {
   fullKey := client.appPrefix + key
   ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
   defer cancel()
   err := client.Del(ctx, fullKey).Err()
   if err != nil {
      return err
   }
   return nil
}

func (client *Client) IncreaseValue(key string, increment int64) (int64, error) {
   fullKey := client.appPrefix + key
   ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
   defer cancel()
   newValue, err := client.IncrBy(ctx, fullKey, increment).Result()
   if err != nil {
      return 0, err
   }
   return newValue, nil
}

func (client *Client) DecreaseValue(key string, decrement int64) (int64, error) {
   fullKey := client.appPrefix + key
   ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
   defer cancel()
   newValue, err := client.DecrBy(ctx, fullKey, decrement).Result()
   if err != nil {
      return 0, err
   }
   return newValue, nil
}
