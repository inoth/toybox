package httpgin

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/singleflight"
)

type Http2GinServer struct {
	ready  bool
	sfg    singleflight.Group
	engine *gin.Engine
	option
}

func NewHttp2(opts ...Option) *Http2GinServer {
	o := defaultOption()
	for _, opt := range opts {
		opt(&o)
	}
	h2gs := Http2GinServer{
		ready:  false,
		sfg:    singleflight.Group{},
		engine: gin.New(),
		option: o,
	}
	return &h2gs
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
