package http

import (
	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/platform/auth"
	httpplatform "lsp/internal/platform/http"
	"lsp/internal/users/domain"
)

func RegisterRoutes(api huma.API, handler *Handler) {
	adminOnly := auth.RequireRoles(string(domain.RoleAdmin))

	huma.Post(api, "/tuks", handler.Create,
		adminOnly,
		httpplatform.WithStatus(201),
		httpplatform.WithDoc("tuks-create", "Create TUK", "Membuat master data Tempat Uji Kompetensi.", "TUK"),
	)
	huma.Get(api, "/tuks", handler.List,
		adminOnly,
		httpplatform.WithDoc("tuks-list", "List TUK", "Daftar master data Tempat Uji Kompetensi.", "TUK"),
	)
	huma.Get(api, "/tuks/{id}", handler.Get,
		adminOnly,
		httpplatform.WithDoc("tuks-get", "Get TUK", "Ambil detail Tempat Uji Kompetensi.", "TUK"),
	)
	huma.Patch(api, "/tuks/{id}", handler.Update,
		adminOnly,
		httpplatform.WithDoc("tuks-update", "Update TUK", "Ubah sebagian master data Tempat Uji Kompetensi.", "TUK"),
	)
	huma.Delete(api, "/tuks/{id}", handler.Delete,
		adminOnly,
		httpplatform.WithStatus(204),
		httpplatform.WithDoc("tuks-delete", "Delete TUK", "Soft-delete master data Tempat Uji Kompetensi.", "TUK"),
	)
}
