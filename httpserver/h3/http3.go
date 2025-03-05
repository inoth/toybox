package ginserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/qlog"
	"golang.org/x/sync/singleflight"
)

const (
	http3name = "http3"
)

type GinHttp3Server struct {
	option
	sfg singleflight.Group
	svr *http3.Server
}

func NewHttp3(opts ...Option) *GinHttp3Server {
	o := defaultOption
	o.engine = gin.New(func(e *gin.Engine) {
		e.ContextWithFallback = true
	})
	for _, opt := range opts {
		opt(&o)
	}
	if o.serverName == "" {
		o.serverName = http3name
	}
	return &GinHttp3Server{
		option: o,
		sfg:    singleflight.Group{},
	}
}

func (h3 *GinHttp3Server) Name() string {
	return h3.serverName
}

func (h3 *GinHttp3Server) Start(ctx context.Context) error {

	if h3.Cert == "" || h3.Key == "" {
		return fmt.Errorf("server %s must be config with tls", http3name)
	}

	for _, h := range h3.handles {
		for _, r := range h.Routers() {
			h3.engine.Handle(
				r.Method,
				h.Prefix()+"/"+r.Path,
				append(h.Middlewares(), r.Handle...)...,
			)
		}
	}

	h3.svr = &http3.Server{
		Addr:           h3.Port,
		Handler:        h3.engine,
		MaxHeaderBytes: 1 << uint(h3.MaxHeaderBytes),
		QUICConfig: &quic.Config{
			Tracer: qlog.DefaultConnectionTracer,
		},
	}
	err := h3.svr.ListenAndServeTLS(h3.Cert, h3.Key)
	if err != nil && err != context.Canceled && err != http.ErrServerClosed {
		return errors.Wrap(err, "start http3 with udp server err")
	}
	return nil
}

func (h3 *GinHttp3Server) Stop(ctx context.Context) error {
	return h3.svr.Close()
}

func (h3 *GinHttp3Server) Do(key string, fn func() (any, error)) (v any, err error, shared bool) {
	return h3.sfg.Do(key, fn)
}
