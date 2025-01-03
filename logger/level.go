package logger

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

var (
	LevelName = []string{
		"DEBUG",
		"INFO",
		"WARN",
		"ERROR",
	}
)
