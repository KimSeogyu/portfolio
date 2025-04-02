package cstore

import (
	"context"
	"fmt"
	"time"

	goredislib "github.com/redis/go-redis/v9"
)

type cstoreImpl struct {
	redisClient goredislib.UniversalClient
}

// Delete implements CStore.
func (c *cstoreImpl) Delete(ctx context.Context, key string) error {
	err := c.redisClient.Del(ctx, key).Err()
	if err != nil {
		if err == goredislib.Nil {
			return nil
		}
		return fmt.Errorf("redisClient.Del: %w", err)
	}
	return nil
}

// Get implements CStore.
func (c *cstoreImpl) Get(ctx context.Context, key string) (string, error) {
	value, err := c.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == goredislib.Nil {
			return "", nil
		}
		return "", fmt.Errorf("redisClient.Get: %w", err)
	}
	return value, nil
}

// Set implements CStore.
func (c *cstoreImpl) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := c.redisClient.Set(ctx, key, value, ttl).Err()
	if err != nil {
		if err == goredislib.Nil {
			return nil
		}
		return fmt.Errorf("redisClient.Set: %w", err)
	}
	return nil
}

// Increment implements CStore.
func (c *cstoreImpl) Increment(ctx context.Context, key string) (int64, error) {
	value, err := c.redisClient.Incr(ctx, key).Result()
	if err != nil {
		if err == goredislib.Nil {
			return 0, nil
		}
		return 0, fmt.Errorf("redisClient.Incr: %w", err)
	}
	return value, nil
}

// Exists implements CStore.
func (c *cstoreImpl) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := c.redisClient.Exists(ctx, key).Result()
	if err != nil {
		if err == goredislib.Nil {
			return false, nil
		}
		return false, fmt.Errorf("redisClient.Exists: %w", err)
	}
	return exists > 0, nil
}
