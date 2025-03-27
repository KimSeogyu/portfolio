package dlock

import "github.com/go-redsync/redsync/v4"

type redsyncImpl struct {
	rs *redsync.Redsync
}

func (m *redsyncImpl) NewMutex(key string, opts ...redsync.Option) DLocker {
	return m.rs.NewMutex(key, opts...)
}
