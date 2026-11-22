package bootstrap

import (
	uploadshttp "lsp/internal/features/uploads/adapters/http"
	uploadsapp "lsp/internal/features/uploads/app"
	"lsp/internal/platform/storage"
)

type UploadsModule struct {
	Service *uploadsapp.Service
	Handler *uploadshttp.Handler
}

func newUploadsModule(objStorage storage.ObjectStorage) *UploadsModule {
	service := uploadsapp.NewService(objStorage)
	handler := uploadshttp.NewHandler(service)
	return &UploadsModule{Service: service, Handler: handler}
}
