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
	client           *minio.Client
	bucketName       string
	publicScheme     string
	publicHost       string
	publicPathPrefix string
}

func NewMinIO(cfg Config) (*MinIOStorage, error) {
	if cfg.Endpoint == "" || cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" || cfg.BucketName == "" {
		return nil, errors.New("minio storage config is incomplete")
	}
	client, err := newMinIOClient(cfg.Endpoint, cfg.AccessKeyID, cfg.AccessKeySecret, cfg.UseSSL)
	if err != nil {
		return nil, err
	}
	return &MinIOStorage{
		client:           client,
		bucketName:       cfg.BucketName,
		publicScheme:     publicEndpointScheme(cfg.PublicEndpoint, cfg.UseSSL),
		publicHost:       publicEndpointHost(cfg.PublicEndpoint),
		publicPathPrefix: publicEndpointPathPrefix(cfg.PublicEndpoint),
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
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	u, err := s.client.PresignedGetObject(ctx, s.bucketName, objectKey, expires, nil)
	if err != nil {
		return "", err
	}
	if s.publicHost != "" {
		if s.publicPathPrefix != "" {
			u.Path = strings.TrimRight(s.publicPathPrefix, "/") + u.Path
		}
		u.Scheme = s.publicScheme
		u.Host = s.publicHost
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
	return minio.New(normalizeEndpoint(endpoint), &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, accessKeySecret, ""),
		Secure: endpointUseSSL(endpoint, useSSL),
	})
}

func endpointUseSSL(endpoint string, fallback bool) bool {
	parsedURL, err := url.Parse(endpoint)
	if err == nil && parsedURL.Scheme != "" {
		return parsedURL.Scheme == "https"
	}
	return fallback
}

func publicEndpointScheme(endpoint string, fallbackUseSSL bool) string {
	if endpointUseSSL(endpoint, fallbackUseSSL) {
		return "https"
	}
	return "http"
}

func publicEndpointHost(endpoint string) string {
	parsedURL, err := url.Parse(endpoint)
	if err == nil && parsedURL.Host != "" {
		return parsedURL.Host
	}
	return normalizeEndpoint(endpoint)
}

func publicEndpointPathPrefix(endpoint string) string {
	parsedURL, err := url.Parse(endpoint)
	if err == nil {
		return strings.TrimRight(parsedURL.EscapedPath(), "/")
	}
	return ""
}
