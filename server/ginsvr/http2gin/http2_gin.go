package http2gin

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/inoth/toybox"
	"github.com/inoth/toybox/server/ginsvr"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/http2"
	"golang.org/x/sync/singleflight"
)

type Http2GinServer struct {
	name     string
	ready    bool
	sfg      singleflight.Group
	engine   *gin.Engine
	handlers []ginsvr.Handler

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

func (hgs *Http2GinServer) IsReady() {
	hgs.ready = true
}

func (h2gs Http2GinServer) Ready() bool {
	return h2gs.ready
}

func (h2gs Http2GinServer) Name() string {
	return h2gs.name
}

func (h2gs *Http2GinServer) Run(ctx context.Context) error {
	if !h2gs.ready {
		return fmt.Errorf("server %s not ready", h2gs.name)
	}
	if !h2gs.TLS || h2gs.Cert == "" || h2gs.Key == "" {
		return fmt.Errorf("server %s must be config with tls", h2gs.name)
	}

	http2svc := &http2.Server{}
	svc := &http.Server{
		Addr:           h2gs.Port,
		Handler:        h2gs.engine,
		ReadTimeout:    time.Second * time.Duration(h2gs.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(h2gs.WriteTimeout),
		MaxHeaderBytes: 1 << uint(h2gs.MaxHeaderBytes),
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			},
		},
	}

	err := http2.ConfigureServer(svc, http2svc)
	if err != nil {
		return err
	}

	return svc.ListenAndServeTLS(h2gs.Cert, h2gs.Key)
}

func (hgs *Http2GinServer) loadRouter() {
	for _, handle := range hgs.handlers {
		for _, router := range handle.Routers() {
			hgs.engine.Handle(
				router.Method,
				handle.Prefix()+"/"+router.Path,
				append(handle.Middlewares(), router.Handle...)...,
			)
		}
	}
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

func (hgs *Http2GinServer) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return hgs.engine.Group(relativePath, handlers...)
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

func (hgs *Http2GinServer) Handle(method, path string, handlers ...gin.HandlerFunc) {
	hgs.engine.Handle(method, path, handlers...)
}
