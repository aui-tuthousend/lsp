package http

import (
	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/platform/auth"
	httpplatform "lsp/internal/platform/http"
)

func RegisterRoutes(api huma.API, handler *Handler) {
	huma.Post(api, "/auth/register", handler.Register,
		auth.PublicOperation,
		httpplatform.WithDoc("auth-register", "Register user", "Register a new user account.", "Auth"),
	)
	huma.Post(api, "/auth/login", handler.Login,
		auth.PublicOperation,
		httpplatform.WithDoc("auth-login", "Login user", "Authenticate user and return JWT token.", "Auth"),
	)
	huma.Get(api, "/auth/me", handler.GetMe,
		httpplatform.WithDoc("auth-me", "Get current user", "Return current authenticated user profile.", "Auth"),
	)
	huma.Get(api, "/users", handler.List,
		httpplatform.WithDoc("users-list", "List users", "Return paginated users list.", "Users"),
	)
}
