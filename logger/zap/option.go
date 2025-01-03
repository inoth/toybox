package zap

import "go.uber.org/zap/zapcore"

type Option func(opt *ZapLogger)

func WithHooks(hooks ...(func(zapcore.Entry) error)) Option {
	return func(opt *ZapLogger) {
		opt.hooks = append(opt.hooks, hooks...)
	}
}
