package http

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/asesi/domain"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return huma.Error404NotFound("asesi profile not found")
	case errors.Is(err, domain.ErrConflict):
		return huma.Error409Conflict("asesi profile already exists")
	case errors.Is(err, domain.ErrNotAsesi):
		return huma.Error422UnprocessableEntity("user is not an asesi")
	case errors.Is(err, domain.ErrInvalidDomicile):
		return huma.Error400BadRequest("domicile_address is required when domicile_same_as_ktp is false")
	default:
		return err
	}
}
