package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatkulnurk/foundation/storage"
)

func main() {
	ctx := context.Background()

	// Choose which example to run
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "local":
			runLocalStorageExample(ctx)
		case "s3":
			runS3StorageExample(ctx)
		case "list":
			runListExample(ctx)
		case "operations":
			runFileOperationsExample(ctx)
		default:
			fmt.Println("Usage: go run main.go [local|s3|list|operations]")
		}
	} else {
		runLocalStorageExample(ctx)
	}
}

// runLocalStorageExample demonstrates local storage operations
func runLocalStorageExample(ctx context.Context) {
	fmt.Println("=== Local Storage Example ===")
	fmt.Println()

	// Create local storage
	cfg := storage.LocalStorageConfig{
		BasePath:              "./storage",
		BaseURL:               "http://localhost:8080/storage",
		DefaultDirPermission:  0755,
		DefaultFilePermission: 0644,
	}

	store, err := storage.NewLocalStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to create local storage: %v", err)
	}

	// Example 1: Upload a text file
	fmt.Println("1. Uploading text file...")
	result, err := store.Upload(ctx, storage.UploadInput{
		FileName:   "documents/hello.txt",
		Content:    "Hello, Local Storage!",
		MimeType:   "text/plain",
		Visibility: storage.VisibilityPublic,
	})
	if err != nil {
		log.Printf("Upload failed: %v\n", err)
	} else {
		fmt.Printf("✓ Uploaded: %s (%s)\n", result.Path, result.SizeHuman)
	}
	fmt.Println()

	// Example 2: Upload a JSON file
	fmt.Println("2. Uploading JSON file...")
	jsonData := `{"name": "John", "age": 30, "city": "New York"}`
	result, err = store.Upload(ctx, storage.UploadInput{
		FileName:   "data/user.json",
		Content:    []byte(jsonData),
		MimeType:   "application/json",
		Visibility: storage.VisibilityPublic,
	})
	if err != nil {
		log.Printf("Upload failed: %v\n", err)
	} else {
		fmt.Printf("✓ Uploaded: %s (%s)\n", result.Path, result.SizeHuman)
	}
	fmt.Println()

	// Example 3: Check if file exists
	fmt.Println("3. Checking if file exists...")
	exists, err := store.Exists(ctx, "documents/hello.txt")
	if err != nil {
		log.Printf("Check failed: %v\n", err)
	} else {
		if exists {
			fmt.Println("✓ File exists!")
		} else {
			fmt.Println("✗ File not found")
		}
	}
	fmt.Println()

	// Example 4: Get file info
	fmt.Println("4. Getting file information...")
	file, err := store.File(ctx, "documents/hello.txt", nil)
	if err != nil {
		log.Printf("Get file info failed: %v\n", err)
	} else {
		fmt.Printf("✓ File Info:\n")
		fmt.Printf("  - Name: %s\n", file.Name)
		fmt.Printf("  - Size: %s\n", file.SizeHuman)
		fmt.Printf("  - Type: %s\n", file.MimeType)
		fmt.Printf("  - URL: %s\n", file.Url)
		fmt.Printf("  - Modified: %s\n", file.LastModified.Format("2006-01-02 15:04:05"))
	}
	fmt.Println()

	// Example 5: Read file content
	fmt.Println("5. Reading file content...")
	content, err := store.Get(ctx, "documents/hello.txt")
	if err != nil {
		log.Printf("Read failed: %v\n", err)
	} else {
		fmt.Printf("✓ Content: %s\n", string(content))
	}
	fmt.Println()

	fmt.Println("✅ Local storage example completed!")
	fmt.Println("Files are stored in: ./storage/")
}

// runS3StorageExample demonstrates S3 storage operations
func runS3StorageExample(ctx context.Context) {
	fmt.Println("=== S3 Storage Example ===")
	fmt.Println()

	// Load S3 configuration from environment
	cfgPtr := storage.LoadS3Config()
	cfg := *cfgPtr // Dereference pointer

	// Check if configuration is set
	if cfg.Bucket == "" {
		fmt.Println("⚠️  S3 configuration not set!")
		fmt.Println("Set these environment variables:")
		fmt.Println("  - STORAGE_S3_REGION")
		fmt.Println("  - STORAGE_S3_BUCKET")
		fmt.Println("  - STORAGE_S3_ACCESS_KEY")
		fmt.Println("  - STORAGE_S3_SECRET_KEY")
		return
	}

	// Create S3 client
	client, err := storage.NewS3Client(cfg)
	if err != nil {
		log.Fatalf("Failed to create S3 client: %v", err)
	}

	// Create S3 storage
	store := storage.NewS3Storage(client, cfg)

	// Example 1: Upload a file
	fmt.Println("1. Uploading file to S3...")
	result, err := store.Upload(ctx, storage.UploadInput{
		FileName:   "uploads/test.txt",
		Content:    "Hello, S3 Storage!",
		MimeType:   "text/plain",
		Visibility: storage.VisibilityPublic,
	})
	if err != nil {
		log.Printf("Upload failed: %v\n", err)
	} else {
		fmt.Printf("✓ Uploaded: %s (%s)\n", result.Path, result.SizeHuman)
	}
	fmt.Println()

	// Example 2: Upload a private file
	fmt.Println("2. Uploading private file...")
	result, err = store.Upload(ctx, storage.UploadInput{
		FileName:   "private/secret.txt",
		Content:    "This is a secret file",
		MimeType:   "text/plain",
		Visibility: storage.VisibilityPrivate,
	})
	if err != nil {
		log.Printf("Upload failed: %v\n", err)
	} else {
		fmt.Printf("✓ Uploaded private file: %s\n", result.Path)
	}
	fmt.Println()

	// Example 3: Get file with temporary URL
	fmt.Println("3. Generating temporary URL...")
	expiry := 15 * time.Minute
	file, err := store.File(ctx, "private/secret.txt", &expiry)
	if err != nil {
		log.Printf("Get file failed: %v\n", err)
	} else {
		fmt.Printf("✓ File Info:\n")
		fmt.Printf("  - Name: %s\n", file.Name)
		fmt.Printf("  - Visibility: %s\n", file.Visibility)
		fmt.Printf("  - Permanent URL: %s\n", file.Url)
		fmt.Printf("  - Temporary URL (15 min): %s\n", file.TempUrl)
	}
	fmt.Println()

	fmt.Println("✅ S3 storage example completed!")
}

