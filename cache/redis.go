package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) Cache {
	return &RedisCache{client: client}
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, ttlSeconds int) error {
	return r.client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Has(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}
