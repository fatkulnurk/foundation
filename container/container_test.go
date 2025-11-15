package container

import "testing"

type Database struct {
	Name string
}

type Logger struct {
	Level string
}

func TestNew(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("Expected non-nil container")
	}
}

func TestNewLocator(t *testing.T) {
	var loc Container = NewContainer()
	if loc == nil {
		t.Fatal("Expected non-nil locator")
	}
}

func TestLocatorInterface(t *testing.T) {
	// Test that Locator implements Container
	var _ Container = (*Locator)(nil)

	// Test using interface
	var loc Container = NewContainer()
	db := &Database{Name: "test"}

	loc.Set("database", db)

	result, err := loc.Get("database")
	if err != nil {
		t.Fatalf("Failed to get service: %v", err)
	}

	retrieved := result.(*Database)
	if retrieved.Name != "test" {
		t.Error("Expected correct database name")
	}

	if !loc.Has("database") {
		t.Error("Expected Has to return true")
	}

	mustResult := loc.MustGet("database").(*Database)
	if mustResult.Name != "test" {
		t.Error("Expected correct database name from MustGet")
	}
}

func TestSetAndGet(t *testing.T) {
	c := New()
	db := &Database{Name: "test"}

	c.Set("database", db)

	result, err := c.Get("database")
	if err != nil {
		t.Fatalf("Failed to get service: %v", err)
	}

	retrieved := result.(*Database)
	if retrieved != db {
		t.Error("Expected same instance")
	}
}

func TestGetNotFound(t *testing.T) {
	c := New()

	_, err := c.Get("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent service")
	}
}

func TestMustGet(t *testing.T) {
	c := New()
	logger := &Logger{Level: "info"}

	c.Set("logger", logger)

	result := c.MustGet("logger")
	retrieved := result.(*Logger)

	if retrieved != logger {
		t.Error("Expected same instance")
	}
}

func TestMustGetPanic(t *testing.T) {
	c := New()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for non-existent service")
		}
	}()

	c.MustGet("nonexistent")
}

func TestHas(t *testing.T) {
	c := New()

	c.Set("logger", &Logger{})

	if !c.Has("logger") {
		t.Error("Expected Has to return true")
	}

	if c.Has("nonexistent") {
		t.Error("Expected Has to return false")
	}
}

func TestConcurrentAccess(t *testing.T) {
	c := New()
	db := &Database{Name: "test"}
	c.Set("database", db)

	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			_, err := c.Get("database")
			if err != nil {
				t.Errorf("Failed to get service: %v", err)
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
