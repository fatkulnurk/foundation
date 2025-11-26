package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fatkulnurk/foundation/workerpool"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "basic":
			runBasicExample()
		case "priority":
			runPriorityExample()
		case "retry":
			runRetryExample()
		case "timeout":
			runTimeoutExample()
		case "scaling":
			runScalingExample()
		case "panic":
			runPanicExample()
		default:
			fmt.Println("Usage: go run main.go [basic|priority|retry|timeout|scaling|panic]")
		}
	} else {
		runBasicExample()
	}
}

// runBasicExample demonstrates basic worker pool usage
func runBasicExample() {
	fmt.Println("=== Basic Worker Pool Example ===")
	fmt.Println()

	// Create pool with 3 workers
	pool := workerpool.NewWorkerPool(3)
	defer pool.Stop()

	fmt.Println("Submitting 10 jobs to pool with 3 workers...")
	fmt.Println()

	// Submit 10 jobs
	for i := 1; i <= 10; i++ {
		jobID := i
		err := pool.Submit(workerpool.Job{
			Task: func(ctx context.Context) error {
				fmt.Printf("📋 Job %d: Starting\n", jobID)
				time.Sleep(1 * time.Second)
				fmt.Printf("✅ Job %d: Completed\n", jobID)
				return nil
			},
			Timeout:  5 * time.Second,
			Priority: workerpool.Normal,
		})

		if err != nil {
			fmt.Printf("❌ Failed to submit job %d: %v\n", jobID, err)
		}
	}

	// Wait for all jobs to complete
	fmt.Println()
	fmt.Println("Waiting for jobs to complete...")
	time.Sleep(5 * time.Second)

	fmt.Println()
	fmt.Println("✅ All jobs completed!")
}

