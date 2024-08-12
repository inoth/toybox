package logger

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

var (
	once sync.Once
)

type Logger interface {
	Debug(ctx context.Context, msg string)
	Info(ctx context.Context, msg string)
	Warn(ctx context.Context, msg string)
	Error(ctx context.Context, msg string)
}

type logger struct {
	server_name string
	log         *ZapComponent
}

func GetLogger(server_name string) Logger {
	once.Do(func() {
		if log == nil {
			log = New()
		}
	})
	return &logger{
		server_name: server_name,
		log:         log,
	}
}

func (l *logger) Debug(ctx context.Context, msg string) {
	if l.log.zp == nil {
		return
	}
	trace_id := ""
	if ctx.Value("trace_id") != nil {
		trace_id = ctx.Value("trace_id").(string)
	}
	l.log.zp.Debug(msg, zap.String("server_name", l.server_name), zap.String("trace_id", trace_id))
}

func (l *logger) Info(ctx context.Context, msg string) {
	if l.log.zp == nil {
		return
	}
	trace_id := ""
	if ctx.Value("trace_id") != nil {
		trace_id = ctx.Value("trace_id").(string)
	}
	l.log.zp.Info(msg, zap.String("server_name", l.server_name), zap.String("trace_id", trace_id))
}

func (l *logger) Warn(ctx context.Context, msg string) {
	if l.log.zp == nil {
		return
	}
	trace_id := ""
	if ctx.Value("trace_id") != nil {
		trace_id = ctx.Value("trace_id").(string)
	}
	l.log.zp.Warn(msg, zap.String("server_name", l.server_name), zap.String("trace_id", trace_id))
}

func (l *logger) Error(ctx context.Context, msg string) {
	if l.log.zp == nil {
		return
	}
	trace_id := ""
	if ctx.Value("trace_id") != nil {
		trace_id = ctx.Value("trace_id").(string)
	}
	l.log.zp.Error(msg, zap.String("server_name", l.server_name), zap.String("trace_id", trace_id))
}
