package httpgin

import (
	"context"
	"crypto/tls"
	"fmt"
	"github/inoth/toybox"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

type HttpGinServer struct {
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

func NewHttpGin(opts ...Option) toybox.Option {
	hgs := defaultOption()
	for _, opt := range opts {
		opt(&hgs)
	}
	return func(tb *toybox.ToyBox) {
		tb.AppendServer(&hgs)
	}
}

func (hgs *HttpGinServer) IsReady() {
	hgs.ready = true
}

func (hgs HttpGinServer) Ready() bool {
	return hgs.ready
}

func (hgs HttpGinServer) Name() string {
	return hgs.name
}

func (hgs *HttpGinServer) Run(ctx context.Context) error {
	if !hgs.ready {
		return fmt.Errorf("server %s not ready", hgs.name)
	}

	svr := &http.Server{
		Addr:           hgs.Port,
		Handler:        hgs.engine,
		ReadTimeout:    time.Second * time.Duration(hgs.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(hgs.WriteTimeout),
		MaxHeaderBytes: 1 << uint(hgs.MaxHeaderBytes),
	}

	if hgs.TLS {
		svr.TLSConfig = &tls.Config{
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
		}
		return svr.ListenAndServeTLS(hgs.Cert, hgs.Key)
	} else {
		return svr.ListenAndServe()
	}
}

func (hgs *HttpGinServer) Engine() *gin.Engine {
	return hgs.engine
}

func (hgs *HttpGinServer) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool) {
	return hgs.sfg.Do(key, fn)
}

func (hgs *HttpGinServer) Use(middleware ...gin.HandlerFunc) {
	hgs.engine.Use(middleware...)
}

func (hgs *HttpGinServer) Group(relativePath string, handlers ...gin.HandlerFunc) {
	hgs.engine.Group(relativePath, handlers...)
}

func (hgs *HttpGinServer) GET(relativePath string, handlers ...gin.HandlerFunc) {
	hgs.engine.GET(relativePath, handlers...)
}

func (hgs *HttpGinServer) POST(relativePath string, handlers ...gin.HandlerFunc) {
	hgs.engine.POST(relativePath, handlers...)
}

func (hgs *HttpGinServer) PUT(relativePath string, handlers ...gin.HandlerFunc) {
	hgs.engine.PUT(relativePath, handlers...)
}

func (hgs *HttpGinServer) DELETE(relativePath string, handlers ...gin.HandlerFunc) {
	hgs.engine.DELETE(relativePath, handlers...)
}
