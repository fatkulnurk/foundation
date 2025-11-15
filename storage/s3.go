package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	appcfg "github.com/fatkulnurk/gostarter/pkg/config"
	"github.com/fatkulnurk/gostarter/pkg/logging"
	"github.com/fatkulnurk/gostarter/pkg/support"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewS3Client(cfg appcfg.S3) (*s3.Client, error) {
	// Load konfigurasi AWS default dari environment, file config, dsb
	awscfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, cfg.Session)),
	)
	if err != nil {
		logging.Error(context.Background(), fmt.Sprintf("unable to load SDK config, %v", err))
		return nil, err
	}

	client := s3.NewFromConfig(awscfg)
	return client, nil
}

type S3Storage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	cfg           appcfg.S3
}

func NewS3Storage(client *s3.Client, cfg appcfg.S3) Storage {
	presignClient := s3.NewPresignClient(client)

	return &S3Storage{
		client:        client,
		presignClient: presignClient,
		cfg:           cfg,
	}
}

func (s *S3Storage) Upload(ctx context.Context, input UploadInput) (*UploadOutput, error) {
	var content io.ReadSeeker

	switch input.Content.(type) {
	case []byte:
		content = bytes.NewReader(input.Content.([]byte))
	case string:
		content = bytes.NewReader([]byte(input.Content.(string)))
	case io.Reader:
		buf := bytes.NewBuffer(nil)
		_, err := buf.ReadFrom(input.Content.(io.Reader))
		if err != nil {
			return nil, err
		}
		content = bytes.NewReader(buf.Bytes())
	default:
		return nil, fmt.Errorf("unsupported content type: %T", input.Content)
	}
	output, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.cfg.Bucket),
		Key:         aws.String(input.FileName),
		Body:        content,
		ContentType: aws.String(input.MimeType),
		ACL:         input.Visibility.ToS3ACL(),
	})

	if err != nil {
		return nil, err
	}

	size := int64(0)
	if output.Size != nil {
		size = *output.Size
	}

	// Extract the filename without the path
	fileName := filepath.Base(input.FileName)

	return &UploadOutput{
		Name:      fileName,
		Path:      input.FileName,
		Size:      size,
		SizeHuman: support.BytesToHumanReadable(size),
	}, nil
}

func (s *S3Storage) Delete(ctx context.Context, path string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(path),
	})
	return err
}

func (s *S3Storage) Copy(ctx context.Context, sourcePath, destinationPath string) error {
	_, err := s.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(s.cfg.Bucket),
		CopySource: aws.String(fmt.Sprintf("%s/%s", s.cfg.Bucket, sourcePath)),
		Key:        aws.String(destinationPath),
	})
	return err
}

func (s *S3Storage) Move(ctx context.Context, sourcePath, destinationPath string) error {
	// First copy the object to the new location
	if err := s.Copy(ctx, sourcePath, destinationPath); err != nil {
		return err
	}

	// Then delete the original object
	return s.Delete(ctx, sourcePath)
}

func (s *S3Storage) Get(ctx context.Context, path string) ([]byte, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

func (s *S3Storage) Exists(ctx context.Context, path string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s *S3Storage) File(ctx context.Context, path string, expiryTempUrl *time.Duration) (*FileStorage, error) {
	result, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	tempUrl := ""
	if expiryTempUrl != nil {
		presignUrl, err := s.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(s.cfg.Bucket),
			Key:    aws.String(path),
		}, s3.WithPresignExpires(*expiryTempUrl))
		if err != nil {
			return nil, err
		}
		tempUrl = presignUrl.URL
	}

	objectUrl := ""
	if s.cfg.Url == "" {
		// Generate a direct URL to the object
		objectUrl = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.cfg.Bucket, s.cfg.Region, path)
		if s.cfg.UseStylePathEndpoint {
			objectUrl = fmt.Sprintf("https://s3.%s.amazonaws.com/%s/%s", s.cfg.Region, s.cfg.Bucket, path)
		}
	} else {
		if s.cfg.UseStylePathEndpoint {
			objectUrl = fmt.Sprintf("%s/%s/%s", s.cfg.Url, s.cfg.Bucket, path)
		} else {
			objectUrl = fmt.Sprintf("%s/%s", s.cfg.Url, path)
		}
	}

	size := int64(0)
	if result.ContentLength != nil {
		size = *result.ContentLength
	}

	lastModified := time.Time{}
	if result.LastModified != nil {
		lastModified = *result.LastModified
	}

	visibility := VisibilityPrivate
	if public, err := s.IsPublicObject(ctx, path); err == nil && public {
		visibility = VisibilityPublic
	}

	file := &FileStorage{
		Name:         filepath.Base(path),
		Path:         path,
		Size:         size,
		SizeHuman:    support.BytesToHumanReadable(size),
		MimeType:     aws.ToString(result.ContentType),
		LastModified: lastModified,
		Visibility:   visibility,
		Url:          objectUrl,
		TempUrl:      tempUrl,
	}

	return file, nil
}

func (s *S3Storage) Files(ctx context.Context, dir string, expiryTempUrl *time.Duration) ([]FileStorage, error) {
	// Ensure directory path ends with a slash
	if dir != "" && !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}

	result, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.cfg.Bucket),
		Prefix:    aws.String(dir),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, err
	}

	var files []FileStorage
	for _, object := range result.Contents {
		// Skip if the object is a directory (ends with /)
		if strings.HasSuffix(*object.Key, "/") {
			continue
		}

		file, err := s.File(ctx, *object.Key, expiryTempUrl)
		if err != nil || file == nil {
			continue
		}

		files = append(files, *file)
	}

	return files, nil
}

func (s *S3Storage) Directories(ctx context.Context, dir string) ([]string, error) {
	// Ensure directory path ends with a slash
	if dir != "" && !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}

	result, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.cfg.Bucket),
		Prefix:    aws.String(dir),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, err
	}

	var directories []string
	for _, prefix := range result.CommonPrefixes {
		// Remove the trailing slash and get the directory name
		dirName := strings.TrimSuffix(*prefix.Prefix, "/")
		if dir != "" {
			// Remove the parent directory prefix to get just the directory name
			dirName = strings.TrimPrefix(dirName, dir)
		}
		directories = append(directories, dirName)
	}
	return directories, nil
}

func (s *S3Storage) IsPublicObject(ctx context.Context, path string) (bool, error) {
	aclOutput, err := s.client.GetObjectAcl(ctx, &s3.GetObjectAclInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return false, err
	}

	for _, grant := range aclOutput.Grants {
		if grant.Grantee != nil && grant.Grantee.URI != nil {
			uri := *grant.Grantee.URI
			if uri == "http://acs.amazonaws.com/groups/global/AllUsers" ||
				uri == "http://acs.amazonaws.com/groups/global/AuthenticatedUsers" {
				if grant.Permission == "READ" || grant.Permission == "FULL_CONTROL" {
					return true, nil
				}
			}
		}
	}

	return false, nil
}
