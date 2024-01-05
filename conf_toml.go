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
	conf      struct {
		Base      map[string]toml.Primitive `toml:"base"`
		Component map[string]toml.Primitive `toml:"component"`
		Server    map[string]toml.Primitive `toml:"server"`
	}
}

func (conf ConfWithToml) Path() string {
	if conf.Env != "" {
		return conf.ConfDir + "/" + conf.Env + "/"
	}
	return conf.ConfDir + "/"
}

func (conf ConfWithToml) Pattern() string {
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

func (conf *ConfWithToml) PrimitiveDecodeComponent(cpts ...Component) error {
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
