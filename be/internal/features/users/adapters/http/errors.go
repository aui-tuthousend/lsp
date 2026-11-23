package http

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/users/domain"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return huma.Error404NotFound("user not found")
	case errors.Is(err, domain.ErrConflict):
		return huma.Error409Conflict("user already exists")
	case errors.Is(err, domain.ErrInvalidCredentials):
		return huma.Error401Unauthorized("invalid credentials")
	case errors.Is(err, domain.ErrInactiveUser):
		return huma.Error403Forbidden("inactive user")
	case errors.Is(err, domain.ErrInvalidRole):
		return huma.Error400BadRequest("invalid role")
	default:
		return err
	}
}

// mapRoleError memetakan error khusus master data roles.
func mapRoleError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return huma.Error404NotFound("role not found")
	case errors.Is(err, domain.ErrConflict):
		return huma.Error409Conflict("role already exists or is in use")
	default:
		return err
	}
}
