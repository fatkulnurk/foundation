# HTTP Router - Modern HTTP Router for Go

Module for building HTTP APIs with routing, middleware, and grouping support.

## What is HTTP Router?

HTTP Router is a lightweight wrapper around Go's standard `net/http` ServeMux (Go 1.22+) that makes it easier to build web applications and APIs. It provides a clean API for defining routes, middleware, and route groups.

**Think of it like:**
- A traffic controller that directs HTTP requests to the right handlers
- A way to organize your API endpoints
- A middleware chain builder

**Use cases:**
- Building REST APIs
- Creating web applications
- Organizing routes with groups (e.g., /api, /admin)
- Adding middleware (logging, auth, CORS, etc.)

## Module Contents

### 1. **router.go** - Main Router
Provides routing with:
- HTTP method helpers (GET, POST, PUT, PATCH, DELETE)
- Route groups with prefixes
- Middleware support (global and per-route)
- Static file serving
- Path parameters (Go 1.22+ native)

### 2. **response.go** - Response Helpers
Helper functions for sending responses:
- `WriteJSON` - Send JSON response
- `WriteHTML` - Send HTML response
- `WriteText` - Send plain text response
- `WriteError` - Send error response

### 3. **middleware/** - Built-in Middleware
- `Logging` - Request logging
- `RecoverMiddleware` - Panic recovery
- `CORS` - Cross-Origin Resource Sharing
- `RateLimit` - Rate limiting
- `RequireAPIKey` - API key authentication

## How to Use

### Basic Usage

```go
import "github.com/fatkulnurk/foundation/httprouter"

func main() {
    r := httprouter.New()

    // Define routes
    r.GET("/", func(w http.ResponseWriter, r *http.Request) {
        httprouter.WriteJSON(w, http.StatusOK, map[string]string{
            "message": "Hello World",
        })
    })

    // Start server
    http.ListenAndServe(":8080", r)
}
```

### HTTP Methods

```go
r := httprouter.New()

// GET request
r.GET("/users", listUsers)

// POST request
r.POST("/users", createUser)

// PUT request
r.PUT("/users/{id}", updateUser)

// PATCH request
r.PATCH("/users/{id}", patchUser)

// DELETE request
r.DELETE("/users/{id}", deleteUser)
```

### Path Parameters

```go
// Go 1.22+ native path parameters
r.GET("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    fmt.Fprintf(w, "User ID: %s", id)
})

// Multiple parameters
r.GET("/posts/{postID}/comments/{commentID}", func(w http.ResponseWriter, r *http.Request) {
    postID := r.PathValue("postID")
    commentID := r.PathValue("commentID")
    fmt.Fprintf(w, "Post: %s, Comment: %s", postID, commentID)
})
```

### Global Middleware

```go
r := httprouter.New()

// Add global middleware (applies to all routes)
r.Use(middleware.Logging)
r.Use(middleware.RecoverMiddleware)

// All routes below will use the middleware
r.GET("/users", listUsers)
r.POST("/users", createUser)
```

### Per-Route Middleware

```go
// Middleware only for this route
r.GET("/admin", adminHandler, middleware.RequireAPIKey)

// Multiple middleware for one route
r.GET("/protected", protectedHandler, 
    middleware.RequireAPIKey,
    middleware.RequireAdmin,
)
```

### Route Groups

```go
r := httprouter.New()

// Group routes with /api prefix
r.Group("/api", func(api httprouter.HttpRouter) {
    // All routes here have /api prefix
    
    // GET /api/users
    api.GET("/users", listUsers)
    
    // POST /api/users
    api.POST("/users", createUser)
    
    // Nested group: /api/admin
    api.Group("/admin", func(admin httprouter.HttpRouter) {
        // GET /api/admin/stats
        admin.GET("/stats", getStats)
    })
})
```

### Group with Middleware

```go
r.Group("/api", func(api httprouter.HttpRouter) {
    // Middleware for all routes in this group
    api.Use(middleware.RequireAPIKey)
    
    // All routes below require API key
    api.GET("/users", listUsers)
    api.POST("/users", createUser)
})
```

### Static Files

```go
// Serve static files from ./public directory
r.Static("/static", "./public")

// Access: http://localhost:8080/static/style.css
// Serves: ./public/style.css
```

### Static Files with Middleware

```go
// Serve with authentication
r.Static("/private", "./private-files", middleware.RequireAPIKey)
```

## Response Helpers

### JSON Response

```go
r.GET("/users", func(w http.ResponseWriter, r *http.Request) {
    users := []User{
        {ID: 1, Name: "John"},
        {ID: 2, Name: "Jane"},
    }
    
    httprouter.WriteJSON(w, http.StatusOK, users)
})
```

### HTML Response

