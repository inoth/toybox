package config

type ConfigSource interface {
	Config(path string) string
}
