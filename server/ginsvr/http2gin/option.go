package http2gin

import (
	"fmt"

	"github.com/inoth/toybox/server/ginsvr"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

var (
	default_name = "gin"
)

type Option func(*Http2GinServer)

func defaultOption() Http2GinServer {
	return Http2GinServer{
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
	return func(o *Http2GinServer) {
		o.name = name
	}
}

func WithPort(port string) Option {
	return func(o *Http2GinServer) {
		o.Port = port
	}
}

func WithReadTimeout(readTimeout int) Option {
	return func(o *Http2GinServer) {
		o.ReadTimeout = readTimeout
	}
}

func WithWriteTimeout(writeTimeout int) Option {
	return func(o *Http2GinServer) {
		o.WriteTimeout = writeTimeout
	}
}

func WithMaxHeaderBytes(maxHeaderBytes int) Option {
	return func(o *Http2GinServer) {
		o.MaxHeaderBytes = maxHeaderBytes
	}
}

func WithTLS(cert, key string) Option {
	if cert == "" || key == "" {
		fmt.Println("WARN: the cert or key is not allowed to be empty.")
	}
	return func(o *Http2GinServer) {
		o.TLS = true
		o.Cert = cert
		o.Key = key
	}
}

func WithMiddleware(mids ...gin.HandlerFunc) Option {
	return func(hgs *Http2GinServer) {
		hgs.Use(mids...)
	}
}

func WithHandlers(cols ...ginsvr.Handler) Option {
	return func(hgs *Http2GinServer) {
		hgs.handlers = append(hgs.handlers, cols...)
	}
}
