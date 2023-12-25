package toybox

type OptionFunc func(*Option)

type Option struct {
	env          string
	confDir      string
	confFileType string
	conf         Conf

	// 组件
	cpts []Component
	// 服务
	svcs []Server
}

func defaultOption() Option {
	return Option{
		confDir:      "config",
		env:          "",
		confFileType: "dev",
		cpts:         make([]Component, 0),
		svcs:         make([]Server, 0),
	}
}

func (o *Option) WithEnv(env string) OptionFunc {
	return func(cfg *Option) {
		cfg.env = env
	}
}

func (o *Option) WithConfDir(confDir string) OptionFunc {
	return func(cfg *Option) {
		cfg.confDir = confDir
	}
}

func (o *Option) WithConfFileType(confFileType string) OptionFunc {
	return func(cfg *Option) {
		cfg.confFileType = confFileType
	}
}

func (o *Option) WithAppendComponent(cpts ...Component) OptionFunc {
	return func(cfg *Option) {
		cfg.cpts = append(cfg.cpts, cpts...)
	}
}

func (o *Option) WithAppendServer(svcs ...Server) OptionFunc {
	return func(cfg *Option) {
		cfg.svcs = append(cfg.svcs, svcs...)
	}
}
