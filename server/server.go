package server

type ServerOption func(sv Server)

type Server interface {
	Name() string
	RequiredComponent() []string
	Start() error
}
