package toybox

import (
	"context"
	"os"

	"github.com/inoth/toybox/config"
	"github.com/inoth/toybox/server"
)

type Option func(opt *option)

type option struct {
	id      string
	name    string
	version string
	ctx     context.Context
	sigs    []os.Signal
	svcs    []server.Server
	cfg     config.ConfigMate
}

func WithName(name string) Option {
	return func(opt *option) {
		opt.name = name
	}
}

func WithVersion(version string) Option {
	return func(opt *option) {
		opt.version = version
	}
}

func WithContext(ctx context.Context) Option {
	return func(opt *option) {
		opt.ctx = ctx
	}
}

func WithServer(svcs ...server.Server) Option {
	return func(opt *option) {
		opt.svcs = append(opt.svcs, svcs...)
	}
}

func WithConfig(cfg config.ConfigMate) Option {
	return func(opt *option) {
		opt.cfg = cfg
	}
}
