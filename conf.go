package toybox

import (
	"fmt"
	"strings"

	"github.com/inoth/toybox/util"

	"github.com/BurntSushi/toml"
)

const (
	Toml = "toml"
	Yaml = "yaml"
	Json = "json"
)

var (
	appconfig Conf
)

var (
	default_config = SetConfig{
		ConfDir:  "config",
		Env:      "",
		FileType: "toml",
	}
)

type Conf interface {
	Path() string
	Pattern() string
	SetConfig(cfg string)
	Decode() error
	GetBaseConf() map[string]interface{}
}

type conf struct {
	Base      map[string]interface{}    `toml:"base"`
	Component map[string]toml.Primitive `toml:"component"`
	Server    map[string]toml.Primitive `toml:"server"`
}

type ConfigMate interface {
	PrimitiveDecodeComponent(cpts ...Component) error
	PrimitiveDecodeServer(svrs ...Server) error
}

type SetConfig struct {
	ConfDir  string
	Env      string
	FileType string
}

func newConfig(cfgs ...SetConfig) Conf {
	cfg := util.First(default_config, cfgs)
	switch cfg.FileType {
	case Toml:
		return &ConfWithToml{
			ConfDir: cfg.ConfDir,
			Env:     cfg.Env,
		}
	default:
		panic(fmt.Errorf("unknown file type"))
	}
}

func loadConfig(cfgs ...SetConfig) Conf {
	cfg := newConfig(cfgs...)
	paths, err := util.PathGlobPattern(cfg.Path() + cfg.Pattern())
	if err != nil {
		panic(fmt.Errorf("no configuration available"))
	}
	var sb strings.Builder
	for _, path := range paths {
		buf, err := util.ReadFile(path)
		if err != nil {
			fmt.Printf("%s read file err: %v", path, err)
			continue
		}
		sb.Write(buf)
		sb.WriteString("\n")
	}

	tomlStr := sb.String()
	cfg.SetConfig(tomlStr)

	if err := cfg.Decode(); err != nil {
		panic(fmt.Errorf("parsing configuration err: %v", err))
	}
	return cfg
}

func NewConfig() ConfigMate {
	return loadConfig().(ConfigMate)
}

func WithConfigMate(mate ConfigMate) Option {
	return func(tb *ToyBox) {
		tb.mate = mate
	}
}

func WithLoadConf(cfgs ...SetConfig) Option {
	cfg := loadConfig(cfgs...)
	appconfig = cfg
	return func(tb *ToyBox) {
		tb.mate = cfg.(ConfigMate)
	}
}

func GetConfMate() ConfigMate {
	if appconfig == nil {
		return nil
	}
	return appconfig.(ConfigMate)
}

func GetString(key string) string {
	if res, ok := util.GetStringValue(appconfig.GetBaseConf(), key); ok {
		return res
	}
	return ""
}

func GetInt(key string) int {
	if res, ok := util.GetIntValue(appconfig.GetBaseConf(), key); ok {
		return res
	}
	return 0
}

func GetFloat(key string) float64 {
	if res, ok := util.GetFloatValue(appconfig.GetBaseConf(), key); ok {
		return res
	}
	return 0
}

func GetBool(key string) bool {
	if res, ok := util.GetBoolValue(appconfig.GetBaseConf(), key); ok {
		return res
	}
	return false
}

func GetStringSlice(key string) []string {
	if res, ok := util.GetInterfaceSlice(appconfig.GetBaseConf(), key); ok {
		r := make([]string, 0, len(res))
		for _, rs := range res {
			if val, ok := rs.(string); ok {
				r = append(r, val)
			}
		}
		return r
	}
	return nil
}
