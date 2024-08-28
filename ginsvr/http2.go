package ginsvr

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/inoth/toybox/validation"
	"github.com/pkg/errors"
	"golang.org/x/net/http2"
	"golang.org/x/sync/singleflight"
)

const (
	http2name = "gin2"
)

type GinHttp2Server struct {
	option
	sfg singleflight.Group
	svr *http.Server
}

func NewHttp2(opts ...Option) *GinHttp2Server {
	o := option{
		Port:    ":9050",
		engine:  gin.New(),
		handles: make([]Handler, 0),
	}
	for _, opt := range opts {
		opt(&o)
	}
	if o.serverName == "" {
		o.serverName = http2name
	}
	return &GinHttp2Server{
		option: o,
		sfg:    singleflight.Group{},
	}
}

func (h2 *GinHttp2Server) Name() string {
	return h2.serverName
}

func (h2 *GinHttp2Server) Start(ctx context.Context) error {

	h2.loadRouter()

	if err := h2.loadValidation(); err != nil {
		return errors.Wrap(err, "loadValidation() failed")
	}

	if h2.Cert == "" || h2.Key == "" {
		return fmt.Errorf("server %s must be config with tls", http2name)
	}

	h2.svr = &http.Server{
		Addr:           h2.Port,
		Handler:        h2.engine,
		ReadTimeout:    time.Second * time.Duration(h2.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(h2.WriteTimeout),
		MaxHeaderBytes: 1 << uint(h2.MaxHeaderBytes),
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

	err := http2.ConfigureServer(h2.svr, &http2.Server{})
	if err != nil {
		return errors.Wrap(err, "ConfigureServer err")
	}

	err = h2.svr.ListenAndServeTLS(h2.Cert, h2.Key)
	if err != nil && err != context.Canceled {
		return errors.Wrap(err, "start http2 server err")
	}
	return nil
}

func (h2 *GinHttp2Server) Stop(ctx context.Context) error {
	return h2.svr.Shutdown(ctx)
}

func (h2 *GinHttp2Server) Do(key string, fn func() (any, error)) (v any, err error, shared bool) {
	return h2.sfg.Do(key, fn)
}

func (h2 *GinHttp2Server) loadRouter() {
	for _, h := range h2.handles {
		for _, r := range h.Routers() {
			h2.engine.Handle(
				r.Method,
				h.Prefix()+"/"+r.Path,
				append(h.Middlewares(), r.Handle...)...,
			)
		}
	}
}

func (h2 *GinHttp2Server) loadValidation() error {
	validation.LoadValidation(h2.validator)
	// trans := validation.GetTranslator()
	// validate := validation.GetDefaultValidator()
	// _ = zh.RegisterDefaultTranslations(validate, trans)
	// validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
	// 	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	// 	if name == "-" {
	// 		return ""
	// 	}
	// 	return name
	// })
	// for _, valid := range h2.validator {
	// 	if valid.Validator() != nil {
	// 		validate.RegisterValidation(valid.Tag(), valid.Validator())
	// 	}
	// 	if valid.RegisterTranslation() != nil && valid.Translation() != nil {
	// 		validate.RegisterTranslation(valid.Tag(), trans, valid.RegisterTranslation(), valid.Translation())
	// 	}
	// }
	return nil
}
