package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/inoth/toybox/util/file"
)

type ConfigWithToml struct {
	mate toml.MetaData
	cfg  struct {
		Server map[string]toml.Primitive `toml:"server"`
	}
}

func (ct *ConfigWithToml) Decode(dir string) error {
	cfgEnv := os.Getenv("CONFIG_ENV")
	if cfgEnv != "" {
		dir = filepath.Join(dir, cfgEnv)
	}
	paths, err := file.PathGlobPattern(filepath.Join(dir, "*.toml"))
	if err != nil {
		panic(fmt.Errorf("no configuration available"))
	}
	cfgStr := loadConfig(paths)
	if cfgStr == "" {
		return fmt.Errorf("failed to load configuration")
	}

	ct.mate, err = toml.Decode(cfgStr, &(ct.cfg))
	if err != nil {
		return err
	}

	return nil
}

func (ct *ConfigWithToml) PrimitiveDecode(vals ...ConfigureMatcher) error {
	for i := 0; i < len(vals); i++ {
		if val, ok := ct.cfg.Server[vals[i].Name()]; ok {
			if err := ct.mate.PrimitiveDecode(val, vals[i]); err != nil {
				return fmt.Errorf("%s -> PrimitiveDecode error: %v", vals[i].Name(), err)
			}
		}
	}
	return nil
}
