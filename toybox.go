package toybox

import (
	"context"
	"errors"
	"fmt"
	"github/inoth/toybox/internal"

	"golang.org/x/sync/errgroup"
)

type ToyBoxOpt func(cfg *ToyBoxConfig)

type ToyBoxConfig struct {
}

type ToyBox struct {
	ctx  context.Context
	eg   *errgroup.Group
	conf ToyBoxConfig

	// 组件
	cpts []Component
	// 服务
	svcs []Server
}

func New(opts ...ToyBoxOpt) *ToyBox {
	tb := &ToyBox{
		conf: ToyBoxConfig{},
		cpts: make([]Component, 0),
		svcs: make([]Server, 0),
	}
	tb.eg, tb.ctx = errgroup.WithContext(context.Background())

	for _, opt := range opts {
		opt(&tb.conf)
	}

	return tb
}

func (tb *ToyBox) AppendComponent(cpt Component) {
	tb.cpts = append(tb.cpts, cpt)
}

func (tb *ToyBox) AppendServer(svc Server) {
	tb.svcs = append(tb.svcs, svc)
}

func (tb *ToyBox) initComponents() error {
	for _, cpt := range tb.cpts {
		if cpt.Status() != ComponentStatusOK {
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
