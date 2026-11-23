package http

import (
	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/users/domain"
	"lsp/internal/platform/auth"
	httpplatform "lsp/internal/platform/http"
)

func RegisterRoutes(api huma.API, handler *Handler) {
	adminOnly := auth.RequireRoles(string(domain.RoleAdmin))

	// Skema sertifikasi.
	huma.Post(api, "/certification-schemas", handler.CreateSchema,
		adminOnly,
		httpplatform.WithStatus(201),
		httpplatform.WithDoc("certification-schemas-create", "Create certification schema", "Membuat master data skema sertifikasi.", "Certification Schemas"),
	)
	huma.Get(api, "/certification-schemas", handler.ListSchemas,
		adminOnly,
		httpplatform.WithDoc("certification-schemas-list", "List certification schemas", "Daftar master data skema sertifikasi.", "Certification Schemas"),
	)
	huma.Get(api, "/certification-schemas/{id}", handler.GetSchema,
		adminOnly,
		httpplatform.WithDoc("certification-schemas-get", "Get certification schema", "Ambil skema beserta struktur unit, elemen, dan kriteria.", "Certification Schemas"),
	)
	huma.Patch(api, "/certification-schemas/{id}", handler.UpdateSchema,
		adminOnly,
		httpplatform.WithDoc("certification-schemas-update", "Update certification schema", "Ubah sebagian master data skema sertifikasi.", "Certification Schemas"),
	)
	huma.Delete(api, "/certification-schemas/{id}", handler.DeleteSchema,
		adminOnly,
		httpplatform.WithStatus(204),
		httpplatform.WithDoc("certification-schemas-delete", "Delete certification schema", "Hapus skema beserta seluruh strukturnya.", "Certification Schemas"),
	)

	// Unit kompetensi.
	huma.Get(api, "/certification-schemas/{schema_id}/units", handler.ListUnits,
		adminOnly,
		httpplatform.WithDoc("competency-units-list", "List competency units", "Daftar unit kompetensi milik sebuah skema.", "Competency Units"),
	)
	huma.Post(api, "/certification-schemas/{schema_id}/units", handler.CreateUnit,
		adminOnly,
		httpplatform.WithStatus(201),
		httpplatform.WithDoc("competency-units-create", "Create competency unit", "Tambah unit kompetensi pada sebuah skema.", "Competency Units"),
	)
	huma.Patch(api, "/competency-units/{id}", handler.UpdateUnit,
		adminOnly,
		httpplatform.WithDoc("competency-units-update", "Update competency unit", "Ubah sebagian data unit kompetensi.", "Competency Units"),
	)
	huma.Delete(api, "/competency-units/{id}", handler.DeleteUnit,
		adminOnly,
		httpplatform.WithStatus(204),
		httpplatform.WithDoc("competency-units-delete", "Delete competency unit", "Hapus unit kompetensi beserta elemen & kriterianya.", "Competency Units"),
	)

	// Elemen unit kompetensi.
	huma.Get(api, "/competency-units/{unit_id}/elements", handler.ListElements,
		adminOnly,
		httpplatform.WithDoc("competency-unit-elements-list", "List competency unit elements", "Daftar elemen milik sebuah unit kompetensi.", "Competency Unit Elements"),
	)
	huma.Post(api, "/competency-units/{unit_id}/elements", handler.CreateElement,
		adminOnly,
		httpplatform.WithStatus(201),
		httpplatform.WithDoc("competency-unit-elements-create", "Create competency unit element", "Tambah elemen pada sebuah unit kompetensi.", "Competency Unit Elements"),
	)
	huma.Patch(api, "/competency-unit-elements/{id}", handler.UpdateElement,
		adminOnly,
		httpplatform.WithDoc("competency-unit-elements-update", "Update competency unit element", "Ubah sebagian data elemen unit kompetensi.", "Competency Unit Elements"),
	)
	huma.Delete(api, "/competency-unit-elements/{id}", handler.DeleteElement,
		adminOnly,
		httpplatform.WithStatus(204),
		httpplatform.WithDoc("competency-unit-elements-delete", "Delete competency unit element", "Hapus elemen beserta kriterianya.", "Competency Unit Elements"),
	)

	// Kriteria unjuk kerja (KUK).
	huma.Get(api, "/competency-unit-elements/{element_id}/criteria", handler.ListCriteria,
		adminOnly,
		httpplatform.WithDoc("competency-unit-criteria-list", "List competency unit criteria", "Daftar kriteria unjuk kerja (KUK) milik sebuah elemen.", "Competency Unit Criteria"),
	)
	huma.Post(api, "/competency-unit-elements/{element_id}/criteria", handler.CreateCriteria,
		adminOnly,
		httpplatform.WithStatus(201),
		httpplatform.WithDoc("competency-unit-criteria-create", "Create competency unit criteria", "Tambah kriteria unjuk kerja (KUK) pada sebuah elemen.", "Competency Unit Criteria"),
	)
	huma.Patch(api, "/competency-unit-criteria/{id}", handler.UpdateCriteria,
		adminOnly,
		httpplatform.WithDoc("competency-unit-criteria-update", "Update competency unit criteria", "Ubah sebagian kriteria unjuk kerja (KUK).", "Competency Unit Criteria"),
	)
	huma.Delete(api, "/competency-unit-criteria/{id}", handler.DeleteCriteria,
		adminOnly,
		httpplatform.WithStatus(204),
		httpplatform.WithDoc("competency-unit-criteria-delete", "Delete competency unit criteria", "Hapus kriteria unjuk kerja (KUK).", "Competency Unit Criteria"),
	)
}
