package logging

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

type slogLogger struct {
	logger *slog.Logger
	files  []*os.File // Keep track of files to close them properly
}

func (s *slogLogger) Close() error {
	for _, f := range s.files {
		if err := f.Sync(); err != nil {
			return err
		}
		if err := f.Close(); err != nil {
			return err
		}
	}
	return nil
}

// MultiHandler wraps multiple slog.Handler to write to all of them
type MultiHandler struct {
	handlers []slog.Handler
}

func NewMultiHandler(handlers ...slog.Handler) *MultiHandler {
	return &MultiHandler{handlers: handlers}
}

func (m *MultiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Enable if any handler is enabled
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *MultiHandler) Handle(ctx context.Context, record slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, record.Level) {
			if err := h.Handle(ctx, record.Clone()); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MultiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithAttrs(attrs)
	}
	return NewMultiHandler(handlers...)
}

func (m *MultiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		handlers[i] = h.WithGroup(name)
	}
	return NewMultiHandler(handlers...)
}

func defaultSlogLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))
}

// NewSlogLoggerWithMultiOutput creates a logger that writes to multiple outputs
func NewSlogLoggerWithMultiOutput(outputs []io.Writer, handlers []slog.Handler, opts *slog.HandlerOptions) Logger {
	if opts == nil {
		opts = &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}
	}

	var allHandlers []slog.Handler

	// Add handlers from outputs
	for _, output := range outputs {
		allHandlers = append(allHandlers, slog.NewJSONHandler(output, opts))
	}

	// Add custom handlers
	allHandlers = append(allHandlers, handlers...)

	multiHandler := NewMultiHandler(allHandlers...)
	logger := slog.New(multiHandler)

	return &slogLogger{logger: logger}
}

// NewSlogLoggerWithFile creates a logger that writes to both stdout and a file
func NewSlogLoggerWithFile(logFilePath string, opts *slog.HandlerOptions) (Logger, error) {
	if opts == nil {
		opts = &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}
	}

	// Create log file
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create handlers
	stdoutHandler := slog.NewTextHandler(os.Stdout, opts)
	fileHandler := slog.NewJSONHandler(logFile, opts)

	multiHandler := NewMultiHandler(stdoutHandler, fileHandler)
	logger := slog.New(multiHandler)

	return &slogLogger{
		logger: logger,
		files:  []*os.File{logFile},
	}, nil
}

// NewSlogLoggerWithRotation creates a logger with multiple outputs including daily rotation
func NewSlogLoggerWithRotation(logDir string, opts *slog.HandlerOptions) (Logger, error) {
	if opts == nil {
		opts = &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}
	}

	// Create logs directory
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log file with date
	logFileName := fmt.Sprintf("%s/log-%s.json", logDir, time.Now().Format("2006-01-02"))
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create handlers: stdout (text), stderr (errors only), file (json)
	stdoutHandler := slog.NewTextHandler(os.Stdout, opts)
	stderrHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelError, // Only errors to stderr
	})
	fileHandler := slog.NewJSONHandler(logFile, opts)

	multiHandler := NewMultiHandler(stdoutHandler, stderrHandler, fileHandler)
	logger := slog.New(multiHandler)

	return &slogLogger{
		logger: logger,
		files:  []*os.File{logFile},
	}, nil
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
