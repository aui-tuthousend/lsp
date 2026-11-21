package http

import (
	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/users/domain"
	"lsp/internal/platform/auth"
	httpplatform "lsp/internal/platform/http"
)

func RegisterRoutes(api huma.API, handler *Handler) {
	huma.Post(api, "/auth/register", handler.Register,
		auth.PublicOperation,
		httpplatform.WithDoc("auth-register", "Register user", "Pendaftaran mandiri; selalu membuat akun dengan role asesi.", "Auth"),
	)
	huma.Post(api, "/auth/login", handler.Login,
		auth.PublicOperation,
		httpplatform.WithDoc("auth-login", "Login user", "Authenticate user and return JWT token.", "Auth"),
	)
	huma.Get(api, "/auth/me", handler.GetMe,
		httpplatform.WithDoc("auth-me", "Get current user", "Return current authenticated user profile.", "Auth"),
	)
	huma.Post(api, "/users", handler.Create,
		auth.RequireRoles(string(domain.RoleAdmin)),
		httpplatform.WithStatus(201),
		httpplatform.WithDoc("users-create", "Create user", "Membuat user baru dengan role apa pun (khusus admin).", "Users"),
	)
	huma.Get(api, "/users", handler.List,
		auth.RequireRoles(string(domain.RoleAdmin)),
		httpplatform.WithDoc("users-list", "List users", "Return paginated users list.", "Users"),
	)
}
