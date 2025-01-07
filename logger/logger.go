package logger

import (
	"context"
	"log"
)

var DefaultLogger Logger

func init() {
	DefaultLogger = &stdLogger{}
}

type Logger interface {
	Log(ctx context.Context, level int, msg string)
}

type stdLogger struct{}

func (l *stdLogger) Log(ctx context.Context, level int, msg string) {
	if level > 3 || level < 0 {
		level = 3
	}
	traceId := ctx.Value("trace_id")
	log.Printf("[%s][%v]%s\n", LevelName[level], traceId, msg)
}

func Log(ctx context.Context, level int, msg string) {
	DefaultLogger.Log(ctx, level, msg)
}
