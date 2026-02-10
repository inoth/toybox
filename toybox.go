package toybox

import (
	"context"
	"log"
	"os"
	"syscall"

	"github.com/inoth/toybox/util"
)

type ToyBox struct {
	option

	id     string
	ctx    context.Context
	cancel context.CancelFunc
}

func New(opts ...Option) *ToyBox {
	o := option{
		sigs: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
	for _, opt := range opts {
		opt(&o)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &ToyBox{
		option: o,
		id:     util.UUID(),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (t *ToyBox) Run() error {
	log.Printf("Starting server ID:%s (PID: %d)\n", t.id, os.Getpid())

	return nil
}
