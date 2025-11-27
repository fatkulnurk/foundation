package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/fatkulnurk/foundation/logging"
)

func main() {
	// Create logs directory
	if err := os.MkdirAll("logs", 0755); err != nil {
		panic(err)
	}

	fmt.Println("=== Example 1: Simple multi-output (stdout + file) ===")
	example1()

	fmt.Println("\n=== Example 2: Multi-output with rotation ===")
	example2()

	fmt.Println("\n=== Example 3: Separate files by level ===")
	example3()
}

func example1() {
	// Creates logger that writes to:
	// - stdout (text format)
	// - logs/app.log (JSON format)
	logger, err := logging.NewSlogLoggerWithFile("logs/app.log", nil)
	if err != nil {
		panic(err)
	}

	logging.InitLogging(logger)
	logging.Info(context.Background(), "Application started with file logging")
	logging.Error(context.Background(), "This is an error", logging.NewField("code", 500))

	// Close logger to flush buffers (if it has Close method)
	if closer, ok := logger.(interface{ Close() error }); ok {
		closer.Close()
	}

	fmt.Println("✓ Logged to stdout and logs/app.log")
}

func example2() {
	// Creates logger that writes to:
	// - stdout (text format, all levels)
	// - stderr (text format, errors only)
	// - logs/log-YYYY-MM-DD.json (JSON format, daily rotation)
	logger, err := logging.NewSlogLoggerWithRotation("logs", nil)
	if err != nil {
		panic(err)
	}

	logging.InitLogging(logger)
	logging.Info(context.Background(), "Using rotation logger")
	logging.Error(context.Background(), "Error goes to stderr and file", logging.NewField("user_id", 123))

	// Close logger to flush buffers
	if closer, ok := logger.(interface{ Close() error }); ok {
		closer.Close()
	}

	fmt.Println("✓ Logged to stdout, stderr (errors only), and logs/log-YYYY-MM-DD.json")
}

func example3() {
	// Create separate log files for different levels
	errorFile, err := os.OpenFile("logs/error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	infoFile, err := os.OpenFile("logs/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// Create multi-output logger with separate files
	logger := logging.NewSlogLoggerWithMultiOutput(
		nil,
		[]slog.Handler{
			// Console output
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: false,
			}),
			// Error file (errors only)
			slog.NewJSONHandler(errorFile, &slog.HandlerOptions{
				Level:     slog.LevelError,
				AddSource: true,
			}),
			// Info file (info and above)
			slog.NewJSONHandler(infoFile, &slog.HandlerOptions{
				Level:     slog.LevelInfo,
				AddSource: true,
			}),
		},
		nil,
	)

	logging.InitLogging(logger)
	logging.Info(context.Background(), "Info message", logging.NewField("status", "ok"))
	logging.Error(context.Background(), "Error message", logging.NewField("status", "failed"))

	// Give slog time to write (slog writes are synchronous but let's be safe)
	time.Sleep(10 * time.Millisecond)

	// Flush and close files
	errorFile.Sync()
	errorFile.Close()
	infoFile.Sync()
	infoFile.Close()

	fmt.Println("Info logged to stdout and logs/info.log")
	fmt.Println("Error logged to stdout, logs/info.log, and logs/error.log")
}
