package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/platform/auth"
)

// requireApplicationOwnerOrRoles mengizinkan operasi bila user memiliki salah
// satu role, atau bila user adalah pemilik permohonan pada path param `id`.
// Tanpa role berarti hanya pemilik yang diizinkan.
func requireApplicationOwnerOrRoles(handler *Handler, roles ...string) func(*huma.Operation) {
	return auth.RequireOwnerOrRoles("id", func(ctx context.Context, applicationID string) (string, bool) {
		detail, err := handler.service.Detail(ctx, applicationID)
		if err != nil {
			return "", false
		}
		return detail.Application.UserID, true
	}, roles...)
}
