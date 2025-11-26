# Container - Service Registry

Module for managing application services and dependencies in one place.

## What is Container?

Container is like a storage box where you keep all your application's services (database, logger, cache, etc.). Instead of creating these services everywhere, you register them once in the container and retrieve them when needed.

**Think of it like:**
- A toolbox where you store all your tools
- A phone book where you look up contacts by name
- A registry where services are stored and retrieved by name

**Use cases:**
- Dependency injection
- Service locator pattern
- Managing application-wide services
- Sharing instances across your application

## Module Contents

### 1. **container.go** - Main Implementation
Provides a thread-safe service container with 4 operations:
- `Set(name, service)` - Register a service
- `Get(name)` - Retrieve a service (returns error if not found)
- `MustGet(name)` - Retrieve a service (panics if not found)
- `Has(name)` - Check if a service exists

## How to Use

### Basic Usage

```go
import "github.com/fatkulnurk/foundation/container"

// Create container
c := container.New()

// Register services
c.Set("database", dbConnection)
c.Set("logger", loggerInstance)
c.Set("cache", cacheInstance)

// Retrieve services
db, err := c.Get("database")
if err != nil {
    log.Fatal(err)
}

// Or use MustGet (panics if not found)
logger := c.MustGet("logger")

// Check if service exists
if c.Has("cache") {
    cache := c.MustGet("cache")
}
```

### Type Assertion

Since container stores `interface{}`, you need to assert the type:

```go
// Register
db := &Database{Host: "localhost"}
c.Set("database", db)

// Retrieve and assert type
dbInterface, _ := c.Get("database")
database := dbInterface.(*Database)

// Or in one line with MustGet
database := c.MustGet("database").(*Database)
```

### Storing Different Types

```go
// String
c.Set("appName", "MyApp")
name := c.MustGet("appName").(string)

// Integer
c.Set("port", 8080)
port := c.MustGet("port").(int)

// Boolean
c.Set("debug", true)
debug := c.MustGet("debug").(bool)

// Struct
c.Set("config", &Config{...})
config := c.MustGet("config").(*Config)

// Map
c.Set("settings", map[string]string{...})
settings := c.MustGet("settings").(map[string]string)

// Slice
c.Set("hosts", []string{"localhost", "127.0.0.1"})
hosts := c.MustGet("hosts").([]string)
```

## Real-World Example

### Dependency Injection Pattern

```go
// Setup container
c := container.New()

// Register core services
c.Set("db", &Database{ConnectionString: "localhost:5432"})
c.Set("logger", &Logger{Level: "info"})
c.Set("cache", &RedisCache{Addr: "localhost:6379"})

// Create service with dependencies
userService := &UserService{
    DB:     c.MustGet("db").(*Database),
    Logger: c.MustGet("logger").(*Logger),
    Cache:  c.MustGet("cache").(*RedisCache),
}
c.Set("userService", userService)

// Use in handlers
func UserHandler(w http.ResponseWriter, r *http.Request) {
    userSvc := c.MustGet("userService").(*UserService)
    users := userSvc.GetAllUsers()
    json.NewEncoder(w).Encode(users)
}
```

### Application Bootstrap

```go
func main() {
    // Create container
    c := container.New()
    
    // Initialize services
    db := initDatabase()
    logger := initLogger()
    cache := initCache()
    
    // Register in container
    c.Set("db", db)
    c.Set("logger", logger)
    c.Set("cache", cache)
    
    // Create application services
    c.Set("userService", NewUserService(c))
    c.Set("authService", NewAuthService(c))
    c.Set("emailService", NewEmailService(c))
    
    // Start application
    startServer(c)
}
```

## Get vs MustGet

### Use `Get()` when:
-  Service might not exist
-  You want to handle errors gracefully
-  Optional dependencies

```go
cache, err := c.Get("cache")
if err != nil {
    log.Println("Cache not available, using fallback")
    cache = NewMemoryCache()
}
```

### Use `MustGet()` when:
-  Service must exist (critical dependency)
-  Application cannot run without it
-  You want to fail fast

