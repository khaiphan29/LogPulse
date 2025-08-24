package unit

import (
   "context"
   "testing"
   "time"
   "strconv"

   "github.com/stretchr/testify/assert"

   "github.com/khaiphan29/logpulse/pkg/logger"
   "github.com/khaiphan29/logpulse/pkg/redis"
   "github.com/khaiphan29/logpulse/internal/constants"
)

const (
   REDIS_PREFIX = "test:" // Prefix for Redis keys
   REDIS_EXPIRATION = 500 * time.Millisecond // Expiration time for Redis keys
)

var (
   redisClient *redis.Client
)

func init() {
   var err error
   redisClient, err = redis.GetClient(constants.REDIS_PORT, "")
   if err != nil {
      logger.Fatal("Failed to initialize Redis client", map[string]any{
         "error": err,
      })
   }
}

func TestRedisClient(t *testing.T) {
   assert.NotNil(t, redisClient, "Redis client should not be nil")
   assert.Equal(t, redisClient.GetPort(), constants.REDIS_PORT, "Redis client port should be 6379")
}

func TestSetAndGetStringRecord(t *testing.T) {
   ctx := context.Background()
   key := REDIS_PREFIX + "test_string"
   value := "test_value"

   // Set a record
   err := redisClient.SetRecord(ctx, key, value, REDIS_EXPIRATION)
   assert.NoError(t, err, "Setting record should not return an error")

   // Get the record
   retrievedValue, err := redisClient.GetRecord(ctx, key)
   assert.NoError(t, err, "Getting record should not return an error")
   assert.Equal(t, value, retrievedValue, "Retrieved value should match the set value")
}

func TestSetAndGetIntRecord(t *testing.T) {
   ctx := context.Background()
   key := REDIS_PREFIX + "test_int"
   value := 42

   // Set a record
   err := redisClient.SetRecord(ctx, key, value, REDIS_EXPIRATION)
   assert.NoError(t, err, "Setting record should not return an error")

   // Get the record
   retrievedValue, err := redisClient.GetRecord(ctx, key)
   assert.NoError(t, err, "Getting record should not return an error")
   assert.Equal(t, strconv.Itoa(value), retrievedValue, "Retrieved value should match the set value")
}

func TestIncrementIntValue(t *testing.T) {
   ctx := context.Background()
   key := REDIS_PREFIX + "test_int"

   initResult, err := redisClient.GetRecord(ctx, key)
   initIntValue, _ := strconv.Atoi(initResult)
   assert.NoError(t, err, "Getting initial record should not return an error")
   assert.NotEqual(t, 0, initIntValue, "Initial value should not be zero")

   // Increment the record
   value, err := redisClient.IncreaseValue(ctx, key, 1)
   assert.NoError(t, err, "Incrementing record should not return an error")
   assert.Equal(t, int64(initIntValue) + 1, value, "Incremented value should be 1 more than initial")
}

func TestDecrementIntValue(t *testing.T) {
   ctx := context.Background()
   key := REDIS_PREFIX + "test_int"

   initResult, err := redisClient.GetRecord(ctx, key)
   initIntValue, _ := strconv.Atoi(initResult)
   assert.NoError(t, err, "Getting initial record should not return an error")
   assert.NotEqual(t, 0, initIntValue, "Initial value should not be zero")

   // Decrement the record
   value, err := redisClient.DecreaseValue(ctx, key, 1)
   assert.NoError(t, err, "Decrementing record should not return an error")
   assert.Equal(t, int64(initIntValue) - 1, value, "Decremented value should be 1 less than initial")
}

func TestIncrementStringValue(t *testing.T) {
   ctx := context.Background()
   key := REDIS_PREFIX + "test_string"

   // Increment the record
   _, err := redisClient.IncreaseValue(ctx, key, 1)
   assert.Error(t, err, "Incrementing a string value should return an error")
}

func TestDecrementStringValue(t *testing.T) {
   ctx := context.Background()
   key := REDIS_PREFIX + "test_string"

   // Decrement the record
   _, err := redisClient.DecreaseValue(ctx, key, 1)
   assert.Error(t, err, "Decrementing a string value should return an error")
}

func TestKeyNotExists(t *testing.T) {
   ctx := context.Background()
   key := REDIS_PREFIX + "non_existent_key"

   // Try to get a non-existent key
   _, err := redisClient.GetRecord(ctx, key)
   assert.True(t, redisClient.ErrKeyNotExists(err), "Should return true for non-existent key")
}

func TestDeleteRecord(t *testing.T) {
   ctx := context.Background()
   key := REDIS_PREFIX + "test_string"

   // Delete the record
   err := redisClient.DeleteRecord(ctx, key)
   assert.NoError(t, err, "Deleting record should not return an error")

   // Try to get the deleted key
   _, err = redisClient.GetRecord(ctx, key)
   assert.True(t, redisClient.ErrKeyNotExists(err), "Should return true for deleted key")
}
