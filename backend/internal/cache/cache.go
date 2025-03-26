package cache

import (
	"context"
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	NewMutex(ctx context.Context, key string, ttl time.Duration) (*redsync.Mutex, error)
	Exists(ctx context.Context, key string) (bool, error)
}

var _ Cache = &redisCache{}

type redisCache struct {
	rs     *redsync.Redsync
	client *goredislib.Client
}

func NewRedisCache(addr string) *redisCache {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: addr,
	})

	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)

	return &redisCache{
		rs:     rs,
		client: client,
	}
}

// Exists implements Cache.
func (r *redisCache) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := r.client.Exists(ctx, key).Result()
	return exists > 0, err
}

// Get implements Cache.
func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

// Set implements Cache.
func (r *redisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

// Lock implements Cache.
func (r *redisCache) NewMutex(ctx context.Context, key string, ttl time.Duration) (*redsync.Mutex, error) {
	mutex := r.rs.NewMutex(key)

	return mutex, nil
}
