# Foundation Packages

A collection of reusable Go packages for building modern web applications. These packages provide common functionality with clean interfaces, making it easy to build scalable and maintainable applications.

## Table of Contents

- [Overview](#overview)
- [Available Packages](#available-packages)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Package Categories](#package-categories)
- [Design Principles](#design-principles)
- [Contributing](#contributing)

---

## Overview

Foundation packages are designed to be:
- **Simple** - Easy to understand and use
- **Modular** - Use only what you need
- **Extensible** - Implement interfaces for custom behavior
- **Well-documented** - Comprehensive documentation and examples
- **Production-ready** - Battle-tested and reliable

## Available Packages

### Core Infrastructure

#### [App](./app)
Application bootstrap and lifecycle management.
```go
app := app.New(app.Config{
    Name:    "MyApp",
    Version: "1.0.0",
})
```

#### [Container](./container)
Dependency injection container for managing services.
```go
c := container.New()
c.Set("db", database)
db := c.MustGet("db").(*Database)
```

#### [Module](./module)
Modular application architecture support.

---

### Data Storage

#### [Cache](./cache)
Temporary data storage with Redis and in-memory support.
```go
cache := cache.NewRedisCache(config, redisClient)
cache.Set(ctx, "key", "value", 60)
```

#### [Storage](./storage)
File storage abstraction for local filesystem and AWS S3.
```go
storage := storage.NewS3Storage(config)
storage.Upload(ctx, storage.UploadInput{...})
```

---

### HTTP & Networking

#### [HTTP Client](./httpclient)
Fluent HTTP client with retry and timeout support.
```go
client := httpclient.New(httpclient.Config{
    BaseURL: "https://api.example.com",
    Timeout: 30 * time.Second,
})
resp, err := client.Get("/users").Send()
```

#### [HTTP Router](./httprouter)
Modern HTTP router with middleware and grouping.
```go
r := httprouter.New()
r.GET("/users", listUsers)
r.POST("/users", createUsers)
```

---

### Background Processing

#### [Queue](./queue)
Task queue and background job processing with Redis.
```go
q, _ := queue.NewQueue(redisClient)
q.Enqueue(ctx, "email:send", payload,
    queue.MaxRetry(3),
    queue.Timeout(30*time.Second),
)
```

#### [Worker Pool](./workerpool)
Concurrent task processing with priority queue.
```go
pool := workerpool.NewWorkerPool(5)
pool.Submit(workerpool.Job{
    Task:     myTask,
    Priority: workerpool.High,
    Timeout:  10 * time.Second,
})
```

---

### Communication

#### [Mailer](./mailer)
Email sending with SMTP and AWS SES support.
```go
mailer := mailer.NewSMTPMailer(smtpClient, from, name)
mailer.Send(ctx, mailer.Message{
    ToAddresses: []string{"user@example.com"},
    Subject:     "Hello",
    HTMLBody:    "<h1>Hello World</h1>",
})
```

---

### Validation & View

#### [Validation](./validation)
Data validation with 20+ built-in rules.
```go
type User struct {
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"required,min=18"`
}

errs := validation.ValidateStruct(user)
```

#### [View](./view)
Template rendering engine with caching.
```go
v := view.New(view.Config{
    ViewsPath: "./templates",
})
v.RenderWithLayout(w, "main", "home", data)
```

---

### Utilities

#### [Logging](./logging)
Structured logging with multiple output formats.
```go
logger := logging.New(logging.Config{
    Level:  "info",
    Format: "json",
})
logger.Info("Application started")
```

#### [Shared](./shared)
Shared utilities and helper functions.

#### [Support](./support)
Support utilities for common operations.

---

## Installation

Install all packages:
```bash
go get github.com/fatkulnurk/foundation
```

Or install specific packages:
```bash
go get github.com/fatkulnurk/foundation/cache
go get github.com/fatkulnurk/foundation/queue
go get github.com/fatkulnurk/foundation/httpclient
```

---

## Quick Start

### Example: Building a Web API

```go
package main

import (
    "context"
    "net/http"
    
    "github.com/fatkulnurk/foundation/container"
    "github.com/fatkulnurk/foundation/httprouter"
    "github.com/fatkulnurk/foundation/cache"
    "github.com/fatkulnurk/foundation/validation"
    "github.com/redis/go-redis/v9"
)

func main() {
    // Setup container
    c := container.New()
    
    // Setup cache
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    cacheService := cache.NewRedisCache(&cache.Config{
        Prefix: "myapp:",
    }, redisClient)
    c.Set("cache", cacheService)
    
    // Setup router
    r := httprouter.New()
    r.Use(middleware.Logging)
    r.Use(middleware.RecoverMiddleware)
    
    // Define routes
    r.GET("/users", func(w http.ResponseWriter, r *http.Request) {
        cache := c.MustGet("cache").(cache.Cache)
        
        // Try to get from cache
        users, err := cache.Get(r.Context(), "users")
        if err == nil {
            httprouter.WriteJSON(w, http.StatusOK, users)
            return
        }
        
        // Fetch from database and cache
        users = fetchUsers()
        cache.Set(r.Context(), "users", users, 300)
        
        httprouter.WriteJSON(w, http.StatusOK, users)
    })
    
    r.POST("/users", func(w http.ResponseWriter, r *http.Request) {
        var user User
        json.NewDecoder(r.Body).Decode(&user)
        
        // Validate
        errs := validation.ValidateStruct(user)
        if errs.HasErrors() {
            httprouter.WriteJSON(w, http.StatusBadRequest, errs)
            return
        }
        
        // Save user
        saveUser(user)
        
        httprouter.WriteJSON(w, http.StatusCreated, user)
    })
    
    // Start server
    http.ListenAndServe(":8080", r)
}
```

### Example: Background Job Processing

```go
package main

import (
    "context"
    "log"
    
    "github.com/fatkulnurk/foundation/queue"
    "github.com/redis/go-redis/v9"
)

func main() {
    redisClient := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
    
    // Producer: Enqueue tasks
    q, _ := queue.NewQueue(redisClient)
    q.Enqueue(context.Background(), "email:send", map[string]string{
        "to":      "user@example.com",
        "subject": "Welcome!",
    }, queue.MaxRetry(3))
    
    // Worker: Process tasks
    w, _ := queue.NewWorker(redisClient, queue.Config{
        Concurrency: 10,
    })
    
    w.RegisterHandler("email:send", func(ctx context.Context, task *queue.Task) error {
        log.Printf("Sending email to: %v", task.Payload)
        return sendEmail(task.Payload)
    })
    
    w.Start()
}
```

---

## Package Categories

### 1. Infrastructure
Core application components:
- **app** - Application lifecycle
- **container** - Dependency injection
- **module** - Modular architecture

### 2. Data Layer
Data storage and caching:
- **cache** - Temporary data storage
- **storage** - File storage abstraction

### 3. HTTP Layer
Web and API functionality:
- **httpclient** - HTTP client
- **httprouter** - HTTP router and middleware

### 4. Background Processing
Asynchronous task processing:
- **queue** - Task queue with Redis
- **workerpool** - In-memory worker pool

### 5. Communication
External communication:
- **mailer** - Email sending

### 6. Presentation
User interface:
- **view** - Template rendering
- **validation** - Data validation

### 7. Utilities
Helper packages:
- **logging** - Structured logging
- **shared** - Shared utilities
- **support** - Support functions

---

## Design Principles

### 1. Interface-Based Design
All packages define clear interfaces, making them easy to mock and extend.

```go
type Cache interface {
    Set(ctx context.Context, key string, value any, ttl int) error
    Get(ctx context.Context, key string) (string, error)
    Delete(ctx context.Context, key string) error
    Has(ctx context.Context, key string) (bool, error)
}
```

### 2. Configuration Structs
Packages use configuration structs for flexibility.

```go
config := httpclient.Config{
    BaseURL:       "https://api.example.com",
    Timeout:       30 * time.Second,
    RetryCount:    3,
    RetryWaitTime: 1 * time.Second,
}
```

### 3. Context Support
All I/O operations accept context for cancellation and timeouts.

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := cache.Get(ctx, "key")
```

### 4. Error Handling
Clear error messages and proper error handling.

```go
if err != nil {
    return fmt.Errorf("failed to fetch user: %w", err)
}
```

### 5. Examples and Documentation
Every package includes:
- Comprehensive README
- Working examples
- API documentation
- Best practices

---

## Contributing

We welcome contributions! Each package has its own directory with:
- Source code
- Tests
- Examples
- Documentation

### Development Setup

1. Clone the repository
```bash
git clone https://github.com/fatkulnurk/foundation.git
cd foundation/pkg
```

2. Install dependencies
```bash
go mod download
```

3. Run tests
```bash
go test ./...
```

### Adding a New Package

1. Create package directory: `pkg/mypackage/`
2. Add source files
3. Write tests
4. Create README.md with:
   - What is it?
   - Features
   - Installation
   - Usage examples
   - API reference
   - Extending section
5. Add example in `mypackage/example/`

---

## Package Documentation

For detailed documentation on each package, see their individual README files:

- [app/README.md](./app/README.md)
- [cache/README.md](./cache/README.md)
- [container/README.md](./container/README.md)
- [httpclient/README.md](./httpclient/README.md)
- [httprouter/README.md](./httprouter/README.md)
- [logging/README.md](./logging/README.md)
- [mailer/README.md](./mailer/README.md)
- [module/README.md](./module/README.md)
- [queue/README.md](./queue/README.md)
- [shared/README.md](./shared/README.md)
- [storage/README.md](./storage/README.md)
- [support/README.md](./support/README.md)
- [validation/README.md](./validation/README.md)
- [view/README.md](./view/README.md)
- [workerpool/README.md](./workerpool/README.md)

---

## Support

For questions, issues, or feature requests, please open an issue on GitHub.

---

**Built with Go | Designed for simplicity and extensibility**
