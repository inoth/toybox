package config

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var (
	Cfg *ViperComponent
)

type ViperComponent struct {
	viperMap      map[string]*viper.Viper
	ConfKeyPrefix string
}

func (m *ViperComponent) Init() error {
	if len(m.ConfKeyPrefix) <= 0 {
		m.ConfKeyPrefix = os.Getenv("GORUNEVN")
		if len(m.ConfKeyPrefix) <= 0 {
			m.ConfKeyPrefix = "config/prod"
		} else {
			m.ConfKeyPrefix = "config/" + m.ConfKeyPrefix
		}
	}
	f, err := os.Open(m.ConfKeyPrefix + "/")
	if err != nil {
		return err
	}
	fileList, err := f.Readdir(1024)
	if err != nil {
		return err
	}
	for _, f0 := range fileList {
		if f0.IsDir() {
			continue
		}
		v := viper.New()
		v.SetConfigType("yaml")
		v.AddConfigPath(m.ConfKeyPrefix)
		pathArr := strings.Split(f0.Name(), ".")
		v.SetConfigName(pathArr[0])
		if err := v.ReadInConfig(); err != nil {
			return err
		}
		if m.viperMap == nil {
			m.viperMap = make(map[string]*viper.Viper)
		}
		m.viperMap[pathArr[0]] = v
	}
	Cfg = m
	return nil
}

//获取get配置信息
func (m *ViperComponent) GetString(key string) string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return ""
	}
	v, ok := m.viperMap[keys[0]]
	if !ok {
		return ""
	}
	confString := v.GetString(strings.Join(keys[1:], "."))
	return confString
}

//获取get配置信息
func (m *ViperComponent) GetStringMap(key string) map[string]interface{} {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := m.viperMap[keys[0]]
	conf := v.GetStringMap(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func (m *ViperComponent) Get(key string) interface{} {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := m.viperMap[keys[0]]
	conf := v.Get(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func (m *ViperComponent) GetBool(key string) bool {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return false
	}
	v := m.viperMap[keys[0]]
	conf := v.GetBool(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func (m *ViperComponent) GetFloat64(key string) float64 {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := m.viperMap[keys[0]]
	conf := v.GetFloat64(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func (m *ViperComponent) GetInt(key string) int {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := m.viperMap[keys[0]]
	conf := v.GetInt(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func (m *ViperComponent) GetStringMapString(key string) map[string]string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := m.viperMap[keys[0]]
	conf := v.GetStringMapString(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func (m *ViperComponent) GetStringSlice(key string) []string {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return nil
	}
	v := m.viperMap[keys[0]]
	conf := v.GetStringSlice(strings.Join(keys[1:], "."))
	return conf
}

//获取get配置信息
func (m *ViperComponent) GetTime(key string) time.Time {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return time.Now()
	}
	v := m.viperMap[keys[0]]
	conf := v.GetTime(strings.Join(keys[1:], "."))
	return conf
}

//获取时间阶段长度
func (m *ViperComponent) GetDuration(key string) time.Duration {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return 0
	}
	v := m.viperMap[keys[0]]
	conf := v.GetDuration(strings.Join(keys[1:], "."))
	return conf
}

//是否设置了key
func (m *ViperComponent) IsSet(key string) bool {
	keys := strings.Split(key, ".")
	if len(keys) < 2 {
		return false
	}
	v := m.viperMap[keys[0]]
	conf := v.IsSet(strings.Join(keys[1:], "."))
	return conf
}
