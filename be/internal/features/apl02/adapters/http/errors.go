package http

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/apl02/domain"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return huma.Error404NotFound("independent assessment not found")
	case errors.Is(err, domain.ErrConflict):
		return huma.Error409Conflict("independent assessment already exists")
	case errors.Is(err, domain.ErrNotDraft):
		return huma.Error409Conflict("independent assessment is not a draft")
	case errors.Is(err, domain.ErrNotSubmitted):
		return huma.Error409Conflict("independent assessment is not submitted")
	case errors.Is(err, domain.ErrNotEligible):
		return huma.Error409Conflict("application is not eligible for independent assessment")
	case errors.Is(err, domain.ErrIncomplete):
		return huma.Error422UnprocessableEntity("independent assessment is incomplete")
	case errors.Is(err, domain.ErrInvalidToken):
		return huma.Error400BadRequest("invalid or expired verification token")
	case errors.Is(err, domain.ErrTokenUsed):
		return huma.Error409Conflict("verification token already used")
	default:
		return err
	}
}
