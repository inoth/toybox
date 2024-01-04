package http2gin

import (
	"context"
	"fmt"
	"github/inoth/toybox"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

type Http2GinServer struct {
	name   string
	ready  bool
	sfg    singleflight.Group
	engine *gin.Engine

	Port           string `toml:"port" json:"port"`
	ReadTimeout    int    `toml:"read_timeout" json:"read_timeout"`
	WriteTimeout   int    `toml:"write_timeout" json:"write_timeout"`
	MaxHeaderBytes int    `toml:"max_header_bytes" json:"max_header_bytes"`

	TLS  bool   `toml:"tls" json:"tls"`
	Cert string `toml:"cert" json:"cert"`
	Key  string `toml:"key" json:"key"`
}

func NewHttp2(opts ...Option) toybox.Option {
	h2gs := defaultOption()
	for _, opt := range opts {
		opt(&h2gs)
	}
	return func(tb *toybox.ToyBox) {
		tb.AppendServer(&h2gs)
	}
}

func (h2gs *Http2GinServer) Ready() bool {
	return h2gs.ready
}

func (h2gs *Http2GinServer) Run(ctx context.Context) error {
	if !h2gs.ready {
		return fmt.Errorf("server %s not ready", h2gs.name)
	}
	if !h2gs.TLS || h2gs.Cert == "" || h2gs.Key == "" {
		return fmt.Errorf("server %s must be config with tls", h2gs.name)
	}
	return nil
}

func (hgs *Http2GinServer) Engine() *gin.Engine {
	return hgs.engine
}

func (hgs *Http2GinServer) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	return hgs.sfg.Do(key, fn)
}

func (hgs *Http2GinServer) Use(middleware ...gin.HandlerFunc) {
	hgs.engine.Use(middleware...)
}

func (hgs *Http2GinServer) Group(relativePath string, handlers ...gin.HandlerFunc) {
	hgs.engine.Group(relativePath, handlers...)
}

func (hgs *Http2GinServer) GET(relativePath string, handlers ...gin.HandlerFunc) {
	hgs.engine.GET(relativePath, handlers...)
}

func (hgs *Http2GinServer) POST(relativePath string, handlers ...gin.HandlerFunc) {
	hgs.engine.POST(relativePath, handlers...)
}

func (hgs *Http2GinServer) PUT(relativePath string, handlers ...gin.HandlerFunc) {
	hgs.engine.PUT(relativePath, handlers...)
}

func (hgs *Http2GinServer) DELETE(relativePath string, handlers ...gin.HandlerFunc) {
	hgs.engine.DELETE(relativePath, handlers...)
}
