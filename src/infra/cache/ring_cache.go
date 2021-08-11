package cache

import "sync"

type RingCache struct {
	list     []interface{}
	keyList  []string
	cacheCap int
	index    int
	sync.RWMutex
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
