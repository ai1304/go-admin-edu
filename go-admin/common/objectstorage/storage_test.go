package objectstorage

import (
	"context"
	"net/url"
	"testing"
	"time"
)

func TestPresignedURLUsesPublicEndpointForSignature(t *testing.T) {
	storage, err := New(Config{
		Driver:          "minio",
		Endpoint:        "minio:9000",
		PublicEndpoint:  "http://localhost:9000",
		AccessKeyID:     "minioadmin",
		AccessKeySecret: "minioadmin",
		BucketName:      "go-admin-edu",
		UseSSL:          false,
	})
	if err != nil {
		t.Fatalf("create storage: %v", err)
	}

	rawURL, err := storage.PresignedGetObject(context.Background(), "tenant/0/resource/example.png", time.Minute)
	if err != nil {
		t.Fatalf("presign object: %v", err)
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("parse presigned url: %v", err)
	}
	if parsedURL.Host != "localhost:9000" {
		t.Fatalf("expected public host localhost:9000, got %s", parsedURL.Host)
	}
}
