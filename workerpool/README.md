# Worker Pool Package

A simple, efficient worker pool implementation for Go with priority queue, retry mechanism, timeout support, and dynamic scaling.

## Table of Contents

- [What is Worker Pool?](#what-is-worker-pool)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Priority Levels](#priority-levels)
- [Job Configuration](#job-configuration)
- [Complete Examples](#complete-examples)
- [API Reference](#api-reference)
- [Best Practices](#best-practices)

---

## What is Worker Pool?

A **worker pool** is like having a team of workers that process jobs from a queue. Instead of creating a new goroutine for every task (which can be expensive), you maintain a fixed number of workers that process tasks one by one.

**Simple Analogy:**
- **Worker Pool** = A team of workers at a factory
- **Worker** = One person who does the work
- **Job** = A task to be done
- **Queue** = Line of tasks waiting to be processed
- **Priority** = Some tasks are more urgent than others
- **Retry** = Try again if the task fails

**Why use a worker pool?**
-  **Control concurrency** - Limit how many tasks run at once
-  **Prevent resource exhaustion** - Don't create unlimited goroutines
-  **Priority handling** - Important tasks processed first
-  **Automatic retry** - Failed tasks are retried automatically
-  **Timeout support** - Tasks don't run forever
-  **Dynamic scaling** - Add or remove workers on the fly

---

## Features

-  **Fixed Worker Count** - Control concurrency with worker limit
-  **Priority Queue** - High, Normal, Low priority levels
-  **FIFO within Priority** - Fair processing within same priority
-  **Automatic Retry** - Configurable retry attempts with delay
-  **Timeout Support** - Per-job timeout configuration
-  **Panic Recovery** - Workers recover from panics
-  **Dynamic Scaling** - Scale workers up or down
-  **Graceful Shutdown** - Wait for running jobs to complete
-  **Thread-Safe** - Safe for concurrent use
-  **Zero Dependencies** - Only uses standard library

---

## Installation

```bash
go get github.com/fatkulnurk/foundation/workerpool
```

**Dependencies:**
- Go 1.25 or higher
- Standard library only

---

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/fatkulnurk/foundation/workerpool"
)

func main() {
    // Create worker pool with 3 workers
    pool := workerpool.NewWorkerPool(3)
    defer pool.Stop()
    
    // Submit a simple job
    pool.Submit(workerpool.Job{
        Task: func(ctx context.Context) error {
            fmt.Println("Processing job...")
            time.Sleep(1 * time.Second)
            fmt.Println("Job completed!")
            return nil
        },
        Retry:      2,                    // Retry 2 times if failed
        Timeout:    5 * time.Second,      // Timeout after 5 seconds
        RetryDelay: 1 * time.Second,      // Wait 1 second between retries
        Priority:   workerpool.Normal,    // Normal priority
    })
    
    // Wait for jobs to complete
    time.Sleep(2 * time.Second)
}
```

### With Priority

```go
pool := workerpool.NewWorkerPool(5)
defer pool.Stop()

// High priority job (processed first)
pool.Submit(workerpool.Job{
    Task: func(ctx context.Context) error {
        fmt.Println("High priority task")
        return nil
    },
    Priority: workerpool.High,
    Timeout:  5 * time.Second,
})

// Normal priority job
pool.Submit(workerpool.Job{
    Task: func(ctx context.Context) error {
        fmt.Println("Normal priority task")
        return nil
    },
    Priority: workerpool.Normal,
    Timeout:  5 * time.Second,
})

// Low priority job (processed last)
pool.Submit(workerpool.Job{
    Task: func(ctx context.Context) error {
        fmt.Println("Low priority task")
        return nil
    },
    Priority: workerpool.Low,
    Timeout:  5 * time.Second,
})
```

---

## Priority Levels

The worker pool supports three priority levels:

```go
const (
    Low    Priority = 0  // Lowest priority
    Normal Priority = 1  // Default priority
    High   Priority = 2  // Highest priority
)
```

**Processing order:**
1. **High priority** jobs are processed first
2. **Normal priority** jobs are processed next
3. **Low priority** jobs are processed last
4. Within the same priority, jobs are processed **FIFO** (First In, First Out)

**Example:**
```
Queue: [High-1, Low-1, High-2, Normal-1, Low-2]
Processing order: High-1 → High-2 → Normal-1 → Low-1 → Low-2
```

---

## Job Configuration

```go
type Job struct {
    // Task is the function to execute
    Task func(ctx context.Context) error
    
    // Retry is the number of retry attempts (0 = no retry)
    Retry int
    
    // Timeout is the maximum duration for the task
    Timeout time.Duration
    
    // RetryDelay is the delay between retry attempts
    RetryDelay time.Duration
    
    // Priority determines processing order
    Priority Priority
    
    // CreatedAt is set automatically (for FIFO ordering)
    CreatedAt time.Time
}
```

### Configuration Examples

#### Simple Job (No Retry)
```go
workerpool.Job{
    Task: func(ctx context.Context) error {
        // Do work
        return nil
    },
    Timeout: 5 * time.Second,
}
```

#### Job with Retry
```go
workerpool.Job{
    Task: func(ctx context.Context) error {
        // Do work that might fail
        return nil
    },
    Retry:      3,                   // Try up to 3 times
    Timeout:    10 * time.Second,    // 10 second timeout per attempt
    RetryDelay: 2 * time.Second,     // Wait 2 seconds between retries
}
```

#### High Priority Job
```go
workerpool.Job{
    Task: func(ctx context.Context) error {
        // Urgent work
        return nil
    },
    Priority: workerpool.High,
    Timeout:  5 * time.Second,
}
```

---

## Complete Examples

### Example 1: Basic Worker Pool

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/fatkulnurk/foundation/workerpool"
)

func main() {
    // Create pool with 3 workers
    pool := workerpool.NewWorkerPool(3)
    defer pool.Stop()
    
    // Submit 10 jobs
    for i := 1; i <= 10; i++ {
        jobID := i
        pool.Submit(workerpool.Job{
            Task: func(ctx context.Context) error {
                fmt.Printf("Processing job %d\n", jobID)
                time.Sleep(1 * time.Second)
                fmt.Printf("Job %d completed\n", jobID)
                return nil
            },
            Timeout: 5 * time.Second,
        })
    }
    
    // Wait for all jobs to complete
    time.Sleep(5 * time.Second)
}
```

### Example 2: Priority Queue

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/fatkulnurk/foundation/workerpool"
)

func main() {
    pool := workerpool.NewWorkerPool(2)
    defer pool.Stop()
    
    // Submit jobs with different priorities
    priorities := []struct {
        name     string
        priority workerpool.Priority
    }{
        {"Low-1", workerpool.Low},
        {"High-1", workerpool.High},
        {"Normal-1", workerpool.Normal},
        {"High-2", workerpool.High},
        {"Low-2", workerpool.Low},
        {"Normal-2", workerpool.Normal},
    }
    
    for _, p := range priorities {
        name := p.name
        pool.Submit(workerpool.Job{
            Task: func(ctx context.Context) error {
                fmt.Printf("Processing: %s\n", name)
                time.Sleep(500 * time.Millisecond)
                return nil
            },
            Priority: p.priority,
            Timeout:  5 * time.Second,
        })
    }
    
    time.Sleep(5 * time.Second)
}
```

**Output:**
```
Processing: High-1
Processing: High-2
Processing: Normal-1
Processing: Normal-2
Processing: Low-1
Processing: Low-2
```

### Example 3: Retry Mechanism

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "math/rand"
    "time"
    "github.com/fatkulnurk/foundation/workerpool"
)

func main() {
    pool := workerpool.NewWorkerPool(2)
    defer pool.Stop()
    
    // Job that might fail
    pool.Submit(workerpool.Job{
        Task: func(ctx context.Context) error {
            // Simulate random failure
            if rand.Float32() < 0.7 {
                return errors.New("random failure")
            }
            fmt.Println("Job succeeded!")
            return nil
        },
        Retry:      3,                   // Retry up to 3 times
        Timeout:    5 * time.Second,     // 5 second timeout
        RetryDelay: 1 * time.Second,     // Wait 1 second between retries
        Priority:   workerpool.Normal,
    })
    
    time.Sleep(10 * time.Second)
}
```

**Output (example):**
```
[Worker 0]  Job failed (Attempt 1): random failure
[Worker 0]  Job failed (Attempt 2): random failure
[Worker 0]  Job success (Attempt 3)
```

### Example 4: Timeout Handling

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/fatkulnurk/foundation/workerpool"
)

func main() {
    pool := workerpool.NewWorkerPool(2)
    defer pool.Stop()
    
    // Job that takes too long
    pool.Submit(workerpool.Job{
        Task: func(ctx context.Context) error {
            select {
            case <-time.After(10 * time.Second):
                fmt.Println("Job completed")
                return nil
            case <-ctx.Done():
                fmt.Println("Job timed out!")
                return ctx.Err()
            }
        },
        Timeout: 2 * time.Second, // Will timeout after 2 seconds
        Retry:   1,
    })
    
    time.Sleep(5 * time.Second)
}
```

### Example 5: Dynamic Scaling

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/fatkulnurk/foundation/workerpool"
)

func main() {
    // Start with 2 workers
    pool := workerpool.NewWorkerPool(2)
    defer pool.Stop()
    
    // Submit some jobs
    for i := 1; i <= 5; i++ {
        jobID := i
        pool.Submit(workerpool.Job{
            Task: func(ctx context.Context) error {
                fmt.Printf("Job %d processing\n", jobID)
                time.Sleep(2 * time.Second)
                return nil
            },
            Timeout: 5 * time.Second,
        })
    }
    
    time.Sleep(1 * time.Second)
    
    // Scale up to 5 workers
    fmt.Println("Scaling up to 5 workers...")
    pool.ScaleTo(5)
    
    // Submit more jobs
    for i := 6; i <= 10; i++ {
        jobID := i
        pool.Submit(workerpool.Job{
            Task: func(ctx context.Context) error {
                fmt.Printf("Job %d processing\n", jobID)
                time.Sleep(2 * time.Second)
                return nil
            },
            Timeout: 5 * time.Second,
        })
    }
    
    time.Sleep(5 * time.Second)
    
    // Scale down to 2 workers
    fmt.Println("Scaling down to 2 workers...")
    pool.ScaleTo(2)
    
    time.Sleep(2 * time.Second)
}
```

### Example 6: Panic Recovery

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/fatkulnurk/foundation/workerpool"
)

func main() {
    pool := workerpool.NewWorkerPool(2)
    defer pool.Stop()
    
    // Job that panics
    pool.Submit(workerpool.Job{
        Task: func(ctx context.Context) error {
            fmt.Println("About to panic...")
            panic("something went wrong!")
        },
        Retry:   2,
        Timeout: 5 * time.Second,
    })
    
    // Normal job (will still be processed)
    pool.Submit(workerpool.Job{
        Task: func(ctx context.Context) error {
            fmt.Println("Normal job processing")
            return nil
        },
        Timeout: 5 * time.Second,
    })
    
    time.Sleep(3 * time.Second)
}
```

**Output:**
```
About to panic...
[Worker 0] PANIC: something went wrong!
[stack trace...]
[Worker 0]  Job failed (Attempt 1): panic: something went wrong!
About to panic...
[Worker 0] PANIC: something went wrong!
[Worker 0]  Job failed (Attempt 2): panic: something went wrong!
[Worker 0]  Job permanently failed after 3 attempts
Normal job processing
[Worker 1]  Job success (Attempt 1)
```

### Example 7: HTTP Request Processing

```go
package main

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "time"
    "github.com/fatkulnurk/foundation/workerpool"
)

