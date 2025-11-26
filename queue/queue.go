package queue

import (
	"context"
	"time"
)

// Queue defines the interface for queueing tasks
type Queue interface {
	// Enqueue adds a task to the queue
	Enqueue(ctx context.Context, taskName string, payload any, opts ...Option) (*OutputEnqueue, error)

	// GetTaskInfo retrieves information about a task by its ID
	GetTaskInfo(ctx context.Context, taskID string) (*TaskInfo, error)

	// Close closes the queue client connection
	Close() error
}

// Worker defines the interface for processing tasks
type Worker interface {
	// Start starts the worker and begins processing tasks
	// This is a blocking call that runs until Stop is called or an error occurs
	Start() error

	// Stop stops the worker gracefully
	// It waits for all in-progress tasks to complete before shutting down
	Stop()

	// Register registers a handler for a specific task type
	// The handler will be called when a task of this type is dequeued
	Register(taskType string, handler Handler) error

	// RegisterWithMiddleware registers a handler with middleware functions
	// Middleware will be executed in the order they are provided
	RegisterWithMiddleware(taskType string, handler Handler, middleware ...MiddlewareFunc) error

	// GetTaskID retrieves the task ID from the context
	// This is useful inside handler functions to get the current task ID
	// Returns the task ID and a boolean indicating if it was found
	GetTaskID(ctx context.Context) (string, bool)

	// GetTaskInfo retrieves information about a task by its ID
	// This allows workers to inspect task details during processing
	GetTaskInfo(ctx context.Context, taskID string) (*TaskInfo, error)
}

// Handler is a function that processes a task
// It receives context and the task payload as []byte
// It should return an error if the task processing fails
type Handler func(ctx context.Context, payload []byte) error

// MiddlewareFunc is a function that wraps a Handler
// It can be used for logging, metrics, error handling, etc.
type MiddlewareFunc func(Handler) Handler

// TaskState represents the state of a task
type TaskState string

const (
	TaskStatePending   TaskState = "pending"
	TaskStateActive    TaskState = "active"
	TaskStateScheduled TaskState = "scheduled"
	TaskStateRetry     TaskState = "retry"
	TaskStateArchived  TaskState = "archived"
	TaskStateCompleted TaskState = "completed"
)

// TaskInfo contains information about a task
type TaskInfo struct {
	ID            string
	Type          string
	Payload       []byte
	State         TaskState
	Queue         string
	MaxRetry      int
	Retried       int
	LastError     string
	CompletedAt   *time.Time
	NextProcessAt *time.Time
}

type OutputEnqueue struct {
	TaskID  string
	Payload []byte
	Options []Option
}

// Option defines a function that configures queue options
type Option func(map[string]any)

// Option keys
const (
	OptMaxRetry  = "max_retry"
	OptQueue     = "queue"
	OptTimeout   = "timeout"
	OptDeadline  = "deadline"
	OptUnique    = "unique"
	OptProcessAt = "process_at"
	OptProcessIn = "process_in"
	OptTaskID    = "task_id"
	OptRetention = "retention"
	OptGroup     = "group"
)

// MaxRetry sets the maximum number of retry attempts for a task
// Example: queue.Enqueue(ctx, "email:send", payload, queue.MaxRetry(3))
func MaxRetry(n int) Option {
	return func(opts map[string]any) {
		opts[OptMaxRetry] = n
	}
}

// QueueName sets the queue name for a task
// Example: queue.Enqueue(ctx, "email:send", payload, queue.QueueName("critical"))
func QueueName(name string) Option {
	return func(opts map[string]any) {
		opts[OptQueue] = name
	}
}

// Timeout sets the maximum execution time for a task
// Example: queue.Enqueue(ctx, "image:process", payload, queue.Timeout(5*time.Minute))
func Timeout(d time.Duration) Option {
	return func(opts map[string]any) {
		opts[OptTimeout] = d
	}
}

// Deadline sets the absolute time after which a task will fail if still running
// Example: queue.Enqueue(ctx, "report:generate", payload, queue.Deadline(time.Now().Add(1*time.Hour)))
func Deadline(t time.Time) Option {
	return func(opts map[string]any) {
		opts[OptDeadline] = t
	}
}

// Unique makes the task unique for the specified duration
// Example: queue.Enqueue(ctx, "notification:send", payload, queue.Unique(30*time.Minute))
func Unique(d time.Duration) Option {
	return func(opts map[string]any) {
		opts[OptUnique] = d
	}
}

// ProcessAt schedules a task to be processed at a specific time
// Example: queue.Enqueue(ctx, "newsletter:send", payload, queue.ProcessAt(tomorrow))
func ProcessAt(t time.Time) Option {
	return func(opts map[string]any) {
		opts[OptProcessAt] = t
	}
}

// ProcessIn schedules a task to be processed after the specified duration
// Example: queue.Enqueue(ctx, "reminder:send", payload, queue.ProcessIn(24*time.Hour))
func ProcessIn(d time.Duration) Option {
	return func(opts map[string]any) {
		opts[OptProcessIn] = d
	}
}

// TaskID assigns a custom ID to a task
// Example: queue.Enqueue(ctx, "order:process", payload, queue.TaskID("order-123"))
func TaskID(id string) Option {
	return func(opts map[string]any) {
		opts[OptTaskID] = id
	}
}

// Retention sets how long task data will be kept after completion
// Example: queue.Enqueue(ctx, "log:cleanup", payload, queue.Retention(7*24*time.Hour))
func Retention(d time.Duration) Option {
	return func(opts map[string]any) {
		opts[OptRetention] = d
	}
}

// Group assigns a task to a specific group
// Example: queue.Enqueue(ctx, "user:sync", payload, queue.Group("user-operations"))
func Group(name string) Option {
	return func(opts map[string]any) {
		opts[OptGroup] = name
	}
}
