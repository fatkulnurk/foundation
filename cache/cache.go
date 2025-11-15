package cache

import (
	"context"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, ttlSeconds int) error
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Has(ctx context.Context, key string) (bool, error)
}
