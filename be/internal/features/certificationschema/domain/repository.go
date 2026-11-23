package domain

import "context"

type ListParams struct {
	Limit  int
	Offset int
}

type SchemaListResult struct {
	Items []CertificationSchema
	Total int
}

type UnitListResult struct {
	Items []CompetencyUnit
	Total int
}

type ElementListResult struct {
	Items []CompetencyUnitElement
	Total int
}

type CriteriaListResult struct {
	Items []CompetencyUnitCriteria
	Total int
}

type Repository interface {
	// Skema sertifikasi.
	CreateSchema(ctx context.Context, schema CertificationSchema) error
	PatchSchema(ctx context.Context, patch CertificationSchemaPatch) (CertificationSchema, error)
	FindSchemaByID(ctx context.Context, id string) (CertificationSchema, error)
	ListSchemas(ctx context.Context, params ListParams) (SchemaListResult, error)
	DeleteSchema(ctx context.Context, id string) error
	SchemaTree(ctx context.Context, id string) (SchemaTree, error)

	// Unit kompetensi.
	CreateUnit(ctx context.Context, unit CompetencyUnit) error
	PatchUnit(ctx context.Context, patch CompetencyUnitPatch) (CompetencyUnit, error)
	FindUnitByID(ctx context.Context, id string) (CompetencyUnit, error)
	ListUnitsBySchema(ctx context.Context, schemaID string, params ListParams) (UnitListResult, error)
	DeleteUnit(ctx context.Context, id string) error

	// Elemen unit kompetensi.
	CreateElement(ctx context.Context, element CompetencyUnitElement) error
	PatchElement(ctx context.Context, patch CompetencyUnitElementPatch) (CompetencyUnitElement, error)
	FindElementByID(ctx context.Context, id string) (CompetencyUnitElement, error)
	ListElementsByUnit(ctx context.Context, unitID string, params ListParams) (ElementListResult, error)
	DeleteElement(ctx context.Context, id string) error

	// Kriteria unjuk kerja (KUK).
	CreateCriteria(ctx context.Context, criteria CompetencyUnitCriteria) error
	PatchCriteria(ctx context.Context, patch CompetencyUnitCriteriaPatch) (CompetencyUnitCriteria, error)
	FindCriteriaByID(ctx context.Context, id string) (CompetencyUnitCriteria, error)
	ListCriteriaByElement(ctx context.Context, elementID string, params ListParams) (CriteriaListResult, error)
	DeleteCriteria(ctx context.Context, id string) error
}
