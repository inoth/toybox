package toybox

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/inoth/toybox/registry"
)

type Option func(opt *options)

type options struct {
	id        string
	name      string
	version   string
	metadata  map[string]string
	endpoints []*url.URL

	ctx  context.Context
	sigs []os.Signal

	registrar        registry.Registrar
	registrarTimeout time.Duration
	stopTimeout      time.Duration
}

func ID(id string) Option {
	return func(opt *options) { opt.id = id }
}

func Name(name string) Option {
	return func(opt *options) { opt.name = name }
}

func Version(ver string) Option {
	return func(opt *options) { opt.version = ver }
}

func Metadata(md map[string]string) Option {
	return func(o *options) { o.metadata = md }
}

func Endpoint(endpoints ...*url.URL) Option {
	return func(o *options) { o.endpoints = endpoints }
}

func Context(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

func Signal(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

// func Logger(logger log.Logger) Option {
// 	return func(o *options) { o.logger = logger }
// }

// func Server(srv ...transport.Server) Option {
// 	return func(o *options) { o.servers = srv }
// }

// func Registrar(r registry.Registrar) Option {
// 	return func(o *options) { o.registrar = r }
// }

// func RegistrarTimeout(t time.Duration) Option {
// 	return func(o *options) { o.registrarTimeout = t }
// }
