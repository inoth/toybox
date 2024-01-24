package httpgin

import (
	"github.com/inoth/toybox/server/ginsvr"
	"github.com/inoth/toybox/server/ginsvr/validaton"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

var (
	default_name = "gin"
)

type Option func(*HttpGinServer)

func defaultOption() HttpGinServer {
	return HttpGinServer{
		name:     default_name,
		ready:    true,
		sfg:      singleflight.Group{},
		engine:   gin.New(),
		handlers: make([]ginsvr.Handler, 0),

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

func WithMiddleware(mids ...gin.HandlerFunc) Option {
	return func(hgs *HttpGinServer) {
		hgs.Use(mids...)
	}
}

func WithHandlers(cols ...ginsvr.Handler) Option {
	return func(hgs *HttpGinServer) {
		hgs.handlers = append(hgs.handlers, cols...)
	}
}

func WithValidator(valid ...validaton.Validaton) Option {
	return func(hgs *HttpGinServer) {
		hgs.validator = append(hgs.validator, valid...)
	}
}
