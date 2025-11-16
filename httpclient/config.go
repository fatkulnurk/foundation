package httpclient

import "time"

// Config untuk HTTP client
type Config struct {
	Timeout        time.Duration
	RetryCount     int
	RetryWaitTime  time.Duration
	BaseURL        string
	DefaultHeaders map[string]string
}
