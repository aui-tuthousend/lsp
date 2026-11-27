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
	// Baca profil/pekerjaan asesi oleh staff (untuk verifikasi permohonan).
	selfOrStaff := auth.RequireSelfOrRoles("user_id",
		string(domain.RoleAdmin), string(domain.RoleAdminLSP), string(domain.RoleManagerSertifikasi))
	ownerOrAdmin := requireProfileOwnerOrRoles(handler, string(domain.RoleAdmin))

	huma.Post(api, "/asesi-profiles", handler.Create,
		adminOnly,
		httpplatform.WithStatus(201),
		httpplatform.WithDoc("asesi-profiles-create", "Create asesi profile", "Membuat profil data diri asesi.", "Asesi Profiles"),
	)
	huma.Get(api, "/asesi-profiles", handler.List,
		adminOnly,
		httpplatform.WithDoc("asesi-profiles-list", "List asesi profiles", "Daftar profil data diri asesi.", "Asesi Profiles"),
	)
	huma.Get(api, "/asesi-profiles/{id}", handler.Get,
		ownerOrAdmin,
		httpplatform.WithDoc("asesi-profiles-get", "Get asesi profile", "Ambil detail profil data diri asesi.", "Asesi Profiles"),
	)
	huma.Get(api, "/users/{user_id}/asesi-profile", handler.GetByUser,
		selfOrStaff,
		httpplatform.WithDoc("asesi-profiles-get-by-user", "Get asesi profile by user", "Ambil profil data diri asesi berdasarkan user ID.", "Asesi Profiles"),
	)
	huma.Patch(api, "/asesi-profiles/{id}", handler.Update,
		ownerOrAdmin,
		httpplatform.WithDoc("asesi-profiles-update", "Update asesi profile", "Ubah sebagian profil data diri asesi.", "Asesi Profiles"),
	)
	huma.Delete(api, "/asesi-profiles/{id}", handler.Delete,
		adminOnly,
		httpplatform.WithStatus(204),
		httpplatform.WithDoc("asesi-profiles-delete", "Delete asesi profile", "Hapus profil data diri asesi.", "Asesi Profiles"),
	)

	// Data pekerjaan asesi (asesi_occupations).
	huma.Post(api, "/asesi-occupations", handler.CreateOccupation,
		adminOnly,
		httpplatform.WithStatus(201),
		httpplatform.WithDoc("asesi-occupations-create", "Create asesi occupation", "Membuat data pekerjaan asesi.", "Asesi Occupations"),
	)
	huma.Get(api, "/asesi-occupations", handler.ListOccupations,
		adminOnly,
		httpplatform.WithDoc("asesi-occupations-list", "List asesi occupations", "Daftar data pekerjaan asesi.", "Asesi Occupations"),
	)
	huma.Get(api, "/asesi-occupations/{id}", handler.GetOccupation,
		adminOnly,
		httpplatform.WithDoc("asesi-occupations-get", "Get asesi occupation", "Ambil detail data pekerjaan asesi.", "Asesi Occupations"),
	)
	huma.Patch(api, "/asesi-occupations/{id}", handler.UpdateOccupation,
		adminOnly,
		httpplatform.WithDoc("asesi-occupations-update", "Update asesi occupation", "Ubah sebagian data pekerjaan asesi.", "Asesi Occupations"),
	)
	huma.Delete(api, "/asesi-occupations/{id}", handler.DeleteOccupation,
		adminOnly,
		httpplatform.WithStatus(204),
		httpplatform.WithDoc("asesi-occupations-delete", "Delete asesi occupation", "Hapus data pekerjaan asesi.", "Asesi Occupations"),
	)
	huma.Get(api, "/users/{user_id}/asesi-occupation", handler.GetOccupationByUser,
		selfOrStaff,
		httpplatform.WithDoc("asesi-occupations-get-by-user", "Get asesi occupation by user", "Ambil data pekerjaan asesi berdasarkan user ID.", "Asesi Occupations"),
	)
	huma.Put(api, "/users/{user_id}/asesi-occupation", handler.UpsertOccupationByUser,
		selfOrAdmin,
		httpplatform.WithDoc("asesi-occupations-upsert-by-user", "Upsert asesi occupation by user", "Isi/perbarui data pekerjaan milik user sendiri (asesi).", "Asesi Occupations"),
	)
}
