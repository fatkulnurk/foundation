package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/fatkulnurk/foundation/container"
)

// Example services
type Database struct {
	ConnectionString string
}

type Logger struct {
	Level string
}

type UserService struct {
	DB     *Database
	Logger *Logger
}

func main() {
	fmt.Println("=== Container Example ===\n")

	// Example 1: Create container
	fmt.Println("1. Creating container...")
	c := container.New()
	fmt.Println("✓ Container created\n")

	// Example 2: Register services
	fmt.Println("2. Registering services...")

	// Register database
	db := &Database{ConnectionString: "localhost:5432"}
	c.Set("database", db)
	fmt.Println("✓ Registered: database")

	// Register logger
	logger := &Logger{Level: "info"}
	c.Set("logger", logger)
	fmt.Println("✓ Registered: logger")

	// Register user service
	userService := &UserService{
		DB:     db,
		Logger: logger,
	}
	c.Set("userService", userService)
	fmt.Println("✓ Registered: userService\n")

	// Example 3: Retrieve services with Get
	fmt.Println("3. Retrieving services with Get()...")

	dbInterface, err := c.Get("database")
	if err != nil {
		log.Fatal(err)
	}
	retrievedDB := dbInterface.(*Database)
	fmt.Printf("✓ Retrieved database: %s\n", retrievedDB.ConnectionString)

	loggerInterface, err := c.Get("logger")
	if err != nil {
		log.Fatal(err)
	}
	retrievedLogger := loggerInterface.(*Logger)
	fmt.Printf("✓ Retrieved logger: level=%s\n", retrievedLogger.Level)

	// Example 4: Retrieve with MustGet (panics if not found)
	fmt.Println("\n4. Retrieving services with MustGet()...")
	userSvc := c.MustGet("userService").(*UserService)
	fmt.Printf("✓ Retrieved userService: DB=%s, Logger=%s\n",
		userSvc.DB.ConnectionString, userSvc.Logger.Level)

	// Example 5: Check if service exists
	fmt.Println("\n5. Checking service existence...")
	if c.Has("database") {
		fmt.Println("✓ database exists")
	}

	if c.Has("cache") {
		fmt.Println("✓ cache exists")
	} else {
		fmt.Println("✗ cache does not exist")
	}

	// Example 6: Error handling
	fmt.Println("\n6. Error handling...")
	_, err = c.Get("nonexistent")
	if err != nil {
		fmt.Printf("✓ Expected error: %v\n", err)
	}

	// Example 7: Using interface type
	fmt.Println("\n7. Using interface types...")
	c.Set("config", map[string]string{
		"app_name": "MyApp",
		"version":  "1.0.0",
	})

	configInterface, _ := c.Get("config")
	config := configInterface.(map[string]string)
	fmt.Printf("✓ Config: %v\n", config)

	// Example 8: Storing different types
	fmt.Println("\n8. Storing different types...")

	// String
	c.Set("appName", "Foundation")
	appName := c.MustGet("appName").(string)
	fmt.Printf("✓ String: %s\n", appName)

	// Integer
	c.Set("port", 8080)
	port := c.MustGet("port").(int)
	fmt.Printf("✓ Integer: %d\n", port)

	// Boolean
	c.Set("debug", true)
	debug := c.MustGet("debug").(bool)
	fmt.Printf("✓ Boolean: %v\n", debug)

	// Slice
	c.Set("allowedHosts", []string{"localhost", "127.0.0.1"})
	hosts := c.MustGet("allowedHosts").([]string)
	fmt.Printf("✓ Slice: %v\n", hosts)

	// Example 9: Real-world scenario - Dependency Injection
	fmt.Println("\n9. Real-world scenario - Dependency Injection...")

	// Setup dependencies
	appContainer := container.New()

	// Register core services
	appContainer.Set("db", &Database{ConnectionString: "prod:5432"})
	appContainer.Set("logger", &Logger{Level: "error"})

	// Create service that depends on other services
	userSvc = &UserService{
		DB:     appContainer.MustGet("db").(*Database),
		Logger: appContainer.MustGet("logger").(*Logger),
	}
	appContainer.Set("userService", userSvc)

	// Use the service
	fmt.Println("✓ Application initialized with dependencies")
	fmt.Printf("  - UserService connected to: %s\n", userSvc.DB.ConnectionString)
	fmt.Printf("  - UserService logging at: %s level\n", userSvc.Logger.Level)

	// Example 10: Using with standard library types
	fmt.Println("\n10. Using with standard library types...")

	// Simulate database connection
	var sqlDB *sql.DB // In real app, this would be sql.Open(...)
	appContainer.Set("sqlDB", sqlDB)

	if appContainer.Has("sqlDB") {
		fmt.Println("✓ SQL database registered")
	}

	fmt.Println("\n✅ All examples completed successfully!")
}
