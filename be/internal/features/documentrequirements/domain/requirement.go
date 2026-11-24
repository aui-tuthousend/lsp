package domain

import "time"

// SchemaDocumentRequirement adalah jenis dokumen yang dipilih untuk sebuah skema
// sertifikasi. IsRequired menandai dokumen wajib/opsional. Menjadi penentu
// "dokumen relevan" pada APL.02.
type SchemaDocumentRequirement struct {
	ID                    string
	CertificationSchemaID string
	DocumentTypeID        string
	IsRequired            bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// RequirementDetail menyertakan master jenis dokumen agar konsumen tidak perlu
// melakukan join sendiri.
type RequirementDetail struct {
	SchemaDocumentRequirement
	DocumentType DocumentType
}
