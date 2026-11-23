package domain

import "time"

// CertificationSchema adalah master data skema sertifikasi (mis. 'Gada Pratama').
type CertificationSchema struct {
	ID         string
	Code       string
	Name       string
	SchemaCode *string
	SortOrder  int
	IsActive   bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CertificationSchemaPatch struct {
	ID         string
	Code       *string
	Name       *string
	SchemaCode *string
	SortOrder  *int
	IsActive   *bool
}

// CompetencyUnit adalah unit kompetensi milik sebuah skema sertifikasi.
type CompetencyUnit struct {
	ID                    string
	CertificationSchemaID string
	Code                  string
	Title                 string
	SortOrder             int
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type CompetencyUnitPatch struct {
	ID        string
	Code      *string
	Title     *string
	SortOrder *int
}

// CompetencyUnitElement adalah elemen milik sebuah unit kompetensi.
type CompetencyUnitElement struct {
	ID               string
	CompetencyUnitID string
	Title            *string
	SortOrder        int
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type CompetencyUnitElementPatch struct {
	ID        string
	Title     *string
	SortOrder *int
}

// CompetencyUnitCriteria adalah kriteria unjuk kerja (KUK) milik sebuah elemen.
type CompetencyUnitCriteria struct {
	ID                      string
	CompetencyUnitElementID string
	Content                 string
	SortOrder               int
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type CompetencyUnitCriteriaPatch struct {
	ID        string
	Content   *string
	SortOrder *int
}

// SchemaTree adalah skema beserta seluruh strukturnya: unit -> elemen -> kriteria.
type SchemaTree struct {
	Schema CertificationSchema
	Units  []UnitNode
}

type UnitNode struct {
	Unit     CompetencyUnit
	Elements []ElementNode
}

type ElementNode struct {
	Element  CompetencyUnitElement
	Criteria []CompetencyUnitCriteria
}
