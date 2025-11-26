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

// EmailPayload is the data for sending an email
// Like an envelope containing: recipient, subject, and message body
type EmailPayload struct {
	To      string `json:"to"`      // Email recipient
	Subject string `json:"subject"` // Email subject
	Body    string `json:"body"`    // Email body
}

// NotificationPayload is the data for sending a notification
// Like a message to be sent to the user
type NotificationPayload struct {
	UserID  string `json:"user_id"` // User ID who will receive the notification
	Message string `json:"message"` // Message content
	Type    string `json:"type"`    // Notification type (order, payment, etc.)
}

func main() {
	// STEP 1: Connect to Redis
	// Redis is where we store the task queue (like a database for queues)
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis server address
	})

	ctx := context.Background()
	// Check if Redis is accessible
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// STEP 2: Choose which mode to run
	// Mode 1: "worker" = Process tasks from the queue
	// Mode 2: default = Add tasks to the queue
	if len(os.Args) > 1 && os.Args[1] == "worker" {
		runWorker(redisClient) // Run as worker
	} else {
		runProducer(redisClient) // Run as producer
	}
}

// runProducer demonstrates how to enqueue tasks
// Analogy: Like putting dirty clothes into a laundry basket
func runProducer(redisClient *redis.Client) {
	fmt.Println("=== Example: Enqueuing Tasks ===")
	fmt.Println()

	// Create Queue instance
	q, err := queue.NewQueue(redisClient)
	if err != nil {
		log.Fatalf("Failed to create queue: %v", err)
	}

	ctx := context.Background()

	// EXAMPLE 1: Simple task - Send email
	// Like putting one piece of clothing into the laundry basket
	fmt.Println("1. Enqueuing simple email task...")
	emailPayload := EmailPayload{
		To:      "user@example.com",
		Subject: "Welcome!",
		Body:    "Welcome to our service!",
	}

	// Enqueue task with name "email:send"
	result, err := q.Enqueue(ctx, "email:send", emailPayload)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Task enqueued successfully: ID=%s\n", result.TaskID)
		fmt.Println()
	}

	// EXAMPLE 2: Task with additional options
	// Like washing with special mode: retry 3x if failed, 30s timeout
	fmt.Println("2. Enqueuing email with retry and timeout...")
	result, err = q.Enqueue(ctx, "email:send", emailPayload,
		queue.MaxRetry(3),             // If failed, retry up to 3 times
		queue.Timeout(30*time.Second), // Maximum 30 seconds to send email
		queue.QueueName("critical"),   // Put in high priority queue
	)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Task enqueued: ID=%s (critical queue, 3 retries, 30s timeout)\n", result.TaskID)
		fmt.Println()
	}

	// EXAMPLE 3: Scheduled task
	// Like setting a timer on the washing machine to start in 5 seconds
	fmt.Println("3. Enqueuing scheduled notification...")
	notifPayload := NotificationPayload{
		UserID:  "user123",
		Message: "Your order has been shipped!",
		Type:    "order",
	}

	result, err = q.Enqueue(ctx, "notification:send", notifPayload,
		queue.ProcessIn(5*time.Second), // Process in 5 seconds from now
		queue.QueueName("default"),     // Put in default queue
	)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Task scheduled: ID=%s (will be processed in 5 seconds)\n", result.TaskID)
		fmt.Println()
	}

	// EXAMPLE 4: Unique task (prevents duplicates)
	// Like: if the same clothes are already in the washing machine, don't add again
	fmt.Println("4. Enqueuing unique task (anti-duplicate)...")
	result, err = q.Enqueue(ctx, "report:generate", map[string]any{
		"report_id": "monthly-2024-11",
		"user_id":   "admin",
	},
		queue.Unique(1*time.Hour),              // This task is unique for 1 hour
		queue.TaskID("report-monthly-2024-11"), // Custom ID for this task
	)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Unique task enqueued: ID=%s (unique for 1 hour)\n", result.TaskID)
		fmt.Println()
	}

	// EXAMPLE 5: Task with deadline
	// Like: laundry must be done before a certain time
	fmt.Println("5. Enqueuing task with deadline...")
	deadline := time.Now().Add(1 * time.Hour)
	result, err = q.Enqueue(ctx, "backup:database", map[string]any{
		"database": "production",
	},
		queue.Deadline(deadline), // Must complete before 1 hour from now
		queue.QueueName("low"),   // Low priority
	)
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Task enqueued: ID=%s (deadline: %s)\n", result.TaskID, deadline.Format("15:04:05"))
		fmt.Println()
	}

	// EXAMPLE 6: Get task information
	// Like checking laundry status: being washed, finished, etc.
	fmt.Println("6. Getting task information...")
	if result != nil {
		taskInfo, err := q.GetTaskInfo(ctx, result.TaskID)
		if err != nil {
			log.Printf("Error getting task info: %v\n", err)
		} else {
			fmt.Printf("✓ Task Info Retrieved Successfully:\n")
			fmt.Printf("  - ID: %s\n", taskInfo.ID)
			fmt.Printf("  - Type: %s\n", taskInfo.Type)
			fmt.Printf("  - State: %s\n", taskInfo.State)
			fmt.Printf("  - Queue: %s\n", taskInfo.Queue)
			fmt.Printf("  - Max Retry: %d\n", taskInfo.MaxRetry)
			fmt.Printf("  - Retried: %d\n", taskInfo.Retried)
			fmt.Println()
		}
	}

	fmt.Println("✅ All tasks enqueued successfully!")
	fmt.Println()
	fmt.Println("Run 'go run main.go worker' to process these tasks.")
}

