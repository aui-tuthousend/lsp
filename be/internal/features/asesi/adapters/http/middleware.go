package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/platform/auth"
)

// requireProfileOwnerOrRoles mengizinkan operasi bila user memiliki salah satu
// role, atau bila user adalah pemilik profil asesi pada path param `id`.
func requireProfileOwnerOrRoles(handler *Handler, roles ...string) func(*huma.Operation) {
	return auth.RequireOwnerOrRoles("id", func(ctx context.Context, profileID string) (string, bool) {
		profile, err := handler.service.Get(ctx, profileID)
		if err != nil {
			return "", false
		}
		return profile.UserID, true
	}, roles...)
}
