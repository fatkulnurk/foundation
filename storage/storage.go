package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Storage interface {
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

type UploadInput struct {
	FileName   string
	Content    any
	MimeType   string
	Visibility Visibility
}

type UploadOutput struct {
	Name      string
	Path      string
	Size      int64
	SizeHuman string
}

type FileStorage struct {
	Name         string
	Path         string
	Size         int64  // in bytes
	SizeHuman    string // in human readable format, like 500kb, 123mb, 4gb, etc
	MimeType     string
	LastModified time.Time
	Visibility   Visibility
	Url          string
	TempUrl      string
}

type Visibility string

const (
	VisibilityPublic  Visibility = "public"
	VisibilityPrivate Visibility = "private"
)

func (v Visibility) String() string {
	return string(v)
}

func (v Visibility) ToS3ACL() types.ObjectCannedACL {
	switch v {
	case VisibilityPublic:
		return types.ObjectCannedACLPublicRead
	case VisibilityPrivate:
		return types.ObjectCannedACLPrivate
	default:
		return types.ObjectCannedACLPrivate
	}
}

func ParseVisibility(v string) Visibility {
	switch v {
	case "public":
		return VisibilityPublic
	case "private":
		return VisibilityPrivate
	default:
		return VisibilityPrivate
	}
}
