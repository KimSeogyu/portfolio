package cstore

import (
	"context"
	"time"

	goredislib "github.com/redis/go-redis/v9"
)

type CacheStore interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Increment(ctx context.Context, key string) (int64, error)
	Exists(ctx context.Context, key string) (bool, error)
}

func NewCacheStore(redisClient goredislib.UniversalClient) CacheStore {
	return &cstoreImpl{
		redisClient: redisClient,
	}
}
