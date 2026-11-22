package http

import asesidomain "lsp/internal/features/asesi/domain"

type RegisterInput struct {
	Body RegisterBody `json:"body"`
}

// RegisterBody adalah pendaftaran mandiri: data akun + profil asesi sekaligus.
type RegisterBody struct {
	Email    string `json:"email" format:"email" maxLength:"255"`
	Password string `json:"password" minLength:"8" maxLength:"72"`

	FullName          string                     `json:"full_name" minLength:"1" maxLength:"255"`
	Phone             string                     `json:"phone" minLength:"5" maxLength:"20"`
	NIK               string                     `json:"nik" minLength:"16" maxLength:"16"`
	Gender            asesidomain.Gender         `json:"gender" enum:"LAKI-LAKI,PEREMPUAN"`
	BirthPlace        string                     `json:"birth_place" minLength:"1" maxLength:"255"`
	BirthDate         string                     `json:"birth_date" pattern:"^\\d{4}-\\d{2}-\\d{2}$" doc:"Tanggal lahir (YYYY-MM-DD)"`
	Nationality       string                     `json:"nationality" minLength:"1" maxLength:"100" default:"Indonesia"`
	Occupation        string                     `json:"occupation" minLength:"1" maxLength:"255"`
	EducationLevel    asesidomain.EducationLevel `json:"education_level" enum:"SD,SMP,SMA_SMK,D1,D2,D3,D4,S1,S2,S3"`
	KTPAddress        string                     `json:"ktp_address" minLength:"1"`
	DomicileSameAsKTP bool                       `json:"domicile_same_as_ktp"`
	DomicileAddress   *string                    `json:"domicile_address,omitempty"`
	Province          string                     `json:"province" minLength:"1" maxLength:"255"`
	City              string                     `json:"city" minLength:"1" maxLength:"255"`
	District          string                     `json:"district" minLength:"1" maxLength:"255"`
	Village           string                     `json:"village" minLength:"1" maxLength:"255"`
	PostalCode        *string                    `json:"postal_code,omitempty" maxLength:"10"`
	Photo             *string                    `json:"photo,omitempty" maxLength:"255"`
	KTPFile           *string                    `json:"ktp_file,omitempty" maxLength:"255"`
	CVFile            *string                    `json:"cv_file,omitempty" maxLength:"255"`
	NPWPFile          *string                    `json:"npwp_file,omitempty" maxLength:"255"`
}

type RegisterOutput struct {
	Body RegisterResponse `json:"body"`
}

type RegisterResponse struct {
	Token string          `json:"token"`
	User  RegisterUserDTO `json:"user"`
}

type RegisterUserDTO struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}
