# Worker Pool Examples

This directory contains working examples demonstrating how to use the worker pool package.

## Prerequisites

1. **Go 1.25 or higher**
2. **No external dependencies** - Uses standard library only

## Running the Examples

### 1. Basic Worker Pool Example

Demonstrates basic worker pool usage with multiple jobs.

```bash
cd pkg/workerpool/example
go run main.go basic
```

**What it does:**
- Creates a pool with 3 workers
- Submits 10 jobs
- Shows how jobs are processed concurrently
- Demonstrates graceful shutdown

**Output:**
```
=== Basic Worker Pool Example ===

Submitting 10 jobs to pool with 3 workers...

📋 Job 1: Starting
📋 Job 2: Starting
📋 Job 3: Starting
✅ Job 1: Completed
📋 Job 4: Starting
✅ Job 2: Completed
📋 Job 5: Starting
...
✅ All jobs completed!
```

### 2. Priority Queue Example

Demonstrates how priority levels affect job processing order.

```bash
go run main.go priority
```

**What it does:**
- Creates a pool with 2 workers
- Submits jobs with High, Normal, and Low priorities
- Shows that High priority jobs are processed first
- Demonstrates FIFO within same priority

**Output:**
```
=== Priority Queue Example ===

Submitting jobs with different priorities...
Expected order: High-1, High-2, Normal-1, Normal-2, Low-1, Low-2

🔄 Processing: High-1
🔄 Processing: High-2
✅ Completed: High-1
🔄 Processing: Normal-1
✅ Completed: High-2
🔄 Processing: Normal-2
...
```

### 3. Retry Mechanism Example

Demonstrates automatic retry with configurable delay.

```bash
go run main.go retry
```

**What it does:**
- Job 1: Succeeds immediately
- Job 2: Fails 2 times, then succeeds
- Job 3: Fails all attempts
- Shows retry delay in action

**Output:**
```
=== Retry Mechanism Example ===

Job 1: Will succeed immediately
  📋 Attempting job 1...
  ✅ Job 1 succeeded!
[Worker 0] ✅ Job success (Attempt 1)

Job 2: Will fail 2 times, then succeed
  📋 Attempting job 2 (attempt 1)...
  ❌ Job 2 failed on attempt 1
[Worker 0] ❌ Job failed (Attempt 1): simulated failure
  📋 Attempting job 2 (attempt 2)...
  ❌ Job 2 failed on attempt 2
[Worker 0] ❌ Job failed (Attempt 2): simulated failure
  📋 Attempting job 2 (attempt 3)...
  ✅ Job 2 succeeded!
[Worker 0] ✅ Job success (Attempt 3)
...
```

### 4. Timeout Example

Demonstrates timeout handling for long-running tasks.

```bash
go run main.go timeout
```

**What it does:**
- Job 1: Completes before timeout (2s task, 5s timeout)
- Job 2: Times out (10s task, 2s timeout)
- Shows context cancellation

**Output:**
```
=== Timeout Example ===

Job 1: Will complete before timeout (2s task, 5s timeout)
  📋 Starting job 1...
  ✅ Job 1 completed successfully
[Worker 0] ✅ Job success (Attempt 1)

Job 2: Will timeout (10s task, 2s timeout)
  📋 Starting job 2...
  ⏱️  Job 2 timed out (as expected)
[Worker 0] ❌ Job failed (Attempt 1): context deadline exceeded
...
```

### 5. Dynamic Scaling Example

Demonstrates scaling workers up and down.

```bash
go run main.go scaling
```

**What it does:**
- Starts with 2 workers
- Submits 5 jobs
- Scales up to 5 workers
- Submits 5 more jobs (processed faster)
- Scales down to 2 workers

