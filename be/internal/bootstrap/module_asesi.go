package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
	asesihttp "lsp/internal/features/asesi/adapters/http"
	asesipg "lsp/internal/features/asesi/adapters/postgres"
	asesiapp "lsp/internal/features/asesi/app"
	userspg "lsp/internal/features/users/adapters/postgres"
	"lsp/internal/platform/storage"
)

type AsesiModule struct {
	Repository *asesipg.Repository
	Service    *asesiapp.Service
	Handler    *asesihttp.Handler
}

func newAsesiModule(pool *pgxpool.Pool, usersRepo *userspg.Repository, objStorage storage.ObjectStorage) *AsesiModule {
	repo := asesipg.NewRepository(pool)
	service := asesiapp.NewService(repo, usersRepo, objStorage)
	handler := asesihttp.NewHandler(service, objStorage)
	return &AsesiModule{Repository: repo, Service: service, Handler: handler}
}
