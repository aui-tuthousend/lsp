package http

import (
	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/platform/auth"
	httpplatform "lsp/internal/platform/http"
)

func RegisterRoutes(api huma.API, handler *Handler) {
	// Publik karena dipakai saat pendaftaran (belum ada sesi). Dibatasi oleh
	// allowlist purpose/content-type dan masa berlaku URL yang singkat.
	huma.Post(api, "/uploads/presign", handler.Presign,
		auth.PublicOperation,
		httpplatform.WithDoc("uploads-presign", "Presign upload", "Membuat presigned URL untuk unggah dokumen langsung ke object storage.", "Uploads"),
	)
}
