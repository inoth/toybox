package yaml

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/inoth/toybox/config"
	"github.com/inoth/toybox/util/encrypt"
	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

const (
	format = "*.yaml"
)

type ConfigWithYaml struct {
	config.Option

	hash string
	p    chan<- struct{}

	cfg struct {
		Server map[string]any `toml:",inline"`
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
	cfg := &ConfigWithYaml{
		Option: o,
	}
	if err := cfg.decode(); err != nil {
		panic(err)
	}
	return cfg
}

func (ct *ConfigWithYaml) decode() (err error) {
	cfgStr, err := ct.Source.Load(format)
	if err != nil {
		return fmt.Errorf("load config err %v\n", err)
	}
	if cfgStr == "" {
		return fmt.Errorf("failed to load configuration")
	}
	ct.hash = encrypt.EncryptMd5(cfgStr)

	err = yaml.Unmarshal([]byte(cfgStr), &(ct.cfg))
	if err != nil {
		return err
	}

	return nil
}

func (ct *ConfigWithYaml) PrimitiveDecode(vals ...config.ConfigureMatcher) error {
	for i := 0; i < len(vals); i++ {
		if val, ok := ct.cfg.Server[vals[i].Name()]; ok {
			if err := mapstructure.Decode(val, vals[i]); err != nil {
				return fmt.Errorf("%s -> PrimitiveDecode error: %v", vals[i].Name(), err)
			}
		}
	}
	return nil
}

func (ct *ConfigWithYaml) Next(ctx context.Context) {
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

			err = yaml.Unmarshal([]byte(cfgStr), &(ct.cfg))
			if err != nil {
				log.Printf("decode configuration error: %v\n", err)
				continue
			}

			ct.p <- struct{}{}
		}
	}
}

func (ct *ConfigWithYaml) Probe(p chan<- struct{}) {
	if ct.p == nil {
		ct.p = p
	}
}
