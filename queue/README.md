# Queue Package

Package queue menyediakan abstraksi untuk task queue dan background job processing menggunakan Redis. Package ini menggunakan [hibiken/asynq](https://github.com/hibiken/asynq) sebagai implementasi internal, tetapi API-nya tidak tergantung pada asynq sehingga implementasi bisa diganti tanpa mengubah kode yang menggunakan package ini.

## Features

- ✅ **Abstracted Interface** - Tidak tergantung pada library queue tertentu
- ✅ **Queue & Worker** - Producer dan consumer pattern
- ✅ **Task Scheduling** - Schedule tasks untuk diproses nanti
- ✅ **Priority Queues** - Multiple queues dengan priority berbeda
- ✅ **Retry Mechanism** - Automatic retry untuk failed tasks
- ✅ **Unique Tasks** - Prevent duplicate tasks
- ✅ **Timeout & Deadline** - Control task execution time
- ✅ **Graceful Shutdown** - Safe worker shutdown
- ✅ **Type-Safe** - Full Go type safety

## Installation

```go
import "github.com/fatkulnurk/foundation/queue"
```

## Quick Start

### 1. Producer (Enqueue Tasks)

```go
// Setup Redis
redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

// Create queue
q, err := queue.NewQueue(redisClient)
if err != nil {
    log.Fatal(err)
}

// Enqueue task
result, err := q.Enqueue(ctx, "email:send", EmailPayload{
    To:      "user@example.com",
    Subject: "Welcome!",
    Body:    "Welcome to our service!",
})
```

### 2. Worker (Process Tasks)

```go
// Create worker config
cfg := &queue.Config{
    Concurrency: 10,
}

// Create worker
worker := queue.NewWorker(cfg, redisClient)

// Register handler
worker.Register("email:send", func(ctx context.Context, payload []byte) error {
    var email EmailPayload
    json.Unmarshal(payload, &email)
    
    // Process email
    return sendEmail(email)
})

// Start worker
worker.Start()
```

## Interfaces

### Queue Interface

```go
type Queue interface {
    Enqueue(ctx context.Context, taskName string, payload any, opts ...Option) (*OutputEnqueue, error)
}
```

### Worker Interface

```go
type Worker interface {
    Start() error
    Stop()
    Register(taskType string, handler Handler) error
}
```

### Handler Type

```go
type Handler func(ctx context.Context, payload []byte) error
```

### Middleware Type

```go
type MiddlewareFunc func(Handler) Handler
```

Middleware allows you to wrap handlers with additional functionality like logging, metrics, recovery, etc.

## Configuration

```go
type Config struct {
    // Concurrency is the maximum number of concurrent processing of tasks
    Concurrency int

    // Queues is a map of queue names to their priority levels
    Queues map[string]int

    // StrictPriority indicates whether the queue priority should be treated strictly
    StrictPriority bool

    // ShutdownTimeout is the duration to wait for workers to finish before forcing shutdown
    ShutdownTimeout int
}
```

## Task Options

### MaxRetry
```go
queue.Enqueue(ctx, "task", payload, queue.MaxRetry(3))
```

### QueueName (Priority)
```go
queue.Enqueue(ctx, "task", payload, queue.QueueName("critical"))
```

### Timeout
```go
queue.Enqueue(ctx, "task", payload, queue.Timeout(30*time.Second))
```

### ProcessIn (Delayed)
```go
queue.Enqueue(ctx, "task", payload, queue.ProcessIn(5*time.Minute))
```

### ProcessAt (Scheduled)
```go
queue.Enqueue(ctx, "task", payload, queue.ProcessAt(tomorrow))
```

### Unique
```go
queue.Enqueue(ctx, "task", payload, queue.Unique(1*time.Hour))
```

### TaskID
```go
queue.Enqueue(ctx, "task", payload, queue.TaskID("custom-id"))
```

### Deadline
```go
queue.Enqueue(ctx, "task", payload, queue.Deadline(time.Now().Add(1*time.Hour)))
```

### Retention
```go
queue.Enqueue(ctx, "task", payload, queue.Retention(7*24*time.Hour))
```

### Group
```go
queue.Enqueue(ctx, "task", payload, queue.Group("user-operations"))
```

## Built-in Middleware

### LoggingMiddleware
Logs task execution start, completion, and errors with duration.
```go
worker.RegisterWithMiddleware("task", handler, 
    queue.LoggingMiddleware("task"),
)
```

### RecoveryMiddleware
Recovers from panics and converts them to errors.
```go
worker.RegisterWithMiddleware("task", handler,
    queue.RecoveryMiddleware("task"),
)
```

### RetryLoggingMiddleware
Logs when a task will be retried.
```go
worker.RegisterWithMiddleware("task", handler,
    queue.RetryLoggingMiddleware("task"),
)
```

### TimeoutMiddleware
Adds a timeout to task execution.
```go
worker.RegisterWithMiddleware("task", handler,
    queue.TimeoutMiddleware(5*time.Minute),
)
```

### MetricsMiddleware
Tracks task metrics (placeholder for actual metrics implementation).
```go
worker.RegisterWithMiddleware("task", handler,
    queue.MetricsMiddleware("task"),
)
```

### ChainMiddleware
Chains multiple middleware functions.
```go
worker.RegisterWithMiddleware("task", handler,
    queue.ChainMiddleware(
        queue.LoggingMiddleware("task"),
        queue.RecoveryMiddleware("task"),
        queue.MetricsMiddleware("task"),
    ),
)
```

## Examples

### Example 1: Simple Task

```go
// Enqueue
q.Enqueue(ctx, "email:send", EmailPayload{
    To:      "user@example.com",
    Subject: "Hello",
    Body:    "Hello World",
})

// Worker
worker.Register("email:send", func(ctx context.Context, payload []byte) error {
    var email EmailPayload
    json.Unmarshal(payload, &email)
    return sendEmail(email)
})
```

### Example 2: Task with Retry and Timeout

```go
q.Enqueue(ctx, "image:process", ImagePayload{
    URL: "https://example.com/image.jpg",
},
    queue.MaxRetry(3),
    queue.Timeout(5*time.Minute),
    queue.QueueName("critical"),
)
```

### Example 3: Scheduled Task

```go
// Send reminder in 24 hours
q.Enqueue(ctx, "reminder:send", ReminderPayload{
    UserID:  "user123",
    Message: "Don't forget!",
},
    queue.ProcessIn(24*time.Hour),
)
```

### Example 4: Unique Task

```go
// Prevent duplicate report generation
q.Enqueue(ctx, "report:generate", ReportPayload{
    ReportID: "monthly-2024-11",
},
    queue.Unique(1*time.Hour),
    queue.TaskID("report-monthly-2024-11"),
)
```

### Example 5: Multiple Queues with Priority

```go
cfg := &queue.Config{
    Concurrency: 10,
    Queues: map[string]int{
        "critical": 6,  // Processed more frequently
        "default":  3,
        "low":      1,  // Processed less frequently
    },
}

worker := queue.NewWorker(cfg, redisClient)
```

### Example 6: Handler with Middleware

```go
// Register handler with logging and recovery middleware
worker.RegisterWithMiddleware("email:send",
    func(ctx context.Context, payload []byte) error {
        var email EmailPayload
        json.Unmarshal(payload, &email)
        return sendEmail(email)
    },
    queue.LoggingMiddleware("email:send"),
    queue.RecoveryMiddleware("email:send"),
    queue.TimeoutMiddleware(30*time.Second),
)
```

## Running the Example

```bash
# Terminal 1: Start worker
cd pkg/queue/example
go run main.go worker

# Terminal 2: Enqueue tasks
go run main.go
```

## Best Practices

### 1. **Use Typed Payloads**
```go
type EmailPayload struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}
```

### 2. **Handle Errors Properly**
```go
worker.Register("task", func(ctx context.Context, payload []byte) error {
    // Return error to trigger retry
    if err := process(); err != nil {
        return fmt.Errorf("processing failed: %w", err)
    }
    return nil
})
```

### 3. **Use Context for Cancellation**
```go
worker.Register("task", func(ctx context.Context, payload []byte) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(1 * time.Second):
        // Do work
    }
    return nil
})
```

### 4. **Set Appropriate Timeouts**
```go
queue.Enqueue(ctx, "long-task", payload,
    queue.Timeout(10*time.Minute),
    queue.MaxRetry(2),
)
```

### 5. **Use Priority Queues**
```go
// Critical tasks
queue.Enqueue(ctx, "payment:process", payload, queue.QueueName("critical"))

// Low priority tasks
queue.Enqueue(ctx, "analytics:update", payload, queue.QueueName("low"))
```

### 6. **Graceful Shutdown**
```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

go func() {
    worker.Start()
}()

<-sigChan
worker.Stop() // Graceful shutdown
```

## Architecture

```
External Code (Your App)
        ↓
    Queue Interface (abstraction)
        ↓
    AsynqQueue (implementation)
        ↓
    hibiken/asynq (internal only)
        ↓
    Redis
```

**Key Point**: External code hanya depend pada `Queue` dan `Worker` interface, bukan pada `asynq` atau implementasi spesifik lainnya.

## Migration Guide

Jika ingin mengganti implementasi dari asynq ke library lain:

1. Buat file baru (e.g., `rabbitmq.go`)
2. Implement `Queue` dan `Worker` interface
3. Update `NewQueue()` dan `NewWorker()` untuk menggunakan implementasi baru
4. Kode external tidak perlu diubah!

## License

MIT
