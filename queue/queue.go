package queue

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
)

// Queue defines the interface for queueing tasks
type Queue interface {
	Enqueue(ctx context.Context, taskName string, payload any, opts ...Option) (*OutputEnqueue, error)
}

type OutputEnqueue struct {
	TaskID  string
	Payload []byte
	Options []Option
}

// Option defines a function that configures queue options
type Option func(*options)

// options holds all configuration for task processing
type options struct {
	maxRetry  int
	queue     string
	timeout   time.Duration
	deadline  time.Time
	unique    time.Duration
	processAt time.Time
	processIn time.Duration
	taskID    string
	retention time.Duration
	group     string
}

// MaxRetry sets the maximum number of retry attempts for a task
// Example: queue.Enqueue(ctx, "email:send", payload, queue.MaxRetry(3))
func MaxRetry(n int) Option {
	return func(o *options) {
		o.maxRetry = n
	}
}

// QueueName sets the queue name for a task
// Example: queue.Enqueue(ctx, "email:send", payload, queue.QueueName("critical"))
func QueueName(name string) Option {
	return func(o *options) {
		o.queue = name
	}
}

// Timeout sets the maximum execution time for a task
// Example: queue.Enqueue(ctx, "image:process", payload, queue.Timeout(5*time.Minute))
func Timeout(d time.Duration) Option {
	return func(o *options) {
		o.timeout = d
	}
}

// Deadline sets the absolute time after which a task will fail if still running
// Example: queue.Enqueue(ctx, "report:generate", payload, queue.Deadline(time.Now().Add(1*time.Hour)))
func Deadline(t time.Time) Option {
	return func(o *options) {
		o.deadline = t
	}
}

// Unique makes the task unique for the specified duration
// Example: queue.Enqueue(ctx, "notification:send", payload, queue.Unique(30*time.Minute))
func Unique(d time.Duration) Option {
	return func(o *options) {
		o.unique = d
	}
}

// ProcessAt schedules a task to be processed at a specific time
// Example: queue.Enqueue(ctx, "newsletter:send", payload, queue.ProcessAt(tomorrow))
func ProcessAt(t time.Time) Option {
	return func(o *options) {
		o.processAt = t
	}
}

// ProcessIn schedules a task to be processed after the specified duration
// Example: queue.Enqueue(ctx, "reminder:send", payload, queue.ProcessIn(24*time.Hour))
func ProcessIn(d time.Duration) Option {
	return func(o *options) {
		o.processIn = d
	}
}

// TaskID assigns a custom ID to a task
// Example: queue.Enqueue(ctx, "order:process", payload, queue.TaskID("order-123"))
func TaskID(id string) Option {
	return func(o *options) {
		o.taskID = id
	}
}

// Retention sets how long task data will be kept after completion
// Example: queue.Enqueue(ctx, "log:cleanup", payload, queue.Retention(7*24*time.Hour))
func Retention(d time.Duration) Option {
	return func(o *options) {
		o.retention = d
	}
}

// Group assigns a task to a specific group
// Example: queue.Enqueue(ctx, "user:sync", payload, queue.Group("user-operations"))
func Group(name string) Option {
	return func(o *options) {
		o.group = name
	}
}

// toAsynqOptions converts our internal options to asynq options
func toAsynqOptions(opts ...Option) []asynq.Option {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	var aOpts []asynq.Option
	if o.maxRetry > 0 {
		aOpts = append(aOpts, asynq.MaxRetry(o.maxRetry))
	}
	if o.queue != "" {
		aOpts = append(aOpts, asynq.Queue(o.queue))
	}
	if o.timeout > 0 {
		aOpts = append(aOpts, asynq.Timeout(o.timeout))
	}
	if !o.deadline.IsZero() {
		aOpts = append(aOpts, asynq.Deadline(o.deadline))
	}
	if o.unique > 0 {
		aOpts = append(aOpts, asynq.Unique(o.unique))
	}
	if !o.processAt.IsZero() {
		aOpts = append(aOpts, asynq.ProcessAt(o.processAt))
	}
	if o.processIn > 0 {
		aOpts = append(aOpts, asynq.ProcessIn(o.processIn))
	}
	if o.taskID != "" {
		aOpts = append(aOpts, asynq.TaskID(o.taskID))
	}
	if o.retention > 0 {
		aOpts = append(aOpts, asynq.Retention(o.retention))
	}
	if o.group != "" {
		aOpts = append(aOpts, asynq.Group(o.group))
	}
	return aOpts
}
