package app

import (
	"context"

	"lsp/internal/features/certificationschema/domain"
	"lsp/internal/platform/id"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

// --- Skema sertifikasi ---

type CreateSchemaCommand struct {
	Code       string
	Name       string
	SchemaCode *string
	SortOrder  int
	IsActive   *bool
}

type PatchSchemaCommand struct {
	ID         string
	Code       *string
	Name       *string
	SchemaCode *string
	SortOrder  *int
	IsActive   *bool
}

func (s *Service) CreateSchema(ctx context.Context, cmd CreateSchemaCommand) (domain.CertificationSchema, error) {
	isActive := true
	if cmd.IsActive != nil {
		isActive = *cmd.IsActive
	}
	schema := domain.CertificationSchema{
		ID:         id.New(),
		Code:       cmd.Code,
		Name:       cmd.Name,
		SchemaCode: cmd.SchemaCode,
		SortOrder:  cmd.SortOrder,
		IsActive:   isActive,
	}
	if err := s.repo.CreateSchema(ctx, schema); err != nil {
		return domain.CertificationSchema{}, err
	}
	return s.repo.FindSchemaByID(ctx, schema.ID)
}

func (s *Service) PatchSchema(ctx context.Context, cmd PatchSchemaCommand) (domain.CertificationSchema, error) {
	return s.repo.PatchSchema(ctx, domain.CertificationSchemaPatch{
		ID:         cmd.ID,
		Code:       cmd.Code,
		Name:       cmd.Name,
		SchemaCode: cmd.SchemaCode,
		SortOrder:  cmd.SortOrder,
		IsActive:   cmd.IsActive,
	})
}

func (s *Service) GetSchemaTree(ctx context.Context, id string) (domain.SchemaTree, error) {
	return s.repo.SchemaTree(ctx, id)
}

func (s *Service) ListSchemas(ctx context.Context, limit, offset int) (domain.SchemaListResult, error) {
	return s.repo.ListSchemas(ctx, domain.ListParams{Limit: limit, Offset: offset})
}

func (s *Service) DeleteSchema(ctx context.Context, id string) error {
	return s.repo.DeleteSchema(ctx, id)
}

// --- Unit kompetensi ---

type CreateUnitCommand struct {
	CertificationSchemaID string
	Code                  string
	Title                 string
	SortOrder             int
}

type PatchUnitCommand struct {
	ID        string
	Code      *string
	Title     *string
	SortOrder *int
}

func (s *Service) CreateUnit(ctx context.Context, cmd CreateUnitCommand) (domain.CompetencyUnit, error) {
	unit := domain.CompetencyUnit{
		ID:                    id.New(),
		CertificationSchemaID: cmd.CertificationSchemaID,
		Code:                  cmd.Code,
		Title:                 cmd.Title,
		SortOrder:             cmd.SortOrder,
	}
	if err := s.repo.CreateUnit(ctx, unit); err != nil {
		return domain.CompetencyUnit{}, err
	}
	return s.repo.FindUnitByID(ctx, unit.ID)
}

func (s *Service) PatchUnit(ctx context.Context, cmd PatchUnitCommand) (domain.CompetencyUnit, error) {
	return s.repo.PatchUnit(ctx, domain.CompetencyUnitPatch{ID: cmd.ID, Code: cmd.Code, Title: cmd.Title, SortOrder: cmd.SortOrder})
}

func (s *Service) ListUnits(ctx context.Context, schemaID string, limit, offset int) (domain.UnitListResult, error) {
	return s.repo.ListUnitsBySchema(ctx, schemaID, domain.ListParams{Limit: limit, Offset: offset})
}

func (s *Service) DeleteUnit(ctx context.Context, id string) error {
	return s.repo.DeleteUnit(ctx, id)
}

// --- Elemen unit kompetensi ---

type CreateElementCommand struct {
	CompetencyUnitID string
	Title            *string
	SortOrder        int
}

type PatchElementCommand struct {
	ID        string
	Title     *string
	SortOrder *int
}

func (s *Service) CreateElement(ctx context.Context, cmd CreateElementCommand) (domain.CompetencyUnitElement, error) {
	element := domain.CompetencyUnitElement{
		ID:               id.New(),
		CompetencyUnitID: cmd.CompetencyUnitID,
		Title:            cmd.Title,
		SortOrder:        cmd.SortOrder,
	}
	if err := s.repo.CreateElement(ctx, element); err != nil {
		return domain.CompetencyUnitElement{}, err
	}
	return s.repo.FindElementByID(ctx, element.ID)
}

func (s *Service) PatchElement(ctx context.Context, cmd PatchElementCommand) (domain.CompetencyUnitElement, error) {
	return s.repo.PatchElement(ctx, domain.CompetencyUnitElementPatch{ID: cmd.ID, Title: cmd.Title, SortOrder: cmd.SortOrder})
}

func (s *Service) ListElements(ctx context.Context, unitID string, limit, offset int) (domain.ElementListResult, error) {
	return s.repo.ListElementsByUnit(ctx, unitID, domain.ListParams{Limit: limit, Offset: offset})
}

func (s *Service) DeleteElement(ctx context.Context, id string) error {
	return s.repo.DeleteElement(ctx, id)
}

// --- Kriteria unjuk kerja (KUK) ---

type CreateCriteriaCommand struct {
	CompetencyUnitElementID string
	Content                 string
	SortOrder               int
}

type PatchCriteriaCommand struct {
	ID        string
	Content   *string
	SortOrder *int
}

func (s *Service) CreateCriteria(ctx context.Context, cmd CreateCriteriaCommand) (domain.CompetencyUnitCriteria, error) {
	criteria := domain.CompetencyUnitCriteria{
		ID:                      id.New(),
		CompetencyUnitElementID: cmd.CompetencyUnitElementID,
		Content:                 cmd.Content,
		SortOrder:               cmd.SortOrder,
	}
	if err := s.repo.CreateCriteria(ctx, criteria); err != nil {
		return domain.CompetencyUnitCriteria{}, err
	}
	return s.repo.FindCriteriaByID(ctx, criteria.ID)
}

func (s *Service) PatchCriteria(ctx context.Context, cmd PatchCriteriaCommand) (domain.CompetencyUnitCriteria, error) {
	return s.repo.PatchCriteria(ctx, domain.CompetencyUnitCriteriaPatch{ID: cmd.ID, Content: cmd.Content, SortOrder: cmd.SortOrder})
}

func (s *Service) ListCriteria(ctx context.Context, elementID string, limit, offset int) (domain.CriteriaListResult, error) {
	return s.repo.ListCriteriaByElement(ctx, elementID, domain.ListParams{Limit: limit, Offset: offset})
}

func (s *Service) DeleteCriteria(ctx context.Context, id string) error {
	return s.repo.DeleteCriteria(ctx, id)
}
