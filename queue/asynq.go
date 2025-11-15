package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fatkulnurk/foundation/config"
	"github.com/fatkulnurk/foundation/logging"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

func NewAsynqClient(cfg *config.Queue, redis *redis.Client) (*asynq.Client, error) {
	client := asynq.NewClientFromRedisClient(redis)

	err := client.Ping()
	if err != nil {
		logging.Error(context.Background(), fmt.Sprintf("failed to ping redis: %v", err))
		return nil, err
	}

	return client, nil
}

type AsynqQueue struct {
	client *asynq.Client
}

func NewAsynqQueue(client *asynq.Client) Queue {
	return &AsynqQueue{client: client}
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
