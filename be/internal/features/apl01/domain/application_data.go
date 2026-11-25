package domain

import "time"

// ApplicationData adalah data sertifikasi yang dipilih asesi: skema sertifikasi
// beserta tujuan asesmen (relasi 1:1 ke permohonan).
type ApplicationData struct {
	ID                    string
	ApplicationID         string
	CertificationSchemaID string
	Purpose               Purpose
	CreatedAt             time.Time
	UpdatedAt             time.Time
}
