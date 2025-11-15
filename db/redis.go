package db

import (
	"context"
	"fmt"

	"github.com/fatkulnurk/gostarter/pkg/config"

	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedis creates a new Redis client with production-ready configuration
func NewRedis(cfg *config.Redis) (*redis.Client, error) {
	// Set default values if not provided
	if cfg.PoolSize == 0 {
		cfg.PoolSize = 10
	}
	if cfg.MinIdleConns == 0 {
		cfg.MinIdleConns = 5
	}
	if cfg.ConnMaxLifetime == 0 {
		cfg.ConnMaxLifetime = time.Hour
	}
	if cfg.PoolTimeout == 0 {
		cfg.PoolTimeout = 4 * time.Second
	}
	if cfg.ConnMaxIdleTime == 0 {
		cfg.ConnMaxIdleTime = 30 * time.Minute
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:            cfg.Addr,
		Password:        cfg.Password,
		DB:              cfg.DB,
		PoolSize:        cfg.PoolSize,
		MinIdleConns:    cfg.MinIdleConns,
		ConnMaxLifetime: cfg.ConnMaxLifetime,
		PoolTimeout:     cfg.PoolTimeout,
		ConnMaxIdleTime: cfg.ConnMaxIdleTime,
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		DialTimeout:     cfg.DialTimeout,
	})

	// Test connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return rdb, nil
}