```go
r.GET("/", func(w http.ResponseWriter, r *http.Request) {
    html := "<h1>Welcome</h1>"
    httprouter.WriteHTML(w, http.StatusOK, html)
})
```

### Text Response

```go
r.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
    httprouter.WriteText(w, http.StatusOK, "pong")
})
```

### Error Response

```go
r.GET("/error", func(w http.ResponseWriter, r *http.Request) {
    httprouter.WriteError(w, http.StatusBadRequest, "Invalid request")
})
```

## Built-in Middleware

### Logging Middleware

```go
r.Use(middleware.Logging)
// Logs: [2024-11-26 18:00:00] GET /users 200 15ms
```

### Recover Middleware

```go
r.Use(middleware.RecoverMiddleware)
// Catches panics and returns 500 error
```

### CORS Middleware

```go
r.Use(middleware.CORS(middleware.CORSOptions{
    AllowedOrigins:   []string{"https://example.com"},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
    AllowedHeaders:   []string{"Content-Type", "Authorization"},
    AllowCredentials: true,
    MaxAge:           600,
}))
```

### Rate Limit Middleware

```go
r.Use(middleware.NewRateLimitMiddleware(middleware.RateLimitConfig{
    Requests: 100,              // 100 requests
    Window:   time.Minute,      // per minute
}))
```

### API Key Middleware

```go
r.Use(middleware.RequireAPIKey)
// Requires X-API-Key header
```

## Real-World Example

### REST API with Groups and Middleware

```go
func main() {
    r := httprouter.New()
    
    // Global middleware
    r.Use(middleware.Logging)
    r.Use(middleware.RecoverMiddleware)
    r.Use(middleware.CORS(middleware.CORSOptions{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
    }))
    
    // Public routes
    r.GET("/", homeHandler)
    r.GET("/health", healthHandler)
    
    // API routes (require API key)
    r.Group("/api", func(api httprouter.HttpRouter) {
        api.Use(middleware.RequireAPIKey)
        
        // User routes
        api.GET("/users", listUsers)
        api.GET("/users/{id}", getUser)
        api.POST("/users", createUser)
        api.PUT("/users/{id}", updateUser)
        api.DELETE("/users/{id}", deleteUser)
        
        // Admin routes (require admin role)
        api.Group("/admin", func(admin httprouter.HttpRouter) {
            admin.Use(requireAdminRole)
            
            admin.GET("/stats", getStats)
            admin.GET("/logs", getLogs)
        })
    })
    
    // Static files
    r.Static("/static", "./public")
    
    log.Println("Server running on :8080")
    http.ListenAndServe(":8080", r)
}
```

### Custom Middleware

```go
// Authentication middleware
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        
        if token == "" {
            httprouter.WriteError(w, http.StatusUnauthorized, "Missing token")
            return
        }
        
        // Validate token...
        
        next.ServeHTTP(w, r)
    })
}

// Use it
r.Use(authMiddleware)
```

### Handler with Error Handling

```go
r.GET("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")
    
    user, err := getUserByID(id)
    if err != nil {
        httprouter.WriteError(w, http.StatusNotFound, "User not found")
        return
    }
    
    httprouter.WriteJSON(w, http.StatusOK, user)
})
```

## Best Practices

### 1. Use Route Groups for Organization

```go
// Good - organized by feature
r.Group("/api", func(api httprouter.HttpRouter) {
    api.Group("/users", func(users httprouter.HttpRouter) {
        users.GET("", listUsers)
        users.POST("", createUser)
        users.GET("/{id}", getUser)
    })
    
    api.Group("/posts", func(posts httprouter.HttpRouter) {
        posts.GET("", listPosts)
        posts.POST("", createPost)
    })
})
```

### 2. Apply Middleware at the Right Level

```go
r := httprouter.New()

// Global middleware (all routes)
r.Use(middleware.Logging)
r.Use(middleware.RecoverMiddleware)

// Group middleware (only /api routes)
r.Group("/api", func(api httprouter.HttpRouter) {
    api.Use(middleware.RequireAPIKey)
    
    // Route middleware (only this specific route)
    api.GET("/admin", adminHandler, requireAdminRole)
})
```

### 3. Use Response Helpers

```go
// Good - consistent responses
r.GET("/users", func(w http.ResponseWriter, r *http.Request) {
    users, err := getUsers()
    if err != nil {
        httprouter.WriteError(w, http.StatusInternalServerError, err.Error())
        return
    }
    httprouter.WriteJSON(w, http.StatusOK, users)
})

// Avoid - manual response writing
r.GET("/users", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(users)
})
```

### 4. Handle Errors Gracefully

```go
r.Use(middleware.RecoverMiddleware)

r.GET("/danger", func(w http.ResponseWriter, r *http.Request) {
    // If this panics, RecoverMiddleware catches it
    result := riskyOperation()
    httprouter.WriteJSON(w, http.StatusOK, result)
})
```

