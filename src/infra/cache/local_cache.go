package cache

import (
	"sync"
	"time"
)

type LocalCache struct {
	st *sync.Map
}

func (l *LocalCache) Put(key string, val interface{}, ept ...time.Duration) error {
	l.st.Store(key, val)
	return nil
}

func (l *LocalCache) Get(key string) (val interface{}, err error) {
	val, ok := l.st.Load(key)
	if !ok {
		return nil, ErrRecordNotFound
	}
	return val, nil
}

func (l *LocalCache) Del(key string) error {
	l.st.Delete(key)
	return nil
}

func NewLocalCache() *LocalCache {
	return &LocalCache{st: &sync.Map{}}
}
