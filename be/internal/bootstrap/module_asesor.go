package bootstrap

import (
	asesorhttp "lsp/internal/features/asesor/adapters/http"
	asesorpg "lsp/internal/features/asesor/adapters/postgres"
	asesorapp "lsp/internal/features/asesor/app"
	userspg "lsp/internal/features/users/adapters/postgres"
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
