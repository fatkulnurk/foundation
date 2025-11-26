package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/fatkulnurk/foundation/logging"
)

func main() {
	ctx := context.Background()

	// Example 1: Using default Slog Logger (recommended)
	logger := logging.NewSlogLogger(nil) // nil will use default config
	logging.InitLogging(logger)

	logging.Info(ctx, "Application started")
	logging.Debug(ctx, "This is a debug message")
	logging.Warning(ctx, "This is a warning")
	logging.Error(ctx, "This is an error message")

	// Example 2: Logging with structured fields
	logging.Info(ctx, "User logged in",
		logging.NewField("user_id", 123),
		logging.NewField("email", "user@example.com"),
		logging.NewField("ip", "192.168.1.1"),
	)

	// Example 3: Using custom Slog logger with JSON format
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	customLogger := logging.NewSlogLogger(slog.New(jsonHandler))
	logging.InitLogging(customLogger)

	logging.Info(ctx, "Switched to JSON format logger")

	// Example 4: Error logging with context
	err := processOrder(ctx, 12345)
	if err != nil {
		logging.Error(ctx, "Failed to process order",
			logging.NewField("order_id", 12345),
			logging.NewField("error", err.Error()),
		)
	}

	// Example 5: Different log levels
	demonstrateLogLevels(ctx)
}

func processOrder(ctx context.Context, orderID int) error {
	logging.Info(ctx, "Processing order",
		logging.NewField("order_id", orderID),
	)

	// Simulate some processing
	// return errors.New("payment failed")

	logging.Info(ctx, "Order processed successfully",
		logging.NewField("order_id", orderID),
		logging.NewField("status", "completed"),
	)

	return nil
}

func demonstrateLogLevels(ctx context.Context) {
	// Simple messages
	logging.Debug(ctx, "Debug: Detailed information for debugging")
	logging.Info(ctx, "Info: General informational messages")
	logging.Warning(ctx, "Warning: Warning messages for potentially harmful situations")
	logging.Error(ctx, "Error: Error messages for error events")

	// With structured fields
	logging.Debug(ctx, "Debug with fields",
		logging.NewField("component", "database"),
		logging.NewField("query", "SELECT * FROM users"),
		logging.NewField("duration", "15ms"),
	)

	logging.Info(ctx, "Info with fields",
		logging.NewField("event", "user_registration"),
		logging.NewField("user_id", 456),
		logging.NewField("username", "john_doe"),
	)

	logging.Warning(ctx, "Warning with fields",
		logging.NewField("issue", "high_memory_usage"),
		logging.NewField("memory_used", "85%"),
		logging.NewField("threshold", "80%"),
	)

	logging.Error(ctx, "Error with fields",
		logging.NewField("error_code", "DB_CONNECTION_FAILED"),
		logging.NewField("host", "localhost:5432"),
		logging.NewField("retry", 3),
	)
}
