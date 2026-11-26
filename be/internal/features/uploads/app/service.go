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
	ErrStorageUnavailable   = errors.New("object storage is not configured")
	ErrInvalidPurpose       = errors.New("invalid upload purpose")
	ErrInvalidContentType   = errors.New("content type not allowed for this purpose")
	ErrDocumentTypeNotFound = errors.New("document type not found")
)

const uploadURLExpiry = 10 * time.Minute

// allowedTypes memetakan purpose (nama field dokumen) ke content-type yang
// diizinkan beserta ekstensi file.
var allowedTypes = map[string]map[string]string{
	"photo":     {"image/jpeg": "jpg", "image/png": "png"},
	"ktp_file":  {"application/pdf": "pdf"},
	"cv_file":   {"application/pdf": "pdf"},
	"npwp_file": {"application/pdf": "pdf"},
	"document":  {"image/jpeg": "jpg", "image/png": "png", "application/pdf": "pdf"},
	// berkas bukti asesmen mandiri (APL-02).
	"evidence": {"image/jpeg": "jpg", "image/png": "png", "application/pdf": "pdf"},
}

// DocumentTypeResolver mengambil daftar MIME yang diizinkan untuk sebuah jenis
// dokumen. Dipakai untuk menegakkan aturan tipe berkas pada purpose 'document'.
type DocumentTypeResolver interface {
	AllowedMimeTypes(ctx context.Context, documentTypeID string) ([]string, error)
}

// Service membuat presigned URL agar file diunggah langsung dari client ke
// object storage; backend hanya menyimpan object key-nya.
type Service struct {
	storage  storage.ObjectStorage
	docTypes DocumentTypeResolver
}

func NewService(storage storage.ObjectStorage, docTypes DocumentTypeResolver) *Service {
	return &Service{storage: storage, docTypes: docTypes}
}

type PresignCommand struct {
	Purpose        string
	ContentType    string
	DocumentTypeID *string
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
	// purpose 'document' dibatasi lebih lanjut oleh tipe berkas jenis dokumen.
	if cmd.Purpose == "document" && cmd.DocumentTypeID != nil && *cmd.DocumentTypeID != "" {
		if err := s.checkDocumentType(ctx, *cmd.DocumentTypeID, cmd.ContentType); err != nil {
			return PresignResult{}, err
		}
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

// checkDocumentType memastikan content type termasuk tipe yang diizinkan jenis
// dokumen. Bila resolver tidak tersedia, aturan dilewati (allowlist global tetap berlaku).
func (s *Service) checkDocumentType(ctx context.Context, documentTypeID, contentType string) error {
	if s.docTypes == nil {
		return nil
	}
	allowed, err := s.docTypes.AllowedMimeTypes(ctx, documentTypeID)
	if err != nil {
		return ErrDocumentTypeNotFound
	}
	for _, mime := range allowed {
		if mime == contentType {
			return nil
		}
	}
	return ErrInvalidContentType
}