func main() {
    pool := workerpool.NewWorkerPool(5)
    defer pool.Stop()
    
    urls := []string{
        "https://api.github.com/users/golang",
        "https://api.github.com/users/microsoft",
        "https://api.github.com/users/google",
    }
    
    for _, url := range urls {
        targetURL := url
        pool.Submit(workerpool.Job{
            Task: func(ctx context.Context) error {
                req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
                if err != nil {
                    return err
                }
                
                resp, err := http.DefaultClient.Do(req)
                if err != nil {
                    return err
                }
                defer resp.Body.Close()
                
                body, err := io.ReadAll(resp.Body)
                if err != nil {
                    return err
                }
                
                fmt.Printf("Fetched %s: %d bytes\n", targetURL, len(body))
                return nil
            },
            Retry:      2,
            Timeout:    10 * time.Second,
            RetryDelay: 1 * time.Second,
            Priority:   workerpool.Normal,
        })
    }
    
    time.Sleep(15 * time.Second)
}
```

### Example 8: Database Operations

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/fatkulnurk/foundation/workerpool"
)

type User struct {
    ID    int
    Name  string
    Email string
}

func main() {
    pool := workerpool.NewWorkerPool(3)
    defer pool.Stop()
    
    users := []User{
        {ID: 1, Name: "John", Email: "john@example.com"},
        {ID: 2, Name: "Jane", Email: "jane@example.com"},
        {ID: 3, Name: "Bob", Email: "bob@example.com"},
    }
    
    // Process each user
    for _, user := range users {
        u := user
        pool.Submit(workerpool.Job{
            Task: func(ctx context.Context) error {
                // Simulate database operation
                fmt.Printf("Saving user %d: %s\n", u.ID, u.Name)
                time.Sleep(1 * time.Second)
                
                // Simulate potential error
                if u.ID == 2 {
                    return fmt.Errorf("database error for user %d", u.ID)
                }
                
                fmt.Printf("User %d saved successfully\n", u.ID)
                return nil
            },
            Retry:      2,
            Timeout:    5 * time.Second,
            RetryDelay: 500 * time.Millisecond,
            Priority:   workerpool.Normal,
        })
    }
    
    time.Sleep(10 * time.Second)
}
```

