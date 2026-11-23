package http

import (
	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/users/domain"
	"lsp/internal/platform/auth"
	httpplatform "lsp/internal/platform/http"
)

func RegisterRoutes(api huma.API, handler *Handler) {
	adminOnly := auth.RequireRoles(string(domain.RoleAdmin))

	huma.Get(api, "/asesor-profiles", handler.List,
		adminOnly,
		httpplatform.WithDoc("asesor-profiles-list", "List asesor profiles", "Daftar profil asesor beserta batas penugasan harian.", "Asesor Profiles"),
	)
	huma.Get(api, "/asesor-profiles/{user_id}", handler.Get,
		adminOnly,
		httpplatform.WithDoc("asesor-profiles-get", "Get asesor profile", "Ambil profil asesor (batas penugasan harian).", "Asesor Profiles"),
	)
	huma.Patch(api, "/asesor-profiles/{user_id}", handler.UpdateProfile,
		adminOnly,
		httpplatform.WithDoc("asesor-profiles-update", "Update asesor profile", "Ubah sebagian profil asesor.", "Asesor Profiles"),
	)

	huma.Get(api, "/asesor-profiles/{user_id}/certifications", handler.ListCertifications,
		adminOnly,
		httpplatform.WithDoc("asesor-certifications-list", "List asesor certifications", "Daftar sertifikasi milik seorang asesor.", "Asesor Certifications"),
	)
	huma.Post(api, "/asesor-profiles/{user_id}/certifications", handler.CreateCertification,
		adminOnly,
		httpplatform.WithStatus(201),
		httpplatform.WithDoc("asesor-certifications-create", "Create asesor certification", "Tambah satu sertifikasi untuk asesor.", "Asesor Certifications"),
	)
	huma.Patch(api, "/asesor-certifications/{cert_id}", handler.UpdateCertification,
		adminOnly,
		httpplatform.WithDoc("asesor-certifications-update", "Update asesor certification", "Ubah sebagian data sertifikasi asesor.", "Asesor Certifications"),
	)
	huma.Delete(api, "/asesor-profiles/{user_id}/certifications/{cert_id}", handler.DeleteCertification,
		adminOnly,
		httpplatform.WithStatus(204),
		httpplatform.WithDoc("asesor-certifications-delete", "Delete asesor certification", "Hapus satu sertifikasi asesor.", "Asesor Certifications"),
	)
}
