package http

import (
	"time"

	"lsp/internal/features/documentrequirements/domain"
	httpplatform "lsp/internal/platform/http"
)

// --- Input ---

type CreateDocumentTypeInput struct {
	Body DocumentTypeBody `json:"body"`
}

type UpdateDocumentTypeInput struct {
	ID   string                `path:"id"`
	Body DocumentTypePatchBody `json:"body"`
}

type GetDocumentTypeInput struct {
	ID string `path:"id"`
}

type DeleteDocumentTypeInput struct {
	ID string `path:"id"`
}

type ListDocumentTypesInput struct {
	httpplatform.PaginationQuery
}

type ListRequirementsInput struct {
	SchemaID string `path:"schema_id"`
}

type ReplaceRequirementsInput struct {
	SchemaID string                  `path:"schema_id"`
	Body     ReplaceRequirementsBody `json:"body"`
}

// --- Body ---

type DocumentTypeBody struct {
	Code             string   `json:"code" minLength:"1" maxLength:"50" doc:"Kode unik, mis. 'ktp'"`
	Label            string   `json:"label" minLength:"1" maxLength:"255" doc:"Label tampilan, mis. 'KTP'"`
	AllowedMimeTypes []string `json:"allowed_mime_types" minItems:"1" maxItems:"5" doc:"MIME yang boleh diunggah, mis. ['application/pdf']"`
}

type DocumentTypePatchBody struct {
	Code             *string  `json:"code,omitempty" minLength:"1" maxLength:"50"`
	Label            *string  `json:"label,omitempty" minLength:"1" maxLength:"255"`
	AllowedMimeTypes []string `json:"allowed_mime_types,omitempty" maxItems:"5"`
}

type ReplaceRequirementsBody struct {
	Items []RequirementItemBody `json:"items" maxItems:"200" doc:"Daftar lengkap dokumen yang dipilih untuk skema ini"`
}

type RequirementItemBody struct {
	DocumentTypeID string `json:"document_type_id" minLength:"1" doc:"ID master jenis dokumen"`
	IsRequired     bool   `json:"is_required" doc:"true = wajib, false = opsional"`
}

// --- Output ---

type DocumentTypeOutput struct {
	Body DocumentTypeDTO `json:"body"`
}

type RequirementsOutput struct {
	Body RequirementsBody `json:"body"`
}

type RequirementsBody struct {
	Data []RequirementDTO `json:"data"`
}

type DeleteOutput struct{}

// --- DTO ---

type DocumentTypeDTO struct {
	ID               string    `json:"id"`
	Code             string    `json:"code"`
	Label            string    `json:"label"`
	AllowedMimeTypes []string  `json:"allowed_mime_types"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type RequirementDTO struct {
	ID                    string          `json:"id"`
	CertificationSchemaID string          `json:"certification_schema_id"`
	DocumentTypeID        string          `json:"document_type_id"`
	IsRequired            bool            `json:"is_required"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
	DocumentType          DocumentTypeDTO `json:"document_type"`
}

// --- Konversi ---

func toDocumentTypeDTO(dt domain.DocumentType) DocumentTypeDTO {
	allowed := dt.AllowedMimeTypes
	if allowed == nil {
		allowed = []string{}
	}
	return DocumentTypeDTO{
		ID:               dt.ID,
		Code:             dt.Code,
		Label:            dt.Label,
		AllowedMimeTypes: allowed,
		CreatedAt:        dt.CreatedAt,
		UpdatedAt:        dt.UpdatedAt,
	}
}

func toRequirementDTO(d domain.RequirementDetail) RequirementDTO {
	return RequirementDTO{
		ID:                    d.ID,
		CertificationSchemaID: d.CertificationSchemaID,
		DocumentTypeID:        d.DocumentTypeID,
		IsRequired:            d.IsRequired,
		CreatedAt:             d.CreatedAt,
		UpdatedAt:             d.UpdatedAt,
		DocumentType:          toDocumentTypeDTO(d.DocumentType),
	}
}

func toRequirementDTOs(items []domain.RequirementDetail) []RequirementDTO {
	result := make([]RequirementDTO, 0, len(items))
	for _, item := range items {
		result = append(result, toRequirementDTO(item))
	}
	return result
}
