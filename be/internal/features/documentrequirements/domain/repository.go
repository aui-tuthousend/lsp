package domain

import "context"

type ListParams struct {
	Limit  int
	Offset int
}

type DocumentTypeListResult struct {
	Items []DocumentType
	Total int
}

type Repository interface {
	// Master jenis dokumen.
	CreateDocumentType(ctx context.Context, docType DocumentType) error
	PatchDocumentType(ctx context.Context, patch DocumentTypePatch) (DocumentType, error)
	FindDocumentTypeByID(ctx context.Context, id string) (DocumentType, error)
	ListDocumentTypes(ctx context.Context, params ListParams) (DocumentTypeListResult, error)
	DeleteDocumentType(ctx context.Context, id string) error

	// Persyaratan dokumen per skema. Penggantian set dilakukan dalam transaksi:
	// hapus yang tidak lagi dipilih lalu upsert yang dipilih.
	ListRequirementsBySchema(ctx context.Context, schemaID string) ([]RequirementDetail, error)
	DeleteRequirementsNotIn(ctx context.Context, schemaID string, documentTypeIDs []string) error
	UpsertRequirement(ctx context.Context, requirement SchemaDocumentRequirement) error
}
