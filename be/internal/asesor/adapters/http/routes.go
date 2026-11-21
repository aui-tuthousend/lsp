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
	huma.Patch(api, "/asesor-profiles/{user_id}", handler.UpdateProfile,
		adminOnly,
		httpplatform.WithDoc("asesor-profiles-update", "Update asesor profile", "Ubah sebagian profil asesor.", "Asesor"),
	)

	huma.Get(api, "/asesor-profiles/{user_id}/certifications", handler.ListCertifications,
		adminOnly,
		httpplatform.WithDoc("asesor-certifications-list", "List asesor certifications", "Daftar sertifikasi milik seorang asesor.", "Asesor"),
	)
	huma.Post(api, "/asesor-profiles/{user_id}/certifications", handler.CreateCertification,
		adminOnly,
		httpplatform.WithStatus(201),
		httpplatform.WithDoc("asesor-certifications-create", "Create asesor certification", "Tambah satu sertifikasi untuk asesor.", "Asesor"),
	)
	huma.Patch(api, "/asesor-certifications/{cert_id}", handler.UpdateCertification,
		adminOnly,
		httpplatform.WithDoc("asesor-certifications-update", "Update asesor certification", "Ubah sebagian data sertifikasi asesor.", "Asesor"),
	)
	huma.Delete(api, "/asesor-profiles/{user_id}/certifications/{cert_id}", handler.DeleteCertification,
		adminOnly,
		httpplatform.WithStatus(204),
		httpplatform.WithDoc("asesor-certifications-delete", "Delete asesor certification", "Hapus satu sertifikasi asesor.", "Asesor"),
	)
}
