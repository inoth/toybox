package toybox

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type ConfWithToml struct {
	ConfDir string
	Env     string

	mate      toml.MetaData
	appConfig string
	conf      conf
}

func (conf *ConfWithToml) Path() string {
	if conf.Env != "" {
		return conf.ConfDir + "/" + conf.Env + "/"
	}
	return conf.ConfDir + "/"
}

func (conf *ConfWithToml) Pattern() string {
	return "*.toml"
}

func (conf *ConfWithToml) SetConfig(cfg string) {
	conf.appConfig = cfg
}

func (conf *ConfWithToml) Decode() error {
	mate, err := toml.Decode(conf.appConfig, &(conf.conf))
	if err != nil {
		return err
	}
	conf.mate = mate
	return nil
}

func (conf *ConfWithToml) GetBaseConf() map[string]interface{} {
	return conf.conf.Base
}

func (conf *ConfWithToml) PrimitiveDecodeComponent(cpts ...Component) error {
	for i := 0; i < len(cpts); i++ {
		if val, ok := conf.conf.Component[cpts[i].Name()]; ok {
			if err := conf.mate.PrimitiveDecode(val, cpts[i]); err != nil {
				return fmt.Errorf("%s -> PrimitiveDecode error: %v", cpts[i].Name(), err)
			}
			cpts[i].IsReady()
		}
	}
	return nil
}

func (conf *ConfWithToml) PrimitiveDecodeServer(svrs ...Server) error {
	for i := 0; i < len(svrs); i++ {
		if val, ok := conf.conf.Server[svrs[i].Name()]; ok {
			if err := conf.mate.PrimitiveDecode(val, svrs[i]); err != nil {
				return fmt.Errorf("%s -> PrimitiveDecode error: %v", svrs[i].Name(), err)
			}
			svrs[i].IsReady()
		}
	}
	return nil
}
