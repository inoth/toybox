package toybox

import (
	"context"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/inoth/toybox/registry"
)

type ToyBoxInfo interface {
	ID() string
	Name() string
	Version() string
	Metadata() map[string]string
	Endpoint() []string
}

type ToyBox struct {
	opt    options
	ctx    context.Context
	cancel context.CancelFunc

	m        sync.Mutex
	instance *registry.ServiceInstance
}

func New(opts ...Option) *ToyBox {
	o := options{
		ctx:              context.Background(),
		sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout: 10 * time.Second,
		stopTimeout:      10 * time.Second,
	}
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(&o)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &ToyBox{
		ctx:    ctx,
		cancel: cancel,
		opt:    o,
	}
}

func (t *ToyBox) ID() string { return t.opt.id }

func (t *ToyBox) Name() string { return t.opt.name }

func (t *ToyBox) Version() string { return t.opt.version }

func (t *ToyBox) Metadata() map[string]string { return t.opt.metadata }

func (t *ToyBox) Endpoint() []string {
	if t.instance != nil {
		return t.instance.Endpoints
	}
	return nil
}

func (t *ToyBox) Run() error {
	return nil
}

func (t *ToyBox) Stop() error {
	return nil
}
