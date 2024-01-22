package toybox

import (
	"context"
	"fmt"

	"github.com/inoth/toybox/internal"
	"github.com/inoth/toybox/util"

	"golang.org/x/sync/errgroup"
)

type ToyBox struct {
	id     string
	ctx    context.Context
	cancel func()
	mate   ConfigMate

	cpts []Component
	svcs []Server
}

func New(opts ...Option) *ToyBox {
	tb := ToyBox{
		id:   util.UUID(),
		cpts: make([]Component, 0),
		svcs: make([]Server, 0),
	}
	for _, opt := range opts {
		opt(&tb)
	}
	tb.ctx, tb.cancel = context.WithCancel(context.Background())
	return &tb
}

func (tb ToyBox) ID() string {
	return tb.id
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
	if err := tb.mate.PrimitiveDecodeComponent(tb.cpts...); err != nil {
		fmt.Printf("load component config err: %v\n", err)
		return err
	}
	if err := tb.mate.PrimitiveDecodeServer(tb.svcs...); err != nil {
		fmt.Printf("load server config err: %v\n", err)
		return err
	}

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
