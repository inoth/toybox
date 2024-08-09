package metric

type Option func(opt *option)

type option struct {
	Port      string `toml:"port"`
	Subsystem string `toml:"subsystem"`
	Namespace string `toml:"namespace"`
}

func WithPort(port string) Option {
	return func(opt *option) {
		opt.Port = port
	}
}

func WithSubsystem(subsystem string) Option {
	return func(opt *option) {
		opt.Subsystem = subsystem
	}
}

func WithNamespace(namespace string) Option {
	return func(opt *option) {
		opt.Namespace = namespace
	}
}
