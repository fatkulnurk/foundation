# Queue Package

A flexible and extensible asynchronous task queue interface for the application with an Asynq implementation.

## Overview

The `queue` package provides a standardized interface for asynchronous task processing within the application. It defines a common interface (`IQueue`) that can be implemented by various queue providers, with Asynq being the default implementation.

## Interface

The package defines the `IQueue` interface with the following method:

```go
type IQueue interface {
	Enqueue(ctx context.Context, taskName string, payload any, opts ...Option) (*OutputEnqueue, error)
}
```

### Methods

- **Enqueue**: Adds a task to the queue with the specified name, payload, and options

## Task Options

The package provides a flexible options system for configuring tasks. Available options include:

- **MaxRetry(n int)**: Sets the maximum number of retry attempts for a task
- **Queue(name string)**: Sets the queue name for a task
- **Timeout(d time.Duration)**: Sets the maximum execution time for a task
- **Deadline(t time.Time)**: Sets the absolute time after which a task will fail if still running
- **Unique(d time.Duration)**: Makes the task unique for the specified duration
- **ProcessAt(t time.Time)**: Schedules a task to be processed at a specific time
- **ProcessIn(d time.Duration)**: Schedules a task to be processed after the specified duration
- **TaskID(id string)**: Assigns a custom ID to a task
- **Retention(d time.Duration)**: Sets how long task data will be kept after completion
- **Group(name string)**: Assigns a task to a specific group

## Implementations

### Asynq Queue

The package includes an implementation using the [Asynq](https://github.com/hibiken/asynq) library, which provides Redis-backed distributed task queue:

```go
type AsynqQueue struct {
	client *asynq.Client
}
```

#### Usage

```go
import (
	"context"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"your-project/pkg/queue"
	"your-project/config"
)

// Initialize Redis client
redisClient := redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

// Initialize Asynq client
asynqClient, err := queue.NewAsynqClient(&config.Queue{}, redisClient)
if err != nil {
	// Handle error
}

// Create queue instance
queueInstance := queue.NewAsynqQueue(asynqClient)

// Use the queue
ctx := context.Background()
payload := map[string]interface{}{
	"user_id": 123,
	"email": "user@example.com",
}

// Enqueue a task with options
result, err := queueInstance.Enqueue(
	ctx,
	"email:send",
	payload,
	queue.MaxRetry(3),
	queue.Queue("critical"),
	queue.ProcessIn(30*time.Minute),
)
```

## Extending

To implement a new queue provider, create a struct that implements the `IQueue` interface.

## Thread Safety

All implementations are designed to be thread-safe and can be safely used concurrently from multiple goroutines.