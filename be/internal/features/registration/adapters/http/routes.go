package http

import (
	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/platform/auth"
	httpplatform "lsp/internal/platform/http"
)

func RegisterRoutes(api huma.API, handler *Handler) {
	huma.Post(api, "/auth/register", handler.Register,
		auth.PublicOperation,
		httpplatform.WithStatus(201),
		httpplatform.WithDoc("auth-register", "Register asesi", "Pendaftaran mandiri: membuat user (role asesi) dan profil asesi dalam satu transaksi.", "Auth"),
	)
}
