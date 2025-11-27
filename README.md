# Foundation Packages

A collection of reusable Go packages for building modern web applications. These packages provide common functionality with clean interfaces, making it easy to build scalable and maintainable applications.
## Available Packages

| Package | Description |
|---------|-------------|
| [cache](./cache) | Temporary data storage (Redis, in-memory) |
| [container](./container) | Dependency injection container |
| [httpclient](./httpclient) | HTTP client with retry and timeout |
| [httprouter](./httprouter) | HTTP router with middleware |
| [logging](./logging) | Structured logging |
| [mailer](./mailer) | Email sending (SMTP, AWS SES) |
| [queue](./queue) | Task queue with Redis |
| [shared](./shared) | Shared utilities |
| [storage](./storage) | File storage (Local, S3) |
| [support](./support) | Support utilities |
| [validation](./validation) | Data validation with 20+ rules |
| [view](./view) | Template rendering engine |
| [workerpool](./workerpool) | Concurrent task processing |

## Installation

Install all packages:
```bash
go get github.com/fatkulnurk/foundation
```

Or install specific packages:
```bash
# Install specific package
go get github.com/fatkulnurk/foundation/cache
go get github.com/fatkulnurk/foundation/queue
```

## Quick Example

```go
// HTTP Router with Cache
r := httprouter.New()
cache := cache.NewRedisCache(config, redisClient)

r.GET("/users", func(w http.ResponseWriter, r *http.Request) {
    users, _ := cache.Get(r.Context(), "users")
    httprouter.WriteJSON(w, http.StatusOK, users)
})

// Background Queue
q, _ := queue.NewQueue(redisClient)
q.Enqueue(ctx, "email:send", payload, queue.MaxRetry(3))
```

## Documentation

Each package has detailed documentation in its README file. Click on the package name in the table above to view its documentation.

---
