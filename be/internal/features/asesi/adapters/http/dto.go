package http

import (
	"context"
	"strings"
	"time"

	"lsp/internal/features/asesi/domain"
	httpplatform "lsp/internal/platform/http"
)

type CreateInput struct {
	Body AsesiProfileBody `json:"body"`
}

type UpdateInput struct {
	ID   string                `path:"id"`
	Body AsesiProfilePatchBody `json:"body"`
}

type GetInput struct {
	ID string `path:"id"`
}

type GetByUserInput struct {
	UserID string `path:"user_id"`
}

type DeleteInput struct {
	ID string `path:"id"`
}

type ListInput struct {
	httpplatform.PaginationQuery
}

type AsesiProfileBody struct {
	UserID            string                `json:"user_id" minLength:"1" maxLength:"32"`
	FullName          string                `json:"full_name" minLength:"1" maxLength:"255"`
	NIK               string                `json:"nik" minLength:"16" maxLength:"16"`
	Phone             string                `json:"phone" minLength:"5" maxLength:"20"`
	Gender            domain.Gender         `json:"gender" enum:"LAKI-LAKI,PEREMPUAN"`
	BirthPlace        string                `json:"birth_place" minLength:"1" maxLength:"255"`
	BirthDate         string                `json:"birth_date" pattern:"^\\d{4}-\\d{2}-\\d{2}$" doc:"Tanggal lahir (YYYY-MM-DD)"`
	Nationality       string                `json:"nationality" minLength:"1" maxLength:"100" default:"Indonesia"`
	Occupation        string                `json:"occupation" minLength:"1" maxLength:"255"`
	EducationLevel    domain.EducationLevel `json:"education_level" enum:"SD,SMP,SMA_SMK,D1,D2,D3,D4,S1,S2,S3"`
	KTPAddress        string                `json:"ktp_address" minLength:"1"`
	DomicileSameAsKTP bool                  `json:"domicile_same_as_ktp"`
	DomicileAddress   *string               `json:"domicile_address,omitempty"`
	Province          string                `json:"province" minLength:"1" maxLength:"255"`
	City              string                `json:"city" minLength:"1" maxLength:"255"`
	District          string                `json:"district" minLength:"1" maxLength:"255"`
	Village           string                `json:"village" minLength:"1" maxLength:"255"`
	PostalCode        *string               `json:"postal_code,omitempty" maxLength:"10"`
	Photo             *string               `json:"photo,omitempty" maxLength:"255"`
	KTPFile           *string               `json:"ktp_file,omitempty" maxLength:"255"`
	CVFile            *string               `json:"cv_file,omitempty" maxLength:"255"`
	NPWPFile          *string               `json:"npwp_file,omitempty" maxLength:"255"`
}

type AsesiProfilePatchBody struct {
	FullName          *string                `json:"full_name,omitempty" minLength:"1" maxLength:"255"`
	NIK               *string                `json:"nik,omitempty" minLength:"16" maxLength:"16"`
	Phone             *string                `json:"phone,omitempty" minLength:"5" maxLength:"20"`
	Gender            *domain.Gender         `json:"gender,omitempty" enum:"LAKI-LAKI,PEREMPUAN"`
	BirthPlace        *string                `json:"birth_place,omitempty" minLength:"1" maxLength:"255"`
	BirthDate         *string                `json:"birth_date,omitempty" pattern:"^\\d{4}-\\d{2}-\\d{2}$" doc:"Tanggal lahir (YYYY-MM-DD)"`
	Nationality       *string                `json:"nationality,omitempty" minLength:"1" maxLength:"100"`
	Occupation        *string                `json:"occupation,omitempty" minLength:"1" maxLength:"255"`
	EducationLevel    *domain.EducationLevel `json:"education_level,omitempty" enum:"SD,SMP,SMA_SMK,D1,D2,D3,D4,S1,S2,S3"`
	KTPAddress        *string                `json:"ktp_address,omitempty" minLength:"1"`
	DomicileSameAsKTP *bool                  `json:"domicile_same_as_ktp,omitempty"`
	DomicileAddress   *string                `json:"domicile_address,omitempty"`
	Province          *string                `json:"province,omitempty" minLength:"1" maxLength:"255"`
	City              *string                `json:"city,omitempty" minLength:"1" maxLength:"255"`
	District          *string                `json:"district,omitempty" minLength:"1" maxLength:"255"`
	Village           *string                `json:"village,omitempty" minLength:"1" maxLength:"255"`
	PostalCode        *string                `json:"postal_code,omitempty" maxLength:"10"`
	Photo             *string                `json:"photo,omitempty" maxLength:"255"`
	KTPFile           *string                `json:"ktp_file,omitempty" maxLength:"255"`
	CVFile            *string                `json:"cv_file,omitempty" maxLength:"255"`
	NPWPFile          *string                `json:"npwp_file,omitempty" maxLength:"255"`
}