// runListExample demonstrates listing files and directories
func runListExample(ctx context.Context) {
	fmt.Println("=== List Files Example ===")
	fmt.Println()

	// Create local storage
	cfg := storage.LocalStorageConfig{
		BasePath:              "./storage",
		BaseURL:               "http://localhost:8080/storage",
		DefaultDirPermission:  0755,
		DefaultFilePermission: 0644,
	}

	store, err := storage.NewLocalStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}

	// Create some test files
	fmt.Println("Creating test files...")
	testFiles := []struct {
		path    string
		content string
	}{
		{"photos/image1.jpg", "fake image 1"},
		{"photos/image2.jpg", "fake image 2"},
		{"photos/vacation/beach.jpg", "beach photo"},
		{"documents/report.pdf", "fake pdf"},
		{"documents/invoice.pdf", "fake invoice"},
	}

	for _, tf := range testFiles {
		store.Upload(ctx, storage.UploadInput{
			FileName:   tf.path,
			Content:    tf.content,
			MimeType:   "application/octet-stream",
			Visibility: storage.VisibilityPublic,
		})
	}
	fmt.Println("✓ Test files created")
	fmt.Println()

	// List files in photos directory
	fmt.Println("1. Files in photos/:")
	files, err := store.Files(ctx, "photos", nil)
	if err != nil {
		log.Printf("List failed: %v\n", err)
	} else {
		for _, file := range files {
			fmt.Printf("  - %s (%s)\n", file.Name, file.SizeHuman)
		}
	}
	fmt.Println()

	// List files in documents directory
	fmt.Println("2. Files in documents/:")
	files, err = store.Files(ctx, "documents", nil)
	if err != nil {
		log.Printf("List failed: %v\n", err)
	} else {
		for _, file := range files {
			fmt.Printf("  - %s (%s) - %s\n", file.Name, file.SizeHuman, file.MimeType)
		}
	}
	fmt.Println()

	// List subdirectories in photos
	fmt.Println("3. Subdirectories in photos/:")
	dirs, err := store.Directories(ctx, "photos")
	if err != nil {
		log.Printf("List directories failed: %v\n", err)
	} else {
		for _, dir := range dirs {
			fmt.Printf("  - %s/\n", dir)
		}
	}
	fmt.Println()

	fmt.Println("✅ List example completed!")
}

// runFileOperationsExample demonstrates copy, move, and delete operations
func runFileOperationsExample(ctx context.Context) {
	fmt.Println("=== File Operations Example ===")
	fmt.Println()

	// Create local storage
	cfg := storage.LocalStorageConfig{
		BasePath:              "./storage",
		BaseURL:               "http://localhost:8080/storage",
		DefaultDirPermission:  0755,
		DefaultFilePermission: 0644,
	}

	store, err := storage.NewLocalStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}

	// Create original file
	fmt.Println("1. Creating original file...")
	result, err := store.Upload(ctx, storage.UploadInput{
		FileName:   "original/document.txt",
		Content:    "This is the original document",
		MimeType:   "text/plain",
		Visibility: storage.VisibilityPublic,
	})
	if err != nil {
		log.Printf("Upload failed: %v\n", err)
		return
	}
	fmt.Printf("✓ Created: %s\n", result.Path)
	fmt.Println()

	// Copy file
	fmt.Println("2. Copying file...")
	err = store.Copy(ctx, "original/document.txt", "backup/document.txt")
	if err != nil {
		log.Printf("Copy failed: %v\n", err)
	} else {
		fmt.Println("✓ File copied to backup/document.txt")

		// Verify copy
		exists, _ := store.Exists(ctx, "backup/document.txt")
		if exists {
			fmt.Println("✓ Copy verified")
		}
	}
	fmt.Println()

	// Move file
	fmt.Println("3. Moving file...")
	err = store.Move(ctx, "original/document.txt", "archive/document.txt")
	if err != nil {
		log.Printf("Move failed: %v\n", err)
	} else {
		fmt.Println("✓ File moved to archive/document.txt")

		// Verify move
		existsOld, _ := store.Exists(ctx, "original/document.txt")
		existsNew, _ := store.Exists(ctx, "archive/document.txt")

		if !existsOld && existsNew {
			fmt.Println("✓ Move verified (old location empty, new location has file)")
		}
	}
	fmt.Println()

	// Delete file
	fmt.Println("4. Deleting file...")
	err = store.Delete(ctx, "backup/document.txt")
	if err != nil {
		log.Printf("Delete failed: %v\n", err)
	} else {
		fmt.Println("✓ File deleted from backup/")

		// Verify deletion
		exists, _ := store.Exists(ctx, "backup/document.txt")
		if !exists {
			fmt.Println("✓ Deletion verified")
		}
	}
	fmt.Println()

	fmt.Println("✅ File operations example completed!")
}
