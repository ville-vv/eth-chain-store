package cache

import (
	"github.com/ville-vv/eth-chain-store/src/common/conf"
	"sync"
	"time"
)

type RingCache struct {
	list     []interface{}
	keyList  []string
	cacheCap int
	index    int
	sync.RWMutex
}

func NewRingCache() *RingCache {
	return &RingCache{
		keyList: make([]string, 20000),
		list:    make([]interface{}, 20000),
	}
}

func (sel *RingCache) Set(key string, val interface{}) error {
	sel.Lock()
	sel.keyList[sel.index] = key
	sel.list[sel.index] = val
	sel.index++
	if sel.index >= sel.cacheCap {
		sel.index = 0
	}
	sel.Unlock()
	return nil
}

func (sel *RingCache) Get(key string) (interface{}, bool) {
	sel.RLock()
	defer sel.RUnlock()
	for i := 0; i < sel.cacheCap; i++ {
		if sel.keyList[i] == key {
			return sel.list[i], true
		}
	}
	return nil, false
}

func (sel *RingCache) Del(key string) bool {
	sel.Lock()
	defer sel.Unlock()
	for i := 0; i < sel.cacheCap; i++ {
		if sel.keyList[i] == key {
			sel.keyList[i] = ""
			sel.list[i] = 0
			return true
		}
	}
	return false
}

type RingStrListV2 struct {
	sync.RWMutex
	list   map[string]int
	length int
}

func NewRingStrListV2() *RingStrListV2 {
	r := &RingStrListV2{
		list:   make(map[string]int),
		length: 0,
	}

	go func() {
		tmr := time.NewTicker(time.Minute * 5)
		for {
			select {
			case <-tmr.C:
				if r.length > 10000000 {
					cp := 0
					list := make(map[string]int)
					for k, val := range r.list {
						if val > 5000 {
							list[k] = 10
						}
						cp++
						if cp > 500000 {
							break
						}
					}
					r.list = list
				}
			case <-conf.GlobalExitSignal:
				return
			}
		}
	}()

	return r
}

func (sel *RingStrListV2) Exist(str string) bool {
	sel.RLock()
	defer sel.RUnlock()
	n, ok := sel.list[str]
	sel.list[str] = n + 1
	return ok
}

func (sel *RingStrListV2) Set(str string) {
	sel.Lock()
	sel.list[str] = 1
	sel.length++
	sel.Unlock()
}

func (sel *RingStrListV2) Del(str string) {
	sel.Lock()
	delete(sel.list, str)
	sel.length--
	sel.Unlock()
}
