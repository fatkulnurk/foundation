# Logging - Unified Logging Package for Go

Module for structured logging with support for multiple backends (slog, zap) and multiple output destinations.

## Table of Contents

- [What is Logging?](#what-is-logging)
- [Module Contents](#module-contents)
- [How to Use](#how-to-use)
- [Multiple Output Support](#multiple-output-support)
- [Real-World Example](#real-world-example)
- [Best Practices](#best-practices)
- [Common Patterns](#common-patterns)
- [Installation](#installation)
- [Dependencies](#dependencies)
- [See Also](#see-also)

---

## What is Logging?

Logging is a unified logging package that provides a common interface for different logging backends. It supports both Go's standard `log/slog` and Uber's `zap` logger, with the ability to write to multiple outputs simultaneously.

**Think of it like:**
- A universal adapter for different logging libraries
- A way to log to multiple destinations (stdout, stderr, files, etc.) at once
- A structured logging solution with consistent API

**Use cases:**
- Application logging with multiple outputs
- Logging to both console and files
- Separate error logs from info logs
- Daily log rotation
- JSON and text format logging simultaneously

## Module Contents

### 1. **logging.go** - Core Interface
Defines the main interface and global functions:
- `Logger` interface - Common interface for all logging backends
- `InitLogging` - Initialize global logger
- `Debug`, `Info`, `Warning`, `Error` - Global logging functions
- `Field` - Structured logging fields

### 2. **slog.go** - Slog Implementation
Go standard library slog implementation:
- `NewSlogLogger` - Create basic slog logger
- `NewSlogLoggerWithFile` - Logger with stdout + file output
- `NewSlogLoggerWithRotation` - Logger with daily rotation
- `NewSlogLoggerWithMultiOutput` - Custom multi-output logger
- `MultiHandler` - Handler that writes to multiple destinations

### 3. **zap.go** - Zap Implementation
Uber zap implementation:
- `NewZapLogger` - Create zap logger
- Default configuration with stdout + file output

## How to Use

### Basic Usage with Slog

```go
import "github.com/fatkulnurk/foundation/logging"

func main() {
    // Create logger
    logger := logging.NewSlogLogger(nil) // nil = use default
    
    // Initialize global logger
    logging.InitLogging(logger)
    
    // Use global logging functions
    logging.Info(context.Background(), "Application started")
    logging.Error(context.Background(), "An error occurred", 
        logging.NewField("code", 500),
        logging.NewField("user_id", 123),
    )
}
```

### Basic Usage with Zap

```go
logger := logging.NewZapLogger(nil) // nil = use default
logging.InitLogging(logger)

logging.Info(context.Background(), "Using zap logger")
```

## Multiple Output Support

### Output to Stdout + File

```go
// Creates logger that writes to:
// - stdout (text format)
// - file (JSON format)
logger, err := logging.NewSlogLoggerWithFile("logs/app.log", nil)
if err != nil {
    panic(err)
}

// Close logger when application exits to flush buffers
defer func() {
    if closer, ok := logger.(interface{ Close() error }); ok {
        closer.Close()
    }
}()

logging.InitLogging(logger)
logging.Info(context.Background(), "This goes to both stdout and file")
```

### Output with Daily Rotation

```go
// Creates logger that writes to:
// - stdout (text format, all levels)
// - stderr (text format, errors only)
// - file (JSON format, daily rotation: logs/log-YYYY-MM-DD.json)
logger, err := logging.NewSlogLoggerWithRotation("logs", nil)
if err != nil {
    panic(err)
}

// Close logger when application exits to flush buffers
defer func() {
    if closer, ok := logger.(interface{ Close() error }); ok {
        closer.Close()
    }
}()

logging.InitLogging(logger)
logging.Info(context.Background(), "Info to stdout and file")
logging.Error(context.Background(), "Error to stdout, stderr, and file")
```

### Custom Multi-Output

```go
// Create custom multi-output logger
logger := logging.NewSlogLoggerWithMultiOutput(
    []io.Writer{
        os.Stdout, // JSON format
    },
    []slog.Handler{
        slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
            Level: slog.LevelError, // Only errors to stderr
        }),
    },
    &slog.HandlerOptions{
        AddSource: true,
        Level:     slog.LevelDebug,
    },
)

logging.InitLogging(logger)
logging.Debug(context.Background(), "Debug message")
logging.Error(context.Background(), "Error to both stdout and stderr")
```

### Separate Files by Level

```go
// Create separate log files for different levels
errorFile, _ := os.OpenFile("logs/error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
infoFile, _ := os.OpenFile("logs/info.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

logger := logging.NewSlogLoggerWithMultiOutput(
    nil,
    []slog.Handler{
        slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
            Level: slog.LevelInfo,
        }),
        slog.NewJSONHandler(errorFile, &slog.HandlerOptions{
            Level: slog.LevelError, // Only errors
        }),
        slog.NewJSONHandler(infoFile, &slog.HandlerOptions{
            Level: slog.LevelInfo, // Info and above
        }),
    },
    nil,
)

logging.InitLogging(logger)
logging.Info(context.Background(), "Goes to stdout and info.log")
logging.Error(context.Background(), "Goes to stdout, info.log, and error.log")

// Important: Flush and close files when done
errorFile.Sync()
errorFile.Close()
infoFile.Sync()
infoFile.Close()
```

### Structured Logging with Fields

```go
logging.Info(context.Background(), "User logged in",
    logging.NewField("user_id", 12345),
    logging.NewField("username", "john_doe"),
    logging.NewField("ip", "192.168.1.1"),
)

logging.Error(context.Background(), "Database connection failed",
    logging.NewField("error", err.Error()),
    logging.NewField("retry_count", 3),
    logging.NewField("database", "postgres"),
)
```

## Real-World Example

### Production Application Logging

```go
package main

import (
    "context"
    "fmt"
    "log/slog"
    "os"
    
    "github.com/fatkulnurk/foundation/logging"
)

func main() {
    // Setup production logging
    logger, err := setupProductionLogger()
    if err != nil {
        panic(err)
    }
    
    logging.InitLogging(logger)
    
    // Application code
    logging.Info(context.Background(), "Application started",
        logging.NewField("version", "1.0.0"),
        logging.NewField("environment", "production"),
    )
    
    // Simulate some operations
    processRequest()
    
    logging.Info(context.Background(), "Application shutdown")
}

func setupProductionLogger() (logging.Logger, error) {
    // Create logs directory
    if err := os.MkdirAll("logs", 0755); err != nil {
        return nil, err
    }
    
    // Open log files
    errorFile, err := os.OpenFile("logs/error.log", 
        os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    
    appFile, err := os.OpenFile("logs/app.log", 
        os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    
    // Create multi-output logger
    logger := logging.NewSlogLoggerWithMultiOutput(
        nil,
        []slog.Handler{
            // Console output (text, info and above)
            slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
                Level:     slog.LevelInfo,
                AddSource: false,
            }),
            // Error console (text, errors only)
            slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
                Level:     slog.LevelError,
                AddSource: true,
            }),
            // Application log (JSON, all levels)
            slog.NewJSONHandler(appFile, &slog.HandlerOptions{
                Level:     slog.LevelDebug,
                AddSource: true,
            }),
            // Error log (JSON, errors only)
            slog.NewJSONHandler(errorFile, &slog.HandlerOptions{
                Level:     slog.LevelError,
                AddSource: true,
            }),
        },
        nil,
    )
    
    return logger, nil
}

func processRequest() {
    ctx := context.Background()
    
    logging.Debug(ctx, "Processing request",
        logging.NewField("request_id", "req-123"),
    )
    
    logging.Info(ctx, "Request processed successfully",
        logging.NewField("request_id", "req-123"),
        logging.NewField("duration_ms", 150),
    )
    
    // Simulate error
    logging.Error(ctx, "Failed to connect to database",
        logging.NewField("request_id", "req-123"),
        logging.NewField("error", "connection timeout"),
        logging.NewField("retry_count", 3),
    )
}
```

## Best Practices

### 1. Initialize Logger Once at Startup

```go
// Good - initialize once
func main() {
    logger, _ := logging.NewSlogLoggerWithRotation("logs", nil)
    logging.InitLogging(logger)
    
    // Use throughout application
    logging.Info(context.Background(), "App started")
}

// Avoid - creating logger multiple times
func handler() {
    logger := logging.NewSlogLogger(nil) // Don't do this
}
```

### 2. Use Structured Fields

```go
// Good - structured fields
logging.Info(ctx, "User action",
    logging.NewField("user_id", userID),
    logging.NewField("action", "login"),
)

// Avoid - string concatenation
logging.Info(ctx, fmt.Sprintf("User %d performed login", userID))
```

### 3. Use Appropriate Log Levels

```go
// Debug - detailed information for debugging
logging.Debug(ctx, "Cache hit", logging.NewField("key", cacheKey))

// Info - general information
logging.Info(ctx, "Server started", logging.NewField("port", 8080))

// Warning - warning messages
logging.Warning(ctx, "High memory usage", logging.NewField("usage", "85%"))

// Error - error messages
logging.Error(ctx, "Database error", logging.NewField("error", err.Error()))
```

### 4. Separate Logs by Purpose

```go
// Production setup with separate files
logger := logging.NewSlogLoggerWithMultiOutput(
    nil,
    []slog.Handler{
        slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
        slog.NewJSONHandler(appLogFile, &slog.HandlerOptions{Level: slog.LevelDebug}),
        slog.NewJSONHandler(errorLogFile, &slog.HandlerOptions{Level: slog.LevelError}),
        slog.NewJSONHandler(auditLogFile, &slog.HandlerOptions{Level: slog.LevelInfo}),
    },
    nil,
)
```

### 5. Always Close File-Based Loggers

```go
// Good - close logger to flush buffers
func main() {
    logger, _ := logging.NewSlogLoggerWithFile("logs/app.log", nil)
    defer func() {
        if closer, ok := logger.(interface{ Close() error }); ok {
            closer.Close()
        }
    }()
    
    logging.InitLogging(logger)
    // ... application code ...
}

// For custom file handlers, close files explicitly
func setupCustomLogger() {
    errorFile, _ := os.OpenFile("logs/error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    
    // ... create logger ...
    
    // Close when done
    defer func() {
        errorFile.Sync()  // Flush buffers
        errorFile.Close() // Close file
    }()
}
```

### How MultiHandler Works

The `MultiHandler` allows you to write logs to multiple destinations simultaneously:

```go
// Create multiple handlers
handler1 := slog.NewJSONHandler(file1, &slog.HandlerOptions{Level: slog.LevelInfo})
handler2 := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
handler3 := slog.NewJSONHandler(file2, &slog.HandlerOptions{Level: slog.LevelError})

// Combine them with MultiHandler
multiHandler := logging.NewMultiHandler(handler1, handler2, handler3)
logger := slog.New(multiHandler)

// Now all logs go to all three destinations (respecting their level filters)
logger.Info("This goes to file1 and stdout")
logger.Error("This goes to all three: file1, stdout, and file2")
```

**Key Features:**
- Each handler can have its own level filter
- Each handler can have its own format (JSON, text, etc.)
- Writes are synchronous to all handlers
- If any handler fails, the error is returned

## Common Patterns

### Development vs Production Logging

```go
func setupLogger(env string) (logging.Logger, error) {
    if env == "development" {
        // Development: text to stdout only
        return logging.NewSlogLogger(nil), nil
    }
    
    // Production: multiple outputs with rotation
    return logging.NewSlogLoggerWithRotation("logs", &slog.HandlerOptions{
        Level:     slog.LevelInfo,
        AddSource: true,
    })
}
```

### Request ID Tracking

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    requestID := generateRequestID()
    ctx := context.WithValue(r.Context(), "request_id", requestID)
    
    logging.Info(ctx, "Request received",
        logging.NewField("request_id", requestID),
        logging.NewField("method", r.Method),
        logging.NewField("path", r.URL.Path),
    )
    
    // Process request...
    
    logging.Info(ctx, "Request completed",
        logging.NewField("request_id", requestID),
        logging.NewField("status", 200),
    )
}
```

### Error Logging with Stack Trace

```go
func processData(data []byte) error {
    if err := validate(data); err != nil {
        logging.Error(context.Background(), "Validation failed",
            logging.NewField("error", err.Error()),
            logging.NewField("data_size", len(data)),
            logging.NewField("stack", string(debug.Stack())),
        )
        return err
    }
    return nil
}
```

### Performance Monitoring

```go
func monitorPerformance(ctx context.Context, operation string, fn func() error) error {
    start := time.Now()
    
    logging.Debug(ctx, "Operation started",
        logging.NewField("operation", operation),
    )
    
    err := fn()
    duration := time.Since(start)
    
    if err != nil {
        logging.Error(ctx, "Operation failed",
            logging.NewField("operation", operation),
            logging.NewField("duration_ms", duration.Milliseconds()),
            logging.NewField("error", err.Error()),
        )
        return err
    }
    
    logging.Info(ctx, "Operation completed",
        logging.NewField("operation", operation),
        logging.NewField("duration_ms", duration.Milliseconds()),
    )
    
    return nil
}
```

## Installation

```bash
go get github.com/fatkulnurk/foundation/logging
```

## Dependencies

- **Slog**: Go standard library `log/slog` (Go 1.21+)
- **Zap**: `go.uber.org/zap`

---

## See Also

- `logging.go` - Core interface and types
- `slog.go` - Slog implementation with multi-output support
- `zap.go` - Zap implementation
- `example/` - Usage examples
