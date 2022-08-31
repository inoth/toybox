package config

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/inoth/ino-toybox/components/cache"
	"github.com/spf13/viper"
)

// const confKeyPrefix = "defaultProject_"

var (
	Cfg  *ViperComponent
	once sync.Once
)

var lastChangeTime time.Time

func init() {
	lastChangeTime = time.Now()
}

type ViperComponent struct {
	defaultValue  map[string]interface{}
	viper         *viper.Viper
	path          string
	confKeyPrefix string
}

func Instance(path ...string) *ViperComponent {
	once.Do(func() {
		Cfg = &ViperComponent{
			defaultValue: make(map[string]interface{}),
			viper:        viper.New(),
		}
		if len(path) > 0 {
			Cfg.path = path[0]
		} else {
			Cfg.path = "config"
		}
	})
	return Cfg
}

func (m *ViperComponent) SetDefaultValue(defaultValue map[string]interface{}) *ViperComponent {
	for k, v := range defaultValue {
		m.defaultValue[k] = v
	}
	return m
}

func (m *ViperComponent) loadDefaultValue() {
	for k, v := range m.defaultValue {
		m.viper.SetDefault(k, v)
	}
}

func (m *ViperComponent) Init() error {
	m.viper.AddConfigPath(m.path)
	m.viper.SetConfigName(selectConfigName())
	m.viper.SetConfigType("yaml")
	if err := m.viper.ReadInConfig(); err != nil {
		return err
	}
	m.loadDefaultValue()

	m.confKeyPrefix = m.GetString("ServerName")
	return nil
}

func selectConfigName() string {
	e := os.Getenv("GORUNEVN")
	if len(e) > 0 {
		return "config." + e
	}
	return "config"
}

// isCached 判断相关键是否已经缓存
func (y *ViperComponent) isCached(key string) bool {
	if _, ok := cache.Cache.IsExist(y.confKeyPrefix + key); ok {
		return true
	}
	return false
}

// cache 对键值进行缓存
func (y *ViperComponent) cache(key string, value interface{}) bool {
	return cache.Cache.Set(y.confKeyPrefix+key, value)
}

// getFromCache 通过键获取缓存的值
func (y *ViperComponent) getFromCache(key string) interface{} {
	return cache.Cache.Get(y.confKeyPrefix + key)
}

// clearCache 清空配置项
func (y *ViperComponent) clearCache() {
	cache.Cache.FuzzyDelete(y.confKeyPrefix)
}

// ConfigFileChangeListen 监听文件变化
func (y *ViperComponent) ConfigFileChangeListen() {
	y.viper.OnConfigChange(func(changeEvent fsnotify.Event) {
		if time.Now().Sub(lastChangeTime).Seconds() >= 1 {
			if changeEvent.Op.String() == "WRITE" {
				y.clearCache()
				lastChangeTime = time.Now()
				fmt.Println("配置文件已更新")
			}
		}
	})
	y.viper.WatchConfig()
}

// Get 获取原始值。先尝试从cache读取，若读取不到，从配置文件读取
func (y *ViperComponent) Get(key string) interface{} {
	if y.isCached(key) {
		return y.getFromCache(key)
	} else {
		value := y.viper.Get(key)
		y.cache(key, value)
		return value
	}
}

// GetString 获取字符串类型的值
func (y *ViperComponent) GetString(key string) string {
	if y.isCached(key) {
		return y.getFromCache(key).(string)
	} else {
		value := y.viper.GetString(key)
		y.cache(key, value)
		return value
	}

}

// GetBool 获取布尔类型的值
func (y *ViperComponent) GetBool(key string) bool {
	if y.isCached(key) {
		return y.getFromCache(key).(bool)
	} else {
		value := y.viper.GetBool(key)
		y.cache(key, value)
		return value
	}
}

// GetInt 获取int类型的值
func (y *ViperComponent) GetInt(key string) int {
	if y.isCached(key) {
		return y.getFromCache(key).(int)
	} else {
		value := y.viper.GetInt(key)
		y.cache(key, value)
		return value
	}
}

// GetInt32 获取int32类型的值
func (y *ViperComponent) GetInt32(key string) int32 {
	if y.isCached(key) {
		return y.getFromCache(key).(int32)
	} else {
		value := y.viper.GetInt32(key)
		y.cache(key, value)
		return value
	}
}

// GetInt64 获取int64类型的值
func (y *ViperComponent) GetInt64(key string) int64 {
	if y.isCached(key) {
		return y.getFromCache(key).(int64)
	} else {
		value := y.viper.GetInt64(key)
		y.cache(key, value)
		return value
	}
}

// GetFloat64 获取浮点数类型的值
func (y *ViperComponent) GetFloat64(key string) float64 {
	if y.isCached(key) {
		return y.getFromCache(key).(float64)
	} else {
		value := y.viper.GetFloat64(key)
		y.cache(key, value)
		return value
	}
}

// GetDuration 获取time.Duration类型的值
func (y *ViperComponent) GetDuration(key string) time.Duration {
	if y.isCached(key) {
		return y.getFromCache(key).(time.Duration)
	} else {
		value := y.viper.GetDuration(key)
		y.cache(key, value)
		return value
	}
}

// GetStringSlice 获取字符串切片的值
func (y *ViperComponent) GetStringSlice(key string) []string {
	if y.isCached(key) {
		return y.getFromCache(key).([]string)
	} else {
		value := y.viper.GetStringSlice(key)
		y.cache(key, value)
		return value
	}
}
