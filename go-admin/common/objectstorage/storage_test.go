package objectstorage

import (
	"context"
	"net/url"
	"testing"
	"time"
)

func TestPresignedURLUsesPublicEndpointForSignature(t *testing.T) {
	objectKey := "tenant/0/resource/example.png"
	expires := time.Minute
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

	rawURL, err := storage.PresignedGetObject(context.Background(), objectKey, expires)
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

	expectedClient, err := newMinIOClient("http://localhost:9000", "minioadmin", "minioadmin", false)
	if err != nil {
		t.Fatalf("create expected public client: %v", err)
	}
	expectedURL, err := expectedClient.PresignedGetObject(context.Background(), "go-admin-edu", objectKey, expires, nil)
	if err != nil {
		t.Fatalf("presign with expected public client: %v", err)
	}
	if got, want := parsedURL.Query().Get("X-Amz-Signature"), expectedURL.Query().Get("X-Amz-Signature"); got != want {
		t.Fatalf("expected signature from public endpoint, got %s want %s", got, want)
	}
}
