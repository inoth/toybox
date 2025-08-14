package logger

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/inoth/toybox/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapRemoteLogger struct {
	writer io.Writer
	log    *zap.Logger
}

func NewZapRemoteLogger(opts ...Option) logger.Logger {
	log := ZapRemoteLogger{}
	for _, opt := range opts {
		opt(&log)
	}
	if log.writer == nil {
		log.writer = os.Stdout
	}
	log.newLogger()
	return &log
}

func (z *ZapRemoteLogger) newLogger() {
	debug := os.Getenv("GORUNEVN")
	if debug == "dev" || debug == "debug" {
		z.log, _ = zap.NewDevelopment()
	}

	encoderConf := genEncoderConf()
	encoder := zapcore.NewConsoleEncoder(encoderConf)
	// if z.Json {
	// 	encoder = zapcore.NewJSONEncoder(encoderConf)
	// }

	z.log = zap.New(zapcore.NewTee(z.newZapcore(encoder, func(lv zapcore.Level) bool {
		if debug == "dev" || debug == "debug" {
			return lv >= zap.DebugLevel
		}
		return lv >= zap.InfoLevel
	})),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zap.ErrorLevel))

	_ = z.log.Sync()
}

func (z *ZapRemoteLogger) newZapcore(encoder zapcore.Encoder, fn func(lv zapcore.Level) bool) zapcore.Core {
	return zapcore.NewCore(encoder, zapcore.AddSync(z.writer), zap.LevelEnablerFunc(fn))
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

func (z *ZapRemoteLogger) LogDebug(ctx context.Context, msg string, args ...any) {
	traceId := ctx.Value("trace_id")
	if traceId == nil {
		traceId = ""
	}
	fields := make([]zap.Field, 0, len(args)+1)
	fields = append(fields, zap.String("trace_id", traceId.(string)))
	for _, v := range args {
		fields = append(fields, zap.Any("arg", v))
	}
	z.log.Debug(msg, fields...)
}

func (z *ZapRemoteLogger) LogInfo(ctx context.Context, msg string, args ...any) {
	traceId := ctx.Value("trace_id")
	if traceId == nil {
		traceId = ""
	}
	fields := make([]zap.Field, 0, len(args)+1)
	fields = append(fields, zap.String("trace_id", traceId.(string)))
	for _, v := range args {
		fields = append(fields, zap.Any("arg", v))
	}
	z.log.Info(msg, fields...)
}

func (z *ZapRemoteLogger) LogWarn(ctx context.Context, msg string, args ...any) {
	traceId := ctx.Value("trace_id")
	if traceId == nil {
		traceId = ""
	}
	fields := make([]zap.Field, 0, len(args)+1)
	fields = append(fields, zap.String("trace_id", traceId.(string)))
	for _, v := range args {
		fields = append(fields, zap.Any("arg", v))
	}
	z.log.Warn(msg, fields...)
}

func (z *ZapRemoteLogger) LogError(ctx context.Context, msg string, args ...any) {
	traceId := ctx.Value("trace_id")
	if traceId == nil {
		traceId = ""
	}
	fields := make([]zap.Field, 0, len(args)+1)
	fields = append(fields, zap.String("trace_id", traceId.(string)))
	for _, v := range args {
		fields = append(fields, zap.Any("arg", v))
	}
	z.log.Error(msg, fields...)
}
