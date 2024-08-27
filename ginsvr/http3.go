package ginsvr

import (
	"context"
	"crypto/tls"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"
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
		Port:    ":9050",
		engine:  gin.New(),
		handles: make([]Handler, 0),
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

	if err := h3.loadValidation(); err != nil {
		return errors.Wrap(err, "loadValidation() failed")
	}

	if h3.Cert == "" || h3.Key == "" {
		return fmt.Errorf("server %s must be config with tls", http2name)
	}

	if h3.IsTCP {
		err := http3.ListenAndServeTLS(h3.Port, h3.Cert, h3.Key, h3.engine)
		if err != nil && err != context.Canceled {
			return errors.Wrap(err, "start http3 with tcp server err")
		}
	} else {
		h3.svr = &http3.Server{
			Addr:           h3.Port,
			Handler:        h3.engine,
			MaxHeaderBytes: 1 << uint(h3.MaxHeaderBytes),
			QUICConfig: &quic.Config{
				Tracer: qlog.DefaultConnectionTracer,
			},
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
				}},
		}
		err := h3.svr.ListenAndServeTLS(h3.Cert, h3.Key)
		if err != nil && err != context.Canceled {
			return errors.Wrap(err, "start http3 with udp server err")
		}
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

func (h3 *GinHttp3Server) loadValidation() error {
	trans := validation.GetTranslator()
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return fmt.Errorf("failed to get validator engine")
	}
	_ = zh.RegisterDefaultTranslations(validate, trans)
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	for _, valid := range h3.validator {
		if valid.Validator() != nil {
			validate.RegisterValidation(valid.Tag(), valid.Validator())
		}
		if valid.RegisterTranslation() != nil && valid.Translation() != nil {
			validate.RegisterTranslation(valid.Tag(), trans, valid.RegisterTranslation(), valid.Translation())
		}
	}
	return nil
}
