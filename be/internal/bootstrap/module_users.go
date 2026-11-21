package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
	asesorpg "lsp/internal/features/asesor/adapters/postgres"
	usershttp "lsp/internal/features/users/adapters/http"
	userspg "lsp/internal/features/users/adapters/postgres"
	"lsp/internal/features/users/app"
	"lsp/internal/platform/auth"
)

type UsersModule struct {
	Repository *userspg.Repository
	Roles      *userspg.RolesRepository
	Service    *app.Service
	Handler    *usershttp.Handler
}

func newUsersModule(pool *pgxpool.Pool, asesorInit *asesorpg.Repository, jwtProvider *auth.JWTProvider) *UsersModule {
	repo := userspg.NewRepository(pool)
	roles := userspg.NewRolesRepository(pool)
	service := app.NewService(repo, roles, asesorInit, jwtProvider)
	handler := usershttp.NewHandler(service)
	return &UsersModule{Repository: repo, Roles: roles, Service: service, Handler: handler}
}