```go
// Application cannot run without database
db := c.MustGet("database").(*Database)
```

## Thread Safety

The container is **thread-safe** and can be used concurrently:

```go
// Safe to use from multiple goroutines
go func() {
    db := c.MustGet("database")
}()

go func() {
    logger := c.MustGet("logger")
}()
```

## Testing

Container makes testing easier by allowing you to swap implementations:

```go
// Production
c.Set("database", &PostgresDB{...})

// Testing
c.Set("database", &MockDB{...})

// Your code works with both!
db := c.MustGet("database").(Database)
```

## Best Practices

### 1. Register Early
Register all services during application startup:
```go
func main() {
    c := container.New()
    registerServices(c)
    startApp(c)
}
```

### 2. Use Interfaces
Store interfaces, not concrete types:
```go
type Logger interface {
    Info(msg string)
    Error(msg string)
}

c.Set("logger", loggerInstance) // loggerInstance implements Logger
logger := c.MustGet("logger").(Logger)
```

### 3. Consistent Naming
Use consistent service names:
```go
// Good
c.Set("database", db)
c.Set("logger", logger)
c.Set("cache", cache)

// Avoid
c.Set("db", db)
c.Set("log", logger)
c.Set("redis", cache)
```

### 4. Check Existence for Optional Services
```go
var cache Cache
if c.Has("cache") {
    cache = c.MustGet("cache").(Cache)
} else {
    cache = NewNullCache() // Fallback
}
```

## Common Patterns

### Singleton Pattern
```go
var appContainer *container.Locator

func GetContainer() *container.Locator {
    if appContainer == nil {
        appContainer = container.New()
        initializeServices(appContainer)
    }
    return appContainer
}
```

### Factory Pattern
```go
func NewUserService(c container.Container) *UserService {
    return &UserService{
        DB:     c.MustGet("database").(*Database),
        Logger: c.MustGet("logger").(*Logger),
    }
}

c.Set("userService", NewUserService(c))
```

## Installation

```bash
go get github.com/fatkulnurk/foundation/container
```

## Dependencies

None - uses only Go standard library.

---

## Extending

The container package is designed to be simple and focused on dependency injection. You can extend it by creating wrapper functions or custom service providers.

### Custom Service Provider

```go
type ServiceProvider interface {
    Register(c *container.Container)
    Boot(c *container.Container)
}

type DatabaseProvider struct{}

func (p *DatabaseProvider) Register(c *container.Container) {
    c.Set("db", func(c *container.Container) interface{} {
        config := c.Get("config").(*Config)
        return connectDatabase(config.DatabaseURL)
    })
}

func (p *DatabaseProvider) Boot(c *container.Container) {
    db := c.Get("db").(*Database)
    db.Migrate()
}
```

### Example: Lazy Loading Services

```go
type LazyContainer struct {
    *container.Container
    loaded map[string]bool
}

func NewLazyContainer() *LazyContainer {
    return &LazyContainer{
        Container: container.New(),
        loaded:    make(map[string]bool),
    }
}

func (lc *LazyContainer) GetLazy(name string) interface{} {
    if !lc.loaded[name] {
        // Initialize service on first access
        service := lc.Get(name)
        if initializer, ok := service.(interface{ Initialize() }); ok {
            initializer.Initialize()
        }
        lc.loaded[name] = true
    }
    return lc.Get(name)
}
```

### Example: Scoped Container

```go
type ScopedContainer struct {
    parent *container.Container
    scope  *container.Container
}

func NewScopedContainer(parent *container.Container) *ScopedContainer {
    return &ScopedContainer{
        parent: parent,
        scope:  container.New(),
    }
}

func (sc *ScopedContainer) Get(name string) interface{} {
    // Try scope first, then parent
    if sc.scope.Has(name) {
        return sc.scope.Get(name)
    }
    return sc.parent.Get(name)
}

func (sc *ScopedContainer) Set(name string, value interface{}) {
    sc.scope.Set(name, value)
}
```

---

## See Also

- `example/main.go` - Complete examples of all features
- `container_test.go` - Unit tests showing usage patterns
