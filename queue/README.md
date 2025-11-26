# Queue Package

A simple, easy-to-understand task queue and background job processing package using Redis. Built on top of [hibiken/asynq](https://github.com/hibiken/asynq) with a clean abstraction layer that allows you to switch implementations without changing your code.

## Table of Contents

- [What is a Queue?](#what-is-a-queue)
- [Why Use a Queue?](#why-use-a-queue)
- [Core Concepts](#core-concepts)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Package Structure](#package-structure)
- [API Reference](#api-reference)
- [Configuration](#configuration)
- [Task Options](#task-options)
- [Task States](#task-states)
- [Middleware](#middleware)
- [Complete Examples](#complete-examples)
- [Real-World Use Cases](#real-world-use-cases)
- [Best Practices](#best-practices)
- [Running the Example](#running-the-example)
- [Troubleshooting](#troubleshooting)

---

## What is a Queue?

A **queue** is like a **line at a bank**. When you have many tasks to do but can't do them all at once, you put them in a queue. The tasks will be processed one by one (or in parallel).

**Simple Analogy:**
- **Queue** = Laundry basket
- **Worker** = Washing machine  
- **Task** = Dirty clothes
- **Handler** = Wash program (normal wash, quick wash, etc.)
- **Redis** = Storage for the basket
- **Middleware** = Extra features (fabric softener, fragrance, etc.)

---

## Why Use a Queue?

Imagine you have a website that needs to:
- Send emails to 1000 people
- Process large images
- Generate time-consuming reports

If you do everything immediately, your website becomes **slow** and users have to **wait a long time**.

**With a queue:**
1. Website immediately responds to user: "OK, will be processed"
2. Heavy work is done in the background
3. Users don't have to wait

**Real-world benefits:**
- ✅ **Faster response times** - Users don't wait for heavy tasks
- ✅ **Better reliability** - Failed tasks are automatically retried
- ✅ **Scalability** - Process multiple tasks in parallel
- ✅ **Scheduling** - Run tasks at specific times
- ✅ **Priority handling** - Important tasks processed first

---

## Core Concepts

### 1. Queue (The Line)
Where you **put tasks** to be processed later.

**Analogy:** Like putting dirty clothes in a laundry basket.

### 2. Worker (The Processor)
**Processes tasks** from the queue.

**Analogy:** Like a washing machine that takes clothes from the basket and washes them.

### 3. Task (The Job)
A job that needs to be done.

**Analogy:** One set of clothes to be washed.

### 4. Handler (The Processor Function)
Code that explains **how to process** a task.

**Analogy:** The wash program in the washing machine (normal, quick, delicate, etc.).

---

## Features

- ✅ **Abstracted Interface** - Not dependent on any specific queue library
- ✅ **Queue & Worker** - Producer and consumer pattern
- ✅ **Task Scheduling** - Schedule tasks to be processed later
- ✅ **Priority Queues** - Multiple queues with different priorities
- ✅ **Retry Mechanism** - Automatic retry for failed tasks
- ✅ **Unique Tasks** - Prevent duplicate tasks
- ✅ **Timeout & Deadline** - Control task execution time
- ✅ **Graceful Shutdown** - Safe worker shutdown
- ✅ **Type-Safe** - Full Go type safety
- ✅ **Middleware Support** - Extensible with middleware functions

---

## Installation

```bash
go get github.com/fatkulnurk/foundation/queue
```

**Dependencies:**
- Redis server (for storing the queue)
- Go 1.19 or higher

---

## Quick Start

### Step 1: Start Redis

```bash
# Using Docker (easiest)
docker run -d -p 6379:6379 redis

# Or using Homebrew (Mac)
brew install redis
brew services start redis

# Or using apt (Linux)
sudo apt install redis-server
sudo systemctl start redis
```

### Step 2: Producer (Enqueue Tasks)

```go
package main

import (
    "context"
    "github.com/fatkulnurk/foundation/queue"
    "github.com/redis/go-redis/v9"
)

func main() {
    // Connect to Redis
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // Create queue
    q, _ := queue.NewQueue(redisClient)

    // Enqueue a task
    q.Enqueue(context.Background(), "email:send", map[string]string{
        "to":      "user@example.com",
        "subject": "Welcome!",
        "body":    "Welcome to our service!",
    })
}
```

### Step 3: Worker (Process Tasks)

```go
package main

import (
    "context"
    "encoding/json"
    "github.com/fatkulnurk/foundation/queue"
    "github.com/redis/go-redis/v9"
)

func main() {
    // Connect to Redis
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // Create worker
    cfg := &queue.Config{
        Concurrency: 10,
    }
    w := queue.NewWorker(cfg, redisClient)

    // Register handler
    w.Register("email:send", func(ctx context.Context, payload []byte) error {
        var email map[string]string
        json.Unmarshal(payload, &email)
        
        // Send email here
        println("Sending email to:", email["to"])
        return nil
    })

    // Start worker
    w.Start()
}
```

---

## Package Structure

```
queue/
├── queue.go          # Interfaces & contracts
├── asynq.go          # Implementation using asynq
├── config.go         # Worker configuration
├── middleware.go     # Additional features for handlers
├── example/          # Usage examples
│   ├── main.go       # Example code
│   └── README.md     # How to run examples
└── go.mod            # Dependencies
```

### File Explanations

#### 1. `queue.go` - The Rules

Contains **interfaces** (contracts/rules) for the queue package.

**Analogy:** Like a manual book explaining what can be done.

**What's inside:**
- `Queue` interface: How to enqueue tasks
- `Worker` interface: How to process tasks
- `Handler`: Function that processes tasks
- `Option`: Additional settings for tasks

#### 2. `asynq.go` - The Implementation

Concrete implementation using the `asynq` library.

**Analogy:** If `queue.go` is the blueprint, then `asynq.go` is the **actual working washing machine**.

**What's inside:**
- `AsynqQueue`: Queue implementation using Redis
- `AsynqWorker`: Worker implementation using Redis

#### 3. `config.go` - Configuration

Settings for the worker.

**Analogy:** Like washing machine settings (spin speed, temperature, etc.).

**What's inside:**
- `Concurrency`: How many tasks can be processed simultaneously
- `Queues`: Queues with different priorities
- `StrictPriority`: Whether priority should be strict
- `ShutdownTimeout`: How long to wait before forcing shutdown

#### 4. `middleware.go` - Extra Features

Functions that wrap handlers to add features.

**Analogy:** Like **extra features in a washing machine** (fragrance, fabric softener, anti-wrinkle).

**What's inside:**
- `LoggingMiddleware`: Records logs
- `RecoveryMiddleware`: Catches errors/panics
- `RetryLoggingMiddleware`: Records retry attempts
- `TimeoutMiddleware`: Limits execution time
- `MetricsMiddleware`: Records statistics
- `ChainMiddleware`: Combines multiple middleware

---

## API Reference

### Queue Interface

```go
type Queue interface {
    // Enqueue adds a task to the queue
    Enqueue(ctx context.Context, taskName string, payload any, opts ...Option) (*OutputEnqueue, error)
    
    // GetTaskInfo retrieves information about a task by its ID
    GetTaskInfo(ctx context.Context, taskID string) (*TaskInfo, error)
    
    // Close closes the queue client connection
    Close() error
}
```

**Usage:**
```go
// Create queue
q, err := queue.NewQueue(redisClient)

// Enqueue task
result, err := q.Enqueue(ctx, "email:send", emailData)

// Get task info
taskInfo, err := q.GetTaskInfo(ctx, result.TaskID)

// Close connection
q.Close()
```

### Worker Interface

```go
type Worker interface {
    // Start starts the worker and begins processing tasks
    Start() error
    
    // Stop stops the worker gracefully
    Stop()
    
    // Register registers a handler for a specific task type
    Register(taskType string, handler Handler) error
    
    // RegisterWithMiddleware registers a handler with middleware functions
    RegisterWithMiddleware(taskType string, handler Handler, middleware ...MiddlewareFunc) error
    
    // GetTaskID retrieves the task ID from the context
    GetTaskID(ctx context.Context) (string, bool)
    
    // GetTaskInfo retrieves information about a task by its ID
    GetTaskInfo(ctx context.Context, taskID string) (*TaskInfo, error)
}
```

**Usage:**
```go
// Create worker
w := queue.NewWorker(cfg, redisClient)

// Register handler
w.Register("email:send", handlerFunc)

// Register with middleware
w.RegisterWithMiddleware("email:send", handlerFunc, 
    queue.LoggingMiddleware("email:send"),
    queue.RecoveryMiddleware("email:send"),
)

// Start worker
w.Start()

// Stop worker
w.Stop()
```

### Handler Type

```go
type Handler func(ctx context.Context, payload []byte) error
```

A function that processes a task. Receives payload (data), returns error if failed.

**Example:**
```go
func emailHandler(ctx context.Context, payload []byte) error {
    var email EmailData
    if err := json.Unmarshal(payload, &email); err != nil {
        return err
    }
    
    // Send email
    return sendEmail(email)
}
```

### TaskInfo

```go
type TaskInfo struct {
    ID            string      // Unique task ID
    Type          string      // Task type (email:send, etc.)
    Payload       []byte      // Task data
    State         TaskState   // Task status
    Queue         string      // Queue name
    MaxRetry      int         // Maximum retries
    Retried       int         // Number of retries so far
    LastError     string      // Last error message
    CompletedAt   *time.Time  // When completed
    NextProcessAt *time.Time  // When will be processed
}
```

---

## Configuration

```go
type Config struct {
    // Concurrency is the maximum number of concurrent processing of tasks
    // Default: 10
    Concurrency int

    // Queues is a map of queue names to their priority levels
    // Higher priority queues will be processed more frequently
    // Example: {"critical": 6, "default": 3, "low": 1}
    Queues map[string]int

    // StrictPriority indicates whether the queue priority should be treated strictly
    // If true, tasks in higher priority queues are processed first
    // Default: false
    StrictPriority bool

    // ShutdownTimeout is the duration to wait for workers to finish before forcing shutdown
    // Default: 8 seconds
    ShutdownTimeout int
}
```

### Configuration Examples

#### Basic Configuration
```go
cfg := &queue.Config{
    Concurrency: 10,
}
```

#### With Priority Queues
```go
cfg := &queue.Config{
    Concurrency: 20,
    Queues: map[string]int{
        "critical": 6,  // Processed 6x more frequently
        "default":  3,  // Processed 3x more frequently
        "low":      1,  // Processed least frequently
    },
}
```

**Analogy:** VIP line is served faster than regular line.

#### Strict Priority
```go
cfg := &queue.Config{
    Concurrency: 10,
    Queues: map[string]int{
        "critical": 6,
        "default":  3,
        "low":      1,
    },
    StrictPriority: true,  // Process critical until empty, then default, then low
}
```

**Analogy:**
- `true` = Serve VIP until empty, then serve regular
- `false` = Serve VIP more often, but still serve regular

---

## Task Options

Options are functions that configure task behavior. They use a flexible key-value pattern.

### Available Options

#### 1. MaxRetry
Retry the task if it fails.

```go
queue.Enqueue(ctx, "email:send", payload, 
    queue.MaxRetry(3),  // Retry up to 3 times
)
```

**Analogy:** If laundry fails, try washing again up to 3 times.

#### 2. QueueName (Priority)
Put task in a specific priority queue.

```go
queue.Enqueue(ctx, "payment:process", payload,
    queue.QueueName("critical"),  // High priority
)
```

**Analogy:** Put in VIP line instead of regular line.

#### 3. Timeout
Maximum time allowed for task execution.

```go
queue.Enqueue(ctx, "image:process", payload,
    queue.Timeout(5*time.Minute),  // Max 5 minutes
)
```

**Analogy:** If washing takes more than 5 minutes, force stop.

#### 4. ProcessIn (Delayed)
Process task after a certain duration.

```go
queue.Enqueue(ctx, "reminder:send", payload,
    queue.ProcessIn(24*time.Hour),  // Process in 24 hours
)
```

**Analogy:** Set timer on washing machine to start in 24 hours.

#### 5. ProcessAt (Scheduled)
Process task at a specific time.

```go
tomorrow := time.Now().Add(24*time.Hour)
queue.Enqueue(ctx, "report:generate", payload,
    queue.ProcessAt(tomorrow),  // Process tomorrow
)
```

**Analogy:** Set washing machine to start at 10 AM tomorrow.

#### 6. Unique
Prevent duplicate tasks for a duration.

```go
queue.Enqueue(ctx, "report:monthly", payload,
    queue.Unique(1*time.Hour),  // Unique for 1 hour
)
```

**Analogy:** If same clothes are already in the machine, don't add again.

#### 7. TaskID
Assign a custom ID to the task.

```go
queue.Enqueue(ctx, "order:process", payload,
    queue.TaskID("order-12345"),  // Custom ID
)
```

#### 8. Deadline
Task must complete before this absolute time.

```go
deadline := time.Now().Add(1*time.Hour)
queue.Enqueue(ctx, "urgent:task", payload,
    queue.Deadline(deadline),  // Must finish within 1 hour
)
```

#### 9. Retention
How long to keep task data after completion.

```go
queue.Enqueue(ctx, "log:cleanup", payload,
    queue.Retention(7*24*time.Hour),  // Keep for 7 days
)
```

#### 10. Group
Assign task to a specific group.

```go
queue.Enqueue(ctx, "user:sync", payload,
    queue.Group("user-operations"),
)
```

### Combining Options

```go
queue.Enqueue(ctx, "important:task", payload,
    queue.MaxRetry(3),
    queue.Timeout(30*time.Second),
    queue.QueueName("critical"),
    queue.Unique(1*time.Hour),
)
```

---

## Task States

Tasks can have different states during their lifecycle:

| State | Description | Analogy |
|-------|-------------|---------|
| `pending` | Waiting to be processed | In the basket, not washed yet |
| `active` | Currently being processed | In the washing machine now |
| `scheduled` | Scheduled for later | Scheduled for 10 AM |
| `retry` | Failed, will be retried | Failed, will be washed again |
| `archived` | Failed permanently | Failed completely, can't be washed |
| `completed` | Successfully completed | Clean and done |

**Checking task state:**
```go
taskInfo, _ := q.GetTaskInfo(ctx, taskID)
fmt.Println(taskInfo.State)  // "pending", "active", "completed", etc.
```

---

## Middleware

Middleware functions wrap handlers to add extra functionality.

### Built-in Middleware

#### 1. LoggingMiddleware
Records logs when task starts and completes.

```go
w.RegisterWithMiddleware("email:send", handler,
    queue.LoggingMiddleware("email:send"),
)
```

**Output:**
```
[email:send] Task started
[email:send] Task completed in 1.2s
```

**Analogy:** Recording in a book "10:00 - Started washing, 10:30 - Finished washing".

#### 2. RecoveryMiddleware
Catches panics so the program doesn't crash.

```go
w.RegisterWithMiddleware("email:send", handler,
    queue.RecoveryMiddleware("email:send"),
)
```

**Analogy:** If washing machine errors, shut down safely, don't explode.

#### 3. RetryLoggingMiddleware
Records when a task will be retried.

```go
w.RegisterWithMiddleware("email:send", handler,
    queue.RetryLoggingMiddleware("email:send"),
)
```

**Output:**
```
[email:send] Task will be retried: connection timeout
```

#### 4. TimeoutMiddleware
Adds a timeout to task execution.

```go
w.RegisterWithMiddleware("email:send", handler,
    queue.TimeoutMiddleware(5*time.Minute),
)
```

**Analogy:** If washing takes more than 5 minutes, force stop.

#### 5. MetricsMiddleware
Tracks task metrics (placeholder for actual metrics implementation).

```go
w.RegisterWithMiddleware("email:send", handler,
    queue.MetricsMiddleware("email:send"),
)
```

**Analogy:** Recording "Today washed 10 times, 9 successful, 1 failed".

#### 6. ChainMiddleware
Combines multiple middleware functions.

```go
w.RegisterWithMiddleware("email:send", handler,
    queue.ChainMiddleware(
        queue.LoggingMiddleware("email:send"),
        queue.RecoveryMiddleware("email:send"),
        queue.TimeoutMiddleware(30*time.Second),
    ),
)
```

**Analogy:** Using logging + recovery + timeout all at once.

---

## Complete Examples

### Example 1: Simple Email Task

```go
// Enqueue
q.Enqueue(ctx, "email:send", map[string]string{
    "to":      "user@example.com",
    "subject": "Hello",
    "body":    "Hello World",
})

// Worker
w.Register("email:send", func(ctx context.Context, payload []byte) error {
    var email map[string]string
    json.Unmarshal(payload, &email)
    return sendEmail(email)
})
```

### Example 2: Task with Retry and Timeout

```go
q.Enqueue(ctx, "image:process", map[string]string{
    "url": "https://example.com/image.jpg",
},
    queue.MaxRetry(3),
    queue.Timeout(5*time.Minute),
    queue.QueueName("critical"),
)
```

### Example 3: Scheduled Task

```go
// Send reminder in 24 hours
q.Enqueue(ctx, "reminder:send", map[string]string{
    "user_id": "user123",
    "message": "Don't forget!",
},
    queue.ProcessIn(24*time.Hour),
)
```

### Example 4: Unique Task

```go
// Prevent duplicate report generation
q.Enqueue(ctx, "report:generate", map[string]string{
    "report_id": "monthly-2024-11",
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
w.RegisterWithMiddleware("email:send",
    func(ctx context.Context, payload []byte) error {
        var email map[string]string
        json.Unmarshal(payload, &email)
        return sendEmail(email)
    },
    queue.LoggingMiddleware("email:send"),
    queue.RecoveryMiddleware("email:send"),
    queue.TimeoutMiddleware(30*time.Second),
)
```

### Example 7: Get Task Info

```go
// Enqueue task
result, _ := q.Enqueue(ctx, "email:send", emailPayload)

// Get task information
taskInfo, err := q.GetTaskInfo(ctx, result.TaskID)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Task ID: %s\n", taskInfo.ID)
fmt.Printf("State: %s\n", taskInfo.State)
fmt.Printf("Queue: %s\n", taskInfo.Queue)
fmt.Printf("Retried: %d/%d\n", taskInfo.Retried, taskInfo.MaxRetry)

if taskInfo.NextProcessAt != nil {
    fmt.Printf("Next process at: %s\n", taskInfo.NextProcessAt)
}
```

### Example 8: Graceful Shutdown

```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

go func() {
    worker.Start()
}()

<-sigChan
worker.Stop()  // Waits for running tasks to finish
```

---

## Real-World Use Cases

### 1. Send Mass Emails

```go
// Enqueue 1000 emails
for _, user := range users {
    q.Enqueue(ctx, "email:send", map[string]string{
        "to":      user.Email,
        "subject": "Newsletter",
        "body":    "Check out our latest news!",
    })
}
// User gets immediate response, emails sent in background
```

### 2. Process Images

```go
// Upload image
q.Enqueue(ctx, "image:resize", map[string]string{
    "url": uploadedImage,
},
    queue.Timeout(5*time.Minute),
    queue.MaxRetry(2),
)
// User doesn't wait for resize to complete
```

### 3. Generate Reports

```go
// Generate monthly report
q.Enqueue(ctx, "report:monthly", map[string]string{
    "month": "November",
},
    queue.ProcessAt(tomorrow),      // Process tomorrow
    queue.Unique(24*time.Hour),     // Prevent duplicates
)
```

### 4. Database Backup

```go
// Schedule nightly backup
q.Enqueue(ctx, "backup:database", map[string]string{
    "database": "production",
},
    queue.ProcessAt(midnight),
    queue.QueueName("low"),
)
```

### 5. Payment Processing

```go
// Critical payment task
q.Enqueue(ctx, "payment:process", paymentData,
    queue.QueueName("critical"),
    queue.MaxRetry(3),
    queue.Timeout(30*time.Second),
    queue.Unique(5*time.Minute),
)
```

---

## Best Practices

### 1. Use Typed Payloads

```go
type EmailPayload struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

// Good
q.Enqueue(ctx, "email:send", EmailPayload{
    To:      "user@example.com",
    Subject: "Welcome",
    Body:    "Hello!",
})
```

### 2. Handle Errors Properly

```go
w.Register("task", func(ctx context.Context, payload []byte) error {
    // Return error to trigger retry
    if err := process(); err != nil {
        return fmt.Errorf("processing failed: %w", err)
    }
    return nil
})
```

### 3. Use Context for Cancellation

```go
w.Register("task", func(ctx context.Context, payload []byte) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(1 * time.Second):
        // Do work
    }
    return nil
})
```

### 4. Set Appropriate Timeouts

```go
queue.Enqueue(ctx, "long-task", payload,
    queue.Timeout(10*time.Minute),
    queue.MaxRetry(2),
)
```

### 5. Use Priority Queues

```go
// Critical tasks
queue.Enqueue(ctx, "payment:process", payload, 
    queue.QueueName("critical"))

// Low priority tasks
queue.Enqueue(ctx, "analytics:update", payload, 
    queue.QueueName("low"))
```

### 6. Implement Graceful Shutdown

```go
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

go func() {
    worker.Start()
}()

<-sigChan
worker.Stop()  // Graceful shutdown
```

### 7. Use Middleware for Cross-Cutting Concerns

```go
w.RegisterWithMiddleware("task", handler,
    queue.LoggingMiddleware("task"),
    queue.RecoveryMiddleware("task"),
    queue.MetricsMiddleware("task"),
)
```

### 8. Monitor Task Status

```go
// After enqueuing
result, _ := q.Enqueue(ctx, "task", payload)

// Check status later
taskInfo, _ := q.GetTaskInfo(ctx, result.TaskID)
if taskInfo.State == queue.TaskStateArchived {
    log.Printf("Task failed permanently: %s", taskInfo.LastError)
}
```

---

## Running the Example

The `example/` directory contains a complete working example.

### Terminal 1: Start Worker

```bash
cd pkg/queue/example
go run main.go worker
```

**Expected output:**
```
=== Example: Worker Processing Tasks ===

Registering handlers for each task type...
✓ All handlers registered successfully

🚀 Starting worker...
Press Ctrl+C to stop
```

### Terminal 2: Enqueue Tasks

```bash
cd pkg/queue/example
go run main.go
```

**Expected output:**
```
=== Example: Enqueuing Tasks ===

1. Enqueuing simple email task...
✓ Task enqueued successfully: ID=abc123

2. Enqueuing email with retry and timeout...
✓ Task enqueued: ID=def456 (critical queue, 3 retries, 30s timeout)

... and so on
```

### Terminal 1: See Worker Processing

```
📋 Processing task ID: abc123
📧 Sending email to user@example.com: Welcome!
✓ Email sent successfully to user@example.com

📋 Task ID: def456
🔔 Sending notification to user user123: Your order has been shipped!
✓ Notification sent successfully

... and so on
```

---

## Troubleshooting

### Error: "Failed to connect to Redis"

**Solution:**
1. Make sure Redis is running: `redis-cli ping`
2. Check port: Default is 6379
3. Check firewall settings

### Tasks Not Being Processed

**Solution:**
1. Make sure worker is running
2. Check handler is registered for the task type
3. Check error logs in worker

### Tasks Keep Failing

**Solution:**
1. Check error message in logs
2. Make sure payload can be unmarshaled
3. Check logic in handler
4. Increase timeout if needed

### Check Queue in Redis

```bash
# Enter Redis CLI
redis-cli

# View all keys
KEYS *

# View queue contents
LRANGE asynq:{default}:pending 0 -1
LRANGE asynq:{critical}:pending 0 -1
LRANGE asynq:{low}:pending 0 -1
```

### Check Active Tasks

```bash
redis-cli
SMEMBERS asynq:{default}:active
```

### Clear All Tasks (Reset)

```bash
redis-cli FLUSHALL
```

---

## Architecture

```
Your Application
        ↓
    Queue Interface (abstraction)
        ↓
    AsynqQueue (implementation)
        ↓
    hibiken/asynq (internal only)
        ↓
    Redis
```

**Key Point:** Your code only depends on the `Queue` and `Worker` interfaces, not on `asynq` or any specific implementation.

---

## Benefits of This Package

### 1. **Abstraction**
Your code doesn't depend on a specific library (asynq). If you want to switch to another library tomorrow, just change the implementation, your code stays the same.

### 2. **Type-Safe**
Uses Go's type system, so it's safer and your IDE can help with autocomplete.

### 3. **Flexible Options**
Can add new options without changing existing code.

### 4. **Middleware Support**
Can add features (logging, metrics, etc.) without changing handlers.

### 5. **Graceful Shutdown**
When shutting down the worker, it waits for running tasks to finish first.

### 6. **Easy to Use**
Simple API that's easy to understand and use.

### 7. **Production Ready**
Built on top of battle-tested asynq library.

---

## Migration Guide

If you want to switch from asynq to another library:

1. Create a new file (e.g., `rabbitmq.go`)
2. Implement `Queue` and `Worker` interfaces
3. Update `NewQueue()` and `NewWorker()` to use the new implementation
4. External code doesn't need to change!

---

## License

MIT

---

## Summary

The queue package is a **task queue system** that:
- Separates heavy work from user responses
- Can process many tasks in parallel
- Has features like retry, scheduling, priority, etc.
- Is easy to use and maintain

**Complete Analogy:**
- Queue = Laundry basket
- Worker = Washing machine
- Task = Dirty clothes
- Handler = Wash program
- Redis = Storage for the basket
- Middleware = Extra features (fragrance, fabric softener, etc.)

Now you can use the queue package for your own applications! 🚀
