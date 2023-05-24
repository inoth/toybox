package toybox

import (
	"github.com/inoth/toybox/component"
)

// 获取本地组件配置
func WithComponentCfgPath(path string) Option {
	return func(tb *ToyBox) {
		tb.cfgSource = "local"
		tb.cfgPath = path
	}
}

// 通过GET获取远程接口的组件配置
func WithComponentCfgURL(url string) Option {
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
