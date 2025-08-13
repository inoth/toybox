package httpserver

import (
	"context"
	"crypto/tls"
	"net/http"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
)

const (
	name = "http"
)

type GinHttpServer struct {
	option
	sfg singleflight.Group
	svr *http.Server
}

func NewHttp(opts ...Option) *GinHttpServer {
	o := defaultOption
	o.engine = gin.New(func(e *gin.Engine) {
		e.ContextWithFallback = true
	})
	for _, opt := range opts {
		opt(&o)
	}
	if o.serverName == "" {
		o.serverName = name
	}
	return &GinHttpServer{
		option: o,
		sfg:    singleflight.Group{},
	}
}

func (h *GinHttpServer) Name() string {
	return h.serverName
}

func (h *GinHttpServer) Start(ctx context.Context) error {

	for _, handle := range h.handles {
		for _, r := range handle.Routers() {
			h.engine.Handle(
				r.Method,
				path.Join(handle.Prefix(), r.Path),
				append(handle.Middlewares(), r.Handle...)...,
			)
		}
	}

	h.svr = &http.Server{
		Addr:           h.Port,
		Handler:        h.engine,
		ReadTimeout:    time.Second * time.Duration(h.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(h.WriteTimeout),
		MaxHeaderBytes: 1 << uint(h.MaxHeaderBytes),
	}
	var err error
	if h.TLS {
		h.svr.TLSConfig = &tls.Config{
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
		err = h.svr.ListenAndServeTLS(h.Cert, h.Key)
	} else {
		err = h.svr.ListenAndServe()
	}
	if err != nil && err != context.Canceled && err != http.ErrServerClosed {
		return errors.Wrap(err, "start http server err")
	}
	return nil
}

func (h *GinHttpServer) Stop(ctx context.Context) error {
	return h.svr.Shutdown(ctx)
}

func (h *GinHttpServer) Do(key string, fn func() (any, error)) (v any, err error, shared bool) {
	return h.sfg.Do(key, fn)
}
