package toybox

import (
	"fmt"
	"github/inoth/toybox/util"
	"strings"
	"sync"
)

const (
	Toml = "toml"
	Yaml = "yaml"
	Json = "json"
)

var (
	conf = sync.Map{}
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

// type ConfWithYaml struct {
// 	ConfDir string
// 	Env     string
// }

// func (conf ConfWithYaml) Path() string {
// 	if conf.Env != "" {
// 		return conf.ConfDir + "/" + conf.Env + "/"
// 	}
// 	return conf.ConfDir + "/"
// }

// func (conf ConfWithYaml) Pattern() string {
// 	return "*.yaml"
// }

// type ConfWithJson struct {
// 	ConfDir string
// 	Env     string
// }

// func (conf ConfWithJson) Path() string {
// 	if conf.Env != "" {
// 		return conf.ConfDir + "/" + conf.Env + "/"
// 	}
// 	return conf.ConfDir + "/"
// }

// func (conf ConfWithJson) Pattern() string {
// 	return "*.json"
// }

func newConfig(cfgs ...SetConfig) Conf {
	cfg := util.First(default_config, cfgs)
	switch cfg.FileType {
	case Toml:
		return &ConfWithToml{
			ConfDir: cfg.ConfDir,
			Env:     cfg.Env,
		}
	// case Yaml:
	// 	return &ConfWithYaml{
	// 		ConfDir: cfg.ConfDir,
	// 		Env:     cfg.Env,
	// 	}
	// case Json:
	// 	return &ConfWithJson{
	// 		ConfDir: cfg.ConfDir,
	// 		Env:     cfg.Env,
	// 	}
	default:
		panic("unknown file type")
	}
}

func WithLoadConf(cfgs ...SetConfig) Option {
	cfg := newConfig(cfgs...)

	paths, err := util.PathGlobPattern(cfg.Path() + cfg.Pattern())
	if err != nil {
		panic("no configuration available")
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
		panic("parsing configuration err: " + err.Error())
	}

	return func(tb *ToyBox) {
		tb.mate = cfg.(ConfigMate)
	}
}
