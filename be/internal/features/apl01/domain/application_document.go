package domain

import "time"

// ApplicationDocument adalah bukti kelengkapan asesi untuk satu jenis dokumen
// (mis. KTP, ijazah). Jenis dokumen yang relevan ditentukan oleh skema.
type ApplicationDocument struct {
	ID             string
	ApplicationID  string
	DocumentTypeID string
	FileURL        *string
	Status         DocumentStatus
	Note           *string
	VerifiedAt     *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
