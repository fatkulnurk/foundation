package container

import (
	"fmt"
	"sync"
)

// Container defines the interface for service locator pattern
type Container interface {
	// Set registers a service with the given name
	Set(name string, service interface{})

	// Get retrieves a service by name
	Get(name string) (interface{}, error)

	// MustGet retrieves a service and panics if not found
	MustGet(name string) interface{}

	// Has checks if a service exists
	Has(name string) bool
}

// Locator is a simple service locator
type Locator struct {
	mu       sync.RWMutex
	services map[string]interface{}
}

// Compile-time check to ensure Locator implements Container
var _ Container = (*Locator)(nil)

// New creates a new Locator
func New() *Locator {
	return &Locator{
		services: make(map[string]interface{}),
	}
}

// NewContainer creates a new Container instance
func NewContainer() Container {
	return New()
}

// Set registers a service with the given name
func (c *Locator) Set(name string, service interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.services[name] = service
}

// Get retrieves a service by name
func (c *Locator) Get(name string) (interface{}, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	service, exists := c.services[name]
	if !exists {
		return nil, fmt.Errorf("service '%s' not found", name)
	}

	return service, nil
}

// MustGet retrieves a service and panics if not found
func (c *Locator) MustGet(name string) interface{} {
	service, err := c.Get(name)
	if err != nil {
		panic(err)
	}
	return service
}

// Has checks if a service exists
func (c *Locator) Has(name string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.services[name]
	return exists
}