---

## API Reference

### WorkerPool

```go
type WorkerPool struct {
    // Internal fields (not exported)
}
```

#### `NewWorkerPool(initialWorkerCount int) *WorkerPool`

Creates a new worker pool with the specified number of workers.

**Parameters:**
- `initialWorkerCount`: Number of workers to start with

**Returns:**
- `*WorkerPool`: New worker pool instance

**Example:**
```go
pool := workerpool.NewWorkerPool(5)
defer pool.Stop()
```

#### `Submit(job Job) error`

Submits a job to the worker pool queue.

**Parameters:**
- `job`: Job to be processed

**Returns:**
- `error`: Error if pool is stopped

**Example:**
```go
err := pool.Submit(workerpool.Job{
    Task: func(ctx context.Context) error {
        // Do work
        return nil
    },
    Timeout: 5 * time.Second,
})
```

#### `ScaleTo(newCount int)`

Scales the worker pool to the specified number of workers.

**Parameters:**
- `newCount`: Target number of workers

**Example:**
```go
// Scale up
pool.ScaleTo(10)

// Scale down
pool.ScaleTo(2)
```

#### `Stop()`

Stops the worker pool gracefully, waiting for all running jobs to complete.

**Example:**
```go
pool.Stop()
```

### Job

```go
type Job struct {
    Task       func(ctx context.Context) error
    Retry      int
    Timeout    time.Duration
    RetryDelay time.Duration
    Priority   Priority
    CreatedAt  time.Time
}
```

