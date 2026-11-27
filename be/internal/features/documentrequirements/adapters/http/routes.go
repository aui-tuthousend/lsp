package http

import (
	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/users/domain"
	"lsp/internal/platform/auth"
	httpplatform "lsp/internal/platform/http"
)

func RegisterRoutes(api huma.API, handler *Handler) {
	adminOnly := auth.RequireRoles(string(domain.RoleAdmin))
	// Persyaratan dokumen boleh dibaca asesi (dipakai saat mengunggah bukti
	// kelengkapan APL-01) dan staff (verifikasi admin LSP/manager); mutasi tetap
	// khusus admin.
	requirementReaders := auth.RequireRoles(string(domain.RoleAdmin), string(domain.RoleAsesi),
		string(domain.RoleAdminLSP), string(domain.RoleManagerSertifikasi))

	// Master jenis dokumen.
	huma.Post(api, "/document-types", handler.CreateDocumentType,
		adminOnly,
		httpplatform.WithStatus(201),
		httpplatform.WithDoc("document-types-create", "Create document type", "Membuat master jenis dokumen.", "Document Types"),
	)
	huma.Get(api, "/document-types", handler.ListDocumentTypes,
		adminOnly,
		httpplatform.WithDoc("document-types-list", "List document types", "Daftar master jenis dokumen.", "Document Types"),
	)
	huma.Get(api, "/document-types/{id}", handler.GetDocumentType,
		adminOnly,
		httpplatform.WithDoc("document-types-get", "Get document type", "Ambil satu jenis dokumen berdasarkan id.", "Document Types"),
	)
	huma.Patch(api, "/document-types/{id}", handler.UpdateDocumentType,
		adminOnly,
		httpplatform.WithDoc("document-types-update", "Update document type", "Ubah sebagian master jenis dokumen.", "Document Types"),
	)
	huma.Delete(api, "/document-types/{id}", handler.DeleteDocumentType,
		adminOnly,
		httpplatform.WithStatus(204),
		httpplatform.WithDoc("document-types-delete", "Delete document type", "Hapus jenis dokumen.", "Document Types"),
	)

	// Persyaratan dokumen per skema (set-semantik: PUT menggantikan seluruh daftar).
	huma.Get(api, "/certification-schemas/{schema_id}/document-requirements", handler.ListRequirements,
		requirementReaders,
		httpplatform.WithDoc("document-requirements-list", "List schema document requirements", "Daftar persyaratan dokumen milik sebuah skema.", "Schema Document Requirements"),
	)
	huma.Put(api, "/certification-schemas/{schema_id}/document-requirements", handler.ReplaceRequirements,
		adminOnly,
		httpplatform.WithDoc("document-requirements-replace", "Replace schema document requirements", "Mengganti seluruh persyaratan dokumen sebuah skema dengan daftar yang dikirim.", "Schema Document Requirements"),
	)
}
