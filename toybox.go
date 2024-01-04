package toybox

import (
	"context"
	"errors"
	"fmt"
	"github/inoth/toybox/internal"

	"golang.org/x/sync/errgroup"
)

type ToyBox struct {
	ctx          context.Context
	eg           *errgroup.Group
	env          string
	confDir      string
	confFileType string

	cpts []Component

	svcs []Server
}

func (tb *ToyBox) AppendComponent(cpts ...Component) *ToyBox {
	tb.cpts = append(tb.cpts, cpts...)
	return tb
}

func (tb *ToyBox) AppendServer(svcs ...Server) *ToyBox {
	tb.svcs = append(tb.svcs, svcs...)
	return tb
}

func New(opts ...Option) *ToyBox {
	tb := defaultOption()
	for _, opt := range opts {
		opt(&tb)
	}
	tb.eg, tb.ctx = errgroup.WithContext(context.Background())
	return &tb
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
		fmt.Printf("init component err: %v\n", err)
		return err
	}
	if err := tb.checkServers(); err != nil {
		fmt.Printf("check servers err: %v\n", err)
		return err
	}

	for _, svc := range tb.svcs {
		eg_svc := svc
		tb.eg.Go(func() error {
			return eg_svc.Run(tb.ctx)
		})
	}

	err := tb.eg.Wait()
	if err != nil {
		fmt.Printf("RunServers Error: %v\n", err)
		return err
	}
	internal.ListenSignal()
	return nil
}
