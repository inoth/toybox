package property

type Option func(opt *option)

type option struct {
	Port string `toml:"port"`
}

func WithPort(port string) Option {
	return func(opt *option) {
		opt.Port = port
	}
}
