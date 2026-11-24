package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
	docreqhttp "lsp/internal/features/documentrequirements/adapters/http"
	docreqpg "lsp/internal/features/documentrequirements/adapters/postgres"
	docreqapp "lsp/internal/features/documentrequirements/app"
	"lsp/internal/platform/database"
)

type DocumentRequirementsModule struct {
	Repository *docreqpg.Repository
	Service    *docreqapp.Service
	Handler    *docreqhttp.Handler
}

func newDocumentRequirementsModule(pool *pgxpool.Pool, tx *database.TxManager) *DocumentRequirementsModule {
	repo := docreqpg.NewRepository(pool)
	service := docreqapp.NewService(repo, tx)
	handler := docreqhttp.NewHandler(service)
	return &DocumentRequirementsModule{Repository: repo, Service: service, Handler: handler}
}
