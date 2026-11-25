package http

import (
	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/users/domain"
	"lsp/internal/platform/auth"
	httpplatform "lsp/internal/platform/http"
)

func RegisterRoutes(api huma.API, handler *Handler) {
	asesiOnly := auth.RequireRoles(string(domain.RoleAsesi))
	staffOnly := auth.RequireRoles(string(domain.RoleAdmin), string(domain.RoleAdminLSP), string(domain.RoleManagerSertifikasi))
	ownerOnly := requireApplicationOwnerOrRoles(handler)
	ownerOrStaff := requireApplicationOwnerOrRoles(handler,
		string(domain.RoleAdmin), string(domain.RoleAdminLSP), string(domain.RoleManagerSertifikasi))

	// Asesi: buat/lanjutkan draft & daftar milik sendiri.
	huma.Post(api, "/apl01", handler.CreateDraft,
		asesiOnly,
		httpplatform.WithStatus(201),
		httpplatform.WithDoc("apl01-create-draft", "Create or resume APL-01 draft", "Membuat draft APL-01 baru atau mengembalikan draft aktif milik asesi.", "APL-01"),
	)
	huma.Get(api, "/apl01/mine", handler.ListMine,
		asesiOnly,
		httpplatform.WithDoc("apl01-list-mine", "List my APL-01 applications", "Daftar permohonan milik asesi yang login (untuk resume draft).", "APL-01"),
	)

	// Admin/staff: daftar seluruh permohonan.
	huma.Get(api, "/apl01", handler.ListAll,
		staffOnly,
		httpplatform.WithDoc("apl01-list", "List APL-01 applications", "Daftar seluruh permohonan APL-01.", "APL-01"),
	)

	// Detail & pengisian step (asesi pemilik; staff boleh membaca).
	huma.Get(api, "/apl01/{id}", handler.Get,
		ownerOrStaff,
		httpplatform.WithDoc("apl01-get", "Get APL-01 application", "Ambil detail permohonan APL-01 beserta seluruh data step.", "APL-01"),
	)
	huma.Patch(api, "/apl01/{id}", handler.UpdateStep,
		ownerOnly,
		httpplatform.WithDoc("apl01-update-step", "Update current step", "Menyimpan step terakhir yang dikerjakan (untuk resume draft).", "APL-01"),
	)
	huma.Put(api, "/apl01/{id}/data", handler.PutData,
		ownerOnly,
		httpplatform.WithDoc("apl01-put-data", "Set certification data", "Menyimpan data sertifikasi (skema + tujuan asesmen).", "APL-01"),
	)
	huma.Put(api, "/apl01/{id}/documents", handler.PutDocuments,
		ownerOnly,
		httpplatform.WithDoc("apl01-put-documents", "Replace evidence documents", "Mengganti seluruh bukti kelengkapan dengan daftar yang dikirim.", "APL-01"),
	)
	huma.Put(api, "/apl01/{id}/schedule", handler.PutSchedule,
		ownerOnly,
		httpplatform.WithDoc("apl01-put-schedule", "Set assessment schedule", "Menyimpan jadwal asesmen yang diajukan asesi.", "APL-01"),
	)
	huma.Post(api, "/apl01/{id}/submit", handler.Submit,
		ownerOnly,
		httpplatform.WithDoc("apl01-submit", "Submit APL-01", "Validasi kelengkapan, submit permohonan, dan terbitkan token QR verifikasi.", "APL-01"),
	)
	huma.Post(api, "/apl01/{id}/verify", handler.Verify,
		ownerOnly,
		httpplatform.WithDoc("apl01-verify", "Verify APL-01 via QR", "Konfirmasi hasil scan QR (tandatangan asesi).", "APL-01"),
	)

	// Verifikasi admin (APPROVED/REJECTED + plotting asesor).
	huma.Post(api, "/apl01/{id}/decision", handler.Decide,
		staffOnly,
		httpplatform.WithDoc("apl01-decide", "Decide APL-01 application", "Menetapkan hasil verifikasi admin beserta plotting asesor.", "APL-01"),
	)
}
