package property

import (
	"context"
	"net/http"
	_ "net/http/pprof"

	"github.com/pkg/errors"
)

const (
	name = "property"
)

type Property struct {
	option

	svr *http.Server
}

func New(opts ...Option) *Property {
	o := option{
		Port: ":9052",
	}
	for _, opt := range opts {
		opt(&o)
	}
	return &Property{
		option: o,
	}
}

func (p *Property) Name() string {
	return name
}

func (p *Property) Start(ctx context.Context) error {
	p.svr = &http.Server{
		Addr: p.Port,
	}
	if err := p.svr.ListenAndServe(); err != nil {
		return errors.Wrap(err, "start pprof err")
	}
	return nil
}

func (p *Property) Stop(ctx context.Context) error {
	return p.svr.Shutdown(ctx)
}
