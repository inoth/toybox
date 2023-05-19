package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/inoth/toybox/common"
	"github.com/inoth/toybox/component"

	"github.com/BurntSushi/toml"
)

var (
	componentName = "config"
	configOnce    sync.Once
	Cfg           *ConfigComponent
)

type Option func(ctx context.Context, cfg *ConfigComponent)

type globalConfig struct {
	Global toml.Primitive `json:"global" toml:"global"`
}

type ConfigComponent struct {
	// 热更重新拉起次数
	count      int
	context    context.Context
	cancelFunc context.CancelFunc
	m          *sync.RWMutex

	Interval int    `toml:"interval" json:"interval" yaml:"interval"`
	CfgPath  string `toml:"cfgPath" json:"cfgPath" yaml:"cfgPath"`

	// 存放从 global 解析出来的数据
	appConfig map[string]interface{}

	currentConfigContentHash string
}

func New(opts ...Option) component.Component {
	cfg := &ConfigComponent{
		count:     1,
		m:         new(sync.RWMutex),
		appConfig: make(map[string]interface{}),
	}

	cfg.context, cfg.cancelFunc = context.WithCancel(context.Background())

	for _, opt := range opts {
		opt(cfg.context, cfg)
	}

	return cfg
}

func (cc *ConfigComponent) Name() string {
	return componentName
}
func (cc *ConfigComponent) String() string {
	// buf, _ := json.Marshal(cc)
	buf, _ := json.Marshal(cc.appConfig)
	return string(buf)
}

func (cc *ConfigComponent) Close() error {
	cc.cancelFunc()
	return nil
}

// TODO: 后续添加配置热更功能
func (cc *ConfigComponent) Init() (err error) {
	// 加载配置
	configOnce.Do(func() {
		var cfgByte []byte
		cfgByte, err = os.ReadFile(cc.CfgPath)
		if err != nil {
			err = fmt.Errorf("读取配置文件失败,Err:%s", err.Error())
			return
		}
		if err = cc.resolveTomlConfig(cfgByte); err != nil {
			return
		}
		cc.watchToml()
		Cfg = cc
	})
	return
}

func (cc *ConfigComponent) resolveTomlConfig(cfgByte []byte) error {
	var globalCfg globalConfig
	mata, err := toml.Decode(string(cfgByte), &globalCfg)
	if err != nil {
		return err
	}

	cc.m.RLock()
	defer cc.m.RUnlock()

	var appConfig map[string]interface{}
	if err := mata.PrimitiveDecode(globalCfg.Global, &appConfig); err != nil {
		return err
	}
	for key := range cc.appConfig {
		delete(cc.appConfig, key)
	}
	for key, value := range appConfig {
		cc.appConfig[key] = value
	}

	cc.currentConfigContentHash = common.Md5(cfgByte)
	return nil
}

func (cc *ConfigComponent) GetString(key string) string {
	if res, ok := common.GetStringValue(cc.appConfig, key); ok {
		return res
	}
	return ""
}

func (cc *ConfigComponent) GetInt(key string) int {
	if res, ok := common.GetIntValue(cc.appConfig, key); ok {
		return res
	}
	return 0
}

func (cc *ConfigComponent) GetFloat(key string) float64 {
	if res, ok := common.GetFloatValue(cc.appConfig, key); ok {
		return res
	}
	return 0
}

func (cc *ConfigComponent) GetBool(key string) bool {
	if res, ok := common.GetBoolValue(cc.appConfig, key); ok {
		return res
	}
	return false
}

func (cc *ConfigComponent) GetStringSlice(key string) []string {
	if res, ok := common.GetStringSlice(cc.appConfig, key); ok {
		return res
	}
	return nil
}

// 监听配置文件进行热更操作
func (cc *ConfigComponent) watchToml() {
	go func() {
		ticker := time.NewTicker(time.Duration(cc.Interval) * time.Second)
		for {
			select {
			case <-cc.context.Done():
				return
			case <-ticker.C:
				cfgByte, err := os.ReadFile(cc.CfgPath)
				if err != nil {
					fmt.Printf("读取配置文件失败,Err:%s\n", err.Error())
					continue
				}
				// 检查hash是否一致
				if cc.currentConfigContentHash == common.Md5(cfgByte) {
					continue
				}
				if err = cc.resolveTomlConfig(cfgByte); err != nil {
					fmt.Printf("更新配置文件失败, Err:%s\n", err.Error())
				}
				fmt.Println("更新配置文件更新成功")
			}
		}
	}()
}
