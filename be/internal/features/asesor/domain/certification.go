package domain

import "time"

// Certification adalah satu record sertifikasi milik asesor (relasi 1:N
// terhadap profil asesor). Record berdiri sendiri tanpa referensi master.
type Certification struct {
	ID        string
	UserID    string
	Name      string
	Issuer    string
	Number    *string
	IssuedAt  *time.Time
	ExpiredAt *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CertificationPatch struct {
	ID        string
	Name      *string
	Issuer    *string
	Number    *string
	IssuedAt  *time.Time
	ExpiredAt *time.Time
}
