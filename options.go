package toybox

import (
	"github.com/inoth/toybox/component"
)

// 获取本地配置
func WithCfgPath(path string) Option {
	return func(tb *ToyBox) {
		tb.cfgSource = "local"
		tb.cfgPath = path
	}
}

// 通过GET获取远程接口的配置
func WithCfgURL(url string) Option {
	return func(tb *ToyBox) {
		tb.cfgSource = "url"
		tb.cfgPath = url
	}
}

func EnableComponents(comps ...component.Component) Option {
	return func(tb *ToyBox) {
		tb.components = comps
	}
}
