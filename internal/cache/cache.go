package cache

import (
	"context"
	"errors"
	"time"
)

const (
	TIMEOUT = 200 * time.Millisecond // Default timeout
	PREFIX  = "cache:"
)

// Client is a wrapper around the Redis (or any satisfied services) client
type cacheClient interface {
	SetRecord(ctx context.Context, key string, value string, exp time.Duration) error
	GetRecord(ctx context.Context, key string) (any, error)
	DeleteRecord(ctx context.Context, key string) error
	ErrKeyNotExists(error) bool
	IncreaseValue(ctx context.Context, key string, value int64) (int64, error)
	DecreaseValue(ctx context.Context, key string, value int64) (int64, error)
}

var client cacheClient

var ErrCacheMiss = errors.New("Cache miss: key does not exist")

// InitCacheClient initializes the cache client with the provided implementation
func InitCacheClient(c cacheClient) error {
	if client != nil {
		return errors.New("Cache client is already initialized")
	} else if c == nil {
		return errors.New("Cache client cannot be nil")
	}
	client = c
	return nil
}

func Get(key string) (any, error) {
	if client == nil {
		return nil, errors.New("Cache client is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	value, err := client.GetRecord(ctx, PREFIX+key)
	if err != nil {
		if client.ErrKeyNotExists(err) {
			return nil, ErrCacheMiss // Key does not exist
		}
		return nil, err // Some other error
	}
	return value, nil
}

func Set(key string, value string, exp time.Duration) error {
	if client == nil {
		return errors.New("Cache client is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	err := client.SetRecord(ctx, PREFIX+key, value, exp)
	if err != nil {
		return err // Some other error
	}
	return nil
}

func Increase(key string, increment int64) (int64, error) {
	if client == nil {
		return 1, errors.New("Cache client is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	newValue, err := client.IncreaseValue(ctx, PREFIX+key, increment)
	if err != nil {
		return 1, err // Some other error
	}
	return newValue, nil
}

func Decrease(key string, decrement int64) (int64, error) {
	if client == nil {
		return 1, errors.New("Cache client is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	newValue, err := client.DecreaseValue(ctx, PREFIX+key, decrement)
	if err != nil {
		return 1, err // Some other error
	}
	return newValue, nil
}

func Delete(key string) error {
	if client == nil {
		return errors.New("Cache client is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()
	err := client.DeleteRecord(ctx, PREFIX+key)
	if err != nil {
		return err // Some other error
	}
	return nil
}
