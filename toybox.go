package toybox

import (
	"context"
	"errors"
	"fmt"
	"github/inoth/toybox/internal"

	"golang.org/x/sync/errgroup"
)

type ToyBox struct {
	ctx context.Context
	eg  *errgroup.Group

	Option
}

func (tb *ToyBox) AppendComponent(cpts ...Component) {
	tb.cpts = append(tb.cpts, cpts...)
}

func (tb *ToyBox) AppendServer(svcs ...Server) {
	tb.svcs = append(tb.svcs, svcs...)
}

func New(opts ...OptionFunc) *ToyBox {
	o := defaultOption()
	for _, opt := range opts {
		opt(&o)
	}
	tb := &ToyBox{
		Option: o,
	}
	tb.eg, tb.ctx = errgroup.WithContext(context.Background())
	return tb
}

func (tb *ToyBox) initComponents() error {
	for _, cpt := range tb.cpts {
		if !cpt.Ready() {
			return errors.New("components that are not yet ready")
		}
		if err := cpt.Init(tb.ctx); err != nil {
			return fmt.Errorf("component init error: %w", err)
		}
	}
	return nil
}

func (tb *ToyBox) checkServers() error {
	for _, svc := range tb.svcs {
		if !svc.Ready() {
			return errors.New("servers that are not yet ready")
		}
	}
	return nil
}

func (tb *ToyBox) Run() error {
	if err := tb.initComponents(); err != nil {
		return err
	}
	if err := tb.checkServers(); err != nil {
		return err
	}

	internal.ListenSignal()
	return nil
}
