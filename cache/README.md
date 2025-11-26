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
-  Application runs on multiple servers
-  Need persistent cache (survives restart)
-  Data needs to be shared across services
-  Production application

### Use **Local** when:
-  For testing
-  Small application (single server)
-  Don't want to install Redis
-  Cache data is not critical

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

---

## Extending

You can create custom cache implementations by implementing the Cache interface.

### Custom Cache Implementation

```go
type Cache interface {
    Set(ctx context.Context, key string, value any, ttlSeconds int) error
    Get(ctx context.Context, key string) (string, error)
    Delete(ctx context.Context, key string) error
    Has(ctx context.Context, key string) (bool, error)
}
```

### Example: Memory Cache with LRU

```go
type LRUCache struct {
    cache    *lru.Cache
    prefix   string
    mu       sync.RWMutex
}

func NewLRUCache(size int, prefix string) *LRUCache {
    c, _ := lru.New(size)
    return &LRUCache{
        cache:  c,
        prefix: prefix,
    }
}

func (c *LRUCache) Set(ctx context.Context, key string, value any, ttlSeconds int) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    fullKey := c.prefix + key
    c.cache.Add(fullKey, value)
    
    // Set expiration timer
    if ttlSeconds > 0 {
        time.AfterFunc(time.Duration(ttlSeconds)*time.Second, func() {
            c.Delete(ctx, key)
        })
    }
    
    return nil
}

func (c *LRUCache) Get(ctx context.Context, key string) (string, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    fullKey := c.prefix + key
    if val, ok := c.cache.Get(fullKey); ok {
        return fmt.Sprintf("%v", val), nil
    }
    
    return "", fmt.Errorf("key not found")
}

func (c *LRUCache) Delete(ctx context.Context, key string) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    fullKey := c.prefix + key
    c.cache.Remove(fullKey)
    return nil
}

func (c *LRUCache) Has(ctx context.Context, key string) (bool, error) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    fullKey := c.prefix + key
    return c.cache.Contains(fullKey), nil
}
```

### Example: File-based Cache

```go
type FileCache struct {
    dir    string
    prefix string
}

func NewFileCache(dir, prefix string) *FileCache {
    os.MkdirAll(dir, 0755)
    return &FileCache{
        dir:    dir,
        prefix: prefix,
    }
}

func (c *FileCache) Set(ctx context.Context, key string, value any, ttlSeconds int) error {
    filename := filepath.Join(c.dir, c.prefix+key)
    data := fmt.Sprintf("%v", value)
    return os.WriteFile(filename, []byte(data), 0644)
}

func (c *FileCache) Get(ctx context.Context, key string) (string, error) {
    filename := filepath.Join(c.dir, c.prefix+key)
    data, err := os.ReadFile(filename)
    if err != nil {
        return "", err
    }
    return string(data), nil
}

func (c *FileCache) Delete(ctx context.Context, key string) error {
    filename := filepath.Join(c.dir, c.prefix+key)
    return os.Remove(filename)
}

func (c *FileCache) Has(ctx context.Context, key string) (bool, error) {
    filename := filepath.Join(c.dir, c.prefix+key)
    _, err := os.Stat(filename)
    return err == nil, nil
}
```

---

## See Also

- `example/main.go` - Complete examples of all features
- Redis documentation: https://redis.io/docs/