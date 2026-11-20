package http

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/tuk/domain"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return huma.Error404NotFound("tuk not found")
	case errors.Is(err, domain.ErrConflict):
		return huma.Error409Conflict("tuk already exists")
	default:
		return err
	}
}
