package cache

import (
	"sync"
	"time"
)

type Cache struct {
	l      sync.Mutex
	cache  map[interface{}]*data
	config *CacheConfig
}

type CacheConfig struct {
	AutoOverride bool
}

func DefaultConfig() *CacheConfig {
	return &CacheConfig{}
}

type data struct {
	d        time.Duration
	item     interface{}
	callback func()
	timer    *time.Timer
}

func New(config *CacheConfig) *Cache {
	return &Cache{
		config: config,
		cache:  make(map[interface{}]*data),
	}
}

func DefaultCache() *Cache {
	return New(DefaultConfig())
}

func (c *Cache) Add(key, value interface{}, duration time.Duration, callback func()) (item interface{}, result bool) {
	c.l.Lock()
	defer c.l.Unlock()

	if !c.config.AutoOverride {
		if v, ok := c.cache[key]; ok {
			return v.item, false
		}
	}

	v := &data{
		callback: callback,
		d:        duration,
		item:     value,
	}
	v.timer = time.AfterFunc(duration, func() {
		c.l.Lock()
		delete(c.cache, key)
		c.l.Unlock()
		v.callback()
	})
	c.cache[key] = v
	return value, true
}

func (c *Cache) Get(key interface{}) (interface{}, bool) {
	v, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	return v.item, true
}

func (c *Cache) Remove(key interface{}) (value *data, result bool) {
	c.l.Lock()
	defer c.l.Unlock()
	v, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	delete(c.cache, key)
	return v, true
}

func (c *Cache) Flush() {
	c.l.Lock()
	defer c.l.Unlock()
	var removeList []*data
	for key := range c.cache {
		if data, ok := c.Remove(key); ok {
			removeList = append(removeList, data)
		}
	}

	for _, v := range removeList {
		v.callback()
	}
}