type AsesiProfileOutput struct {
	Body AsesiProfileDTO `json:"body"`
}

type AsesiProfileDTO struct {
	ID                string                `json:"id"`
	UserID            string                `json:"user_id"`
	FullName          string                `json:"full_name"`
	NIK               string                `json:"nik"`
	Phone             string                `json:"phone"`
	Gender            domain.Gender         `json:"gender"`
	BirthPlace        string                `json:"birth_place"`
	BirthDate         string                `json:"birth_date"`
	Nationality       string                `json:"nationality"`
	Occupation        string                `json:"occupation"`
	EducationLevel    domain.EducationLevel `json:"education_level"`
	KTPAddress        string                `json:"ktp_address"`
	DomicileSameAsKTP bool                  `json:"domicile_same_as_ktp"`
	DomicileAddress   *string               `json:"domicile_address,omitempty"`
	Province          string                `json:"province"`
	City              string                `json:"city"`
	District          string                `json:"district"`
	Village           string                `json:"village"`
	PostalCode        *string               `json:"postal_code,omitempty"`
	Photo             *string               `json:"photo,omitempty"`
	KTPFile           *string               `json:"ktp_file,omitempty"`
	CVFile            *string               `json:"cv_file,omitempty"`
	NPWPFile          *string               `json:"npwp_file,omitempty"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
}

type DeleteOutput struct{}

// fileURLExpiry adalah masa berlaku presigned GET URL untuk file privat.
const fileURLExpiry = 15 * time.Minute

// toDTO mengubah profil menjadi DTO. Object key dokumen ditandatangani menjadi
// presigned GET URL agar bisa diakses langsung (bucket privat).
func (h *Handler) toDTO(ctx context.Context, profile domain.Profile) AsesiProfileDTO {
	return AsesiProfileDTO{
		ID:                profile.ID,
		UserID:            profile.UserID,
		FullName:          profile.FullName,
		NIK:               profile.NIK,
		Phone:             profile.Phone,
		Gender:            profile.Gender,
		BirthPlace:        profile.BirthPlace,
		BirthDate:         profile.BirthDate.Format(time.DateOnly),
		Nationality:       profile.Nationality,
		Occupation:        profile.Occupation,
		EducationLevel:    profile.EducationLevel,
		KTPAddress:        profile.KTPAddress,
		DomicileSameAsKTP: profile.DomicileSameAsKTP,
		DomicileAddress:   profile.DomicileAddress,
		Province:          profile.Province,
		City:              profile.City,
		District:          profile.District,
		Village:           profile.Village,
		PostalCode:        profile.PostalCode,
		Photo:             h.fileURL(ctx, profile.Photo),
		KTPFile:           h.fileURL(ctx, profile.KTPFile),
		CVFile:            h.fileURL(ctx, profile.CVFile),
		NPWPFile:          h.fileURL(ctx, profile.NPWPFile),
		CreatedAt:         profile.CreatedAt,
		UpdatedAt:         profile.UpdatedAt,
	}
}

// fileURL menandatangani object key menjadi presigned GET URL. Bila storage
// tidak dikonfigurasi, nilai sudah berupa URL penuh, atau penandatanganan
// gagal, nilai asli dikembalikan apa adanya.
func (h *Handler) fileURL(ctx context.Context, key *string) *string {
	if key == nil || *key == "" {
		return key
	}
	if h.files == nil || strings.HasPrefix(*key, "http://") || strings.HasPrefix(*key, "https://") {
		return key
	}
	url, err := h.files.PresignedGetURL(ctx, *key, fileURLExpiry)
	if err != nil {
		return key
	}
	return &url
}

func parseDate(value string) (time.Time, error) {
	return time.Parse(time.DateOnly, value)
}

func parseOptionalDate(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	t, err := parseDate(*value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
