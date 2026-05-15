package objectstorage

import (
	"context"
	"errors"
	"io"
	"net/url"
	"strings"
	"time"

	appconfig "go-admin/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type ObjectStorage interface {
	PutObject(ctx context.Context, objectKey string, reader io.Reader, size int64, contentType string) error
	PresignedGetObject(ctx context.Context, objectKey string, expires time.Duration) (string, error)
	DeleteObject(ctx context.Context, objectKey string) error
	EnsureBucket(ctx context.Context) error
	BucketName() string
}

type Config struct {
	Driver          string
	Endpoint        string
	PublicEndpoint  string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
	UseSSL          bool
}

func NewFromExtend() (ObjectStorage, error) {
	cfg := appconfig.ExtConfig.Storage
	return New(Config{
		Driver:          cfg.Driver,
		Endpoint:        cfg.Endpoint,
		PublicEndpoint:  cfg.PublicEndpoint,
		AccessKeyID:     cfg.AccessKeyID,
		AccessKeySecret: cfg.AccessKeySecret,
		BucketName:      cfg.BucketName,
		UseSSL:          cfg.UseSSL,
	})
}

func New(cfg Config) (ObjectStorage, error) {
	switch strings.ToLower(cfg.Driver) {
	case "minio", "s3":
		return NewMinIO(cfg)
	case "":
		return nil, errors.New("object storage driver is empty")
	default:
		return nil, errors.New("unsupported object storage driver: " + cfg.Driver)
	}
}

type MinIOStorage struct {
	client          *minio.Client
	presignedClient *minio.Client
	bucketName      string
}

func NewMinIO(cfg Config) (*MinIOStorage, error) {
	if cfg.Endpoint == "" || cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" || cfg.BucketName == "" {
		return nil, errors.New("minio storage config is incomplete")
	}
	client, err := newMinIOClient(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret, cfg.UseSSL)
	if err != nil {
		return nil, err
	}
	presignedClient := client
	if cfg.PublicEndpoint != "" {
		publicClient, err := newMinIOClient(cfg.PublicEndpoint, cfg.AccessKeyID, cfg.AccessKeySecret, cfg.UseSSL)
		if err != nil {
			return nil, err
		}
		presignedClient = publicClient
	}
	return &MinIOStorage{
		client:          client,
		presignedClient: presignedClient,
		bucketName:      cfg.BucketName,
	}, nil
}

func (s *MinIOStorage) BucketName() string {
	return s.bucketName
}

func (s *MinIOStorage) EnsureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
}

func (s *MinIOStorage) PutObject(ctx context.Context, objectKey string, reader io.Reader, size int64, contentType string) error {
	_, err := s.client.PutObject(ctx, s.bucketName, objectKey, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (s *MinIOStorage) PresignedGetObject(ctx context.Context, objectKey string, expires time.Duration) (string, error) {
	u, err := s.presignedClient.PresignedGetObject(ctx, s.bucketName, objectKey, expires, nil)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

func (s *MinIOStorage) DeleteObject(ctx context.Context, objectKey string) error {
	return s.client.RemoveObject(ctx, s.bucketName, objectKey, minio.RemoveObjectOptions{})
}

func normalizeEndpoint(endpoint string) string {
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	return strings.TrimRight(endpoint, "/")
}

func newMinIOClient(endpoint, accessKeyID, accessKeySecret string, useSSL bool) (*minio.Client, error) {
	secure := useSSL
	parsedURL, err := url.Parse(endpoint)
	if err == nil && parsedURL.Scheme != "" {
		secure = parsedURL.Scheme == "https"
	}
	return minio.New(normalizeEndpoint(endpoint), &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, accessKeySecret, ""),
		Secure: secure,
	})
}
