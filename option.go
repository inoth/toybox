package toybox

type Option func(*ToyBox)

func WithEnv(env string) Option {
	return func(cfg *ToyBox) {
		cfg.env = env
	}
}

func WithConfDir(confDir string) Option {
	return func(cfg *ToyBox) {
		cfg.confDir = confDir
	}
}

func WithConfFileType(confFileType string) Option {
	return func(cfg *ToyBox) {
		cfg.confFileType = confFileType
	}
}
