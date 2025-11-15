package logging

import (
	"context"
	"sync"
)

type Logger interface {
	Debug(ctx context.Context, msg string, fields ...Field)
	Info(ctx context.Context, msg string, fields ...Field)
	Warning(ctx context.Context, msg string, fields ...Field)
	Error(ctx context.Context, msg string, fields ...Field)
}

// LogLevel represents the log level
type LogLevel int

const (
	LevelDebug LogLevel = 1
	LevelInfo  LogLevel = 2
	LevelWarn  LogLevel = 3
	LevelError LogLevel = 4
)

// Field represents a key-value pair for structured logging
type Field struct {
	Key   string
	Value any
}

func NewField(key string, value any) Field {
	return Field{Key: key, Value: value}
}

var (
	once sync.Once
	l    Logger
)

func InitLogging(logger Logger) {
	once.Do(func() {
		l = logger
	})
}

// Debug logs a message at the debug level
func Debug(ctx context.Context, msg string, fields ...Field) {
	l.Debug(ctx, msg, fields...)
}

// Info logs a message at the info level
func Info(ctx context.Context, msg string, fields ...Field) {
	l.Info(ctx, msg, fields...)
}

// Warning logs a message at the warning level
func Warning(ctx context.Context, msg string, fields ...Field) {
	l.Warning(ctx, msg, fields...)
}

// Error logs a message at the error level
func Error(ctx context.Context, msg string, fields ...Field) {
	l.Error(ctx, msg, fields...)
}
