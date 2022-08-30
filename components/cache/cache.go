package cache

import (
	"strings"
	"sync"
)

var (
	m     sync.Map
	Cache *LocalCache
)

// 本地缓存
type LocalCache struct{}

type CacheComponents struct{}

func (CacheComponents) Init() error {
	Cache = &LocalCache{}
	return nil
}

// IsExist 判断key是否存在
func (c *LocalCache) IsExist(key string) (interface{}, bool) {
	return m.Load(key)
}

// Set 设置缓存
func (c *LocalCache) Set(key string, value interface{}) bool {
	if _, ok := c.IsExist(key); ok {
		return false
	}
	m.Store(key, value)
	return true
}

// Get 获取缓存数据
func (c *LocalCache) Get(key string) interface{} {
	value, ok := c.IsExist(key)
	if ok {
		return value
	}
	return nil
}

// Delete 删除缓存
func (c *LocalCache) Delete(key string) bool {
	return c.Delete(key)
}

// FuzzyDelete 清空全部缓存
func (c *LocalCache) FuzzyDelete(prefix string) {
	m.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			if strings.HasPrefix(k, prefix) {
				m.Delete(k)
			}
		}
		return true
	})
}
