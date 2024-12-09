package config

import (
	"fmt"
	"strings"

	"github.com/inoth/toybox/util/file"
)

const (
	Toml = "toml"
	Yaml = "yaml"
	Json = "json"
)

const (
	DefaultDir = "config"
)

type ConfigMate interface {
	PrimitiveDecode(vals ...ConfigureMatcher) error
}

type Configuration interface {
	Decode(dir string) error
}

type ConfigureMatcher interface {
	Name() string
}

func NewDefaultConfig() ConfigMate {
	tomlCfg := &ConfigWithToml{}
	if err := tomlCfg.Decode(DefaultDir); err != nil {
		panic(err)
	}
	return tomlCfg
}
func NewConfig(opts ...Option) ConfigMate {
	o := &option{
		cfg: &ConfigWithToml{},
		dir: DefaultDir,
	}

	for _, opt := range opts {
		opt(o)
	}
	if err := o.cfg.Decode(o.dir); err != nil {
		panic(err)
	}

	return o.cfg.(ConfigMate)
}
func loadConfig(paths []string) string {
	var sb strings.Builder
	for _, path := range paths {
		buf, err := file.ReadFile(path)
		if err != nil {
			fmt.Printf("%s read file err: %v", path, err)
			continue
		}
		sb.Write(buf)
		sb.WriteString("\n")
	}
	return sb.String()
}
