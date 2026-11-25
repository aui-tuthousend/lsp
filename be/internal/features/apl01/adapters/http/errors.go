package http

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/apl01/domain"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return huma.Error404NotFound("certification application not found")
	case errors.Is(err, domain.ErrConflict):
		return huma.Error409Conflict("certification application already exists")
	case errors.Is(err, domain.ErrNotDraft):
		return huma.Error409Conflict("certification application is not a draft")
	case errors.Is(err, domain.ErrIncomplete):
		return huma.Error422UnprocessableEntity("certification application is incomplete")
	case errors.Is(err, domain.ErrInvalidToken):
		return huma.Error400BadRequest("invalid or expired verification token")
	case errors.Is(err, domain.ErrTokenUsed):
		return huma.Error409Conflict("verification token already used")
	case errors.Is(err, domain.ErrScheduleConflict):
		return huma.Error409Conflict("you already have an application on this date")
	default:
		return err
	}
}
