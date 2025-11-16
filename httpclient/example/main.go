package main

import (
	"context"
	"os"
	"time"

	"github.com/fatkulnurk/foundation/httpclient"
	"github.com/google/uuid"
)

func main() {
	client := httpclient.NewDefault()

	// 1. GET dengan headers
	resp, err := client.Get("/users").
		WithHeader("Authorization", "Bearer token").
		WithHeader("X-Request-ID", "abc123").
		Send()

	// 2. POST dengan JSON
	resp, err = client.Post("/users").
		WithJSON(map[string]string{
			"name":  "John",
			"email": "john@example.com",
		}).
		Send()

	// 3. POST dengan Form URL Encoded
	resp, err = client.Post("/login").
		WithFormURLEncoded(map[string]string{
			"username": "john",
			"password": "secret",
		}).
		Send()

	// 4. POST dengan Multipart Form (file upload)
	fileData, _ := os.ReadFile("photo.jpg")
	resp, err = client.Post("/upload").
		WithMultipartForm(
			map[string]string{"title": "My Photo"},
			map[string][]byte{"photo.jpg": fileData},
		).
		Send()

	// 5. PUT dengan multiple headers
	resp, err = client.Put("/users/123").
		WithHeaders(map[string]string{
			"Authorization": "Bearer token",
			"X-Request-ID":  "abc123",
			"Content-Type":  "application/json",
		}).
		WithJSON(map[string]interface{}{
			"name": "John Updated",
			"age":  30,
		}).
		Send()

	// 6. DELETE dengan context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err = client.Delete("/users/123").
		WithContext(ctx).
		WithHeader("Authorization", "Bearer token").
		Send()

	// 7. Custom request dengan raw body
	resp, err = client.NewRequest().
		WithMethod("PATCH").
		WithURL("/api/data").
		WithRaw([]byte("custom data"), "application/octet-stream").
		Send()

	// 8. Plain text request
	resp, err = client.Post("/webhook").
		WithHeader("X-Webhook-Secret", "secret123").
		WithText("Event payload text").
		Send()

	// 9. Complex chaining
	resp, err = client.Post("/api/users").
		WithContext(ctx).
		WithHeaders(map[string]string{
			"Authorization": "Bearer token",
			"X-API-Version": "v2",
		}).
		WithHeader("X-Request-ID", uuid.New().String()).
		WithJSON(userData).
		Send()
}
