package logger

import (
	"context"
	"fmt"
	"log"
)

var defaultLogger Logger

func init() {
	defaultLogger = &stdLogger{}
}

func LogDebug(ctx context.Context, msg string, args ...any) {
	defaultLogger.LogDebug(ctx, fmt.Sprintf(msg, args...))
}
func LogInfo(ctx context.Context, msg string, args ...any) {
	defaultLogger.LogInfo(ctx, fmt.Sprintf(msg, args...))
}
func LogWarn(ctx context.Context, msg string, args ...any) {
	defaultLogger.LogWarn(ctx, fmt.Sprintf(msg, args...))
}
func LogError(ctx context.Context, msg string, args ...any) {
	defaultLogger.LogError(ctx, fmt.Sprintf(msg, args...))
}

type Logger interface {
	LogDebug(ctx context.Context, msg string, args ...any)
	LogInfo(ctx context.Context, msg string, args ...any)
	LogWarn(ctx context.Context, msg string, args ...any)
	LogError(ctx context.Context, msg string, args ...any)
}

type stdLogger struct{}

func (l *stdLogger) doLog(ctx context.Context, level int, msg string) {
	if level > LevelError || level < LevelDebug {
		level = LevelError
	}
	traceId := ctx.Value("trace_id")
	log.Printf("[%s][%v] %s\n", LevelName[level], traceId, msg)
}

func (l *stdLogger) LogDebug(ctx context.Context, msg string, args ...any) {
	l.doLog(ctx, LevelDebug, fmt.Sprintf(msg, args...))
}
func (l *stdLogger) LogInfo(ctx context.Context, msg string, args ...any) {
	l.doLog(ctx, LevelInfo, fmt.Sprintf(msg, args...))
}
func (l *stdLogger) LogWarn(ctx context.Context, msg string, args ...any) {
	l.doLog(ctx, LevelWarn, fmt.Sprintf(msg, args...))
}
func (l *stdLogger) LogError(ctx context.Context, msg string, args ...any) {
	l.doLog(ctx, LevelError, fmt.Sprintf(msg, args...))
}
