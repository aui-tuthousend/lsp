package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/apl02/domain"
	"lsp/internal/platform/auth"
)

// assessmentOwnerOf memetakan path param `id` (asesmen mandiri) menjadi id user
// pemilik sesuai `pick` (asesi pemilik permohonan atau asesor yang di-plotting).
func assessmentOwnerOf(handler *Handler, pick func(domain.ApplicationContext) string) func(context.Context, string) (string, bool) {
	return func(ctx context.Context, assessmentID string) (string, bool) {
		appCtx, err := handler.service.ApplicationContext(ctx, assessmentID)
		if err != nil {
			return "", false
		}
		owner := pick(appCtx)
		if owner == "" {
			return "", false
		}
		return owner, true
	}
}

// requireAssessmentRead: asesi pemilik, asesor, atau staff boleh membaca.
func requireAssessmentRead(handler *Handler, roles ...string) func(*huma.Operation) {
	return auth.RequireOwnerOrRoles("id", assessmentOwnerOf(handler, func(c domain.ApplicationContext) string {
		return c.UserID
	}), roles...)
}

// requireAssessmentOwner: hanya asesi pemilik permohonan (untuk mutasi).
func requireAssessmentOwner(handler *Handler) func(*huma.Operation) {
	return auth.RequireOwnerOrRoles("id", assessmentOwnerOf(handler, func(c domain.ApplicationContext) string {
		return c.UserID
	}))
}

// requireAssessmentAsesor: asesor yang di-plotting, atau staff, untuk tinjauan.
func requireAssessmentAsesor(handler *Handler, roles ...string) func(*huma.Operation) {
	return auth.RequireOwnerOrRoles("id", assessmentOwnerOf(handler, func(c domain.ApplicationContext) string {
		if c.AsesorID == nil {
			return ""
		}
		return *c.AsesorID
	}), roles...)
}
