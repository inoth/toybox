package config

type Option func(opt *option)

type option struct {
	dir string
	cfg Configuration
}

func WithConfiguration(cfg Configuration) Option {
	return func(opt *option) {
		opt.cfg = cfg
	}
}

func WithConfigDir(dir string) Option {
	return func(opt *option) {
		opt.dir = dir
	}
}
