package http

import (
	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/platform/auth"
	httpplatform "lsp/internal/platform/http"
	"lsp/internal/users/domain"
)

func RegisterRoutes(api huma.API, handler *Handler) {
	adminOnly := auth.RequireRoles(string(domain.RoleAdmin))

	huma.Get(api, "/asesor-profiles", handler.List,
		adminOnly,
		httpplatform.WithDoc("asesor-profiles-list", "List asesor profiles", "Daftar profil asesor beserta batas penugasan harian.", "Asesor"),
	)
	huma.Get(api, "/asesor-profiles/{user_id}", handler.Get,
		adminOnly,
		httpplatform.WithDoc("asesor-profiles-get", "Get asesor profile", "Ambil profil asesor (batas penugasan harian).", "Asesor"),
	)
	huma.Put(api, "/asesor-profiles/{user_id}", handler.Upsert,
		adminOnly,
		httpplatform.WithDoc("asesor-profiles-upsert", "Set asesor daily limit", "Set/ubah batas penugasan harian asesor.", "Asesor"),
	)
}