### Priority

```go
type Priority int

const (
    Low    Priority = 0
    Normal Priority = 1
    High   Priority = 2
)
```

---

## Best Practices

### 1. Choose Appropriate Worker Count

```go
// Too few workers - jobs will queue up
pool := workerpool.NewWorkerPool(1)

// Good - based on CPU cores or I/O operations
pool := workerpool.NewWorkerPool(runtime.NumCPU())

// For I/O-bound tasks, can use more workers
pool := workerpool.NewWorkerPool(runtime.NumCPU() * 2)
```

### 2. Set Reasonable Timeouts

```go
// Good - appropriate timeout for task
workerpool.Job{
    Task:    quickTask,
    Timeout: 5 * time.Second,
}

// Good - longer timeout for slow operations
workerpool.Job{
    Task:    slowTask,
    Timeout: 30 * time.Second,
}
```

### 3. Use Retry for Transient Failures

```go
// Good - retry for network operations
workerpool.Job{
    Task:       fetchFromAPI,
    Retry:      3,
    RetryDelay: 1 * time.Second,
    Timeout:    10 * time.Second,
}

// Bad - don't retry for permanent errors
workerpool.Job{
    Task:  validateInput, // Will always fail if input is invalid
    Retry: 3,             // Wasteful
}
```

### 4. Use Priority Wisely

