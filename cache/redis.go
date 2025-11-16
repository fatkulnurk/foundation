package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	cfg    *Config
	client *redis.Client
}

func NewRedisCache(cfg *Config, client *redis.Client) Cache {
	return &RedisCache{cfg: cfg, client: client}
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, ttlSeconds int) error {
	key = r.cfg.Prefix + key
	return r.client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second).Err()
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	key = r.cfg.Prefix + key
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	key = r.cfg.Prefix + key
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Has(ctx context.Context, key string) (bool, error) {
	key = r.cfg.Prefix + key
	count, err := r.client.Exists(ctx, key).Result()
	return count > 0, err
}