// runWorker demonstrates how to process tasks from the queue
// Analogy: Like a washing machine that takes clothes from the basket and washes them
func runWorker(redisClient *redis.Client) {
	fmt.Println("=== Example: Worker Processing Tasks ===")
	fmt.Println()

	// Create worker configuration
	cfg := &queue.Config{
		Concurrency: 5, // Can process 5 tasks simultaneously (parallel)
		Queues: map[string]int{
			"critical": 6, // High priority queue (processed more frequently)
			"default":  3, // Normal queue
			"low":      1, // Low priority queue (processed less frequently)
		},
		StrictPriority:  false, // false = all queues are still processed
		ShutdownTimeout: 10,    // Wait 10 seconds before forcing shutdown
	}

	// Create Worker instance
	w := queue.NewWorker(cfg, redisClient)

	// Register handlers (how to process each type of task)
	fmt.Println("Registering handlers for each task type...")

	// HANDLER 1: email:send with middleware
	// Middleware = additional features like logging and error handling
	w.RegisterWithMiddleware("email:send",
		func(ctx context.Context, payload []byte) error {
			// Get task ID from context (for tracking)
			taskID, ok := w.GetTaskID(ctx)
			if ok {
				fmt.Printf("📋 Processing task ID: %s\n", taskID)
			}

			// Parse email data from payload
			var email EmailPayload
			if err := json.Unmarshal(payload, &email); err != nil {
				return fmt.Errorf("failed to parse email data: %w", err)
			}

			// Send email (simulated with sleep)
			fmt.Printf("📧 Sending email to %s: %s\n", email.To, email.Subject)
			time.Sleep(1 * time.Second) // Simulate email sending process
			fmt.Printf("✓ Email sent successfully to %s\n", email.To)
			fmt.Println()
			return nil
		},
		queue.LoggingMiddleware("email:send"),  // Middleware for logging
		queue.RecoveryMiddleware("email:send"), // Middleware to catch panics
	)

	// HANDLER 2: notification:send
	// Example handler that also retrieves full task info
	w.Register("notification:send", func(ctx context.Context, payload []byte) error {
		// Get task ID and full info
		if taskID, ok := w.GetTaskID(ctx); ok {
			fmt.Printf("📋 Task ID: %s\n", taskID)

			// Get full task information (status, retry count, etc.)
			if taskInfo, err := w.GetTaskInfo(ctx, taskID); err == nil {
				fmt.Printf("   State: %s, Retry: %d/%d\n",
					taskInfo.State, taskInfo.Retried, taskInfo.MaxRetry)
			}
		}

		// Parse notification data
		var notif NotificationPayload
		if err := json.Unmarshal(payload, &notif); err != nil {
			return fmt.Errorf("failed to parse notification data: %w", err)
		}

		// Send notification (simulated)
		fmt.Printf("🔔 Sending notification to user %s: %s\n", notif.UserID, notif.Message)
		time.Sleep(500 * time.Millisecond) // Simulate notification sending process
		fmt.Printf("✓ Notification sent successfully\n")
		fmt.Println()
		return nil
	})

	// HANDLER 3: report:generate
	// Handler for generating reports
	w.Register("report:generate", func(ctx context.Context, payload []byte) error {
		// Parse report data
		var data map[string]any
		if err := json.Unmarshal(payload, &data); err != nil {
			return fmt.Errorf("failed to parse report data: %w", err)
		}

		// Generate report (simulated longer process)
		fmt.Printf("📊 Generating report: %v\n", data["report_id"])
		time.Sleep(2 * time.Second) // Simulate report generation process
		fmt.Printf("✓ Report generated successfully\n")
		fmt.Println()
		return nil
	})

	// HANDLER 4: backup:database
	// Handler for database backup
	w.Register("backup:database", func(ctx context.Context, payload []byte) error {
		// Parse backup data
		var data map[string]any
		if err := json.Unmarshal(payload, &data); err != nil {
			return fmt.Errorf("failed to parse backup data: %w", err)
		}

		// Backup database (simulated long process)
		fmt.Printf("💾 Backing up database: %v\n", data["database"])
		time.Sleep(3 * time.Second) // Simulate backup process
		fmt.Printf("✓ Database backup completed\n")
		fmt.Println()
		return nil
	})

	fmt.Println("✓ All handlers registered successfully")
	fmt.Println()

	// Setup graceful shutdown
	// If Ctrl+C is pressed, wait for running tasks to finish first
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run worker in a separate goroutine
	go func() {
		fmt.Println("🚀 Starting worker...")
		fmt.Println("Press Ctrl+C to stop")
		fmt.Println()
		if err := w.Start(); err != nil {
			log.Fatalf("Worker error: %v", err)
		}
	}()

	// Wait for shutdown signal (Ctrl+C)
	<-sigChan
	fmt.Println()
	fmt.Println()
	fmt.Println("⏳ Shutting down worker gracefully...")
	w.Stop() // Stop worker, wait for running tasks to finish
	fmt.Println("✅ Worker stopped successfully")
}
