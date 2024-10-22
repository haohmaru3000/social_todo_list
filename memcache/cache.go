package memcache

import (
	"sync"
)

type Cache interface {
	Write(k string, value interface{})
	Read(k string) interface{}
}

type caching struct {
	store  map[string]interface{}
	locker *sync.RWMutex
}

func NewCaching() *caching {
	return &caching{
		store:  make(map[string]interface{}),
		locker: new(sync.RWMutex),
	}
}

func (c *caching) Write(key string, value interface{}) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.store[key] = value
}

func (c *caching) Read(key string) interface{} {
	c.locker.RLock()
	defer c.locker.RUnlock()
	return c.store[key]
}
