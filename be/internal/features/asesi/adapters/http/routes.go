package http

import (
	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/users/domain"
	"lsp/internal/platform/auth"
	httpplatform "lsp/internal/platform/http"
)

func RegisterRoutes(api huma.API, handler *Handler) {
	adminOnly := auth.RequireRoles(string(domain.RoleAdmin))
	selfOrAdmin := auth.RequireSelfOrRoles("user_id", string(domain.RoleAdmin))
	ownerOrAdmin := requireProfileOwnerOrRoles(handler, string(domain.RoleAdmin))

	huma.Post(api, "/asesi-profiles", handler.Create,
		adminOnly,
		httpplatform.WithStatus(201),
		httpplatform.WithDoc("asesi-profiles-create", "Create asesi profile", "Membuat profil data diri asesi.", "Asesi"),
	)
	huma.Get(api, "/asesi-profiles", handler.List,
		adminOnly,
		httpplatform.WithDoc("asesi-profiles-list", "List asesi profiles", "Daftar profil data diri asesi.", "Asesi"),
	)
	huma.Get(api, "/asesi-profiles/{id}", handler.Get,
		ownerOrAdmin,
		httpplatform.WithDoc("asesi-profiles-get", "Get asesi profile", "Ambil detail profil data diri asesi.", "Asesi"),
	)
	huma.Get(api, "/users/{user_id}/asesi-profile", handler.GetByUser,
		selfOrAdmin,
		httpplatform.WithDoc("asesi-profiles-get-by-user", "Get asesi profile by user", "Ambil profil data diri asesi berdasarkan user ID.", "Asesi"),
	)
	huma.Patch(api, "/asesi-profiles/{id}", handler.Update,
		ownerOrAdmin,
		httpplatform.WithDoc("asesi-profiles-update", "Update asesi profile", "Ubah sebagian profil data diri asesi.", "Asesi"),
	)
	huma.Delete(api, "/asesi-profiles/{id}", handler.Delete,
		adminOnly,
		httpplatform.WithStatus(204),
		httpplatform.WithDoc("asesi-profiles-delete", "Delete asesi profile", "Hapus profil data diri asesi.", "Asesi"),
	)
}
