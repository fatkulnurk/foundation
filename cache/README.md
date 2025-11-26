# Cache - Temporary Data Storage

Module for storing temporary data to make your application faster.

## What is Cache?

Cache is a temporary data storage. Think of it like a notepad on your desk - you keep frequently used information there so you don't have to walk to the filing cabinet every time you need the same data.

**Use cases:**
- Store frequently accessed user data
- Store heavy computation results
- Store database query results to avoid repeated queries

## Module Contents

### 1. **cache.go** - Main Interface
Defines 4 basic cache operations:
- `Set` - Store data (with expiration time)
- `Get` - Retrieve data
- `Delete` - Remove data
- `Has` - Check if data exists

### 2. **redis.go** - Redis Cache
Cache implementation using Redis (fast in-memory database).

**How to use:**
```go
// Create Redis connection
redisClient := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
})

// Create cache with prefix
cfg := &cache.Config{Prefix: "myapp:"}
c := cache.NewRedisCache(cfg, redisClient)

// Store data for 60 seconds
c.Set(ctx, "user:1", "John Doe", 60)

// Retrieve data
value, _ := c.Get(ctx, "user:1")
```

### 3. **local.go** - In-Memory Cache
Cache implementation stored in application memory (no Redis needed).

**When to use:**
- For testing
- Small applications that don't need Redis
- Data that doesn't need to be shared across servers

**How to use:**
```go
cfg := &cache.Config{Prefix: "myapp:"}
c := cache.NewLocalCache(cfg)

// Same API as Redis
c.Set(ctx, "session:abc", "data", 300)
value, _ := c.Get(ctx, "session:abc")
```

### 4. **config.go** - Configuration
Simple cache configuration:
- `Prefix` - Prefix for all keys (example: "myapp:")

Can be loaded from environment variable:
```bash
CACHE_PREFIX=myapp: go run main.go
```

### 5. **example/** - Usage Examples
Folder contains complete examples of how to use cache.

**Run:**
```bash
# Make sure Redis is running
redis-server

# Run example
go run ./example
```

## When to Use Redis vs Local?

### Use **Redis** when:
- ✅ Application runs on multiple servers
- ✅ Need persistent cache (survives restart)
- ✅ Data needs to be shared across services
- ✅ Production application

### Use **Local** when:
- ✅ For testing
- ✅ Small application (single server)
- ✅ Don't want to install Redis
- ✅ Cache data is not critical

## Important Features

### TTL (Time To Live)
Every data has an expiration time. After the time expires, data is automatically removed.

```go
// Data expires after 60 seconds
c.Set(ctx, "key", "value", 60)

// Data expires after 1 hour (3600 seconds)
c.Set(ctx, "key", "value", 3600)
```

### Prefix
Prefix helps organize data and avoid key conflicts.

```go
cfg := &cache.Config{Prefix: "myapp:"}
c := cache.NewRedisCache(cfg, redisClient)

// Actual key: "myapp:user:1"
c.Set(ctx, "user:1", "John", 60)
```

## Installation

```bash
go get github.com/fatkulnurk/foundation/cache
```

## Dependencies

- Redis cache: `github.com/redis/go-redis/v9`
- Support utilities: `github.com/fatkulnurk/foundation/support`

## See Also

- `example/main.go` - Complete examples of all features
- Redis documentation: https://redis.io/docs/