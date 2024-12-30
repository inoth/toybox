package config

type Source interface {
	Load() (string, error)
}
