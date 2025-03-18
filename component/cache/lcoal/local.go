package cache

import (
	"sync"
	"time"
)

type LocalCacheComponent struct {
	items sync.Map
}

func NewLocalCache() *LocalCacheComponent {
	return &LocalCacheComponent{
		items: sync.Map{},
	}
}

func (c *LocalCacheComponent) Set(key string, value any, expiration time.Duration) {
	if expiration > 0 {
		go func() {
			time.Sleep(expiration * time.Second)
			c.items.Delete(key)
		}()
	}
	c.items.Store(key, value)
}

func (c *LocalCacheComponent) Get(key string) (any, bool) {
	return c.items.Load(key)
}

func (c *LocalCacheComponent) Delete(key string) {
	c.items.Delete(key)
}
