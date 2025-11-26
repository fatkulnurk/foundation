# Storage Package

A simple, flexible file storage interface for Go applications with Local filesystem and AWS S3 implementations.

## Table of Contents

- [What is Storage?](#what-is-storage)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Interface](#interface)
- [Configuration](#configuration)
- [Local Storage](#local-storage)
- [S3 Storage](#s3-storage)
- [Complete Examples](#complete-examples)
- [Best Practices](#best-practices)

---

## What is Storage?

The storage package provides a **unified interface** for file operations, whether you're storing files locally or in the cloud (S3). You can easily switch between storage providers without changing your code.

**Simple Analogy:**
- **Storage Interface** = Universal remote control
- **Local Storage** = Storing files in your computer's hard drive
- **S3 Storage** = Storing files in the cloud (like Dropbox)
- **Visibility** = Public (anyone can access) or Private (only you can access)

---

## Features

-  **Unified Interface** - Same code works for local and S3 storage
-  **Multiple Implementations** - Local filesystem and AWS S3
-  **File Operations** - Upload, download, delete, copy, move
-  **Directory Operations** - List files and directories
-  **Visibility Control** - Public or private files
-  **Temporary URLs** - Generate time-limited access URLs (S3)
-  **MIME Type Detection** - Automatic content type detection
-  **Human-Readable Sizes** - File sizes in KB, MB, GB format
-  **Thread-Safe** - Safe for concurrent use
-  **Easy Configuration** - Environment variables or code

---

## Installation

```bash
go get github.com/fatkulnurk/foundation/storage
```

**Dependencies:**
- Go 1.25 or higher
- For S3: AWS SDK for Go v2
- For Local: Standard library only

---

## Quick Start

### Local Storage

```go
package main

import (
    "context"
    "github.com/fatkulnurk/foundation/storage"
)

func main() {
    // Create local storage
    cfg := storage.LocalStorageConfig{
        BasePath: "./uploads",
        BaseURL:  "http://localhost:8080/uploads",
        DefaultDirPermission:  0755,
        DefaultFilePermission: 0644,
    }
    
    store, _ := storage.NewLocalStorage(cfg)
    
    // Upload a file
    result, _ := store.Upload(context.Background(), storage.UploadInput{
        FileName:   "photos/image.jpg",
        Content:    []byte("image data here"),
        MimeType:   "image/jpeg",
        Visibility: storage.VisibilityPublic,
    })
    
    println("File uploaded:", result.Path)
}
```

### S3 Storage

```go
package main

import (
    "context"
    "github.com/fatkulnurk/foundation/storage"
)

func main() {
    // Create S3 storage
    cfg := storage.S3Config{
        Region:    "us-east-1",
        Bucket:    "my-bucket",
        AccessKey: "YOUR_ACCESS_KEY",
        SecretKey: "YOUR_SECRET_KEY",
    }
    
    client, _ := storage.NewS3Client(cfg)
    store := storage.NewS3Storage(client, cfg)
    
    // Upload a file
    result, _ := store.Upload(context.Background(), storage.UploadInput{
        FileName:   "documents/file.pdf",
        Content:    []byte("pdf data here"),
        MimeType:   "application/pdf",
        Visibility: storage.VisibilityPrivate,
    })
    
    println("File uploaded:", result.Path)
}
```

---

## Interface

The `Storage` interface defines all file operations:

```go
type Storage interface {
    // Upload stores a file
    Upload(ctx context.Context, input UploadInput) (*UploadOutput, error)
    
    // Delete removes a file
    Delete(ctx context.Context, path string) error
    
    // Copy creates a copy of a file
    Copy(ctx context.Context, sourcePath, destinationPath string) error
    
    // Move moves a file to a new location
    Move(ctx context.Context, sourcePath, destinationPath string) error
    
    // Get retrieves file content
    Get(ctx context.Context, path string) ([]byte, error)
    
    // File gets information about a single file
    File(ctx context.Context, path string, expiryTempUrl *time.Duration) (*FileStorage, error)
    
    // Files lists files in a directory
    Files(ctx context.Context, dir string, expiryTempUrl *time.Duration) ([]FileStorage, error)
    
    // Directories lists subdirectories
    Directories(ctx context.Context, dir string) ([]string, error)
    
    // Exists checks if a file exists
    Exists(ctx context.Context, path string) (bool, error)
}
```

### Data Structures

#### UploadInput

```go
type UploadInput struct {
    FileName   string      // Path where file will be stored
    Content    any         // File content ([]byte, string, or io.Reader)
    MimeType   string      // MIME type (e.g., "image/jpeg")
    Visibility Visibility  // Public or Private
}
```

#### UploadOutput

```go
type UploadOutput struct {
    Name      string  // File name
    Path      string  // Full path
    Size      int64   // Size in bytes
    SizeHuman string  // Human-readable size (e.g., "1.5 MB")
}
```

#### FileStorage

```go
type FileStorage struct {
    Name         string      // File name
    Path         string      // Full path
    Size         int64       // Size in bytes
    SizeHuman    string      // Human-readable size
    MimeType     string      // MIME type
    LastModified time.Time   // Last modification time
    Visibility   Visibility  // Public or Private
    Url          string      // Permanent URL
    TempUrl      string      // Temporary URL (if requested)
}
```

#### Visibility

```go
type Visibility string

const (
    VisibilityPublic  Visibility = "public"   // Anyone can access
    VisibilityPrivate Visibility = "private"  // Authentication required
)
```

---

## Configuration

### Local Storage Configuration

```go
type LocalStorageConfig struct {
    BasePath              string      // Base directory path
    BaseURL               string      // Base URL for file access
    DefaultDirPermission  os.FileMode // Directory permissions (default: 0755)
    DefaultFilePermission os.FileMode // File permissions (default: 0644)
}
```

**Load from environment:**
```go
cfg := storage.LoadLocalStorageConfig()
// Uses these environment variables:
// - STORAGE_LOCAL_BASE_PATH (default: "./storage")
// - STORAGE_LOCAL_BASE_URL (default: "http://localhost:8080/storage")
// - STORAGE_LOCAL_DEFAULT_DIR_PERMISSION (default: 0755)
// - STORAGE_LOCAL_DEFAULT_FILE_PERMISSION (default: 0644)
```

### S3 Configuration

```go
type S3Config struct {
    Region               string  // AWS region (e.g., "us-east-1")
    Bucket               string  // S3 bucket name
    AccessKey            string  // AWS access key
    SecretKey            string  // AWS secret key
    Session              string  // AWS session token (optional)
    Url                  string  // Custom URL (for MinIO, etc.)
    UseStylePathEndpoint bool    // Path style vs virtual hosted style
}
```

**Load from environment:**
```go
cfg := storage.LoadS3Config()
// Uses these environment variables:
// - STORAGE_S3_REGION (default: "us-east-1")
// - STORAGE_S3_BUCKET
// - STORAGE_S3_ACCESS_KEY
// - STORAGE_S3_SECRET_KEY
// - STORAGE_S3_SESSION
// - STORAGE_S3_URL
// - STORAGE_S3_USE_STYLE_PATH_ENDPOINT (default: false)
```

**URL Formats:**

When `UseStylePathEndpoint` is:
- `false` (default): `https://bucket.s3.region.amazonaws.com/path/to/file`
- `true`: `https://s3.region.amazonaws.com/bucket/path/to/file`

For MinIO or custom S3:
```go
cfg := storage.S3Config{
    Region:               "us-east-1",
    Bucket:               "my-bucket",
    AccessKey:            "minioadmin",
    SecretKey:            "minioadmin",
    Url:                  "http://localhost:9000",
    UseStylePathEndpoint: true,
}
```

---

## Local Storage

Local storage saves files to your server's filesystem.

### Creating Local Storage

```go
cfg := storage.LocalStorageConfig{
    BasePath:              "./storage",
    BaseURL:               "http://localhost:8080/storage",
    DefaultDirPermission:  0755,
    DefaultFilePermission: 0644,
}

store, err := storage.NewLocalStorage(cfg)
if err != nil {
    log.Fatal(err)
}
```

### Features

-  Automatic directory creation
-  Configurable file/directory permissions
-  MIME type detection from file extension
-  Content-based MIME detection fallback
-  Human-readable file sizes
-  URL generation for file access

### Example: Upload File

```go
// Upload from bytes
result, err := store.Upload(ctx, storage.UploadInput{
    FileName:   "uploads/photo.jpg",
    Content:    imageBytes,
    MimeType:   "image/jpeg",
    Visibility: storage.VisibilityPublic,
})

// Upload from string
result, err := store.Upload(ctx, storage.UploadInput{
    FileName:   "text/note.txt",
    Content:    "Hello, World!",
    MimeType:   "text/plain",
    Visibility: storage.VisibilityPublic,
})

// Upload from io.Reader
file, _ := os.Open("document.pdf")
defer file.Close()

result, err := store.Upload(ctx, storage.UploadInput{
    FileName:   "docs/document.pdf",
    Content:    file,
    MimeType:   "application/pdf",
    Visibility: storage.VisibilityPrivate,
})
```

---

## S3 Storage

S3 storage uses Amazon S3 or compatible services (MinIO, DigitalOcean Spaces, etc.).

### Creating S3 Storage

```go
// Configure S3
cfg := storage.S3Config{
    Region:    "us-east-1",
    Bucket:    "my-bucket",
    AccessKey: "YOUR_ACCESS_KEY",
    SecretKey: "YOUR_SECRET_KEY",
}

// Create S3 client
client, err := storage.NewS3Client(cfg)
if err != nil {
    log.Fatal(err)
}

// Create storage
store := storage.NewS3Storage(client, cfg)
```

### Features

-  AWS S3 and compatible services
-  Public/Private ACL support
-  Presigned URLs for temporary access
-  Automatic ACL detection
-  Efficient copy/move operations
-  Directory listing with delimiters

### Example: Temporary URLs

```go
// Get file with temporary URL (valid for 1 hour)
expiry := 1 * time.Hour
file, err := store.File(ctx, "private/document.pdf", &expiry)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Permanent URL:", file.Url)
fmt.Println("Temporary URL:", file.TempUrl) // Valid for 1 hour
```

---

## Complete Examples

### Example 1: Upload and Download

```go
package main

import (
    "context"
    "fmt"
    "github.com/fatkulnurk/foundation/storage"
)

func main() {
    ctx := context.Background()
    
    // Create storage
    cfg := storage.LoadLocalStorageConfig()
    store, _ := storage.NewLocalStorage(cfg)
    
    // Upload file
    content := []byte("Hello, Storage!")
    result, err := store.Upload(ctx, storage.UploadInput{
        FileName:   "test/hello.txt",
        Content:    content,
        MimeType:   "text/plain",
        Visibility: storage.VisibilityPublic,
    })
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Uploaded: %s (%s)\n", result.Path, result.SizeHuman)
    
    // Download file
    data, err := store.Get(ctx, "test/hello.txt")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Downloaded: %s\n", string(data))
}
```

### Example 2: Copy and Move Files

```go
package main

import (
    "context"
    "github.com/fatkulnurk/foundation/storage"
)

func main() {
    ctx := context.Background()
    store, _ := storage.NewLocalStorage(storage.LoadLocalStorageConfig())
    
    // Upload original file
    store.Upload(ctx, storage.UploadInput{
        FileName:   "original/file.txt",
        Content:    []byte("Original content"),
        MimeType:   "text/plain",
        Visibility: storage.VisibilityPublic,
    })
    
    // Copy file
    err := store.Copy(ctx, "original/file.txt", "backup/file.txt")
    if err != nil {
        panic(err)
    }
    println("File copied to backup/")
    
    // Move file
    err = store.Move(ctx, "original/file.txt", "archive/file.txt")
    if err != nil {
        panic(err)
    }
    println("File moved to archive/")
}
```

### Example 3: List Files and Directories

```go
package main

import (
    "context"
    "fmt"
    "github.com/fatkulnurk/foundation/storage"
)

func main() {
    ctx := context.Background()
    store, _ := storage.NewLocalStorage(storage.LoadLocalStorageConfig())
    
    // List files in directory
    files, err := store.Files(ctx, "uploads", nil)
    if err != nil {
        panic(err)
    }
    
    fmt.Println("Files in uploads/:")
    for _, file := range files {
        fmt.Printf("  - %s (%s) - %s\n", 
            file.Name, file.SizeHuman, file.MimeType)
    }
    
    // List subdirectories
    dirs, err := store.Directories(ctx, "uploads")
    if err != nil {
        panic(err)
    }
    
    fmt.Println("\nSubdirectories:")
    for _, dir := range dirs {
        fmt.Printf("  - %s\n", dir)
    }
}
```

### Example 4: Check File Existence

```go
package main

import (
    "context"
    "fmt"
    "github.com/fatkulnurk/foundation/storage"
)

func main() {
    ctx := context.Background()
    store, _ := storage.NewLocalStorage(storage.LoadLocalStorageConfig())
    
    // Check if file exists
    exists, err := store.Exists(ctx, "uploads/photo.jpg")
    if err != nil {
        panic(err)
    }
    
    if exists {
        fmt.Println("File exists!")
        
        // Get file info
        file, _ := store.File(ctx, "uploads/photo.jpg", nil)
        fmt.Printf("Size: %s\n", file.SizeHuman)
        fmt.Printf("Type: %s\n", file.MimeType)
        fmt.Printf("URL: %s\n", file.Url)
    } else {
        fmt.Println("File not found")
    }
}
```

### Example 5: Delete Files

```go
package main

import (
    "context"
    "github.com/fatkulnurk/foundation/storage"
)

func main() {
    ctx := context.Background()
    store, _ := storage.NewLocalStorage(storage.LoadLocalStorageConfig())
    
    // Delete a file
    err := store.Delete(ctx, "temp/old-file.txt")
    if err != nil {
        panic(err)
    }
    
    println("File deleted successfully")
    
    // Verify deletion
    exists, _ := store.Exists(ctx, "temp/old-file.txt")
    if !exists {
        println("Confirmed: file no longer exists")
    }
}
```

### Example 6: S3 with Presigned URLs

```go
package main

import (
    "context"
    "fmt"
    "time"
    "github.com/fatkulnurk/foundation/storage"
)

func main() {
    ctx := context.Background()
    
    // Create S3 storage
    cfg := storage.LoadS3Config()
    client, _ := storage.NewS3Client(cfg)
    store := storage.NewS3Storage(client, cfg)
    
    // Upload private file
    store.Upload(ctx, storage.UploadInput{
        FileName:   "private/secret.pdf",
        Content:    []byte("Secret content"),
        MimeType:   "application/pdf",
        Visibility: storage.VisibilityPrivate,
    })
    
    // Generate temporary URL (valid for 15 minutes)
    expiry := 15 * time.Minute
    file, _ := store.File(ctx, "private/secret.pdf", &expiry)
    
    fmt.Println("Share this temporary URL:")
    fmt.Println(file.TempUrl)
    fmt.Println("Valid for 15 minutes")
}
```

### Example 7: Switch Between Storage Providers

```go
package main

import (
    "context"
    "os"
    "github.com/fatkulnurk/foundation/storage"
)

func getStorage() storage.Storage {
    storageType := os.Getenv("STORAGE_TYPE") // "local" or "s3"
    
    if storageType == "s3" {
        cfg := storage.LoadS3Config()
        client, _ := storage.NewS3Client(cfg)
        return storage.NewS3Storage(client, cfg)
    }
    
    // Default to local
    cfg := storage.LoadLocalStorageConfig()
    store, _ := storage.NewLocalStorage(cfg)
    return store
}

func main() {
    ctx := context.Background()
    
    // Get storage (works with both local and S3)
    store := getStorage()
    
    // Use storage (same code for both!)
    result, _ := store.Upload(ctx, storage.UploadInput{
        FileName:   "uploads/file.txt",
        Content:    []byte("Hello!"),
        MimeType:   "text/plain",
        Visibility: storage.VisibilityPublic,
    })
    
    println("File uploaded:", result.Path)
}
```

### Example 8: Upload Multiple Files

```go
package main

import (
    "context"
    "fmt"
    "github.com/fatkulnurk/foundation/storage"
)

func main() {
    ctx := context.Background()
    store, _ := storage.NewLocalStorage(storage.LoadLocalStorageConfig())
    
    files := []struct {
        name    string
        content string
    }{
        {"file1.txt", "Content 1"},
        {"file2.txt", "Content 2"},
        {"file3.txt", "Content 3"},
    }
    
    for _, f := range files {
        result, err := store.Upload(ctx, storage.UploadInput{
            FileName:   "batch/" + f.name,
            Content:    f.content,
            MimeType:   "text/plain",
            Visibility: storage.VisibilityPublic,
        })
        
        if err != nil {
            fmt.Printf("Failed to upload %s: %v\n", f.name, err)
            continue
        }
        
        fmt.Printf("Uploaded: %s (%s)\n", result.Name, result.SizeHuman)
    }
}
```

---

## Best Practices

### 1. Use Environment Variables for Configuration

```go
// Good - Easy to change without code changes
cfg := storage.LoadLocalStorageConfig()
store, _ := storage.NewLocalStorage(cfg)

// Or for S3
cfg := storage.LoadS3Config()
client, _ := storage.NewS3Client(cfg)
store := storage.NewS3Storage(client, cfg)
```

### 2. Always Check Errors

```go
result, err := store.Upload(ctx, input)
if err != nil {
    log.Printf("Upload failed: %v", err)
    return err
}
```

### 3. Use Context for Cancellation

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := store.Upload(ctx, input)
```

### 4. Set Appropriate Visibility

```go
// Public files (images, CSS, JS)
storage.VisibilityPublic

// Private files (user documents, sensitive data)
storage.VisibilityPrivate
```

### 5. Use Temporary URLs for Private Files

```go
// Instead of making files public, use temporary URLs
expiry := 1 * time.Hour
file, _ := store.File(ctx, "private/document.pdf", &expiry)
// Share file.TempUrl (valid for 1 hour)
```

### 6. Organize Files in Directories

```go
// Good structure
"uploads/users/123/profile.jpg"
"uploads/products/456/image1.jpg"
"documents/invoices/2024/invoice-001.pdf"

// Bad structure
"profile.jpg"
"image1.jpg"
"invoice-001.pdf"
```

### 7. Handle Large Files Efficiently

```go
// Use io.Reader for large files
file, _ := os.Open("large-video.mp4")
defer file.Close()

store.Upload(ctx, storage.UploadInput{
    FileName:   "videos/large-video.mp4",
    Content:    file, // Streams the file
    MimeType:   "video/mp4",
    Visibility: storage.VisibilityPublic,
})
```

### 8. Clean Up Temporary Files

```go
// Delete temporary files after use
defer store.Delete(ctx, "temp/processing-file.tmp")
```

### 9. Use Proper MIME Types

```go
// Common MIME types
"image/jpeg"           // .jpg, .jpeg
"image/png"            // .png
"image/gif"            // .gif
"application/pdf"      // .pdf
"text/plain"           // .txt
"text/html"            // .html
"application/json"     // .json
"video/mp4"            // .mp4
"audio/mpeg"           // .mp3
```

### 10. Test with Both Storage Providers

```go
func TestStorage(t *testing.T, store storage.Storage) {
    ctx := context.Background()
    
    // Test upload
    result, err := store.Upload(ctx, storage.UploadInput{
        FileName:   "test/file.txt",
        Content:    []byte("test"),
        MimeType:   "text/plain",
        Visibility: storage.VisibilityPublic,
    })
    assert.NoError(t, err)
    assert.NotEmpty(t, result.Path)
    
    // Test exists
    exists, err := store.Exists(ctx, "test/file.txt")
    assert.NoError(t, err)
    assert.True(t, exists)
    
    // Test delete
    err = store.Delete(ctx, "test/file.txt")
    assert.NoError(t, err)
}
```

---

## Thread Safety

All storage implementations are **thread-safe** and can be used concurrently from multiple goroutines:

```go
var wg sync.WaitGroup
store, _ := storage.NewLocalStorage(cfg)

for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        
        store.Upload(ctx, storage.UploadInput{
            FileName:   fmt.Sprintf("concurrent/file-%d.txt", id),
            Content:    fmt.Sprintf("Content %d", id),
            MimeType:   "text/plain",
            Visibility: storage.VisibilityPublic,
        })
    }(i)
}

wg.Wait()
```

---

## Extending

To create a custom storage provider, implement the `Storage` interface:

```go
type MyCustomStorage struct {
    // Your fields
}

func (s *MyCustomStorage) Upload(ctx context.Context, input storage.UploadInput) (*storage.UploadOutput, error) {
    // Your implementation
}

func (s *MyCustomStorage) Delete(ctx context.Context, path string) error {
    // Your implementation
}

// Implement all other methods...
```

---

## Troubleshooting

### Local Storage Issues

**Problem:** Permission denied
```
Solution: Check file/directory permissions
- Ensure BasePath is writable
- Check DefaultDirPermission and DefaultFilePermission values
```

**Problem:** File not found
```
Solution: Check file path
- Paths are relative to BasePath
- Use forward slashes (/) even on Windows
```

### S3 Storage Issues

**Problem:** Access denied
```
Solution: Check AWS credentials and permissions
- Verify AccessKey and SecretKey
- Ensure IAM user has s3:PutObject, s3:GetObject, etc.
```

**Problem:** Bucket not found
```
Solution: Check bucket configuration
- Verify bucket name
- Ensure bucket exists in the specified region
```

**Problem:** Invalid presigned URL
```
Solution: Check URL generation
- Verify Url and UseStylePathEndpoint settings
- For MinIO, set UseStylePathEndpoint to true
```

---

## Extending

You can create custom storage providers by implementing the Storage interface.

### Custom Storage Implementation

```go
type MyCustomStorage struct {
    // Your fields
}

func (s *MyCustomStorage) Upload(ctx context.Context, input UploadInput) (*UploadOutput, error) {
    // Your implementation
    return &UploadOutput{}, nil
}

func (s *MyCustomStorage) Delete(ctx context.Context, path string) error {
    // Your implementation
    return nil
}

func (s *MyCustomStorage) Copy(ctx context.Context, sourcePath, destinationPath string) error {
    // Your implementation
    return nil
}

func (s *MyCustomStorage) Move(ctx context.Context, sourcePath, destinationPath string) error {
    // Your implementation
    return nil
}

func (s *MyCustomStorage) Get(ctx context.Context, path string) ([]byte, error) {
    // Your implementation
    return []byte{}, nil
}

func (s *MyCustomStorage) File(ctx context.Context, path string, expiryTempUrl *time.Duration) (*FileStorage, error) {
    // Your implementation
    return &FileStorage{}, nil
}

func (s *MyCustomStorage) Files(ctx context.Context, dir string, expiryTempUrl *time.Duration) ([]FileStorage, error) {
    // Your implementation
    return []FileStorage{}, nil
}

func (s *MyCustomStorage) Directories(ctx context.Context, dir string) ([]string, error) {
    // Your implementation
    return []string{}, nil
}

func (s *MyCustomStorage) Exists(ctx context.Context, path string) (bool, error) {
    // Your implementation
    return false, nil
}
```

### Example: FTP Storage

```go
type FTPStorage struct {
    host     string
    username string
    password string
}

func NewFTPStorage(host, username, password string) Storage {
    return &FTPStorage{
        host:     host,
        username: username,
        password: password,
    }
}

func (s *FTPStorage) Upload(ctx context.Context, input UploadInput) (*UploadOutput, error) {
    // Connect to FTP server
    // Upload file
    // Return output
    return &UploadOutput{
        Name:      input.FileName,
        Path:      input.FileName,
        Size:      0,
        SizeHuman: "0 B",
    }, nil
}

// Implement other methods...
```

---

## Summary

The storage package provides a **simple, unified interface** for file storage:
- Works with local filesystem and AWS S3
- Easy to switch between providers
- Supports all common file operations
- Thread-safe and production-ready

**Key Features:**
- Upload, download, delete, copy, move files
- List files and directories
- Public/private visibility control
- Temporary URLs for secure access
- Automatic MIME type detection
- Human-readable file sizes

Now you can easily manage file storage in your Go applications! 
