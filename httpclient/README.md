# HTTP Client - Simple HTTP Request Wrapper

Module for making HTTP requests with a clean and fluent API.

## What is HTTP Client?

HTTP Client is a wrapper around Go's standard `net/http` that makes it easier to make HTTP requests. Think of it like a simplified way to talk to APIs and web services without writing repetitive code.

**Use cases:**
- Calling REST APIs
- Consuming third-party services
- Making HTTP requests with retry logic
- Sending data to webhooks

## Module Contents

### 1. **httpclient.go** - Main Client
Creates HTTP client instances with configuration.

### 2. **request.go** - Request Builder
Fluent API for building HTTP requests with:
- Headers
- JSON body
- Form data
- Multipart form (file uploads)
- Context support

### 3. **response.go** - Response Wrapper
Wraps HTTP responses with helper methods:
- `StatusCode` - HTTP status code
- `Body` - Response body as bytes
- `String()` - Response body as string
- `JSON(target)` - Parse JSON response
- `IsSuccess()` - Check if 2xx status

### 4. **config.go** - Configuration
Client configuration options:
- `BaseURL` - Base URL for all requests
- `Timeout` - Request timeout
- `RetryCount` - Number of retries on failure
- `RetryWaitTime` - Wait time between retries
- `DefaultHeaders` - Headers added to all requests

## How to Use

### Basic Usage

```go
import "github.com/fatkulnurk/foundation/httpclient"

// Create client
client := httpclient.New(httpclient.Config{
    BaseURL: "https://api.example.com",
    Timeout: 30 * time.Second,
})

// GET request
resp, err := client.Get("/users").Send()
if err != nil {
    log.Fatal(err)
}

fmt.Println(resp.String())
```

### GET with Headers

```go
resp, err := client.Get("/users").
    WithHeader("Authorization", "Bearer token").
    WithHeader("X-Request-ID", "abc123").
    Send()
```

### POST with JSON

```go
resp, err := client.Post("/users").
    WithJSON(map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
        "age":   30,
    }).
    Send()
```

### POST with Form Data

```go
resp, err := client.Post("/login").
    WithFormURLEncoded(map[string]string{
        "username": "john",
        "password": "secret",
    }).
    Send()
```

### File Upload (Multipart Form)

```go
fileData, _ := os.ReadFile("photo.jpg")

resp, err := client.Post("/upload").
    WithMultipartForm(
        map[string]string{"title": "My Photo"},
        map[string][]byte{"photo.jpg": fileData},
    }).
    Send()
```

### PUT Request

```go
resp, err := client.Put("/users/123").
    WithJSON(map[string]interface{}{
        "name": "John Updated",
        "age":  31,
    }).
    Send()
```

### PATCH Request

```go
resp, err := client.Patch("/users/123").
    WithJSON(map[string]interface{}{
        "age": 32,
    }).
    Send()
```

### DELETE Request

```go
resp, err := client.Delete("/users/123").
    WithHeader("Authorization", "Bearer token").
    Send()
```

### Multiple Headers

```go
resp, err := client.Get("/users").
    WithHeaders(map[string]string{
        "Authorization": "Bearer token",
        "X-API-Version": "v2",
        "Accept":        "application/json",
    }).
    Send()
```

### With Context Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := client.Get("/users").
    WithContext(ctx).
    Send()
```

### Custom Request Method

```go
resp, err := client.NewRequest().
    WithMethod("PATCH").
    WithURL("/api/data").
    WithHeader("Content-Type", "application/json").
    WithJSON(data).
    Send()
```

## Configuration Options

### Basic Configuration

```go
client := httpclient.New(httpclient.Config{
    BaseURL: "https://api.example.com",
    Timeout: 30 * time.Second,
})
```

### With Retry

```go
client := httpclient.New(httpclient.Config{
    BaseURL:       "https://api.example.com",
    Timeout:       30 * time.Second,
    RetryCount:    3,                // Retry 3 times on failure
    RetryWaitTime: 1 * time.Second,  // Wait 1 second between retries
})
```

### With Default Headers

```go
client := httpclient.New(httpclient.Config{
    BaseURL: "https://api.example.com",
    Timeout: 30 * time.Second,
    DefaultHeaders: map[string]string{
        "User-Agent":    "MyApp/1.0",
        "Authorization": "Bearer default-token",
        "Accept":        "application/json",
    },
})
```

### Default Client

```go
// Uses default configuration
client := httpclient.NewDefault()
```

## Response Handling

### Check Status Code

```go
resp, err := client.Get("/users").Send()
if err != nil {
    log.Fatal(err)
}

if resp.StatusCode == 200 {
    fmt.Println("Success!")
}

