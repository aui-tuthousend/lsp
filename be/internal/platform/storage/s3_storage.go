package storage

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	appconfig "lsp/internal/platform/config"
)

type S3Storage struct {
	client        *s3.Client
	presigner     *s3.PresignClient
	bucket        string
	publicBaseURL string
}

func NewS3Storage(ctx context.Context, cfg appconfig.Config) (*S3Storage, error) {
	if err := cfg.ValidateR2(); err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2AccountID)
	awsCfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(cfg.R2Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.R2AccessKeyID, cfg.R2SecretAccessKey, "")),
	)
	if err != nil {
		return nil, err
	}
	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})
	return &S3Storage{client: client, presigner: s3.NewPresignClient(client), bucket: cfg.R2Bucket, publicBaseURL: cfg.R2PublicBaseURL}, nil
}

func (s *S3Storage) PresignedPutURL(ctx context.Context, objectKey string, contentType string, expires time.Duration) (string, error) {
	out, err := s.presigner.PresignPutObject(ctx, &s3.PutObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(objectKey), ContentType: aws.String(contentType)}, s3.WithPresignExpires(expires))
	if err != nil {
		return "", err
	}
	return out.URL, nil
}

func (s *S3Storage) PresignedGetURL(ctx context.Context, objectKey string, expires time.Duration) (string, error) {
	out, err := s.presigner.PresignGetObject(ctx, &s3.GetObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(objectKey)}, s3.WithPresignExpires(expires))
	if err != nil {
		return "", err
	}
	return out.URL, nil
}

func (s *S3Storage) PublicURL(objectKey string) string {
	if s.publicBaseURL == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s", s.publicBaseURL, objectKey)
}

func (s *S3Storage) Delete(ctx context.Context, objectKey string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(objectKey)})
	return err
}

func (s *S3Storage) Exists(ctx context.Context, objectKey string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: aws.String(s.bucket), Key: aws.String(objectKey)})
	if err == nil {
		return true, nil
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) && (apiErr.ErrorCode() == "NotFound" || apiErr.ErrorCode() == "NoSuchKey") {
		return false, nil
	}
	var respErr interface{ HTTPStatusCode() int }
	if errors.As(err, &respErr) && respErr.HTTPStatusCode() == http.StatusNotFound {
		return false, nil
	}
	return false, err
}
