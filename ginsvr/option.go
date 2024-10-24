package ginsvr

import (
	"github.com/gin-gonic/gin"
	"github.com/inoth/toybox/validation"
)

type Option func(opt *option)

type option struct {
	serverName string

	ReadTimeout    int    `toml:"read_timeout" json:"read_timeout"`
	WriteTimeout   int    `toml:"write_timeout" json:"write_timeout"`
	MaxHeaderBytes int    `toml:"max_header_bytes" json:"max_header_bytes"`
	TLS            bool   `toml:"tls" json:"tls"`
	Cert           string `toml:"cert" json:"cert"`
	Key            string `toml:"key" json:"key"`
	Port           string `toml:"port" json:"port"`

	engine    *gin.Engine
	handles   []Handler
	validator []validation.Validation
}

func WithName(name string) Option {
	return func(opt *option) {
		opt.serverName = name
	}
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
		opt.handles = append(opt.handles, handles...)
	}
}

func WithMiddleware(handles ...gin.HandlerFunc) Option {
	return func(opt *option) {
		opt.engine.Use(handles...)
	}
}

func WithValidator(v ...validation.Validation) Option {
	return func(opt *option) {
		opt.validator = append(opt.validator, v...)
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
