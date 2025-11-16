package httpclient

import (
	"net/http"
	"time"
)

// HttpClient interface untuk HTTP client operations
type HttpClient interface {
	Get(url string) *Request
	Post(url string) *Request
	Put(url string) *Request
	Patch(url string) *Request
	Delete(url string) *Request
	NewRequest() *Request
}

// ContentType constants
const (
	ContentTypeJSON           = "application/json"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
	ContentTypeMultipartForm  = "multipart/form-data"
	ContentTypeXML            = "application/xml"
	ContentTypeText           = "text/plain"
)

// client implements HttpClient interface
type client struct {
	httpClient *http.Client
	config     Config
}

// New creates a new HTTP client with config
func New(config Config) HttpClient {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.RetryWaitTime == 0 {
		config.RetryWaitTime = 1 * time.Second
	}

	return &client{
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		config: config,
	}
}

// NewDefault creates HTTP client with default config
func NewDefault() HttpClient {
	return New(Config{
		Timeout:    30 * time.Second,
		RetryCount: 0,
	})
}

// Get creates a GET request builder
func (c *client) Get(url string) *Request {
	return c.NewRequest().WithMethod(http.MethodGet).WithURL(url)
}

// Post creates a POST request builder
func (c *client) Post(url string) *Request {
	return c.NewRequest().WithMethod(http.MethodPost).WithURL(url)
}

// Put creates a PUT request builder
func (c *client) Put(url string) *Request {
	return c.NewRequest().WithMethod(http.MethodPut).WithURL(url)
}

// Patch creates a PATCH request builder
func (c *client) Patch(url string) *Request {
	return c.NewRequest().WithMethod(http.MethodPatch).WithURL(url)
}

// Delete creates a DELETE request builder
func (c *client) Delete(url string) *Request {
	return c.NewRequest().WithMethod(http.MethodDelete).WithURL(url)
}
