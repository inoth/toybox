package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/inoth/toybox/config"
	"github.com/inoth/toybox/util/encrypt"
)

const (
	format = "*.toml"
)

type ConfigWithToml struct {
	config.Option

	hash string
	p    chan<- struct{}

	mate toml.MetaData
	cfg  struct {
		Server map[string]toml.Primitive `toml:"server"`
	}
}

func NewConfiguration(opts ...config.Options) config.ConfigMate {
	o := config.Option{
		Interval: 0,
	}
	for _, opt := range opts {
		opt(&o)
	}
	if o.Source == nil {
		panic(fmt.Errorf("the configuration source is not set"))
	}
	cfg := &ConfigWithToml{
		Option: o,
	}
	if err := cfg.decode(); err != nil {
		panic(err)
	}
	return cfg
}

func (ct *ConfigWithToml) decode() (err error) {
	cfgStr, err := ct.Source.Load(format)
	if err != nil {
		return fmt.Errorf("load config err %v\n", err)
	}
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

func (ct *ConfigWithToml) PrimitiveDecode(vals ...config.ConfigureMatcher) error {
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
	if ct.Interval <= 0 {
		return
	}
	ticker := time.NewTicker(time.Second * time.Duration(ct.Interval))
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Println("checking configuration...")

			cfgStr, err := ct.Source.Load(format)
			if err != nil {
				log.Printf("load config err %v\n", err)
				continue
			}
			if cfgStr == "" {
				log.Println("configuration is empty")
				continue
			}

			hash := encrypt.EncryptMd5(cfgStr)
			if hash == ct.hash {
				continue
			}
			ct.hash = hash

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
