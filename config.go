package toybox

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type appConfig struct {
	Base       toml.Primitive              `json:"base" toml:"base"`
	Components map[string][]toml.Primitive `json:"components" toml:"components"`
}

func (tb *ToyBox) resolveConfig(cfgByte []byte) error {
	switch tb.cfgType {
	case Toml:
		return tb.resolveTomlConfig(cfgByte)
	case Json:
		return fmt.Errorf("%v format is not supported yet", tb.cfgType)
	case Yaml:
		return fmt.Errorf("%v format is not supported yet", tb.cfgType)
	default:
		return fmt.Errorf("unknown format is not supported yet")
	}
}

// TODO: 存在全局变量，之后考虑添加上，优先吧组件的配置注入完成
func (tb *ToyBox) resolveTomlConfig(cfgByte []byte) error {
	var cfg appConfig
	mata, err := toml.Decode(string(cfgByte), &cfg)
	if err != nil {
		return err
	}
	for i := 0; i < len(tb.components); i++ {
		prs := cfg.Components[tb.components[i].Name()]
		for _, pr := range prs {
			if err := mata.PrimitiveDecode(pr, tb.components[i]); err != nil {
				return fmt.Errorf("cannot load component correctly %v: %v", tb.components[i].Name(), err)
			}
			// TODO: 加入校验逻辑
		}
	}
	return nil
}
