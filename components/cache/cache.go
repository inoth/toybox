package cache

import (
	"strings"
	"sync"
)

var (
	m     sync.Map
	Cache *CacheComponent
)

// 本地缓存
type CacheComponent struct{}

func (c *CacheComponent) Init() error {
	Cache = c
	return nil
}

// IsExist 判断key是否存在
func (c *CacheComponent) IsExist(key string) (interface{}, bool) {
	return m.Load(key)
}

// Set 设置缓存
func (c *CacheComponent) Set(key string, value interface{}) bool {
	if _, ok := c.IsExist(key); ok {
		return false
	}
	m.Store(key, value)
	return true
}

// Get 获取缓存数据
func (c *CacheComponent) Get(key string) interface{} {
	value, ok := c.IsExist(key)
	if ok {
		return value
	}
	return nil
}

// Delete 删除缓存
func (c *CacheComponent) Delete(key string) bool {
	return c.Delete(key)
}

// FuzzyDelete 清空全部缓存
func (c *CacheComponent) FuzzyDelete(prefix string) {
	m.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			if strings.HasPrefix(k, prefix) {
				m.Delete(k)
			}
		}
		return true
	})
}

func (c *CacheComponent) GetCacheList() map[string]interface{} {
	res := make(map[string]interface{}, 0)
	m.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			res[k] = value
		}
		return true
	})
	return res
}
