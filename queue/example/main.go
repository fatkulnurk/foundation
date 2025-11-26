package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatkulnurk/foundation/queue"
	"github.com/redis/go-redis/v9"
)

// EmailPayload represents the data for sending an email
type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// NotificationPayload represents the data for sending a notification
type NotificationPayload struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

func main() {
	// Setup Redis connection
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Choose which example to run
	if len(os.Args) > 1 && os.Args[1] == "worker" {
		runWorker(redisClient)
	} else {
		runProducer(redisClient)
	}
}

// runProducer demonstrates how to enqueue tasks
func runProducer(redisClient *redis.Client) {
	fmt.Println("=== Queue Producer Example ===\n")

	// Create queue instance
	q, err := queue.NewQueue(redisClient)
	if err != nil {
		log.Fatalf("Failed to create queue: %v", err)
	}

	ctx := context.Background()

	// Example 1: Simple email task
	fmt.Println("1. Enqueuing simple email task...")
	emailPayload := EmailPayload{
		To:      "user@example.com",
		Subject: "Welcome!",
		Body:    "Welcome to our service!",
	}

	result, err := q.Enqueue(ctx, "email:send", emailPayload)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Task enqueued: ID=%s\n\n", result.TaskID)
	}

	// Example 2: Email with retry and timeout
	fmt.Println("2. Enqueuing email with retry and timeout...")
	result, err = q.Enqueue(ctx, "email:send", emailPayload,
		queue.MaxRetry(3),
		queue.Timeout(30*time.Second),
		queue.QueueName("critical"),
	)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Task enqueued: ID=%s (critical queue, 3 retries, 30s timeout)\n\n", result.TaskID)
	}

	// Example 3: Scheduled notification
	fmt.Println("3. Enqueuing scheduled notification...")
	notifPayload := NotificationPayload{
		UserID:  "user123",
		Message: "Your order has been shipped!",
		Type:    "order",
	}

	result, err = q.Enqueue(ctx, "notification:send", notifPayload,
		queue.ProcessIn(5*time.Second),
		queue.QueueName("default"),
	)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Task scheduled: ID=%s (will process in 5 seconds)\n\n", result.TaskID)
	}

	// Example 4: Unique task (prevents duplicates)
	fmt.Println("4. Enqueuing unique task...")
	result, err = q.Enqueue(ctx, "report:generate", map[string]any{
		"report_id": "monthly-2024-11",
		"user_id":   "admin",
	},
		queue.Unique(1*time.Hour),
		queue.TaskID("report-monthly-2024-11"),
	)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Unique task enqueued: ID=%s (unique for 1 hour)\n\n", result.TaskID)
	}

	// Example 5: Task with deadline
	fmt.Println("5. Enqueuing task with deadline...")
	deadline := time.Now().Add(1 * time.Hour)
	result, err = q.Enqueue(ctx, "backup:database", map[string]any{
		"database": "production",
	},
		queue.Deadline(deadline),
		queue.QueueName("low"),
	)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Task enqueued: ID=%s (deadline: %s)\n\n", result.TaskID, deadline.Format("15:04:05"))
	}

	// Example 6: Get task info
	fmt.Println("6. Getting task info...")
	if result != nil {
		taskInfo, err := q.GetTaskInfo(ctx, result.TaskID)
		if err != nil {
			log.Printf("Error getting task info: %v\n", err)
		} else {
			fmt.Printf("✓ Task Info Retrieved:\n")
			fmt.Printf("  - ID: %s\n", taskInfo.ID)
			fmt.Printf("  - Type: %s\n", taskInfo.Type)
			fmt.Printf("  - State: %s\n", taskInfo.State)
			fmt.Printf("  - Queue: %s\n", taskInfo.Queue)
			fmt.Printf("  - Max Retry: %d\n", taskInfo.MaxRetry)
			fmt.Printf("  - Retried: %d\n\n", taskInfo.Retried)
		}
	}

	fmt.Println("✅ All tasks enqueued successfully!")
	fmt.Println("\nRun 'go run main.go worker' to start the worker and process these tasks.")
}

