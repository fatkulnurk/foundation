package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fatkulnurk/foundation/logging"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

// AsynqQueue implements Queue interface using asynq
type AsynqQueue struct {
	client *asynq.Client
}

// NewQueue creates a new Queue instance using asynq
// This is the main constructor that external code should use
func NewQueue(redis *redis.Client) (Queue, error) {
	client := asynq.NewClientFromRedisClient(redis)

	err := client.Ping()
	if err != nil {
		logging.Error(context.Background(), fmt.Sprintf("failed to ping redis: %v", err))
		return nil, err
	}

	return &AsynqQueue{client: client}, nil
}

func (q *AsynqQueue) Enqueue(ctx context.Context, taskName string, payload any, opts ...Option) (*OutputEnqueue, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	task := asynq.NewTask(taskName, data)
	aOpts := toAsynqOptions(opts...)
	tInfo, err := q.client.EnqueueContext(ctx, task, aOpts...)
	if err != nil {
		return nil, err
	}
	return &OutputEnqueue{TaskID: tInfo.ID, Payload: data, Options: opts}, nil
}

// AsynqWorker implements Worker interface using asynq
type AsynqWorker struct {
	server   *asynq.Server
	mux      *asynq.ServeMux
	handlers map[string]Handler
}

// NewWorker creates a new Worker instance using asynq
// This is the main constructor that external code should use
func NewWorker(cfg *Config, redis *redis.Client) Worker {
	// Set defaults
	if cfg.Concurrency == 0 {
		cfg.Concurrency = 10
	}
	if cfg.Queues == nil {
		cfg.Queues = map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		}
	}
	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = 8
	}

	serverCfg := asynq.Config{
		Concurrency:     cfg.Concurrency,
		Queues:          cfg.Queues,
		StrictPriority:  cfg.StrictPriority,
		ShutdownTimeout: time.Duration(cfg.ShutdownTimeout) * time.Second,
	}

	server := asynq.NewServerFromRedisClient(redis, serverCfg)
	mux := asynq.NewServeMux()

	return &AsynqWorker{
		server:   server,
		mux:      mux,
		handlers: make(map[string]Handler),
	}
}

func (w *AsynqWorker) Register(taskType string, handler Handler) error {
	// Store handler for reference
	w.handlers[taskType] = handler

	// Wrap our Handler to asynq.Handler
	w.mux.HandleFunc(taskType, func(ctx context.Context, task *asynq.Task) error {
		return handler(ctx, task.Payload())
	})

	return nil
}

func (w *AsynqWorker) Start() error {
	logging.Info(context.Background(), fmt.Sprintf("Starting worker with %d registered handlers", len(w.handlers)))
	return w.server.Run(w.mux)
}

func (w *AsynqWorker) Stop() {
	w.server.Shutdown()
	logging.Info(context.Background(), "Worker stopped")
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
