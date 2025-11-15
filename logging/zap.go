package logging

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
	"time"
)

func defaultZapLogger() *zap.Logger {
	var once sync.Once
	var logger *zap.Logger
	once.Do(func() {
		logFileName := fmt.Sprintf("logs/log-%s.log", time.Now().Format("2006-01-02"))

		if err := os.MkdirAll("logs", 0755); err != nil {
			panic(fmt.Sprintf("failed to create logs directory: %v", err))
		}

		logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Sprintf("failed to open log file: %v", err))
		}

		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "timestamp"
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder := zapcore.NewConsoleEncoder(encoderCfg)

		fileSyncer := zapcore.AddSync(logFile)
		stdoutSyncer := zapcore.AddSync(os.Stdout)

		core := zapcore.NewTee(
			zapcore.NewCore(encoder, stdoutSyncer, zap.InfoLevel),
			zapcore.NewCore(encoder, fileSyncer, zap.InfoLevel),
		)

		logger = zap.New(core, zap.AddCaller())
	})
	return logger
}

type zapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(logger *zap.Logger) Logger {
	if logger == nil {
		logger = defaultZapLogger()
	}

	return &zapLogger{logger: logger}
}

func (z zapLogger) Debug(ctx context.Context, msg string, fields ...Field) {
	z.logWithZap(ctx, LevelDebug, msg, fields...)
}

func (z zapLogger) Info(ctx context.Context, msg string, fields ...Field) {
	z.logWithZap(ctx, LevelInfo, msg, fields...)
}

func (z zapLogger) Warning(ctx context.Context, msg string, fields ...Field) {
	z.logWithZap(ctx, LevelWarn, msg, fields...)
}

func (z zapLogger) Error(ctx context.Context, msg string, fields ...Field) {
	z.logWithZap(ctx, LevelError, msg, fields...)
}

func (z zapLogger) logWithZap(ctx context.Context, level LogLevel, msg string, fields ...Field) {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}

	switch {
	case level <= LevelDebug:
		z.logger.Debug(msg, zapFields...)
	case level <= LevelInfo:
		z.logger.Info(msg, zapFields...)
	case level <= LevelWarn:
		z.logger.Warn(msg, zapFields...)
	default:
		z.logger.Error(msg, zapFields...)
	}
}