// Or use helper
if resp.IsSuccess() {
    fmt.Println("2xx status!")
}
```

### Get Response as String

```go
resp, err := client.Get("/users").Send()
body := resp.String()
fmt.Println(body)
```

### Parse JSON Response

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

resp, err := client.Get("/users/1").Send()
if err != nil {
    log.Fatal(err)
}

var user User
if err := resp.JSON(&user); err != nil {
    log.Fatal(err)
}

fmt.Printf("User: %+v\n", user)
```

### Parse JSON Array

```go
var users []User
resp, err := client.Get("/users").Send()
if err != nil {
    log.Fatal(err)
}

if err := resp.JSON(&users); err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d users\n", len(users))
```

## Real-World Examples

### API Client

```go
type APIClient struct {
    client *httpclient.Client
    token  string
}

func NewAPIClient(baseURL, token string) *APIClient {
    return &APIClient{
        client: httpclient.New(httpclient.Config{
            BaseURL: baseURL,
            Timeout: 30 * time.Second,
            DefaultHeaders: map[string]string{
                "Authorization": "Bearer " + token,
                "Accept":        "application/json",
            },
        }),
        token: token,
    }
}

func (c *APIClient) GetUser(id int) (*User, error) {
    resp, err := c.client.Get(fmt.Sprintf("/users/%d", id)).Send()
    if err != nil {
        return nil, err
    }

    var user User
    if err := resp.JSON(&user); err != nil {
        return nil, err
    }

    return &user, nil
}

func (c *APIClient) CreateUser(user *User) error {
    resp, err := c.client.Post("/users").
        WithJSON(user).
        Send()
    if err != nil {
        return err
    }

    if !resp.IsSuccess() {
        return fmt.Errorf("failed to create user: %d", resp.StatusCode)
    }

    return nil
}
```

### Webhook Sender

```go
func SendWebhook(url string, data interface{}) error {
    client := httpclient.New(httpclient.Config{
        Timeout:    10 * time.Second,
        RetryCount: 3,
    })

    resp, err := client.Post(url).
        WithHeader("X-Webhook-Secret", "secret123").
        WithJSON(data).
        Send()
    if err != nil {
        return err
    }

    if !resp.IsSuccess() {
        return fmt.Errorf("webhook failed: %d", resp.StatusCode)
    }

    return nil
}
```

## Best Practices

### 1. Reuse Client Instances

```go
// Good - reuse client
var apiClient = httpclient.New(httpclient.Config{
    BaseURL: "https://api.example.com",
    Timeout: 30 * time.Second,
})

func GetUser(id int) (*User, error) {
    resp, err := apiClient.Get(fmt.Sprintf("/users/%d", id)).Send()
    // ...
}

// Avoid - creating new client every time
func GetUser(id int) (*User, error) {
    client := httpclient.New(httpclient.Config{...})
    resp, err := client.Get(fmt.Sprintf("/users/%d", id)).Send()
    // ...
}
```

### 2. Use Context for Timeouts

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := client.Get("/users").
    WithContext(ctx).
    Send()
```

### 3. Handle Errors Properly

```go
resp, err := client.Get("/users").Send()
if err != nil {
    log.Printf("Request failed: %v", err)
    return err
}

if !resp.IsSuccess() {
    log.Printf("API error: %d - %s", resp.StatusCode, resp.String())
    return fmt.Errorf("API returned %d", resp.StatusCode)
}
```

### 4. Use Retry for Unreliable APIs

```go
client := httpclient.New(httpclient.Config{
    BaseURL:       "https://unreliable-api.com",
    RetryCount:    3,
    RetryWaitTime: 2 * time.Second,
})
```

## Common Patterns

### Pagination

```go
func GetAllUsers(client *httpclient.Client) ([]User, error) {
    var allUsers []User
    page := 1

    for {
        resp, err := client.Get(fmt.Sprintf("/users?page=%d", page)).Send()
        if err != nil {
            return nil, err
        }

        var users []User
        if err := resp.JSON(&users); err != nil {
            return nil, err
        }

        if len(users) == 0 {
            break
        }

        allUsers = append(allUsers, users...)
        page++
    }

    return allUsers, nil
}
```

### Rate Limiting

```go
type RateLimitedClient struct {
    client  *httpclient.Client
    limiter *rate.Limiter
}

func (c *RateLimitedClient) Get(url string) (*httpclient.Response, error) {
    c.limiter.Wait(context.Background())
    return c.client.Get(url).Send()
}
```

## Installation

```bash
go get github.com/fatkulnurk/foundation/httpclient
```

## Dependencies

None - uses only Go standard library.

## See Also

- `example/main.go` - Complete examples with real API calls
- `request.go` - Full request builder API
- `response.go` - Response helper methods
