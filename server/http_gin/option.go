package httpgin

var (
	default_name = "gin"
)

type OptionFunc func(*Option)

type Option struct {
	name string

	Port           string `toml:"port" json:"port"`
	ReadTimeout    int    `toml:"read_timeout" json:"read_timeout"`
	WriteTimeout   int    `toml:"write_timeout" json:"write_timeout"`
	MaxHeaderBytes int    `toml:"max_header_bytes" json:"max_header_bytes"`

	TLS  bool   `toml:"tls" json:"tls"`
	Cert string `toml:"cert" json:"cert"`
	Key  string `toml:"key" json:"key"`
}

func defaultOption() Option {
	return Option{
		name:           default_name,
		Port:           ":8080",
		ReadTimeout:    10,
		WriteTimeout:   10,
		MaxHeaderBytes: 10,
		TLS:            false,
	}
}

func (o *Option) WithName(name string) OptionFunc {
	return func(o *Option) {
		o.name = name
	}
}

func (o *Option) WithPort(port string) OptionFunc {
	return func(o *Option) {
		o.Port = port
	}
}