**Output:**
```
=== Dynamic Scaling Example ===

Starting with 2 workers...
Submitting 5 jobs...
  📋 Job 1 processing (2 workers)
  📋 Job 2 processing (2 workers)
  ✅ Job 1 completed
  📋 Job 3 processing (2 workers)
...

⬆️  Scaling up to 5 workers...
⬆️ Scaled up to 5 workers
Submitting 5 more jobs...
  📋 Job 6 processing (5 workers)
  📋 Job 7 processing (5 workers)
  📋 Job 8 processing (5 workers)
  📋 Job 9 processing (5 workers)
  📋 Job 10 processing (5 workers)
...
```

### 6. Panic Recovery Example

Demonstrates that workers recover from panics.

```bash
go run main.go panic
```

**What it does:**
- Job 1: Panics (worker recovers)
- Job 2: Normal job (proves worker is still alive)
- Shows stack trace on panic

**Output:**
```
=== Panic Recovery Example ===

Job 1: Will panic (worker should recover)
  📋 Job 1 starting...
  💥 Job 1 about to panic!
[Worker 0] PANIC: something went terribly wrong!
[stack trace...]
[Worker 0] ❌ Job failed (Attempt 1): panic: something went terribly wrong!
...

Job 2: Normal job (to verify worker recovered)
  📋 Job 2 starting...
  ✅ Job 2 completed successfully!
  ✅ Worker is still alive after panic!
[Worker 1] ✅ Job success (Attempt 1)
```

## Code Examples

### Example 1: Basic Usage

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
    
    // Submit jobs
    for i := 1; i <= 10; i++ {
        jobID := i
        pool.Submit(workerpool.Job{
            Task: func(ctx context.Context) error {
                fmt.Printf("Processing job %d\n", jobID)
                time.Sleep(1 * time.Second)
                return nil
            },
            Timeout: 5 * time.Second,
        })
    }
    
    // Wait for completion
    time.Sleep(5 * time.Second)
}
```

### Example 2: With Priority

```go
// High priority - processed first
pool.Submit(workerpool.Job{
    Task: func(ctx context.Context) error {
        fmt.Println("Urgent task")
        return nil
    },
    Priority: workerpool.High,
    Timeout:  5 * time.Second,
})

// Normal priority
pool.Submit(workerpool.Job{
    Task: func(ctx context.Context) error {
        fmt.Println("Regular task")
        return nil
    },
    Priority: workerpool.Normal,
    Timeout:  5 * time.Second,
})

// Low priority - processed last
pool.Submit(workerpool.Job{
    Task: func(ctx context.Context) error {
        fmt.Println("Background task")
        return nil
    },
    Priority: workerpool.Low,
    Timeout:  5 * time.Second,
})
```

### Example 3: With Retry

```go
pool.Submit(workerpool.Job{
    Task: func(ctx context.Context) error {
        // Task that might fail
        if err := doSomething(); err != nil {
            return err
        }
        return nil
    },
    Retry:      3,                   // Retry up to 3 times
    Timeout:    10 * time.Second,    // 10 second timeout per attempt
    RetryDelay: 2 * time.Second,     // Wait 2 seconds between retries
    Priority:   workerpool.Normal,
})
```

### Example 4: With Timeout

```go
pool.Submit(workerpool.Job{
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
    Timeout: 3 * time.Second, // Will timeout after 3 seconds
})
```

### Example 5: Dynamic Scaling

```go
pool := workerpool.NewWorkerPool(2)
defer pool.Stop()

// Start with 2 workers
fmt.Println("Starting with 2 workers")

// Scale up when load increases
pool.ScaleTo(10)
fmt.Println("Scaled up to 10 workers")

// Scale down when load decreases
pool.ScaleTo(2)
fmt.Println("Scaled down to 2 workers")
```

## Common Patterns

### Pattern 1: Batch Processing

```go
pool := workerpool.NewWorkerPool(5)
defer pool.Stop()

items := []string{"item1", "item2", "item3", "item4", "item5"}

for _, item := range items {
    i := item
    pool.Submit(workerpool.Job{
        Task: func(ctx context.Context) error {
            return processItem(i)
        },
        Retry:   2,
        Timeout: 10 * time.Second,
    })
}
```

### Pattern 2: HTTP Requests

```go
pool := workerpool.NewWorkerPool(10)
defer pool.Stop()