```go
// High priority for critical operations
pool.Submit(workerpool.Job{
    Task:     processPayment,
    Priority: workerpool.High,
    Timeout:  10 * time.Second,
})

// Normal priority for regular operations
pool.Submit(workerpool.Job{
    Task:     sendEmail,
    Priority: workerpool.Normal,
    Timeout:  5 * time.Second,
})

// Low priority for background tasks
pool.Submit(workerpool.Job{
    Task:     generateReport,
    Priority: workerpool.Low,
    Timeout:  30 * time.Second,
})
```

### 5. Always Stop the Pool

```go
// Good - using defer
pool := workerpool.NewWorkerPool(5)
defer pool.Stop()

// Good - explicit stop
pool := workerpool.NewWorkerPool(5)
// ... do work ...
pool.Stop()
```

### 6. Handle Context Cancellation

```go
workerpool.Job{
    Task: func(ctx context.Context) error {
        select {
        case <-time.After(5 * time.Second):
            // Work completed
            return nil
        case <-ctx.Done():
            // Context cancelled or timeout
            return ctx.Err()
        }
    },
    Timeout: 3 * time.Second,
}
```

### 7. Return Meaningful Errors

```go
workerpool.Job{
    Task: func(ctx context.Context) error {
        if err := doSomething(); err != nil {
            return fmt.Errorf("failed to do something: %w", err)
        }
        return nil
    },
}
```

### 8. Scale Based on Load

```go
// Monitor queue length and scale accordingly
if queueLength > threshold {
    pool.ScaleTo(currentWorkers * 2)
} else if queueLength < lowThreshold {
    pool.ScaleTo(currentWorkers / 2)
}
```

### 9. Avoid Blocking Operations

```go
// Bad - blocking without context
workerpool.Job{
    Task: func(ctx context.Context) error {
        time.Sleep(10 * time.Second) // Ignores context
        return nil
    },
}

// Good - respects context
workerpool.Job{
    Task: func(ctx context.Context) error {
        select {
        case <-time.After(10 * time.Second):
            return nil
        case <-ctx.Done():
            return ctx.Err()
        }
    },
}
```

### 10. Test Your Jobs

```go
func TestJob(t *testing.T) {
    pool := workerpool.NewWorkerPool(1)
    defer pool.Stop()
    
    done := make(chan bool)
    
    err := pool.Submit(workerpool.Job{
        Task: func(ctx context.Context) error {
            defer func() { done <- true }()
            // Test logic
            return nil
        },
        Timeout: 5 * time.Second,
    })
    
    if err != nil {
        t.Fatalf("Submit failed: %v", err)
    }
    
    select {
    case <-done:
        // Success
    case <-time.After(10 * time.Second):
        t.Fatal("Job timed out")
    }
}
```

---

## Performance

### Worker Count

- **CPU-bound tasks**: Use `runtime.NumCPU()` workers
- **I/O-bound tasks**: Use `runtime.NumCPU() * 2` or more
- **Mixed workload**: Start with `runtime.NumCPU()` and scale as needed

### Memory

- Each worker uses minimal memory (goroutine stack)
- Job queue grows with pending jobs
- Consider queue size limits for memory-constrained environments

### Concurrency

- Thread-safe with mutex protection
- Workers process jobs concurrently
- Priority sorting happens on job retrieval

---

## Troubleshooting

### Jobs Not Processing

**Problem:** Jobs submitted but not processing

**Solution:**
- Ensure pool is not stopped
- Check if workers are created (`NewWorkerPool` called)
- Verify timeout is not too short

### Jobs Timing Out

**Problem:** Jobs always timeout

**Solution:**
- Increase `Timeout` duration
- Check if task respects context cancellation
- Verify task doesn't block indefinitely

### Too Many Retries

**Problem:** Jobs retry unnecessarily

