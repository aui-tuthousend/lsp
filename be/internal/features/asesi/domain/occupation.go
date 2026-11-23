package domain

import "time"

// Occupation adalah data pekerjaan terkini asesi (relasi 1:1 ke asesi_profiles).
type Occupation struct {
	ID                string
	AsesiProfileID    string
	CompanyName       string
	Position          string
	TenureMonths      *int
	SubordinatesCount *int
	Address           string
	Province          string
	City              string
	District          string
	PostalCode        *string
	OfficeEmail       *string
	OfficePhone       *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// OccupationPatch adalah perubahan parsial data pekerjaan (kirim field yang berubah saja).
type OccupationPatch struct {
	ID                string
	CompanyName       *string
	Position          *string
	TenureMonths      *int
	SubordinatesCount *int
	Address           *string
	Province          *string
	City              *string
	District          *string
	PostalCode        *string
	OfficeEmail       *string
	OfficePhone       *string
}
