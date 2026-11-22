package app

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"lsp/internal/platform/storage"
)

var (
	ErrStorageUnavailable = errors.New("object storage is not configured")
	ErrInvalidPurpose     = errors.New("invalid upload purpose")
	ErrInvalidContentType = errors.New("content type not allowed for this purpose")
	ErrInvalidNIK         = errors.New("invalid NIK")
)

const uploadURLExpiry = 10 * time.Minute

var nikPattern = regexp.MustCompile(`^[0-9]{16}$`)

// fileRule menentukan nama file tetap dan content-type yang diizinkan per
// purpose (nama field dokumen).
type fileRule struct {
	basename string
	types    map[string]string // content-type -> ekstensi
}

var allowedTypes = map[string]fileRule{
	"photo":     {"foto", map[string]string{"image/jpeg": "jpg", "image/png": "png"}},
	"ktp_file":  {"ktp", map[string]string{"image/jpeg": "jpg", "image/png": "png", "application/pdf": "pdf"}},
	"cv_file":   {"cv", map[string]string{"application/pdf": "pdf"}},
	"npwp_file": {"npwp", map[string]string{"image/jpeg": "jpg", "image/png": "png", "application/pdf": "pdf"}},
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
	NIK         string
	ContentType string
}

type PresignResult struct {
	ObjectKey string
	UploadURL string
	ExpiresIn int
}

func (s *Service) Presign(ctx context.Context, cmd PresignCommand) (PresignResult, error) {
	rule, ok := allowedTypes[cmd.Purpose]
	if !ok {
		return PresignResult{}, ErrInvalidPurpose
	}
	ext, ok := rule.types[cmd.ContentType]
	if !ok {
		return PresignResult{}, ErrInvalidContentType
	}
	if !nikPattern.MatchString(cmd.NIK) {
		return PresignResult{}, ErrInvalidNIK
	}
	if s.storage == nil {
		return PresignResult{}, ErrStorageUnavailable
	}
	// Nama file dibuat tetap (foto/ktp/cv/npwp) sehingga berkas tersimpan
	// terprediksi dalam satu folder: uploads/{tahun}/asesi/{nik}/dokumen-pribadi/.
	objectKey := fmt.Sprintf("uploads/%s/asesi/%s/dokumen-pribadi/%s.%s",
		time.Now().UTC().Format("2006"), cmd.NIK, rule.basename, ext)
	url, err := s.storage.PresignedPutURL(ctx, objectKey, cmd.ContentType, uploadURLExpiry)
	if err != nil {
		return PresignResult{}, err
	}
	return PresignResult{ObjectKey: objectKey, UploadURL: url, ExpiresIn: int(uploadURLExpiry.Seconds())}, nil
}
