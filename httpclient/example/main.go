package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fatkulnurk/foundation/httpclient"
	"github.com/google/uuid"
)

func main() {
	fmt.Println("=== HTTP Client Example ===\n")

	// Example 1: Simple GET request
	fmt.Println("1. Simple GET request")
	client := httpclient.New(httpclient.Config{
		BaseURL: "https://jsonplaceholder.typicode.com",
		Timeout: 30 * time.Second,
	})

	resp, err := client.Get("/posts/1").Send()
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Status: %d\n", resp.StatusCode)
		fmt.Printf("✓ Body: %s\n\n", resp.String()[:100]+"...")
	}

	// Example 2: GET with headers
	fmt.Println("2. GET with custom headers")
	resp, err = client.Get("/posts/1").
		WithHeader("Authorization", "Bearer token123").
		WithHeader("X-Request-ID", uuid.New().String()).
		Send()
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Status: %d\n\n", resp.StatusCode)
	}

	// Example 3: POST with JSON
	fmt.Println("3. POST with JSON body")
	resp, err = client.Post("/posts").
		WithJSON(map[string]interface{}{
			"title":  "Hello World",
			"body":   "This is a test post",
			"userId": 1,
		}).
		Send()
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Status: %d\n", resp.StatusCode)
		fmt.Printf("✓ Response: %s\n\n", resp.String()[:80]+"...")
	}

	// Example 4: POST with Form URL Encoded
	fmt.Println("4. POST with Form URL Encoded")
	resp, err = client.Post("/posts").
		WithFormURLEncoded(map[string]string{
			"title":  "Form Post",
			"body":   "Posted via form",
			"userId": "1",
		}).
		Send()
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Status: %d\n\n", resp.StatusCode)
	}

	// Example 5: PUT request
	fmt.Println("5. PUT request to update resource")
	resp, err = client.Put("/posts/1").
		WithJSON(map[string]interface{}{
			"id":     1,
			"title":  "Updated Title",
			"body":   "Updated body content",
			"userId": 1,
		}).
		Send()
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Status: %d\n\n", resp.StatusCode)
	}

	// Example 6: PATCH request
	fmt.Println("6. PATCH request for partial update")
	resp, err = client.Patch("/posts/1").
		WithJSON(map[string]interface{}{
			"title": "Patched Title",
		}).
		Send()
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Status: %d\n\n", resp.StatusCode)
	}

	// Example 7: DELETE request
	fmt.Println("7. DELETE request")
	resp, err = client.Delete("/posts/1").
		WithHeader("Authorization", "Bearer token123").
		Send()
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Status: %d\n\n", resp.StatusCode)
	}

	// Example 8: Request with context timeout
	fmt.Println("8. Request with context timeout")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err = client.Get("/posts").
		WithContext(ctx).
		Send()
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Status: %d\n", resp.StatusCode)
		fmt.Printf("✓ Retrieved posts\n\n")
	}

	// Example 9: Multiple headers
	fmt.Println("9. Request with multiple headers")
	resp, err = client.Get("/posts/1").
		WithHeaders(map[string]string{
			"Authorization": "Bearer token123",
			"X-API-Version": "v2",
			"X-Request-ID":  uuid.New().String(),
			"Accept":        "application/json",
		}).
		Send()
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Status: %d\n\n", resp.StatusCode)
	}

	// Example 10: Custom request method
	fmt.Println("10. Custom request method")
	resp, err = client.NewRequest().
		WithMethod("GET").
		WithURL("/posts/1").
		WithHeader("User-Agent", "CustomClient/1.0").
		Send()
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Status: %d\n\n", resp.StatusCode)
	}

	// Example 11: Client with retry configuration
	fmt.Println("11. Client with retry configuration")
	retryClient := httpclient.New(httpclient.Config{
		BaseURL:       "https://jsonplaceholder.typicode.com",
		Timeout:       10 * time.Second,
		RetryCount:    3,
		RetryWaitTime: 1 * time.Second,
	})

	resp, err = retryClient.Get("/posts/1").Send()
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Status: %d (with retry support)\n\n", resp.StatusCode)
	}

	// Example 12: Default headers
	fmt.Println("12. Client with default headers")
	clientWithDefaults := httpclient.New(httpclient.Config{
		BaseURL: "https://jsonplaceholder.typicode.com",
		Timeout: 30 * time.Second,
		DefaultHeaders: map[string]string{
			"User-Agent":    "MyApp/1.0",
			"Accept":        "application/json",
			"Authorization": "Bearer default-token",
		},
	})

	resp, err = clientWithDefaults.Get("/posts/1").Send()
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("✓ Status: %d (with default headers)\n\n", resp.StatusCode)
	}

	fmt.Println("✅ All examples completed successfully!")
}
