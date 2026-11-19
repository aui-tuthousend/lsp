package storage

import (
	"context"
	"time"
)

type ObjectStorage interface {
	PresignedPutURL(ctx context.Context, objectKey string, contentType string, expires time.Duration) (string, error)
	PresignedGetURL(ctx context.Context, objectKey string, expires time.Duration) (string, error)
	PublicURL(objectKey string) string
	Delete(ctx context.Context, objectKey string) error
	Exists(ctx context.Context, objectKey string) (bool, error)
}