## Common Patterns

### Versioned API

```go
r.Group("/api/v1", func(v1 httprouter.HttpRouter) {
    v1.GET("/users", listUsersV1)
})

r.Group("/api/v2", func(v2 httprouter.HttpRouter) {
    v2.GET("/users", listUsersV2)
})
```

### Protected Routes

```go
r.Group("/admin", func(admin httprouter.HttpRouter) {
    admin.Use(middleware.RequireAPIKey)
    admin.Use(requireAdminRole)
    
    admin.GET("/dashboard", dashboardHandler)
    admin.GET("/users", adminUsersHandler)
})
```

### File Upload

```go
r.POST("/upload", func(w http.ResponseWriter, r *http.Request) {
    file, header, err := r.FormFile("file")
    if err != nil {
        httprouter.WriteError(w, http.StatusBadRequest, "Invalid file")
        return
    }
    defer file.Close()
    
    // Save file...
    
    httprouter.WriteJSON(w, http.StatusOK, map[string]string{
        "filename": header.Filename,
        "size":     fmt.Sprintf("%d", header.Size),
    })
})
```

## Testing

```go
func TestRouter(t *testing.T) {
    r := httprouter.New()
    
    r.GET("/ping", func(w http.ResponseWriter, r *http.Request) {
        httprouter.WriteText(w, http.StatusOK, "pong")
    })
    
    req := httptest.NewRequest("GET", "/ping", nil)
    w := httptest.NewRecorder()
    
    r.ServeHTTP(w, req)
    
    if w.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d", w.Code)
    }
    
    if w.Body.String() != "pong" {
        t.Errorf("Expected 'pong', got %s", w.Body.String())
    }
}
```

## Installation

```bash
go get github.com/fatkulnurk/foundation/httprouter
```

## Dependencies

None - uses only Go standard library (requires Go 1.22+ for path parameters).

---

## Extending

You can extend the HTTP router by creating custom middleware and response helpers.

### Custom Middleware

```go
// Timing middleware
func TimingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        // Call next handler
        next.ServeHTTP(w, r)
        
        // Log timing
        duration := time.Since(start)
        log.Printf("[%s] %s - %v", r.Method, r.URL.Path, duration)
    })
}

// Use it
r.Use(TimingMiddleware)
```

### Custom Response Writer

```go
type ResponseWriter struct {
    http.ResponseWriter
    statusCode int
    bytes      int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
    return &ResponseWriter{
        ResponseWriter: w,
        statusCode:     http.StatusOK,
    }
}

func (rw *ResponseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
    n, err := rw.ResponseWriter.Write(b)
    rw.bytes += n
    return n, err
}

// Logging middleware using custom response writer
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        rw := NewResponseWriter(w)
        start := time.Now()
        
        next.ServeHTTP(rw, r)
        
        log.Printf("%s %s - %d (%d bytes) - %v",
            r.Method, r.URL.Path, rw.statusCode, rw.bytes, time.Since(start))
    })
}
```

### Example: Request ID Middleware

```go
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        // Add to response header
        w.Header().Set("X-Request-ID", requestID)
        
        // Add to context
        ctx := context.WithValue(r.Context(), "request_id", requestID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func generateRequestID() string {
    return fmt.Sprintf("%d", time.Now().UnixNano())
}
```

### Example: Custom Error Handler

```go
type ErrorResponse struct {
    Error     string `json:"error"`
    RequestID string `json:"request_id,omitempty"`
    Timestamp string `json:"timestamp"`
}

func WriteErrorWithContext(w http.ResponseWriter, r *http.Request, status int, message string) {
    requestID, _ := r.Context().Value("request_id").(string)
    
    response := ErrorResponse{
        Error:     message,
        RequestID: requestID,
        Timestamp: time.Now().Format(time.RFC3339),
    }
    
    httprouter.WriteJSON(w, status, response)
}
```

### Example: Rate Limiter per IP

```go
type IPRateLimiter struct {
    limiters map[string]*rate.Limiter
    mu       sync.RWMutex
    rate     rate.Limit
    burst    int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
    return &IPRateLimiter{
        limiters: make(map[string]*rate.Limiter),
        rate:     r,
        burst:    b,
    }
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
    i.mu.Lock()
    defer i.mu.Unlock()
    
    limiter, exists := i.limiters[ip]
    if !exists {
        limiter = rate.NewLimiter(i.rate, i.burst)
        i.limiters[ip] = limiter
    }
    
    return limiter
}

func (i *IPRateLimiter) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        limiter := i.GetLimiter(ip)
        
        if !limiter.Allow() {
            httprouter.WriteError(w, http.StatusTooManyRequests, "Rate limit exceeded")
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

---

## See Also

- `example/main.go` - Complete example with all features
- `middleware/` - Built-in middleware implementations
- `router_test.go` - Unit tests and usage examples
