package queue

import (
	"context"
	"fmt"
	"time"

	"github.com/fatkulnurk/foundation/logging"
)

// LoggingMiddleware logs task execution
func LoggingMiddleware(taskType string) MiddlewareFunc {
	return func(next Handler) Handler {
		return func(ctx context.Context, payload []byte) error {
			start := time.Now()
			logging.Info(ctx, fmt.Sprintf("[%s] Task started", taskType))

			err := next(ctx, payload)

			duration := time.Since(start)
			if err != nil {
				logging.Error(ctx, fmt.Sprintf("[%s] Task failed after %v: %v", taskType, duration, err))
			} else {
				logging.Info(ctx, fmt.Sprintf("[%s] Task completed in %v", taskType, duration))
			}

			return err
		}
	}
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(taskType string) MiddlewareFunc {
	return func(next Handler) Handler {
		return func(ctx context.Context, payload []byte) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("[%s] panic recovered: %v", taskType, r)
					logging.Error(ctx, err.Error())
				}
			}()

			return next(ctx, payload)
		}
	}
}

// RetryLoggingMiddleware logs retry attempts
func RetryLoggingMiddleware(taskType string) MiddlewareFunc {
	return func(next Handler) Handler {
		return func(ctx context.Context, payload []byte) error {
			// Check if this is a retry (asynq adds retry count to context)
			// For now, just execute the handler
			err := next(ctx, payload)

			if err != nil {
				logging.Info(ctx, fmt.Sprintf("[%s] Task will be retried: %v", taskType, err))
			}

			return err
		}
	}
}

// TimeoutMiddleware adds a timeout to task execution
func TimeoutMiddleware(timeout time.Duration) MiddlewareFunc {
	return func(next Handler) Handler {
		return func(ctx context.Context, payload []byte) error {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			done := make(chan error, 1)
			go func() {
				done <- next(ctx, payload)
			}()

			select {
			case err := <-done:
				return err
			case <-ctx.Done():
				return fmt.Errorf("task timeout after %v: %w", timeout, ctx.Err())
			}
		}
	}
}

// MetricsMiddleware tracks task metrics (placeholder for actual metrics implementation)
func MetricsMiddleware(taskType string) MiddlewareFunc {
	return func(next Handler) Handler {
		return func(ctx context.Context, payload []byte) error {
			start := time.Now()
			err := next(ctx, payload)
			duration := time.Since(start)

			// Here you would send metrics to your metrics system
			// For example: prometheus, datadog, etc.
			_ = duration // Use duration for metrics

			if err != nil {
				// Increment error counter
				logging.Debug(ctx, fmt.Sprintf("[%s] Task error metric recorded", taskType))
			} else {
				// Increment success counter
				logging.Debug(ctx, fmt.Sprintf("[%s] Task success metric recorded", taskType))
			}

			return err
		}
	}
}

// ChainMiddleware chains multiple middleware functions
func ChainMiddleware(middleware ...MiddlewareFunc) MiddlewareFunc {
	return func(next Handler) Handler {
		for i := len(middleware) - 1; i >= 0; i-- {
			next = middleware[i](next)
		}
		return next
	}
}
