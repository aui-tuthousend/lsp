package bootstrap

import (
	"github.com/jackc/pgx/v5/pgxpool"
	schemahttp "lsp/internal/features/certificationschema/adapters/http"
	schemapg "lsp/internal/features/certificationschema/adapters/postgres"
	"lsp/internal/features/certificationschema/app"
)

type SchemaModule struct {
	Repository *schemapg.Repository
	Service    *app.Service
	Handler    *schemahttp.Handler
}

func newSchemaModule(pool *pgxpool.Pool) *SchemaModule {
	repo := schemapg.NewRepository(pool)
	service := app.NewService(repo)
	handler := schemahttp.NewHandler(service)
	return &SchemaModule{Repository: repo, Service: service, Handler: handler}
}
