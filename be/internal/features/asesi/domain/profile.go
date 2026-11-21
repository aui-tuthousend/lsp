package domain

import "time"

type Gender string

const (
	GenderMale   Gender = "LAKI-LAKI"
	GenderFemale Gender = "PEREMPUAN"
)

type EducationLevel string

const (
	EducationSD     EducationLevel = "SD"
	EducationSMP    EducationLevel = "SMP"
	EducationSMASMK EducationLevel = "SMA_SMK"
	EducationD1     EducationLevel = "D1"
	EducationD2     EducationLevel = "D2"
	EducationD3     EducationLevel = "D3"
	EducationD4     EducationLevel = "D4"
	EducationS1     EducationLevel = "S1"
	EducationS2     EducationLevel = "S2"
	EducationS3     EducationLevel = "S3"
)

type Profile struct {
	ID                string
	UserID            string
	FullName          string
	NIK               string
	Phone             string
	Gender            Gender
	BirthPlace        string
	BirthDate         time.Time
	Nationality       string
	Occupation        string
	EducationLevel    EducationLevel
	KTPAddress        string
	DomicileSameAsKTP bool
	DomicileAddress   *string
	Province          string
	City              string
	District          string
	Village           string
	PostalCode        *string
	Photo             *string
	KTPFile           *string
	CVFile            *string
	NPWPFile          *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Patch struct {
	ID                string
	FullName          *string
	NIK               *string
	Phone             *string
	Gender            *Gender
	BirthPlace        *string
	BirthDate         *time.Time
	Nationality       *string
	Occupation        *string
	EducationLevel    *EducationLevel
	KTPAddress        *string
	DomicileSameAsKTP *bool
	DomicileAddress   *string
	Province          *string
	City              *string
	District          *string
	Village           *string
	PostalCode        *string
	Photo             *string
	KTPFile           *string
	CVFile            *string
	NPWPFile          *string
}
