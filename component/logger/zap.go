package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	Name        = "zap"
	DefaultPath = "/toybox/log/"
)

type Option func(opt *ZapComponent)

type ZapComponent struct {
	Debug     string `toml:"debug"`
	Info      string `toml:"info"`
	Warn      string `toml:"warn"`
	Err       string `toml:"err"`
	MaxSize   int    `toml:"max_size"`
	MaxAge    int    `toml:"max_age"`
	MaxBackup int    `toml:"max_backup"`
	Compress  bool   `toml:"compress"`
	Json      bool   `toml:"json"`

	hooks [](func(zapcore.Entry) error)
	zp    *zap.Logger
}

func SetHooks(hooks ...(func(zapcore.Entry) error)) Option {
	return func(opt *ZapComponent) {
		opt.hooks = append(opt.hooks, hooks...)
	}
}

func new(opts ...Option) *ZapComponent {
	o := ZapComponent{
		Debug:     DefaultPath + "debug/debug.log",
		Info:      DefaultPath + "info/info.log",
		Warn:      DefaultPath + "warn/warn.log",
		Err:       DefaultPath + "err/err.log",
		MaxSize:   10 << 10,
		MaxAge:    15,
		MaxBackup: 30,
		Compress:  true,
		Json:      true,
	}
	for _, opt := range opts {
		opt(&o)
	}
	o.zp, _ = o.newLogger()
	return &o
}

func (zc *ZapComponent) Name() string {
	return Name
}

func (zc *ZapComponent) newLogger() (*zap.Logger, error) {
	debug := os.Getenv("GORUNEVN")
	if debug == "dev" || debug == "debug" {
		return zap.NewDevelopment()
	}

	encoderConf := genEncoderConf()
	encoder := zapcore.NewConsoleEncoder(encoderConf)
	if zc.Json {
		encoder = zapcore.NewJSONEncoder(encoderConf)
	}

	writers := make([]zapcore.Core, 0)
	if zc.Err != "" {
		writers = append(writers, zc.newZapcore(zc.Err, encoder, func(lv zapcore.Level) bool {
			return lv == zap.ErrorLevel
		}))
	}
	if zc.Warn != "" {
		writers = append(writers, zc.newZapcore(zc.Warn, encoder, func(lv zapcore.Level) bool {
			return lv == zap.WarnLevel
		}))
	}
	if zc.Info != "" {
		writers = append(writers, zc.newZapcore(zc.Info, encoder, func(lv zapcore.Level) bool {
			return lv == zap.InfoLevel
		}))
	}
	if zc.Debug != "" {
		writers = append(writers, zc.newZapcore(zc.Debug, encoder, func(lv zapcore.Level) bool {
			return lv == zap.DebugLevel
		}))
	}
	logger := zap.New(zapcore.NewTee(writers...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel))
	defer logger.Sync()
	return logger, nil
}

func (zc *ZapComponent) newZapcore(path string, encoder zapcore.Encoder, fn func(lv zapcore.Level) bool) zapcore.Core {
	return zapcore.NewCore(encoder, zapcore.AddSync(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    zc.MaxSize,
		MaxAge:     zc.MaxAge,
		MaxBackups: zc.MaxBackup,
		Compress:   zc.Compress,
	}), zap.LevelEnablerFunc(fn))
}

func genEncoderConf() zapcore.EncoderConfig {
	encoderConf := zap.NewProductionEncoderConfig()
	encoderConf.TimeKey = "created_at"
	encoderConf.MessageKey = "message"
	encoderConf.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConf.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05:000"))
	}
	return encoderConf
}
