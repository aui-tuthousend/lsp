package http

import (
	"time"

	"lsp/internal/features/certificationschema/domain"
	httpplatform "lsp/internal/platform/http"
)

// --- Input ---

type CreateSchemaInput struct {
	Body SchemaBody `json:"body"`
}

type UpdateSchemaInput struct {
	ID   string          `path:"id"`
	Body SchemaPatchBody `json:"body"`
}

type GetSchemaInput struct {
	ID string `path:"id"`
}

type DeleteSchemaInput struct {
	ID string `path:"id"`
}

type ListSchemasInput struct {
	httpplatform.PaginationQuery
}

type ListUnitsInput struct {
	SchemaID string `path:"schema_id"`
	httpplatform.PaginationQuery
}

type CreateUnitInput struct {
	SchemaID string   `path:"schema_id"`
	Body     UnitBody `json:"body"`
}

type UpdateUnitInput struct {
	ID   string        `path:"id"`
	Body UnitPatchBody `json:"body"`
}

type DeleteUnitInput struct {
	ID string `path:"id"`
}

type ListElementsInput struct {
	UnitID string `path:"unit_id"`
	httpplatform.PaginationQuery
}

type CreateElementInput struct {
	UnitID string      `path:"unit_id"`
	Body   ElementBody `json:"body"`
}

type UpdateElementInput struct {
	ID   string           `path:"id"`
	Body ElementPatchBody `json:"body"`
}

type DeleteElementInput struct {
	ID string `path:"id"`
}

type ListCriteriaInput struct {
	ElementID string `path:"element_id"`
	httpplatform.PaginationQuery
}

type CreateCriteriaInput struct {
	ElementID string       `path:"element_id"`
	Body      CriteriaBody `json:"body"`
}

type UpdateCriteriaInput struct {
	ID   string            `path:"id"`
	Body CriteriaPatchBody `json:"body"`
}

type DeleteCriteriaInput struct {
	ID string `path:"id"`
}

// --- Body ---

type SchemaBody struct {
	Code       string  `json:"code" minLength:"1" maxLength:"50" doc:"Slug internal, mis. 'gada_pratama'"`
	Name       string  `json:"name" minLength:"1" maxLength:"255" doc:"Nama skema, mis. 'Gada Pratama'"`
	SchemaCode *string `json:"schema_code,omitempty" maxLength:"100" doc:"Nomor resmi BNSP (opsional)"`
	SortOrder  int     `json:"sort_order,omitempty" minimum:"0"`
	IsActive   *bool   `json:"is_active,omitempty" doc:"Default true bila tidak diisi"`
}

