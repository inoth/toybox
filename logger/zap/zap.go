package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/inoth/toybox/config"
	"github.com/inoth/toybox/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	name = "zap"
)

type ZapLogger struct {
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

	log *zap.Logger
}

func NewZapLogger(conf config.ConfigMate) logger.Logger {
	log := ZapLogger{
		hooks: make([](func(zapcore.Entry) error), 0),
	}
	if err := conf.PrimitiveDecode(&log); err != nil {
		panic(fmt.Errorf("init logger err: %v", err))
	}
	log.newLogger()
	return &log
}

func (z *ZapLogger) Name() string {
	return name
}

func (z *ZapLogger) Log(ctx context.Context, level int, msg string) {
	if level > 3 || level < 0 {
		level = 3
	}
	traceId := ctx.Value("trace_id")
	if traceId == nil {
		traceId = ""
	}
	switch level {
	case logger.LevelDebug:
		z.log.Debug(msg, zap.String("trace_id", traceId.(string)))
	case logger.LevelInfo:
		z.log.Info(msg, zap.String("trace_id", traceId.(string)))
	case logger.LevelWarn:
		z.log.Warn(msg, zap.String("trace_id", traceId.(string)))
	case logger.LevelError:
		z.log.Error(msg, zap.String("trace_id", traceId.(string)))
	}
}

func (z *ZapLogger) newLogger() {
	debug := os.Getenv("GORUNEVN")
	if debug == "dev" || debug == "debug" {
		z.log, _ = zap.NewDevelopment()
	}

	encoderConf := genEncoderConf()
	encoder := zapcore.NewConsoleEncoder(encoderConf)
	if z.Json {
		encoder = zapcore.NewJSONEncoder(encoderConf)
	}

	writers := make([]zapcore.Core, 0)
	if z.Err != "" {
		writers = append(writers, z.newZapcore(z.Err, encoder, func(lv zapcore.Level) bool {
			return lv == zap.ErrorLevel
		}))
	}
	if z.Warn != "" {
		writers = append(writers, z.newZapcore(z.Warn, encoder, func(lv zapcore.Level) bool {
			return lv == zap.WarnLevel
		}))
	}
	if z.Info != "" {
		writers = append(writers, z.newZapcore(z.Info, encoder, func(lv zapcore.Level) bool {
			return lv == zap.InfoLevel
		}))
	}
	if z.Debug != "" {
		writers = append(writers, z.newZapcore(z.Debug, encoder, func(lv zapcore.Level) bool {
			return lv == zap.DebugLevel
		}))
	}
	z.log = zap.New(zapcore.NewTee(writers...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel))

	_ = z.log.Sync()
}

func (z *ZapLogger) newZapcore(path string, encoder zapcore.Encoder, fn func(lv zapcore.Level) bool) zapcore.Core {
	return zapcore.NewCore(encoder, zapcore.AddSync(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    z.MaxSize,
		MaxAge:     z.MaxAge,
		MaxBackups: z.MaxBackup,
		Compress:   z.Compress,
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
