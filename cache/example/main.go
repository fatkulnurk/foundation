package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fatkulnurk/foundation/cache"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	// Example 1: Using Redis Cache
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// Test Redis connection
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Create cache config
	cfg := &cache.Config{
		Prefix: "myapp:",
	}

	// Create cache instance
	c := cache.NewRedisCache(cfg, redisClient)

	// Example 2: Set values with TTL
	fmt.Println("=== Setting cache values ===")

	err := c.Set(ctx, "user:1", "John Doe", 60) // 60 seconds TTL
	if err != nil {
		log.Printf("Error setting cache: %v", err)
	} else {
		fmt.Println("✓ Set user:1 = John Doe (TTL: 60s)")
	}

	err = c.Set(ctx, "user:2", "Jane Smith", 120) // 120 seconds TTL
	if err != nil {
		log.Printf("Error setting cache: %v", err)
	} else {
		fmt.Println("✓ Set user:2 = Jane Smith (TTL: 120s)")
	}

	// Example 3: Get values
	fmt.Println("\n=== Getting cache values ===")

	value, err := c.Get(ctx, "user:1")
	if err != nil {
		log.Printf("Error getting cache: %v", err)
	} else {
		fmt.Printf("✓ Get user:1 = %s\n", value)
	}

	value, err = c.Get(ctx, "user:2")
	if err != nil {
		log.Printf("Error getting cache: %v", err)
	} else {
		fmt.Printf("✓ Get user:2 = %s\n", value)
	}

	// Example 4: Check if key exists
	fmt.Println("\n=== Checking key existence ===")

	exists, err := c.Has(ctx, "user:1")
	if err != nil {
		log.Printf("Error checking cache: %v", err)
	} else {
		fmt.Printf("✓ user:1 exists: %v\n", exists)
	}

	exists, err = c.Has(ctx, "user:999")
	if err != nil {
		log.Printf("Error checking cache: %v", err)
	} else {
		fmt.Printf("✓ user:999 exists: %v\n", exists)
	}

	// Example 5: Delete a key
	fmt.Println("\n=== Deleting cache values ===")

	err = c.Delete(ctx, "user:1")
	if err != nil {
		log.Printf("Error deleting cache: %v", err)
	} else {
		fmt.Println("✓ Deleted user:1")
	}

	// Verify deletion
	exists, err = c.Has(ctx, "user:1")
	if err != nil {
		log.Printf("Error checking cache: %v", err)
	} else {
		fmt.Printf("✓ user:1 exists after deletion: %v\n", exists)
	}

	// Example 6: Using Local Cache (in-memory)
	fmt.Println("\n=== Using Local Cache ===")

	localCache := cache.NewLocalCache(cfg)

	err = localCache.Set(ctx, "session:abc123", "user_data", 300)
	if err != nil {
		log.Printf("Error setting local cache: %v", err)
	} else {
		fmt.Println("✓ Set session:abc123 in local cache")
	}

	value, err = localCache.Get(ctx, "session:abc123")
	if err != nil {
		log.Printf("Error getting local cache: %v", err)
	} else {
		fmt.Printf("✓ Get session:abc123 = %s\n", value)
	}

	// Example 7: Caching complex data (JSON)
	fmt.Println("\n=== Caching complex data ===")

	userData := `{"id":123,"name":"John Doe","email":"john@example.com"}`
	err = c.Set(ctx, "user:json:123", userData, 300)
	if err != nil {
		log.Printf("Error setting JSON cache: %v", err)
	} else {
		fmt.Println("✓ Set user JSON data")
	}

	jsonValue, err := c.Get(ctx, "user:json:123")
	if err != nil {
		log.Printf("Error getting JSON cache: %v", err)
	} else {
		fmt.Printf("✓ Get user JSON data: %s\n", jsonValue)
	}

	// Example 8: Demonstrating TTL expiration
	fmt.Println("\n=== Demonstrating TTL expiration ===")

	err = c.Set(ctx, "temp:key", "temporary value", 3) // 3 seconds TTL
	if err != nil {
		log.Printf("Error setting temp cache: %v", err)
	} else {
		fmt.Println("✓ Set temp:key with 3 seconds TTL")
	}

	// Check immediately
	exists, _ = c.Has(ctx, "temp:key")
	fmt.Printf("✓ temp:key exists immediately: %v\n", exists)

	// Wait 4 seconds
	fmt.Println("⏳ Waiting 4 seconds...")
	time.Sleep(4 * time.Second)

	// Check after expiration
	exists, _ = c.Has(ctx, "temp:key")
	fmt.Printf("✓ temp:key exists after 4 seconds: %v\n", exists)

	// Cleanup
	fmt.Println("\n=== Cleanup ===")
	c.Delete(ctx, "user:2")
	c.Delete(ctx, "user:json:123")
	fmt.Println("✓ Cleanup completed")

	fmt.Println("\n✅ All examples completed successfully!")
}
