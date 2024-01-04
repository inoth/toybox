package toybox

type Option func(*ToyBox)

func defaultOption() ToyBox {
	return ToyBox{
		confDir:      "config",
		env:          "",
		confFileType: "dev",
		cpts:         make([]Component, 0),
		svcs:         make([]Server, 0),
	}
}

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
