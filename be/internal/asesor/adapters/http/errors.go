package http

import (
	"errors"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/asesor/domain"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return huma.Error404NotFound("asesor profile not found")
	case errors.Is(err, domain.ErrNotAsesor):
		return huma.Error422UnprocessableEntity("user is not an asesor")
	case errors.Is(err, domain.ErrTUKNotFound):
		return huma.Error422UnprocessableEntity("tuk not found")
	case errors.Is(err, domain.ErrInvalidLimit):
		return huma.Error400BadRequest("daily assignment limit must be >= 0")
	case errors.Is(err, domain.ErrCertificationNotFound):
		return huma.Error404NotFound("asesor certification not found")
	case errors.Is(err, domain.ErrInvalidCertificationPeriod):
		return huma.Error400BadRequest("certification expired_at must be after issued_at")
	default:
		return err
	}
}
