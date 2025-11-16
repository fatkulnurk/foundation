package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrNotFound = errors.New("key not found")
)

type item struct {
	value     string
	expiresAt time.Time // zero value (time.Time{}) berarti tidak ada TTL
}

type LocalCache struct {
	cfg   *Config
	mu    sync.RWMutex
	items map[string]item
}

func NewLocalCache(cfg *Config) Cache {
	return &LocalCache{
		cfg:   cfg,
		items: make(map[string]item),
	}
}

// Set menyimpan key dengan TTL dalam detik.
// ttlSeconds <= 0 => tidak ada expiry.
func (c *LocalCache) Set(ctx context.Context, key string, value any, ttlSeconds int) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	key = c.cfg.Prefix + key

	stringValue, err := toString(value)
	if err != nil {
		return err
	}

	var expiresAt time.Time
	if ttlSeconds > 0 {
		expiresAt = time.Now().Add(time.Duration(ttlSeconds) * time.Second)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item{
		value:     stringValue,
		expiresAt: expiresAt,
	}

	return nil
}

func (c *LocalCache) Get(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	key = c.cfg.Prefix + key

	c.mu.RLock()
	it, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return "", ErrNotFound
	}

	// Cek expired
	if isExpired(it.expiresAt) {
		// Hapus kalau expired
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return "", ErrNotFound
	}

	return it.value, nil
}

func (c *LocalCache) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	key = c.cfg.Prefix + key

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
	return nil
}

func (c *LocalCache) Has(ctx context.Context, key string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}

	key = c.cfg.Prefix + key

	c.mu.RLock()
	it, ok := c.items[key]
	c.mu.RUnlock()

	if !ok {
		return false, nil
	}

	if isExpired(it.expiresAt) {
		// Bersihkan kalau sudah kadaluwarsa
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return false, nil
	}

	return true, nil
}

// Helper: cek apakah sudah expired
func isExpired(expiresAt time.Time) bool {
	if expiresAt.IsZero() {
		return false // tidak ada TTL
	}
	return time.Now().After(expiresAt)
}

// Helper: convert any -> string
func toString(v any) (string, error) {
	switch t := v.(type) {
	case string:
		return t, nil
	case fmt.Stringer:
		return t.String(), nil
	default:
		// fallback pakai fmt.Sprint
		return fmt.Sprint(v), nil
	}
}
