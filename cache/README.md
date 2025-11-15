# Cache Package

A flexible and extensible caching interface for the application with a Redis implementation.

## Overview

The `cache` package provides a standardized interface for caching operations within the application. It defines a common interface (`ICache`) that can be implemented by various cache providers, with Redis being the default implementation.

## Interface

The package defines the `ICache` interface with the following methods:

```go
type ICache interface {
	Set(ctx context.Context, key string, value any, ttlSeconds int) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Has(ctx context.Context, key string) (bool, error)
}
```

### Methods

- **Set**: Stores a value in the cache with a specified TTL (Time-To-Live) in seconds
- **Get**: Retrieves a value from the cache by its key
- **Delete**: Removes a value from the cache by its key
- **Has**: Checks if a key exists in the cache

## Implementations

### Redis Cache

The package includes a Redis implementation of the cache interface:

```go
type RedisCache struct {
	client *redis.Client
}
```

#### Usage

```go
import (
	"github.com/redis/go-redis/v9"
	"your-project/pkg/cache"
)

// Initialize Redis client
redisClient := redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

// Create cache instance
cacheInstance := cache.NewRedisCache(redisClient)

// Use the cache
ctx := context.Background()
cacheInstance.Set(ctx, "key", "value", 3600) // Cache for 1 hour
value, err := cacheInstance.Get(ctx, "key")
```

## Extending

To implement a new cache provider, create a struct that implements all methods of the `ICache` interface.

## Thread Safety

All implementations are designed to be thread-safe and can be safely used concurrently from multiple goroutines.