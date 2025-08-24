package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
   *redis.Client
   port string
}

var client *Client

const (
   DB = 0 // Default Redis database
   PREFIX = "logpulse:" // Prefix for Redis keys
)

func GetClient(port string, password string) (*Client, error) {
   if client != nil {
      return client, nil
   }

   reidsClient := redis.NewClient(&redis.Options{
      Addr:     "localhost:" + port,
      Password: password,
      DB:       DB, // Use default DB
   })

   client = &Client{
      Client: reidsClient,
      port:   port,
   }

   return client, nil
}

func (client *Client) Close() error {
   err := client.Client.Close()
   if err != nil {
      return err
   }
   return nil
}

func (client *Client) GetPort() string {
   if client == nil {
      return ""
   }
   return client.port
}

func (client *Client) SetRecord(ctx context.Context, key string, value any, exp time.Duration) error {
   fullKey := PREFIX + key
   err := client.Set(ctx, fullKey, value, exp).Err()
   if err != nil {
      return err
   }
   return nil
}

// Normal GetRecord function that returns a string value
func (client *Client) GetRecord(ctx context.Context, key string) (string, error) {
   fullKey := PREFIX + key
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

func (client *Client) DeleteRecord(ctx context.Context, key string) error {
   fullKey := PREFIX + key
   err := client.Del(ctx, fullKey).Err()
   if err != nil {
      return err
   }
   return nil
}

func (client *Client) IncreaseValue(ctx context.Context, key string, increment int64) (int64, error) {
   fullKey := PREFIX + key
   newValue, err := client.IncrBy(ctx, fullKey, increment).Result()
   if err != nil {
      return 0, err
   }
   return newValue, nil
}

func (client *Client) DecreaseValue(ctx context.Context, key string, decrement int64) (int64, error) {
   fullKey := PREFIX + key
   newValue, err := client.DecrBy(ctx, fullKey, decrement).Result()
   if err != nil {
      return 0, err
   }
   return newValue, nil
}
