package cstore

import (
	"context"
	"time"

	goredislib "github.com/redis/go-redis/v9"
)

type cstoreImpl struct {
	redisClient goredislib.UniversalClient
}

// Delete implements CStore.
func (c *cstoreImpl) Delete(ctx context.Context, key string) error {
	return c.redisClient.Del(ctx, key).Err()
}

// Get implements CStore.
func (c *cstoreImpl) Get(ctx context.Context, key string) (string, error) {
	return c.redisClient.Get(ctx, key).Result()
}

// Set implements CStore.
func (c *cstoreImpl) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.redisClient.Set(ctx, key, value, ttl).Err()
}

// Increment implements CStore.
func (c *cstoreImpl) Increment(ctx context.Context, key string) (int64, error) {
	return c.redisClient.Incr(ctx, key).Result()
}

// Exists implements CStore.
func (c *cstoreImpl) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := c.redisClient.Exists(ctx, key).Result()
	return exists > 0, err
}
