package storage

import (
	"os"

	"github.com/fatkulnurk/foundation/support"
)

type S3Config struct {
	Region               string
	Bucket               string
	AccessKey            string
	SecretKey            string
	Session              string
	Url                  string // url for generate url, if fill this field, it will be used to generate url for file, example https://minio.example.com for usePathStyleEndpoint = true, and https://bucket.minio.example.com for usePathStyleEndpoint = false
	UseStylePathEndpoint bool   // if true, format will be s3.amazonaws.com/bucket, if false, format will be bucket.s3.amazonaws.com
}

func LoadS3Config() *S3Config {
	return &S3Config{
		Region:               support.GetEnv("STORAGE_S3_REGION", "us-east-1"),
		Bucket:               support.GetEnv("STORAGE_S3_BUCKET", ""),
		AccessKey:            support.GetEnv("STORAGE_S3_ACCESS_KEY", ""),
		SecretKey:            support.GetEnv("STORAGE_S3_SECRET_KEY", ""),
		Session:              support.GetEnv("STORAGE_S3_SESSION", ""),
		Url:                  support.GetEnv("STORAGE_S3_URL", ""),
		UseStylePathEndpoint: support.GetBoolEnv("STORAGE_S3_USE_STYLE_PATH_ENDPOINT", false),
	}
}

type LocalStorageConfig struct {
	BasePath              string
	BaseURL               string
	DefaultDirPermission  os.FileMode // default 0755
	DefaultFilePermission os.FileMode // default 0644
}

func LoadLocalStorageConfig() *LocalStorageConfig {
	return &LocalStorageConfig{
		BasePath:              support.GetEnv("STORAGE_LOCAL_BASE_PATH", "./storage"),
		BaseURL:               support.GetEnv("STORAGE_LOCAL_BASE_URL", "http://localhost:8080/storage"),
		DefaultDirPermission:  os.FileMode(support.GetIntEnv("STORAGE_LOCAL_DEFAULT_DIR_PERMISSION", 0755)),
		DefaultFilePermission: os.FileMode(support.GetIntEnv("STORAGE_LOCAL_DEFAULT_FILE_PERMISSION", 0644)),
	}
}
