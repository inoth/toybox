package ginsvr

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/qlog"

	"github.com/inoth/toybox/validation"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
)

const (
	http3name = "gin3"
)

type GinHttp3Server struct {
	option
	sfg singleflight.Group
	svr *http3.Server
}

func NewHttp3(opts ...Option) *GinHttp3Server {
	o := option{
		ReadTimeout:    10,
		WriteTimeout:   10,
		MaxHeaderBytes: 20,
		Port:           ":9050",
		engine:         gin.New(),
		handles:        make([]Handler, 0),
	}
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

	h3.loadRouter()
	h3.loadValidation()

	if h3.Cert == "" || h3.Key == "" {
		return fmt.Errorf("server %s must be config with tls", http2name)
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
	if err != nil && err != context.Canceled {
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

func (h3 *GinHttp3Server) loadRouter() {
	for _, h := range h3.handles {
		for _, r := range h.Routers() {
			h3.engine.Handle(
				r.Method,
				h.Prefix()+"/"+r.Path,
				append(h.Middlewares(), r.Handle...)...,
			)
		}
	}
}

func (h3 *GinHttp3Server) loadValidation() {
	validation.LoadValidation(h3.validator)
}
