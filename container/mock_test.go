package container

import (
	"fmt"
	"testing"
)

// MockLocator is a simple mock implementation of Container interface
type MockLocator struct {
	services map[string]interface{}
	getCalls int
	setCalls int
}

func NewMockLocator() *MockLocator {
	return &MockLocator{
		services: make(map[string]interface{}),
	}
}

func (m *MockLocator) Set(name string, service interface{}) {
	m.setCalls++
	m.services[name] = service
}

func (m *MockLocator) Get(name string) (interface{}, error) {
	m.getCalls++
	if s, ok := m.services[name]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("service '%s' not found", name)
}

func (m *MockLocator) MustGet(name string) interface{} {
	service, err := m.Get(name)
	if err != nil {
		panic(err)
	}
	return service
}

func (m *MockLocator) Has(name string) bool {
	_, ok := m.services[name]
	return ok
}

// Example service that depends on Container
type EmailService struct {
	logger *Logger
}

func NewEmailService(loc Container) *EmailService {
	return &EmailService{
		logger: loc.MustGet("logger").(*Logger),
	}
}

func (e *EmailService) SendEmail(to, subject string) error {
	// Use logger level for logging
	fmt.Printf("[%s] Sending email to %s: %s\n", e.logger.Level, to, subject)
	return nil
}

// Test using mock
func TestEmailServiceWithMock(t *testing.T) {
	// Create mock locator
	mock := NewMockLocator()

	// Setup mock data
	logger := &Logger{Level: "INFO"}
	mock.Set("logger", logger)

	// Create service with mock
	emailService := NewEmailService(mock)

	// Test service
	err := emailService.SendEmail("test@example.com", "Test Subject")
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}

	// Verify mock was called
	if mock.getCalls == 0 {
		t.Error("Expected Get to be called")
	}

	if mock.setCalls == 0 {
		t.Error("Expected Set to be called")
	}
}

// Test using real container
func TestEmailServiceWithRealContainer(t *testing.T) {
	// Create real container
	var loc Container = NewContainer()

	// Setup real data
	logger := &Logger{Level: "DEBUG"}
	loc.Set("logger", logger)

	// Create service with real container
	emailService := NewEmailService(loc)

	// Test service
	err := emailService.SendEmail("real@example.com", "Real Subject")
	if err != nil {
		t.Fatalf("Failed to send email: %v", err)
	}
}

// Test that mock implements Container interface
func TestMockImplementsLocator(t *testing.T) {
	var _ Container = (*MockLocator)(nil)
	var _ Container = NewMockLocator()
}

// Benchmark comparison
func BenchmarkRealContainer(b *testing.B) {
	loc := NewContainer()
	loc.Set("test", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = loc.Get("test")
	}
}

func BenchmarkMockLocator(b *testing.B) {
	mock := NewMockLocator()
	mock.Set("test", "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mock.Get("test")
	}
}
