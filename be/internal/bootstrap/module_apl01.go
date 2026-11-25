package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
	apl01http "lsp/internal/features/apl01/adapters/http"
	apl01pg "lsp/internal/features/apl01/adapters/postgres"
	apl01app "lsp/internal/features/apl01/app"
	"lsp/internal/platform/database"
	"lsp/internal/platform/events"
)

type APL01Module struct {
	Repository *apl01pg.Repository
	Service    *apl01app.Service
	Handler    *apl01http.Handler
}

func newAPL01Module(pool *pgxpool.Pool, tx *database.TxManager, bus events.Bus) *APL01Module {
	repo := apl01pg.NewRepository(pool)
	service := apl01app.NewService(repo, tx, bus)
	handler := apl01http.NewHandler(service)
	return &APL01Module{Repository: repo, Service: service, Handler: handler}
}
