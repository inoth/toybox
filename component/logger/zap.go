package logger

import (
	"context"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	log  *zap.Logger
	once sync.Once
)

type ZapComponent struct {
	ready bool
	name  string
	hooks [](func(zapcore.Entry) error)

	Debug     string `toml:"debug"`
	Info      string `toml:"info"`
	Warn      string `toml:"warn"`
	Err       string `toml:"err"`
	MaxSize   int    `toml:"max_size"`
	MaxAge    int    `toml:"max_age"`
	MaxBackup int    `toml:"max_backup"`
	Compress  bool   `toml:"compress"`
	Json      bool   `toml:"json"`
}

func (zc ZapComponent) Name() string {
	return zc.name
}

func (zc ZapComponent) Ready() bool {
	return zc.ready
}

func (zc *ZapComponent) IsReady() {
	zc.ready = true
}

func (zc *ZapComponent) Init(ctx context.Context) (err error) {
	once.Do(func() {
		log, err = zc.newLogger()
	})
	return
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
