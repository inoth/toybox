package profile

import (
	"context"
	"net/http"
	_ "net/http/pprof"

	"github.com/pkg/errors"
)

const (
	name = "property"
)

type Profile struct {
	option

	svr *http.Server
}

func New(opts ...Option) *Profile {
	o := option{
		Port: ":9001",
	}
	for _, opt := range opts {
		opt(&o)
	}
	return &Profile{
		option: o,
	}
}

func (p *Profile) Name() string {
	return name
}

func (p *Profile) Start(ctx context.Context) error {
	p.svr = &http.Server{
		Addr: p.Port,
	}
	if err := p.svr.ListenAndServe(); err != nil && err != context.Canceled && err != http.ErrServerClosed {
		return errors.Wrap(err, "start pprof err")
	}
	return nil
}

func (p *Profile) Stop(ctx context.Context) error {
	return p.svr.Shutdown(ctx)
}