**Solution:**
- Reduce `Retry` count
- Return `nil` for non-retriable errors
- Check if error is transient or permanent

### Memory Usage Growing

**Problem:** Memory usage increases over time

**Solution:**
- Limit job submission rate
- Increase worker count
- Implement queue size limits

---

## Extending

The worker pool is designed to be simple and focused. However, you can extend it by wrapping it or creating custom job types.

### Custom Job Wrapper

```go
type CustomJob struct {
    ID       string
    Name     string
    Priority workerpool.Priority
    Task     func(ctx context.Context) error
}

func (cj CustomJob) ToWorkerPoolJob() workerpool.Job {
    return workerpool.Job{
        Task: func(ctx context.Context) error {
            log.Printf("Starting custom job: %s (ID: %s)", cj.Name, cj.ID)
            err := cj.Task(ctx)
            if err != nil {
                log.Printf("Custom job failed: %s", cj.Name)
            } else {
                log.Printf("Custom job completed: %s", cj.Name)
            }
            return err
        },
        Priority:   cj.Priority,
        Retry:      3,
        Timeout:    30 * time.Second,
        RetryDelay: 2 * time.Second,
    }
}

// Usage
pool := workerpool.NewWorkerPool(5)
defer pool.Stop()

customJob := CustomJob{
    ID:       "job-001",
    Name:     "Process Data",
    Priority: workerpool.High,
    Task: func(ctx context.Context) error {
        // Your task logic
        return nil
    },
}

pool.Submit(customJob.ToWorkerPoolJob())
```

### Job Queue Wrapper

```go
type JobQueue struct {
    pool    *workerpool.WorkerPool
    metrics struct {
        submitted int
        completed int
        failed    int
    }
    mu sync.Mutex
}

func NewJobQueue(workerCount int) *JobQueue {
    return &JobQueue{
        pool: workerpool.NewWorkerPool(workerCount),
    }
}

func (jq *JobQueue) Submit(task func(ctx context.Context) error, priority workerpool.Priority) error {
    jq.mu.Lock()
    jq.metrics.submitted++
    jq.mu.Unlock()
    
    return jq.pool.Submit(workerpool.Job{
        Task: func(ctx context.Context) error {
            err := task(ctx)
            
            jq.mu.Lock()
            if err != nil {
                jq.metrics.failed++
            } else {
                jq.metrics.completed++
            }
            jq.mu.Unlock()
            
            return err
        },
        Priority: priority,
        Timeout:  30 * time.Second,
    })
}

func (jq *JobQueue) GetMetrics() (submitted, completed, failed int) {
    jq.mu.Lock()
    defer jq.mu.Unlock()
    return jq.metrics.submitted, jq.metrics.completed, jq.metrics.failed
}

func (jq *JobQueue) Stop() {
    jq.pool.Stop()
}
```

### Example: Rate Limited Pool

```go
type RateLimitedPool struct {
    pool    *workerpool.WorkerPool
    limiter *rate.Limiter
}

func NewRateLimitedPool(workerCount int, rps int) *RateLimitedPool {
    return &RateLimitedPool{
        pool:    workerpool.NewWorkerPool(workerCount),
        limiter: rate.NewLimiter(rate.Limit(rps), rps),
    }
}

func (rlp *RateLimitedPool) Submit(job workerpool.Job) error {
    // Wait for rate limiter
    if err := rlp.limiter.Wait(context.Background()); err != nil {
        return err
    }
    
    return rlp.pool.Submit(job)
}

func (rlp *RateLimitedPool) Stop() {
    rlp.pool.Stop()
}
```

---

## Summary

The worker pool package provides **efficient concurrent task processing**:
- Fixed number of workers to control concurrency
- Priority queue for important tasks
- Automatic retry with configurable delay
- Timeout support for long-running tasks
- Panic recovery to keep workers alive
- Dynamic scaling for load management
- Graceful shutdown for clean termination

**Key Features:**
- Simple API with `Submit()` and `Stop()`
- Three priority levels (High, Normal, Low)
- FIFO within same priority
- Thread-safe for concurrent use
- Zero external dependencies

Now you can efficiently process concurrent tasks in your Go applications! 
