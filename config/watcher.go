package config

type Watcher interface {
	Next() ([]byte, error)
	Stop() error
}
