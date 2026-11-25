package app

import (
	"context"

	"lsp/internal/features/documentrequirements/domain"
	"lsp/internal/platform/database"
	"lsp/internal/platform/id"
)

type Service struct {
	repo domain.Repository
	tx   *database.TxManager
}

func NewService(repo domain.Repository, tx *database.TxManager) *Service {
	return &Service{repo: repo, tx: tx}
}

// --- Master jenis dokumen ---

type CreateDocumentTypeCommand struct {
	Code             string
	Label            string
	AllowedMimeTypes []string
}

type PatchDocumentTypeCommand struct {
	ID               string
	Code             *string
	Label            *string
	AllowedMimeTypes []string
}

func (s *Service) CreateDocumentType(ctx context.Context, cmd CreateDocumentTypeCommand) (domain.DocumentType, error) {
	docType := domain.DocumentType{
		ID:               id.New(),
		Code:             cmd.Code,
		Label:            cmd.Label,
		AllowedMimeTypes: cmd.AllowedMimeTypes,
	}
	if err := s.repo.CreateDocumentType(ctx, docType); err != nil {
		return domain.DocumentType{}, err
	}
	return s.repo.FindDocumentTypeByID(ctx, docType.ID)
}

func (s *Service) PatchDocumentType(ctx context.Context, cmd PatchDocumentTypeCommand) (domain.DocumentType, error) {
	return s.repo.PatchDocumentType(ctx, domain.DocumentTypePatch{
		ID:               cmd.ID,
		Code:             cmd.Code,
		Label:            cmd.Label,
		AllowedMimeTypes: cmd.AllowedMimeTypes,
	})
}

func (s *Service) GetDocumentType(ctx context.Context, id string) (domain.DocumentType, error) {
	return s.repo.FindDocumentTypeByID(ctx, id)
}

func (s *Service) ListDocumentTypes(ctx context.Context, limit, offset int) (domain.DocumentTypeListResult, error) {
	return s.repo.ListDocumentTypes(ctx, domain.ListParams{Limit: limit, Offset: offset})
}

func (s *Service) DeleteDocumentType(ctx context.Context, id string) error {
	return s.repo.DeleteDocumentType(ctx, id)
}

// --- Persyaratan dokumen per skema ---

// ReplaceRequirementItem mewakili satu jenis dokumen yang dipilih untuk skema.
type ReplaceRequirementItem struct {
	DocumentTypeID string
	IsRequired     bool
}

// ReplaceRequirementsCommand menggantikan seluruh set persyaratan dokumen sebuah
// skema dengan daftar yang dikirim (set-semantik, idempoten).
type ReplaceRequirementsCommand struct {
	CertificationSchemaID string
	Items                 []ReplaceRequirementItem
}

func (s *Service) ReplaceRequirements(ctx context.Context, cmd ReplaceRequirementsCommand) ([]domain.RequirementDetail, error) {
	documentTypeIDs := make([]string, 0, len(cmd.Items))
	for _, item := range cmd.Items {
		documentTypeIDs = append(documentTypeIDs, item.DocumentTypeID)
	}

	err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.DeleteRequirementsNotIn(txCtx, cmd.CertificationSchemaID, documentTypeIDs); err != nil {
			return err
		}
		for _, item := range cmd.Items {
			if err := s.repo.UpsertRequirement(txCtx, domain.SchemaDocumentRequirement{
				ID:                    id.New(),
				CertificationSchemaID: cmd.CertificationSchemaID,
				DocumentTypeID:        item.DocumentTypeID,
				IsRequired:            item.IsRequired,
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return s.repo.ListRequirementsBySchema(ctx, cmd.CertificationSchemaID)
}

func (s *Service) ListRequirements(ctx context.Context, schemaID string) ([]domain.RequirementDetail, error) {
	return s.repo.ListRequirementsBySchema(ctx, schemaID)
}
