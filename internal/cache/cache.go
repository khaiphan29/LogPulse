package cache

import (
	"context"
	"errors"
	"time"
   "sync"

   "github.com/redis/go-redis/v9"
   "github.com/khaiphan29/logpulse/internal/constants"
)

const (
	TIMEOUT = 200 * time.Millisecond // Default timeout
   PREFIX  = "logpulse:cache:"
   DB      = 0                     // Default DB
   DEFAULT_EXPIRATION = 24 * time.Hour // Default expiration time for cache records
   DEFAULT_CONTEXT_TIMEOUT = 100 * time.Millisecond // Default context timeout
)

type Client struct {
   *redis.Client
}

var (
   instance *Client
   once     sync.Once
)

func init() {
   once.Do(func() {
      reidsClient := redis.NewClient(&redis.Options{
         Addr:     "localhost:" + constants.REDIS_PORT,
         Password: "",
         DB:       DB, // Use default DB
      })
      instance = &Client{
         Client: reidsClient,
      }
   })
}

func Instance() (*Client, error) {
   return instance, nil
}

func (client *Client) SetRecord(key string, value any, exp time.Duration) error {
   fullKey := PREFIX + key
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
   fullKey := PREFIX + key
   ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
   defer cancel()
   value, err := client.Get(ctx, fullKey).Result()
   if err != nil {
      return "", err
   }
   return value, nil
}

func (client *Client) ErrKeyNotExists(err error) bool {
   if errors.Is(err, redis.Nil) {
      return true
   }
   return false
}

func (client *Client) DeleteRecord(key string) error {
   fullKey := PREFIX + key
   ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
   defer cancel()
   err := client.Del(ctx, fullKey).Err()
   if err != nil {
      return err
   }
   return nil
}

func (client *Client) IncreaseValue(key string, increment int64) (int64, error) {
   fullKey := PREFIX + key
   ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
   defer cancel()
   newValue, err := client.IncrBy(ctx, fullKey, increment).Result()
   if err != nil {
      return 0, err
   }
   return newValue, nil
}

func (client *Client) DecreaseValue(key string, decrement int64) (int64, error) {
   fullKey := PREFIX + key
   ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_CONTEXT_TIMEOUT)
   defer cancel()
   newValue, err := client.DecrBy(ctx, fullKey, decrement).Result()
   if err != nil {
      return 0, err
   }
   return newValue, nil
}
