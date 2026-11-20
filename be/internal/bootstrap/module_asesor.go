package bootstrap

import (
	asesorhttp "lsp/internal/asesor/adapters/http"
	asesorpg "lsp/internal/asesor/adapters/postgres"
	asesorapp "lsp/internal/asesor/app"
	userspg "lsp/internal/users/adapters/postgres"
)

type AsesorModule struct {
	Repository *asesorpg.Repository
	Service    *asesorapp.Service
	Handler    *asesorhttp.Handler
}

func newAsesorModule(repo *asesorpg.Repository, usersRepo *userspg.Repository) *AsesorModule {
	service := asesorapp.NewService(repo, usersRepo)
	handler := asesorhttp.NewHandler(service)
	return &AsesorModule{Repository: repo, Service: service, Handler: handler}
}
