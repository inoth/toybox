package config

const (
	DefaultDir = "config"
)

type ConfigMate interface {
	PrimitiveDecode(vals ...ConfigureMatcher) error
}

type ConfigureMatcher interface {
	Name() string
}

// type Configuration interface {
// 	Decode() error
// }
