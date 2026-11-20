package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
	tukhttp "lsp/internal/tuk/adapters/http"
	tukpg "lsp/internal/tuk/adapters/postgres"
	"lsp/internal/tuk/app"
)

type TUKModule struct {
	Repository *tukpg.Repository
	Service    *app.Service
	Handler    *tukhttp.Handler
}

func newTUKModule(pool *pgxpool.Pool) *TUKModule {
	repo := tukpg.NewRepository(pool)
	service := app.NewService(repo)
	handler := tukhttp.NewHandler(service)
	return &TUKModule{Repository: repo, Service: service, Handler: handler}
}
