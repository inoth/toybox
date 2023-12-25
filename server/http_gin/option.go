package httpgin

var (
	default_name = "gin"
)

type Option func(*option)

type option struct {
	name string

	Port           string `toml:"port" json:"port"`
	ReadTimeout    int    `toml:"read_timeout" json:"read_timeout"`
	WriteTimeout   int    `toml:"write_timeout" json:"write_timeout"`
	MaxHeaderBytes int    `toml:"max_header_bytes" json:"max_header_bytes"`

	TLS  bool   `toml:"tls" json:"tls"`
	Cert string `toml:"cert" json:"cert"`
	Key  string `toml:"key" json:"key"`
}

func defaultOption() option {
	return option{
		name:           default_name,
		Port:           ":8080",
		ReadTimeout:    10,
		WriteTimeout:   10,
		MaxHeaderBytes: 10,
		TLS:            false,
	}
}

func (o *option) WithName(name string) Option {
	return func(o *option) {
		o.name = name
	}
}

func (o *option) WithPort(port string) Option {
	return func(o *option) {
		o.Port = port
	}
}

func (o *option) WithReadTimeout(readTimeout int) Option {
	return func(o *option) {
		o.ReadTimeout = readTimeout
	}
}

func (o *option) WithWriteTimeout(writeTimeout int) Option {
	return func(o *option) {
		o.WriteTimeout = writeTimeout
	}
}

func (o *option) WithMaxHeaderBytes(maxHeaderBytes int) Option {
	return func(o *option) {
		o.MaxHeaderBytes = maxHeaderBytes
	}
}

func (o *option) WithTLS(tls bool, cert, key string) Option {
	return func(o *option) {
		o.TLS = tls
		o.Cert = cert
		o.Key = key
	}
}
