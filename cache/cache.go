package cache

import (
	"sync"
)

type Cache struct {
	mutex      sync.Mutex
	lru        *lru
	CacheBytes int64
}

func (c *Cache) Set(key string, value ByteView) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		c.lru = New(c.CacheBytes, nil)
	}
	c.lru.Set(key, value)
}

func (c *Cache) Get(key string) (v ByteView, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		return
	}
	value, ok := c.lru.Get(key)
	if ok {
		return value.(ByteView), true
	}
	return
}
