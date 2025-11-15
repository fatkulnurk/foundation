package logging

import (
	"context"
	"log/slog"
	"os"
)

type slogLogger struct {
	logger *slog.Logger
}

func defaultSlogLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
}

func NewSlogLogger(logger *slog.Logger) Logger {
	if logger == nil {
		logger = defaultSlogLogger()
	}

	return &slogLogger{logger: logger}
}

func (s slogLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	s.logWithSlog(ctx, LevelDebug, msg, fields...)
}

func (s slogLogger) Info(ctx context.Context, msg string, fields ...Field) {
	s.logWithSlog(ctx, LevelInfo, msg, fields...)
}

func (s slogLogger) Warning(ctx context.Context, msg string, fields ...Field) {
	s.logWithSlog(ctx, LevelWarn, msg, fields...)
}

func (s slogLogger) Error(ctx context.Context, msg string, fields ...Field) {
	s.logWithSlog(ctx, LevelError, msg, fields...)
}

func (s slogLogger) logWithSlog(ctx context.Context, level LogLevel, msg string, fields ...Field) {
	slogLevel := func(level LogLevel) slog.Level {
		switch level {
		case LevelDebug:
			return slog.LevelDebug
		case LevelInfo:
			return slog.LevelInfo
		case LevelWarn:
			return slog.LevelWarn
		case LevelError:
			return slog.LevelError
		default:
			return slog.LevelInfo
		}
	}(level)

	if !s.logger.Enabled(ctx, slogLevel) {
		return
	}

	var attrs []slog.Attr
	for _, field := range fields {
		attrs = append(attrs, slog.Any(field.Key, field.Value))
	}

	s.logger.LogAttrs(ctx, slogLevel, msg, attrs...)
}