// runWorker demonstrates how to process tasks
func runWorker(redisClient *redis.Client) {
	fmt.Println("=== Queue Worker Example ===\n")

	// Create worker config
	cfg := &queue.Config{
		Concurrency: 5,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
		StrictPriority:  false,
		ShutdownTimeout: 10,
	}

	// Create worker instance
	w := queue.NewWorker(cfg, redisClient)

	// Register handlers
	fmt.Println("Registering task handlers...")

	// Handler for email:send with middleware
	w.RegisterWithMiddleware("email:send",
		func(ctx context.Context, payload []byte) error {
			// Get task ID from context using worker
			taskID, ok := w.GetTaskIDFromContext(ctx)
			if ok {
				fmt.Printf("📋 Processing task ID: %s\n", taskID)
			}

			var email EmailPayload
			if err := json.Unmarshal(payload, &email); err != nil {
				return fmt.Errorf("failed to unmarshal email payload: %w", err)
			}

			fmt.Printf("📧 Sending email to %s: %s\n", email.To, email.Subject)
			time.Sleep(1 * time.Second) // Simulate work
			fmt.Printf("✓ Email sent successfully to %s\n\n", email.To)
			return nil
		},
		queue.LoggingMiddleware("email:send"),
		queue.RecoveryMiddleware("email:send"),
	)

	// Handler for notification:send
	w.Register("notification:send", func(ctx context.Context, payload []byte) error {
		// Get task ID from context
		if taskID, ok := w.GetTaskIDFromContext(ctx); ok {
			fmt.Printf("📋 Task ID: %s\n", taskID)

			// Get full task info from worker
			if taskInfo, err := w.GetTaskInfo(ctx, taskID); err == nil {
				fmt.Printf("   State: %s, Retry: %d/%d\n",
					taskInfo.State, taskInfo.Retried, taskInfo.MaxRetry)
			}
		}

		var notif NotificationPayload
		if err := json.Unmarshal(payload, &notif); err != nil {
			return fmt.Errorf("failed to unmarshal notification payload: %w", err)
		}

		fmt.Printf("🔔 Sending notification to user %s: %s\n", notif.UserID, notif.Message)
		time.Sleep(500 * time.Millisecond) // Simulate work
		fmt.Printf("✓ Notification sent successfully\n\n")
		return nil
	})

	// Handler for report:generate
	w.Register("report:generate", func(ctx context.Context, payload []byte) error {
		var data map[string]any
		if err := json.Unmarshal(payload, &data); err != nil {
			return fmt.Errorf("failed to unmarshal report payload: %w", err)
		}

		fmt.Printf("📊 Generating report: %v\n", data["report_id"])
		time.Sleep(2 * time.Second) // Simulate work
		fmt.Printf("✓ Report generated successfully\n\n")
		return nil
	})

	// Handler for backup:database
	w.Register("backup:database", func(ctx context.Context, payload []byte) error {
		var data map[string]any
		if err := json.Unmarshal(payload, &data); err != nil {
			return fmt.Errorf("failed to unmarshal backup payload: %w", err)
		}

		fmt.Printf("💾 Backing up database: %v\n", data["database"])
		time.Sleep(3 * time.Second) // Simulate work
		fmt.Printf("✓ Database backup completed\n\n")
		return nil
	})

	fmt.Println("✓ All handlers registered\n")

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start worker in goroutine
	go func() {
		fmt.Println("🚀 Starting worker...")
		fmt.Println("Press Ctrl+C to stop\n")
		if err := w.Start(); err != nil {
			log.Fatalf("Worker error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	fmt.Println("\n\n⏳ Shutting down gracefully...")
	w.Stop()
	fmt.Println("✅ Worker stopped")
}
