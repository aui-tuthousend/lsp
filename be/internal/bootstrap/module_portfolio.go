package bootstrap

import "lsp/internal/platform/storage"

type PortfolioModule struct {
	Storage storage.ObjectStorage
}

func newPortfolioModule(objStorage storage.ObjectStorage) *PortfolioModule {
	return &PortfolioModule{Storage: objStorage}
}
