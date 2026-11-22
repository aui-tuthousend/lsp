package bootstrap

import (
	asesiapp "lsp/internal/features/asesi/app"
	registrationhttp "lsp/internal/features/registration/adapters/http"
	registrationapp "lsp/internal/features/registration/app"
	usersapp "lsp/internal/features/users/app"
	"lsp/internal/platform/database"
)

type RegistrationModule struct {
	Service *registrationapp.Service
	Handler *registrationhttp.Handler
}

func newRegistrationModule(users *usersapp.Service, asesi *asesiapp.Service, tx *database.TxManager) *RegistrationModule {
	service := registrationapp.NewService(users, asesi, tx)
	handler := registrationhttp.NewHandler(service)
	return &RegistrationModule{Service: service, Handler: handler}
}