urls := []string{"url1", "url2", "url3"}

for _, url := range urls {
    targetURL := url
    pool.Submit(workerpool.Job{
        Task: func(ctx context.Context) error {
            req, _ := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
            resp, err := http.DefaultClient.Do(req)
            if err != nil {
                return err
            }
            defer resp.Body.Close()
            // Process response
            return nil
        },
        Retry:      2,
        Timeout:    30 * time.Second,
        RetryDelay: 1 * time.Second,
    })
}
```

### Pattern 3: Database Operations

```go
pool := workerpool.NewWorkerPool(5)
defer pool.Stop()

users := []User{{ID: 1}, {ID: 2}, {ID: 3}}

for _, user := range users {
    u := user
    pool.Submit(workerpool.Job{
        Task: func(ctx context.Context) error {
            return db.SaveUser(ctx, u)
        },
        Retry:      3,
        Timeout:    5 * time.Second,
        RetryDelay: 500 * time.Millisecond,
    })
}
```

### Pattern 4: Priority-Based Processing

```go
pool := workerpool.NewWorkerPool(5)
defer pool.Stop()

// Critical operations
pool.Submit(workerpool.Job{
    Task:     processCritical,
    Priority: workerpool.High,
    Timeout:  10 * time.Second,
})

// Regular operations
pool.Submit(workerpool.Job{
    Task:     processRegular,
    Priority: workerpool.Normal,
    Timeout:  10 * time.Second,
})

// Background operations
pool.Submit(workerpool.Job{
    Task:     processBackground,
    Priority: workerpool.Low,
    Timeout:  30 * time.Second,
})
```

## Modifying the Examples

### Add Your Own Example

```go
func runMyExample() {
    fmt.Println("=== My Custom Example ===")
    
    pool := workerpool.NewWorkerPool(3)
    defer pool.Stop()
    
    // Your code here
    pool.Submit(workerpool.Job{
        Task: func(ctx context.Context) error {
            // Your task logic
            return nil
        },
        Timeout: 5 * time.Second,
    })
    
    time.Sleep(2 * time.Second)
    fmt.Println("✅ My example completed!")
}
```

Then add to main():
```go
case "myexample":
    runMyExample()
```

## Troubleshooting

### Jobs Not Processing

**Problem:** Jobs submitted but not processing

**Solution:**
- Ensure `defer pool.Stop()` is called
- Check if pool was created with `NewWorkerPool()`
- Verify timeout is not too short

### Jobs Always Failing

**Problem:** Jobs fail even with retry

**Solution:**
- Check task logic for errors
- Increase timeout duration
- Verify retry count is appropriate
- Check if error is retriable

### Workers Not Scaling

**Problem:** `ScaleTo()` doesn't seem to work

**Solution:**
- Scaling up creates new workers immediately
- Scaling down happens after current jobs complete
- Check logs for scaling messages

## Performance Tips

1. **Choose appropriate worker count**
   - CPU-bound: `runtime.NumCPU()`
   - I/O-bound: `runtime.NumCPU() * 2`

2. **Set reasonable timeouts**
   - Too short: Jobs timeout unnecessarily
   - Too long: Resources held too long

3. **Use retry wisely**
   - Only for transient failures
   - Set appropriate retry delay

4. **Monitor and scale**
   - Scale up when queue grows
   - Scale down when queue shrinks

## Next Steps

After running these examples:

1. **Integrate into your application**
   - Copy patterns from examples
   - Adapt to your use case

2. **Add monitoring**
   - Track job success/failure rates
   - Monitor queue length
   - Measure processing time

3. **Implement auto-scaling**
   - Scale based on queue length
   - Scale based on system load

## Learn More

See the main [README.md](../README.md) for:
- Complete API documentation
- All configuration options
- Best practices
- Advanced usage

---

Happy concurrent processing! 🚀