// runPriorityExample demonstrates priority queue
func runPriorityExample() {
	fmt.Println("=== Priority Queue Example ===")
	fmt.Println()

	// Create pool with 2 workers (so we can see priority in action)
	pool := workerpool.NewWorkerPool(2)
	defer pool.Stop()

	fmt.Println("Submitting jobs with different priorities...")
	fmt.Println("Expected order: High-1, High-2, Normal-1, Normal-2, Low-1, Low-2")
	fmt.Println()

	// Submit jobs with different priorities
	jobs := []struct {
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

	for _, j := range jobs {
		name := j.name
		pool.Submit(workerpool.Job{
			Task: func(ctx context.Context) error {
				fmt.Printf("🔄 Processing: %s\n", name)
				time.Sleep(500 * time.Millisecond)
				fmt.Printf("✅ Completed: %s\n", name)
				return nil
			},
			Priority: j.priority,
			Timeout:  5 * time.Second,
		})
	}

	// Wait for all jobs to complete
	time.Sleep(5 * time.Second)

	fmt.Println()
	fmt.Println("✅ Priority example completed!")
}

// runRetryExample demonstrates retry mechanism
func runRetryExample() {
	fmt.Println("=== Retry Mechanism Example ===")
	fmt.Println()

	pool := workerpool.NewWorkerPool(2)
	defer pool.Stop()

	// Job 1: Will succeed on first try
	fmt.Println("Job 1: Will succeed immediately")
	pool.Submit(workerpool.Job{
		Task: func(ctx context.Context) error {
			fmt.Println("  📋 Attempting job 1...")
			time.Sleep(500 * time.Millisecond)
			fmt.Println("  ✅ Job 1 succeeded!")
			return nil
		},
		Retry:      3,
		Timeout:    5 * time.Second,
		RetryDelay: 1 * time.Second,
		Priority:   workerpool.Normal,
	})

	time.Sleep(2 * time.Second)

	// Job 2: Will fail a few times then succeed
	fmt.Println()
	fmt.Println("Job 2: Will fail 2 times, then succeed")
	attempt := 0
	pool.Submit(workerpool.Job{
		Task: func(ctx context.Context) error {
			attempt++
			fmt.Printf("  📋 Attempting job 2 (attempt %d)...\n", attempt)
			time.Sleep(500 * time.Millisecond)

			if attempt < 3 {
				fmt.Printf("  ❌ Job 2 failed on attempt %d\n", attempt)
				return errors.New("simulated failure")
			}

			fmt.Println("  ✅ Job 2 succeeded!")
			return nil
		},
		Retry:      3,
		Timeout:    5 * time.Second,
		RetryDelay: 1 * time.Second,
		Priority:   workerpool.Normal,
	})

	time.Sleep(8 * time.Second)

	// Job 3: Will always fail
	fmt.Println()
	fmt.Println("Job 3: Will fail all attempts")
	pool.Submit(workerpool.Job{
		Task: func(ctx context.Context) error {
			fmt.Println("  📋 Attempting job 3...")
			time.Sleep(500 * time.Millisecond)
			fmt.Println("  ❌ Job 3 failed")
			return errors.New("permanent failure")
		},
		Retry:      2,
		Timeout:    5 * time.Second,
		RetryDelay: 1 * time.Second,
		Priority:   workerpool.Normal,
	})

	time.Sleep(8 * time.Second)

	fmt.Println()
	fmt.Println("✅ Retry example completed!")
}

// runTimeoutExample demonstrates timeout handling
func runTimeoutExample() {
	fmt.Println("=== Timeout Example ===")
	fmt.Println()

	pool := workerpool.NewWorkerPool(2)
	defer pool.Stop()

	// Job 1: Completes before timeout
	fmt.Println("Job 1: Will complete before timeout (2s task, 5s timeout)")
	pool.Submit(workerpool.Job{
		Task: func(ctx context.Context) error {
			fmt.Println("  📋 Starting job 1...")
			select {
			case <-time.After(2 * time.Second):
				fmt.Println("  ✅ Job 1 completed successfully")
				return nil
			case <-ctx.Done():
				fmt.Println("  ⏱️  Job 1 timed out")
				return ctx.Err()
			}
		},
		Timeout:  5 * time.Second,
		Priority: workerpool.Normal,
	})

	time.Sleep(3 * time.Second)

	// Job 2: Will timeout
	fmt.Println()
	fmt.Println("Job 2: Will timeout (10s task, 2s timeout)")
	pool.Submit(workerpool.Job{
		Task: func(ctx context.Context) error {
			fmt.Println("  📋 Starting job 2...")
			select {
			case <-time.After(10 * time.Second):
				fmt.Println("  ✅ Job 2 completed successfully")
				return nil
			case <-ctx.Done():
				fmt.Println("  ⏱️  Job 2 timed out (as expected)")
				return ctx.Err()
			}
		},
		Timeout:  2 * time.Second,
		Retry:    1,
		Priority: workerpool.Normal,
	})

	time.Sleep(6 * time.Second)

	fmt.Println()
	fmt.Println("✅ Timeout example completed!")
}

// runScalingExample demonstrates dynamic scaling
func runScalingExample() {
	fmt.Println("=== Dynamic Scaling Example ===")
	fmt.Println()

	// Start with 2 workers
	fmt.Println("Starting with 2 workers...")
	pool := workerpool.NewWorkerPool(2)
	defer pool.Stop()

	// Submit 5 jobs
	fmt.Println("Submitting 5 jobs...")
	for i := 1; i <= 5; i++ {
		jobID := i
		pool.Submit(workerpool.Job{
			Task: func(ctx context.Context) error {
				fmt.Printf("  📋 Job %d processing (2 workers)\n", jobID)
				time.Sleep(2 * time.Second)
				fmt.Printf("  ✅ Job %d completed\n", jobID)
				return nil
			},
			Timeout:  5 * time.Second,
			Priority: workerpool.Normal,
		})
	}

	time.Sleep(1 * time.Second)

	// Scale up to 5 workers
	fmt.Println()
	fmt.Println("⬆️  Scaling up to 5 workers...")
	pool.ScaleTo(5)

	// Submit more jobs
	fmt.Println("Submitting 5 more jobs...")
	for i := 6; i <= 10; i++ {
		jobID := i
		pool.Submit(workerpool.Job{
			Task: func(ctx context.Context) error {
				fmt.Printf("  📋 Job %d processing (5 workers)\n", jobID)
				time.Sleep(2 * time.Second)
				fmt.Printf("  ✅ Job %d completed\n", jobID)
				return nil
			},
			Timeout:  5 * time.Second,
			Priority: workerpool.Normal,
		})
	}

	time.Sleep(5 * time.Second)

	// Scale down to 2 workers
	fmt.Println()
	fmt.Println("⬇️  Scaling down to 2 workers...")
	pool.ScaleTo(2)

	time.Sleep(2 * time.Second)

	fmt.Println()
	fmt.Println("✅ Scaling example completed!")
}

// runPanicExample demonstrates panic recovery
func runPanicExample() {
	fmt.Println("=== Panic Recovery Example ===")
	fmt.Println()

	pool := workerpool.NewWorkerPool(2)
	defer pool.Stop()

	// Job 1: Will panic
	fmt.Println("Job 1: Will panic (worker should recover)")
	pool.Submit(workerpool.Job{
		Task: func(ctx context.Context) error {
			fmt.Println("  📋 Job 1 starting...")
			time.Sleep(500 * time.Millisecond)
			fmt.Println("  💥 Job 1 about to panic!")
			panic("something went terribly wrong!")
		},
		Retry:      2,
		Timeout:    5 * time.Second,
		RetryDelay: 1 * time.Second,
		Priority:   workerpool.Normal,
	})

	time.Sleep(5 * time.Second)

	// Job 2: Normal job (to prove worker is still alive)
	fmt.Println()
	fmt.Println("Job 2: Normal job (to verify worker recovered)")
	pool.Submit(workerpool.Job{
		Task: func(ctx context.Context) error {
			fmt.Println("  📋 Job 2 starting...")
			time.Sleep(500 * time.Millisecond)
			fmt.Println("  ✅ Job 2 completed successfully!")
			fmt.Println("  ✅ Worker is still alive after panic!")
			return nil
		},
		Timeout:  5 * time.Second,
		Priority: workerpool.Normal,
	})

	time.Sleep(2 * time.Second)

	fmt.Println()
	fmt.Println("✅ Panic recovery example completed!")
}
