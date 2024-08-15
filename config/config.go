package config

import (
	"fmt"
	"strings"

	"github.com/inoth/toybox/util"
)

const (
	Toml = "toml"
	Yaml = "yaml"
	Json = "json"
)

var (
	defaultCfg = CfgBasic{
		Remote:   false,
		CfgDir:   "config",
		FileType: "toml",
		Env:      "",
	}
)

type ConfigMate interface {
	PrimitiveDecode(vals ...ConfigureMatcher) error
}

type Configuration interface {
	Decode() error
}

type ConfigureMatcher interface {
	Name() string
}

type CfgBasic struct {
	Remote   bool
	CfgDir   string
	FileType string
	Env      string
}

func NewDefaultConfig() ConfigMate {
	return newConfig()
}

func NewConfig(cb CfgBasic) ConfigMate {
	return newConfig(cb)
}

func newConfig(cbs ...CfgBasic) ConfigMate {
	cb := util.First(defaultCfg, cbs)
	if cb.CfgDir == "" {
		cb.CfgDir = defaultCfg.CfgDir
	}
	if cb.FileType == "" {
		cb.FileType = defaultCfg.FileType
	}
	if cb.Env == "" {
		cb.Env = defaultCfg.Env
	}
	switch cb.FileType {
	case Toml:
		tomlCfg := &ConfigWithToml{
			basic: cb,
		}
		if err := tomlCfg.Decode(); err != nil {
			fmt.Printf("%v", err)
		}
		return tomlCfg
	default:
		panic(fmt.Errorf("unknown file type"))
	}
}

func loadConfig(paths []string) string {
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
	return sb.String()
}
