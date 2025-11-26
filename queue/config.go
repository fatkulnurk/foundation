package queue

// Config holds configuration for queue worker
type Config struct {
	// Concurrency is the maximum number of concurrent processing of tasks
	// If not set or set to 0, defaults to 10
	Concurrency int

	// Queues is a map of queue names to their priority levels
	// Higher priority queues will be processed more frequently
	// Example: {"critical": 6, "default": 3, "low": 1}
	Queues map[string]int

	// StrictPriority indicates whether the queue priority should be treated strictly
	// If true, tasks in higher priority queues are processed first
	StrictPriority bool

	// ShutdownTimeout is the duration to wait for workers to finish before forcing shutdown
	// Default: 8 seconds
	ShutdownTimeout int
}
