package config

type Source interface {
	LoadConfig() (string, error)
}
