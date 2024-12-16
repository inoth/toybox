package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/inoth/toybox/util/encrypt"
	"github.com/inoth/toybox/util/file"
)

type ConfigWithToml struct {
	paths    []string
	hash     string
	interval int
	p        chan<- struct{}
	source   Source

	mate toml.MetaData
	cfg  struct {
		Server map[string]toml.Primitive `toml:"server"`
	}
}

func (ct *ConfigWithToml) Decode(dir string) (err error) {
	cfgEnv := os.Getenv("CONFIG_ENV")
	if cfgEnv != "" {
		dir = filepath.Join(dir, cfgEnv)
	}
	ct.paths, err = file.PathGlobPattern(filepath.Join(dir, "*.toml"))
	if err != nil {
		panic(fmt.Errorf("no configuration available"))
	}
	cfgStr := loadConfig(ct.paths)
	if cfgStr == "" {
		return fmt.Errorf("failed to load configuration")
	}
	ct.hash = encrypt.EncryptMd5(cfgStr)

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

func (ct *ConfigWithToml) Next(ctx context.Context) {
	ticker := time.NewTicker(time.Second * time.Duration(ct.interval))
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Println("checking configuration...")

			cfgStr := loadConfig(ct.paths)
			if cfgStr == "" {
				log.Println("configuration is empty")
				continue
			}

			hash := encrypt.EncryptMd5(cfgStr)
			if hash == ct.hash {
				continue
			}
			ct.hash = hash

			var err error
			ct.mate, err = toml.Decode(cfgStr, &(ct.cfg))
			if err != nil {
				log.Printf("decode configuration error: %v\n", err)
				continue
			}

			ct.p <- struct{}{}
		}
	}
}

func (ct *ConfigWithToml) Probe(p chan<- struct{}) {
	if ct.p == nil {
		ct.p = p
	}
}
