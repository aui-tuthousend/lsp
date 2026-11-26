package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
	apl02http "lsp/internal/features/apl02/adapters/http"
	apl02pg "lsp/internal/features/apl02/adapters/postgres"
	apl02app "lsp/internal/features/apl02/app"
	"lsp/internal/platform/database"
	"lsp/internal/platform/events"
	"lsp/internal/platform/storage"
)

type APL02Module struct {
	Repository *apl02pg.Repository
	Service    *apl02app.Service
	Handler    *apl02http.Handler
}

func newAPL02Module(pool *pgxpool.Pool, tx *database.TxManager, bus events.Bus, objStorage storage.ObjectStorage) *APL02Module {
	repo := apl02pg.NewRepository(pool)
	service := apl02app.NewService(repo, tx, bus)
	handler := apl02http.NewHandler(service, objStorage)
	return &APL02Module{Repository: repo, Service: service, Handler: handler}
}
