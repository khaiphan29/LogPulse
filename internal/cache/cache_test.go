package cache_test

import (
	"testing"
	"time"

	"github.com/khaiphan29/logpulse/pkg/logger"

	"github.com/khaiphan29/logpulse/internal/cache"
	"github.com/stretchr/testify/assert"
)

var cacheClient *cache.Client

func init() {
   var err error
   cacheClient, err = cache.Instance()
   if err != nil {
      logger.Fatal("Failed to initialize cache client", map[string]any{
         "error": err,
      })
   }
}

func TestSetCacheStringValue(t *testing.T) {
	key := "test_string"
	value := "value"
	expiration := cache.DEFAULT_EXPIRATION

	err := cacheClient.SetRecord(key, value, expiration)
	assert.NoError(t, err, "Setting cacheClient should not return an error")
}

func TestSetIntCacheValue(t *testing.T) {
	key := "test_int"
	value := "100"
	expiration := 1 * time.Second

	err := cacheClient.SetRecord(key, value, expiration)
	assert.NoError(t, err, "Setting int cacheClient should not return an error")
}

func TestGetCacheString(t *testing.T) {
	key := "test_string"
	expectedValue := "value"

	value, err := cacheClient.GetRecord(key)
	assert.NoError(t, err, "Getting cacheClient should not return an error")
	assert.Equal(t, expectedValue, value, "Retrieved value should match the set value")
}

func TestGetIntCacheValue(t *testing.T) {
	key := "test_int"
	expectedValue := "100"

	value, err := cacheClient.GetRecord(key)
	assert.NoError(t, err, "Getting int cacheClient should not return an error")
	assert.Equal(t, expectedValue, value, "Retrieved int value should match the set value")
}

func TestCacheValueIncrement(t *testing.T) {
	key := "test_int"
	incrementBy := 50
	expectedValue := int64(150) // 100 + 50

	newValue, err := cacheClient.IncreaseValue(key, int64(incrementBy))
	assert.NoError(t, err, "Incrementing cacheClient should not return an error")
	assert.Equal(t, expectedValue, newValue, "Incremented value should match the expected value")
}

func TestCacheValueDecrement (t *testing.T) {
	key := "test_int"
	decrementBy := 30
	expectedValue := int64(120) // 150 - 30

	newValue, err := cacheClient.DecreaseValue(key, int64(decrementBy))
	assert.NoError(t, err, "Decrementing cacheClient should not return an error")
	assert.Equal(t, expectedValue, newValue, "Decremented value should match the expected value")
}

func TestCacheDeletion(t *testing.T) {
	key := "test_string"

	err := cacheClient.DeleteRecord(key)
	assert.NoError(t, err, "Deleting cacheClient should not return an error")

	// Try to get the deleted key
	value, err := cacheClient.GetRecord(key)
	assert.Error(t, err, "Getting deleted cacheClient should return an error")
   assert.True(t, cacheClient.ErrKeyNotExists(err), "Error should indicate that the key does not exist")
   assert.Equal(t, "", value, "Value of deleted key should be empty")
}
