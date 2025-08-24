package unit_test

import (
	"testing"
	"time"

	"github.com/khaiphan29/logpulse/pkg/logger"
	"github.com/khaiphan29/logpulse/tests/helpers"

	"github.com/khaiphan29/logpulse/internal/cache"
	"github.com/stretchr/testify/assert"
)

func init() {
	mockClient := &helpers.MockCacheClient{}
	err := cache.InitCacheClient(mockClient)
	if err != nil {
		logger.Fatal("Failed to initialize Cache client", map[string]any{
			"error": err,
		})
	}
}

func TestSetCache(t *testing.T) {
	key := "test_string"
	value := "value"
	expiration := 1 * time.Second

	err := cache.Set(key, value, expiration)
	assert.NoError(t, err, "Setting cache should not return an error")
}

func TestSetIntCache(t *testing.T) {
	key := "test_int"
	value := "100"
	expiration := 1 * time.Second

	err := cache.Set(key, value, expiration)
	assert.NoError(t, err, "Setting int cache should not return an error")
}

func TestGetCache(t *testing.T) {
	key := "test_string"
	expectedValue := "value"

	value, err := cache.Get(key)
	assert.NoError(t, err, "Getting cache should not return an error")
	assert.Equal(t, expectedValue, value, "Retrieved value should match the set value")
}

func TestGetIntCache(t *testing.T) {
	key := "test_int"
	expectedValue := "100"

	value, err := cache.Get(key)
	assert.NoError(t, err, "Getting int cache should not return an error")
	assert.Equal(t, expectedValue, value, "Retrieved int value should match the set value")
}

func TestIncrementCache(t *testing.T) {
	key := "test_int"
	incrementBy := 50
	expectedValue := int64(150) // 100 + 50

	newValue, err := cache.Increase(key, int64(incrementBy))
	assert.NoError(t, err, "Incrementing cache should not return an error")
	assert.Equal(t, expectedValue, newValue, "Incremented value should match the expected value")
}

func TestDecrementCache(t *testing.T) {
	key := "test_int"
	decrementBy := 30
	expectedValue := int64(120) // 150 - 30

	newValue, err := cache.Decrease(key, int64(decrementBy))
	assert.NoError(t, err, "Decrementing cache should not return an error")
	assert.Equal(t, expectedValue, newValue, "Decremented value should match the expected value")
}

func TestCacheDelete(t *testing.T) {
	key := "test_string"

	err := cache.Delete(key)
	assert.NoError(t, err, "Deleting cache should not return an error")

	// Try to get the deleted key
	value, err := cache.Get(key)
	assert.Error(t, err, "Getting deleted cache should return an error")
	assert.Equal(t, cache.ErrCacheMiss, err, "Error should be ErrCacheMiss for deleted key")
	assert.Nil(t, value, "Value for deleted key should be nil")
}
