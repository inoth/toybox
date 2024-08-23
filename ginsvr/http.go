package ginsvr

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"

	"github.com/inoth/toybox/validation"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
)

const (
	name = "gin"
)

type GinHttpServer struct {
	option
	sfg singleflight.Group
	svr *http.Server
}

func New(opts ...Option) *GinHttpServer {
	o := option{
		Port:    ":9050",
		engine:  gin.New(),
		handles: make([]Handler, 0),
	}
	for _, opt := range opts {
		opt(&o)
	}
	return &GinHttpServer{
		option: o,
		sfg:    singleflight.Group{},
	}
}

func (h *GinHttpServer) Name() string {
	return name
}

func (e *GinHttpServer) Start(ctx context.Context) error {

	e.loadRouter()

	if err := e.loadValidation(); err != nil {
		return errors.Wrap(err, "loadValidation() failed")
	}

	e.svr = &http.Server{
		Addr:           e.Port,
		Handler:        e.engine,
		ReadTimeout:    time.Second * time.Duration(e.ReadTimeout),
		WriteTimeout:   time.Second * time.Duration(e.WriteTimeout),
		MaxHeaderBytes: 1 << uint(e.MaxHeaderBytes),
	}
	var err error
	if e.TLS {
		e.svr.TLSConfig = &tls.Config{
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
		err = e.svr.ListenAndServeTLS(e.Cert, e.Key)
	} else {
		err = e.svr.ListenAndServe()
	}
	if err != nil && err != context.Canceled {
		return errors.Wrap(err, "start http server err")
	}
	return nil
}

func (e *GinHttpServer) Stop(ctx context.Context) error {
	return e.svr.Shutdown(ctx)
}

func (e *GinHttpServer) Do(key string, fn func() (any, error)) (v any, err error, shared bool) {
	return e.sfg.Do(key, fn)
}

func (e *GinHttpServer) loadRouter() {
	for _, h := range e.handles {
		for _, r := range h.Routers() {
			e.engine.Handle(
				r.Method,
				h.Prefix()+"/"+r.Path,
				append(h.Middlewares(), r.Handle...)...,
			)
		}
	}
}

func (e *GinHttpServer) loadValidation() error {
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
	for _, valid := range e.validator {
		if valid.Validator() != nil {
			validate.RegisterValidation(valid.Tag(), valid.Validator())
		}
		if valid.RegisterTranslation() != nil && valid.Translation() != nil {
			validate.RegisterTranslation(valid.Tag(), trans, valid.RegisterTranslation(), valid.Translation())
		}
	}
	return nil
}
