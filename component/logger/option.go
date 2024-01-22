package logger

import (
	"github.com/inoth/toybox"

	"go.uber.org/zap/zapcore"
)

const (
	default_name = "logger"
)

type Option func(*ZapComponent)

func defaultOption() ZapComponent {
	return ZapComponent{
		name:      default_name,
		ready:     true,
		Debug:     "log/debug.log",
		Info:      "log/info.log",
		Warn:      "log/warn.log",
		Err:       "log/err.log",
		MaxSize:   10,
		MaxAge:    15,
		MaxBackup: 30,
		Compress:  true,
		Json:      true,
	}
}

func SetHooks(hooks ...(func(zapcore.Entry) error)) Option {
	return func(zc *ZapComponent) {
		hooks = append(hooks, hooks...)
	}
}

func WithLogger(opts ...Option) toybox.Option {
	o := defaultOption()
	for _, opt := range opts {
		opt(&o)
	}
	return func(tb *toybox.ToyBox) {
		tb.AppendComponent(&o)
	}
}