type SchemaPatchBody struct {
	Code       *string `json:"code,omitempty" minLength:"1" maxLength:"50"`
	Name       *string `json:"name,omitempty" minLength:"1" maxLength:"255"`
	SchemaCode *string `json:"schema_code,omitempty" maxLength:"100"`
	SortOrder  *int    `json:"sort_order,omitempty" minimum:"0"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

type UnitBody struct {
	Code      string `json:"code" minLength:"1" maxLength:"50" doc:"Kode unit, mis. 'N.80PAM00.001.2'"`
	Title     string `json:"title" minLength:"1" maxLength:"255"`
	SortOrder int    `json:"sort_order,omitempty" minimum:"0"`
}

type UnitPatchBody struct {
	Code      *string `json:"code,omitempty" minLength:"1" maxLength:"50"`
	Title     *string `json:"title,omitempty" minLength:"1" maxLength:"255"`
	SortOrder *int    `json:"sort_order,omitempty" minimum:"0"`
}

type ElementBody struct {
	Title     *string `json:"title,omitempty" maxLength:"255"`
	SortOrder int     `json:"sort_order,omitempty" minimum:"0"`
}

type ElementPatchBody struct {
	Title     *string `json:"title,omitempty" maxLength:"255"`
	SortOrder *int    `json:"sort_order,omitempty" minimum:"0"`
}

type CriteriaBody struct {
	Content   string `json:"content" minLength:"1" doc:"Teks kriteria unjuk kerja (KUK)"`
	SortOrder int    `json:"sort_order,omitempty" minimum:"0"`
}

type CriteriaPatchBody struct {
	Content   *string `json:"content,omitempty" minLength:"1"`
	SortOrder *int    `json:"sort_order,omitempty" minimum:"0"`
}

// --- Output ---

type SchemaOutput struct {
	Body SchemaDTO `json:"body"`
}

type SchemaTreeOutput struct {
	Body SchemaTreeDTO `json:"body"`
}

type UnitOutput struct {
	Body UnitDTO `json:"body"`
}

type ElementOutput struct {
	Body ElementDTO `json:"body"`
}

type CriteriaOutput struct {
	Body CriteriaDTO `json:"body"`
}

type DeleteOutput struct{}

// --- DTO ---

type SchemaDTO struct {
	ID         string    `json:"id"`
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	SchemaCode *string   `json:"schema_code,omitempty"`
	SortOrder  int       `json:"sort_order"`
	IsActive   bool      `json:"is_active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UnitDTO struct {
	ID                    string    `json:"id"`
	CertificationSchemaID string    `json:"certification_schema_id"`
	Code                  string    `json:"code"`
	Title                 string    `json:"title"`
	SortOrder             int       `json:"sort_order"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type ElementDTO struct {
	ID               string    `json:"id"`
	CompetencyUnitID string    `json:"competency_unit_id"`
	Title            *string   `json:"title,omitempty"`
	SortOrder        int       `json:"sort_order"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CriteriaDTO struct {
	ID                      string    `json:"id"`
	CompetencyUnitElementID string    `json:"competency_unit_element_id"`
	Content                 string    `json:"content"`
	SortOrder               int       `json:"sort_order"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

// SchemaTreeDTO meng-inline field skema dan menambahkan struktur unit.
type SchemaTreeDTO struct {
	SchemaDTO
	Units []UnitNodeDTO `json:"units"`
}

type UnitNodeDTO struct {
	UnitDTO
	Elements []ElementNodeDTO `json:"elements"`
}

type ElementNodeDTO struct {
	ElementDTO
	Criteria []CriteriaDTO `json:"criteria"`
}

// --- Konversi ---

func toSchemaDTO(s domain.CertificationSchema) SchemaDTO {
	return SchemaDTO{
		ID:         s.ID,
		Code:       s.Code,
		Name:       s.Name,
		SchemaCode: s.SchemaCode,
		SortOrder:  s.SortOrder,
		IsActive:   s.IsActive,
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
	}
}

func toUnitDTO(u domain.CompetencyUnit) UnitDTO {
	return UnitDTO{
		ID:                    u.ID,
		CertificationSchemaID: u.CertificationSchemaID,
		Code:                  u.Code,
		Title:                 u.Title,
		SortOrder:             u.SortOrder,
		CreatedAt:             u.CreatedAt,
		UpdatedAt:             u.UpdatedAt,
	}
}

func toElementDTO(e domain.CompetencyUnitElement) ElementDTO {
	return ElementDTO{
		ID:               e.ID,
		CompetencyUnitID: e.CompetencyUnitID,
		Title:            e.Title,
		SortOrder:        e.SortOrder,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
	}
}

func toCriteriaDTO(c domain.CompetencyUnitCriteria) CriteriaDTO {
	return CriteriaDTO{
		ID:                      c.ID,
		CompetencyUnitElementID: c.CompetencyUnitElementID,
		Content:                 c.Content,
		SortOrder:               c.SortOrder,
		CreatedAt:               c.CreatedAt,
		UpdatedAt:               c.UpdatedAt,
	}
}

func toSchemaTreeDTO(tree domain.SchemaTree) SchemaTreeDTO {
	units := make([]UnitNodeDTO, 0, len(tree.Units))
	for _, unitNode := range tree.Units {
		elements := make([]ElementNodeDTO, 0, len(unitNode.Elements))
		for _, elementNode := range unitNode.Elements {
			criteria := make([]CriteriaDTO, 0, len(elementNode.Criteria))
			for _, item := range elementNode.Criteria {
				criteria = append(criteria, toCriteriaDTO(item))
			}
			elements = append(elements, ElementNodeDTO{ElementDTO: toElementDTO(elementNode.Element), Criteria: criteria})
		}
		units = append(units, UnitNodeDTO{UnitDTO: toUnitDTO(unitNode.Unit), Elements: elements})
	}
	return SchemaTreeDTO{SchemaDTO: toSchemaDTO(tree.Schema), Units: units}
}
