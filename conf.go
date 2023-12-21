package toybox

import "github.com/BurntSushi/toml"

type Conf interface {
}

type ConfWithToml struct {
	Server map[string]toml.Primitive `toml:"server" json:"server"`
}
