package httpapi

import "github.com/gin-gonic/gin"

type Option func(opt *option)

type option struct {
	ReadTimeout    int    `toml:"read_timeout" json:"read_timeout"`
	WriteTimeout   int    `toml:"write_timeout" json:"write_timeout"`
	MaxHeaderBytes int    `toml:"max_header_bytes" json:"max_header_bytes"`
	TLS            bool   `toml:"tls" json:"tls"`
	Cert           string `toml:"cert" json:"cert"`
	Key            string `toml:"key" json:"key"`
	Port           string `toml:"port" json:"port"`

	engine  *gin.Engine
	handles []Handler
}

func WithPort(port string) Option {
	return func(opt *option) {
		opt.Port = port
	}
}

func WithTLS(cert, key string) Option {
	return func(opt *option) {
		opt.TLS = true
		opt.Key = key
		opt.Cert = cert
	}
}

func WithHandlers(handles ...Handler) Option {
	return func(opt *option) {
		opt.handles = handles
	}
}

func WithGET(path string, hs ...gin.HandlerFunc) Option {
	return func(opt *option) {
		opt.engine.GET(path, hs...)
	}
}

func WithPOST(path string, hs ...gin.HandlerFunc) Option {
	return func(opt *option) {
		opt.engine.POST(path, hs...)
	}
}

func WithPUT(path string, hs ...gin.HandlerFunc) Option {
	return func(opt *option) {
		opt.engine.PUT(path, hs...)
	}
}

func WithDELETE(path string, hs ...gin.HandlerFunc) Option {
	return func(opt *option) {
		opt.engine.DELETE(path, hs...)
	}
}
