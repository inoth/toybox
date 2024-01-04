package toybox

import (
	"context"
	"fmt"
	"github/inoth/toybox/internal"

	"golang.org/x/sync/errgroup"
)

type ToyBox struct {
	ctx          context.Context
	cancel       func()
	env          string
	confDir      string
	confFileType string

	cpts []Component

	svcs []Server
}

func New(opts ...Option) *ToyBox {
	tb := ToyBox{
		confDir:      "config",
		env:          "",
		confFileType: "",
		cpts:         make([]Component, 0),
		svcs:         make([]Server, 0),
	}
	for _, opt := range opts {
		opt(&tb)
	}
	tb.ctx, tb.cancel = context.WithCancel(context.Background())
	return &tb
}

func (tb *ToyBox) initComponents() error {
	for _, cpt := range tb.cpts {
		if !cpt.Ready() {
			return fmt.Errorf("components that are not yet ready")
		}
		if err := cpt.Init(tb.ctx); err != nil {
			return fmt.Errorf("component init error: %w", err)
		}
	}
	return nil
}

func (tb *ToyBox) AppendComponent(cpts ...Component) { tb.cpts = append(tb.cpts, cpts...) }

func (tb *ToyBox) AppendServer(svcs ...Server) { tb.svcs = append(tb.svcs, svcs...) }

func (tb *ToyBox) checkServers() error {
	for _, svc := range tb.svcs {
		if !svc.Ready() {
			return fmt.Errorf("servers that are not yet ready")
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

	eg, _ := errgroup.WithContext(tb.ctx)

	for _, svc := range tb.svcs {
		svc := svc
		eg.Go(func() error {
			return svc.Run(tb.ctx)
		})
	}

	if err := eg.Wait(); err != nil {
		tb.cancel()
		fmt.Printf("run servers err: %v\n", err)
		return err
	}
	internal.ListenSignal()
	return nil
}
