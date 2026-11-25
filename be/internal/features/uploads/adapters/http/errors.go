package http

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/uploads/app"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, app.ErrStorageUnavailable):
		return huma.Error503ServiceUnavailable("object storage is not configured")
	case errors.Is(err, app.ErrInvalidPurpose):
		return huma.Error400BadRequest("invalid upload purpose")
	case errors.Is(err, app.ErrInvalidContentType):
		return huma.Error400BadRequest("content type not allowed for this purpose")
	case errors.Is(err, app.ErrDocumentTypeNotFound):
		return huma.Error400BadRequest("document type not found")
	default:
		return err
	}
}
