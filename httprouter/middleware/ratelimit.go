package middleware

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RateLimitConfig struct {
	Requests int           // max request
	Window   time.Duration // dalam durasi ini
}

func NewRateLimitMiddleware(cfg RateLimitConfig) func(http.Handler) http.Handler {
	if cfg.Requests <= 0 {
		cfg.Requests = 100
	}
	if cfg.Window <= 0 {
		cfg.Window = time.Minute
	}

	type client struct {
		count       int
		windowStart time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)

			now := time.Now()

			mu.Lock()
			c, ok := clients[ip]
			if !ok {
				c = &client{count: 0, windowStart: now}
				clients[ip] = c
			}

			// reset kalau sudah lewat window
			if now.Sub(c.windowStart) > cfg.Window {
				c.windowStart = now
				c.count = 0
			}

			c.count++
			currentCount := c.count
			windowStart := c.windowStart
			mu.Unlock()

			if currentCount > cfg.Requests {
				retryAfter := cfg.Window - now.Sub(windowStart)
				if retryAfter < 0 {
					retryAfter = 0
				}
				w.Header().Set("Retry-After", strconv.Itoa(int(retryAfter.Seconds())))
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// clientIP mencoba ambil IP dari header proxy, lalu fallback RemoteAddr
func clientIP(r *http.Request) string {
	// prioritas: X-Real-IP
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	// lalu: X-Forwarded-For (ambil IP pertama)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	// fallback RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
