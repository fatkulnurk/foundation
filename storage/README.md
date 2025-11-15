# Storage Package

A flexible and extensible file storage interface for the application with Local and S3 implementations.

## Overview

The `storage` package provides a standardized interface for file storage operations within the application. It defines a common interface (`IStorage`) that can be implemented by various storage providers, with Local filesystem and AWS S3 being the default implementations.

## Interface

The package defines the `IStorage` interface with the following methods:

```go
type IStorage interface {
	Upload(ctx context.Context, input UploadInput) (*UploadOutput, error)
	Delete(ctx context.Context, path string) error
	Copy(ctx context.Context, sourcePath, destinationPath string) error
	Move(ctx context.Context, sourcePath, destinationPath string) error
	Get(ctx context.Context, path string) ([]byte, error)
	File(ctx context.Context, path string, expiryTempUrl *time.Duration) (*FileStorage, error)
	Files(ctx context.Context, dir string, expiryTempUrl *time.Duration) ([]FileStorage, error)
	Directories(ctx context.Context, dir string) ([]string, error)
	Exists(ctx context.Context, path string) (bool, error)
}
```

### Methods

- **Upload**: Stores a file in the storage with the specified content and options
- **Delete**: Removes a file from the storage by its path
- **Copy**: Creates a copy of a file at a new location
- **Move**: Moves a file from one location to another
- **Get**: Retrieves the content of a file as a byte array
- **File**: Gets detailed information about a single file
- **Files**: Lists files in a directory
- **Directories**: Lists subdirectories in a directory
- **Exists**: Checks if a file exists at the specified path

## File Visibility

The package supports two visibility levels for files:

```go
type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
)
```

- **Public**: Files are accessible without authentication
- **Private**: Files require authentication to access

## Implementations

### Local Storage

The package includes a local filesystem implementation:

```go
type LocalStorage struct {
	cfg config.LocalStorage
}
```

Local storage saves files to the local filesystem with configurable paths and permissions.

#### Usage

```go
import (
	"context"
	"your-project/config"
	"your-project/pkg/storage"
)

// Initialize local storage
cfg := config.LocalStorage{
	BasePath:             "/path/to/storage",
	BaseURL:              "http://localhost:8080/storage",
	DefaultDirPermission: 0755,
	DefaultFilePermission: 0644,
}

localStorage, err := storage.NewLocalStorage(cfg)
if err != nil {
	// Handle error
}

// Use the storage
ctx := context.Background()
output, err := localStorage.Upload(ctx, storage.UploadInput{
	FileName:   "uploads/image.jpg",
	Content:    imageBytes,
	MimeType:   "image/jpeg",
	Visibility: storage.VisibilityPublic,
})
```

### S3 Storage

The package includes an AWS S3 implementation:

```go
type S3Storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	cfg           config.S3
}
```

S3 storage uses Amazon S3 or compatible services for file storage.

#### Usage

```go
import (
	"context"
	"your-project/config"
	"your-project/pkg/storage"
)

// Initialize S3 client
cfg := config.S3{
	Region:    "us-west-2",
	Bucket:    "my-bucket",
	AccessKey: "YOUR_ACCESS_KEY",
	SecretKey: "YOUR_SECRET_KEY",
}

s3Client, err := storage.NewS3Client(cfg)
if err != nil {
	// Handle error
}

// Create storage instance
s3Storage := storage.NewS3Storage(s3Client, cfg)

// Use the storage
ctx := context.Background()
output, err := s3Storage.Upload(ctx, storage.UploadInput{
	FileName:   "uploads/document.pdf",
	Content:    pdfBytes,
	MimeType:   "application/pdf",
	Visibility: storage.VisibilityPrivate,
})
```

## Extending

To implement a new storage provider, create a struct that implements all methods of the `IStorage` interface.

## Thread Safety

All implementations are designed to be thread-safe and can be safely used concurrently from multiple goroutines.