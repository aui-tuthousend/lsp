package http

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/documentrequirements/domain"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return huma.Error404NotFound("document requirement not found")
	case errors.Is(err, domain.ErrConflict):
		return huma.Error409Conflict("document requirement already exists")
	default:
		return err
	}
}
