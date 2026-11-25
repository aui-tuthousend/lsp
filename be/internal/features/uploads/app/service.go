package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"lsp/internal/platform/id"
	"lsp/internal/platform/storage"
)

var (
	ErrStorageUnavailable = errors.New("object storage is not configured")
	ErrInvalidPurpose     = errors.New("invalid upload purpose")
	ErrInvalidContentType = errors.New("content type not allowed for this purpose")
)

const uploadURLExpiry = 10 * time.Minute

// allowedTypes memetakan purpose (nama field dokumen) ke content-type yang
// diizinkan beserta ekstensi file.
var allowedTypes = map[string]map[string]string{
	"photo":     {"image/jpeg": "jpg", "image/png": "png"},
	"ktp_file":  {"image/jpeg": "jpg", "image/png": "png", "application/pdf": "pdf"},
	"cv_file":   {"application/pdf": "pdf"},
	"npwp_file": {"image/jpeg": "jpg", "image/png": "png", "application/pdf": "pdf"},
	"document":  {"image/jpeg": "jpg", "image/png": "png", "application/pdf": "pdf"},
}

// Service membuat presigned URL agar file diunggah langsung dari client ke
// object storage; backend hanya menyimpan object key-nya.
type Service struct {
	storage storage.ObjectStorage
}

func NewService(storage storage.ObjectStorage) *Service {
	return &Service{storage: storage}
}

type PresignCommand struct {
	Purpose     string
	ContentType string
}

type PresignResult struct {
	ObjectKey string
	UploadURL string
	ExpiresIn int
}

func (s *Service) Presign(ctx context.Context, cmd PresignCommand) (PresignResult, error) {
	types, ok := allowedTypes[cmd.Purpose]
	if !ok {
		return PresignResult{}, ErrInvalidPurpose
	}
	ext, ok := types[cmd.ContentType]
	if !ok {
		return PresignResult{}, ErrInvalidContentType
	}
	if s.storage == nil {
		return PresignResult{}, ErrStorageUnavailable
	}
	objectKey := fmt.Sprintf("uploads/%s/%s/%s.%s", cmd.Purpose, time.Now().UTC().Format("2006/01"), id.New(), ext)
	url, err := s.storage.PresignedPutURL(ctx, objectKey, cmd.ContentType, uploadURLExpiry)
	if err != nil {
		return PresignResult{}, err
	}
	return PresignResult{ObjectKey: objectKey, UploadURL: url, ExpiresIn: int(uploadURLExpiry.Seconds())}, nil
}
