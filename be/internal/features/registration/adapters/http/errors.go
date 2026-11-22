package http

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	asesidomain "lsp/internal/features/asesi/domain"
	usersdomain "lsp/internal/features/users/domain"
)

// mapError menyatukan error domain dari modul users (pembuatan akun) dan
// modul asesi (pembuatan profil) karena satu request menyentuh keduanya.
func mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, usersdomain.ErrConflict):
		return huma.Error409Conflict("email already registered")
	case errors.Is(err, usersdomain.ErrInvalidRole):
		return huma.Error400BadRequest("invalid role")
	case errors.Is(err, asesidomain.ErrConflict):
		return huma.Error409Conflict("asesi profile already exists")
	case errors.Is(err, asesidomain.ErrNotAsesi):
		return huma.Error422UnprocessableEntity("user is not an asesi")
	case errors.Is(err, asesidomain.ErrInvalidDomicile):
		return huma.Error400BadRequest("domicile_address is required when domicile_same_as_ktp is false")
	default:
		return err
	}
}
