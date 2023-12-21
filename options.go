package toybox

type ToyBoxOptFunc func(cfg *ToyBoxOption)

type ToyBoxOption struct {
	env          string
	confDir      string
	confFileType string
	conf         interface{}

	// 组件
	cpts []Component
	// 服务
	svcs []Server
}

func defaultOption() ToyBoxOption {
	return ToyBoxOption{
		confDir:      "config",
		env:          "",
		confFileType: "dev",
		cpts:         make([]Component, 0),
		svcs:         make([]Server, 0),
	}
}

func (o *ToyBoxOption) WithEnv(env string) ToyBoxOptFunc {
	return func(cfg *ToyBoxOption) {
		cfg.env = env
	}
}

func (o *ToyBoxOption) WithConfDir(confDir string) ToyBoxOptFunc {
	return func(cfg *ToyBoxOption) {
		cfg.confDir = confDir
	}
}

func (o *ToyBoxOption) WithConfFileType(confFileType string) ToyBoxOptFunc {
	return func(cfg *ToyBoxOption) {
		cfg.confFileType = confFileType
	}
}

func (o *ToyBoxOption) WithAppendComponent(cpts ...Component) ToyBoxOptFunc {
	return func(cfg *ToyBoxOption) {
		cfg.cpts = append(cfg.cpts, cpts...)
	}
}

func (o *ToyBoxOption) WithAppendServer(svcs ...Server) ToyBoxOptFunc {
	return func(cfg *ToyBoxOption) {
		cfg.svcs = append(cfg.svcs, svcs...)
	}
}
