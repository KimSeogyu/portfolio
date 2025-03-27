package dlock

import (
	"context"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

type DLockerFactory interface {
	NewMutex(key string, opts ...redsync.Option) DLocker
}

type DLocker interface {
	LockContext(ctx context.Context) error
	UnlockContext(ctx context.Context) (bool, error)
}

var ErrLocked = redsync.ErrFailed

func NewDLockerFactory(redisClient goredislib.UniversalClient) DLockerFactory {
	pool := goredis.NewPool(redisClient) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)

	return &redsyncImpl{
		rs: rs,
	}
}
