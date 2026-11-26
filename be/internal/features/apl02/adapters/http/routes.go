package http

import (
	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/users/domain"
	"lsp/internal/platform/auth"
	httpplatform "lsp/internal/platform/http"
)

func RegisterRoutes(api huma.API, handler *Handler) {
	asesiOnly := auth.RequireRoles(string(domain.RoleAsesi))
	readAccess := requireAssessmentRead(handler,
		string(domain.RoleAdmin), string(domain.RoleAdminLSP), string(domain.RoleManagerSertifikasi), string(domain.RoleAsesor))
	ownerOnly := requireAssessmentOwner(handler)
	asesorOrStaff := requireAssessmentAsesor(handler,
		string(domain.RoleAdmin), string(domain.RoleAdminLSP), string(domain.RoleManagerSertifikasi))

	// Asesi: buat/ambil asesmen mandiri untuk permohonan SUBMITTED/APPROVED.
	huma.Post(api, "/apl02", handler.Create,
		asesiOnly,
		httpplatform.WithStatus(201),
		httpplatform.WithDoc("apl02-create", "Create or resume APL-02", "Membuat atau mengembalikan asesmen mandiri (APL-02) untuk sebuah permohonan.", "APL-02"),
	)
	huma.Get(api, "/apl02/mine", handler.ListMine,
		asesiOnly,
		httpplatform.WithDoc("apl02-list-mine", "List my APL-02", "Daftar asesmen mandiri milik asesi yang login.", "APL-02"),
	)
	huma.Get(api, "/apl02/{id}", handler.Get,
		readAccess,
		httpplatform.WithDoc("apl02-get", "Get APL-02", "Ambil detail asesmen mandiri beserta seluruh jawaban unit.", "APL-02"),
	)
	huma.Put(api, "/apl02/{id}/units", handler.PutUnits,
		ownerOnly,
		httpplatform.WithDoc("apl02-put-units", "Save APL-02 answers", "Menyimpan jawaban per unit kompetensi beserta berkas bukti.", "APL-02"),
	)
	huma.Post(api, "/apl02/{id}/submit", handler.Submit,
		ownerOnly,
		httpplatform.WithDoc("apl02-submit", "Submit APL-02", "Validasi kelengkapan, submit asesmen mandiri, dan terbitkan token QR.", "APL-02"),
	)
	huma.Post(api, "/apl02/verify", handler.Verify,
		asesiOnly,
		httpplatform.WithDoc("apl02-verify", "Verify APL-02 via QR", "Konfirmasi hasil scan QR (tandatangan asesi) untuk asesmen mandiri.", "APL-02"),
	)
	huma.Get(api, "/apl02/verify/{token}", handler.VerifyStatus,
		asesiOnly,
		httpplatform.WithDoc("apl02-verify-status", "Get APL-02 QR status", "Status token verifikasi asesmen mandiri tanpa memakainya.", "APL-02"),
	)

	// Asesor (yang di-plotting) atau staff: tinjauan + jadwal zoom.
	huma.Post(api, "/apl02/{id}/review", handler.Review,
		asesorOrStaff,
		httpplatform.WithDoc("apl02-review", "Review APL-02", "Menetapkan rekomendasi asesor, catatan, dan jadwal zoom.", "APL-02"),
	)
}
