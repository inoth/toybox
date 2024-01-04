package toybox

type Option func(*option)

type option struct {
	env          string
	confDir      string
	confFileType string
	conf         Conf

	// 组件
	cpts []Component
	// 服务
	svcs []Server
}

func defaultOption() option {
	return option{
		confDir:      "config",
		env:          "",
		confFileType: "dev",
		cpts:         make([]Component, 0),
		svcs:         make([]Server, 0),
	}
}

func WithEnv(env string) Option {
	return func(cfg *option) {
		cfg.env = env
	}
}

func WithConfDir(confDir string) Option {
	return func(cfg *option) {
		cfg.confDir = confDir
	}
}

func WithConfFileType(confFileType string) Option {
	return func(cfg *option) {
		cfg.confFileType = confFileType
	}
}

func WithAppendComponent(cpts ...Component) Option {
	return func(cfg *option) {
		cfg.cpts = append(cfg.cpts, cpts...)
	}
}

func WithAppendServer(svcs ...Server) Option {
	return func(cfg *option) {
		cfg.svcs = append(cfg.svcs, svcs...)
	}
}
