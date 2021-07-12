package cache

import (
	"errors"
	"time"
)

type ICache interface {
	Put(key string, val interface{}, ept ...time.Duration) error
	Get(key string) (val interface{}, err error)
	Del(key string) error
}

var (
	ErrRecordNotFound = errors.New("record not found error")
)
