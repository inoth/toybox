package httpgin

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

var (
	default_name = "gin"
)

type Option func(*HttpGinServer)

func defaultOption() HttpGinServer {
	return HttpGinServer{
		name:   default_name,
		ready:  true,
		sfg:    singleflight.Group{},
		engine: gin.New(),

		Port:           ":8080",
		ReadTimeout:    10,
		WriteTimeout:   10,
		MaxHeaderBytes: 10,
		TLS:            false,
	}
}

func WithName(name string) Option {
	return func(hgs *HttpGinServer) {
		hgs.name = name
	}
}

func WithPort(port string) Option {
	return func(hgs *HttpGinServer) {
		hgs.Port = port
	}
}

func WithReadTimeout(readTimeout int) Option {
	return func(hgs *HttpGinServer) {
		hgs.ReadTimeout = readTimeout
	}
}

func WithWriteTimeout(writeTimeout int) Option {
	return func(hgs *HttpGinServer) {
		hgs.WriteTimeout = writeTimeout
	}
}

func WithMaxHeaderBytes(maxHeaderBytes int) Option {
	return func(hgs *HttpGinServer) {
		hgs.MaxHeaderBytes = maxHeaderBytes
	}
}

func WithTLS(cert, key string) Option {
	return func(hgs *HttpGinServer) {
		hgs.TLS = true
		hgs.Cert = cert
		hgs.Key = key
	}
}
