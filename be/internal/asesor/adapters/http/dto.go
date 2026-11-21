package http

import (
	"time"

	"lsp/internal/asesor/domain"
	httpplatform "lsp/internal/platform/http"
)

type UpsertProfileInput struct {
	UserID string            `path:"user_id"`
	Body   UpsertProfileBody `json:"body"`
}

type UpdateProfileInput struct {
	UserID string           `path:"user_id"`
	Body   PatchProfileBody `json:"body"`
}

type UpsertProfileBody struct {
	TUKID                *string `json:"tuk_id,omitempty" doc:"ID TUK tempat asesor berafiliasi"`
	Province             *string `json:"province,omitempty" maxLength:"120" doc:"Domisili asesor (provinsi)"`
	City                 *string `json:"city,omitempty" maxLength:"120" doc:"Domisili asesor (kota)"`
	DailyAssignmentLimit int     `json:"daily_assignment_limit" minimum:"0" doc:"Batas penugasan harian (0 = tanpa batas)"`
}

type PatchProfileBody struct {
	TUKID                *string `json:"tuk_id,omitempty" doc:"ID TUK tempat asesor berafiliasi"`
	Province             *string `json:"province,omitempty" maxLength:"120" doc:"Domisili asesor (provinsi)"`
	City                 *string `json:"city,omitempty" maxLength:"120" doc:"Domisili asesor (kota)"`
	DailyAssignmentLimit *int    `json:"daily_assignment_limit,omitempty" minimum:"0" doc:"Batas penugasan harian (0 = tanpa batas)"`
}

type GetProfileInput struct {
	UserID string `path:"user_id"`
}

type ListProfilesInput struct {
	httpplatform.PaginationQuery
}

type ProfileOutput struct {
	Body ProfileDTO `json:"body"`
}

type ProfileDTO struct {
	UserID               string     `json:"user_id"`
	Name                 string     `json:"name,omitempty"`
	Email                string     `json:"email,omitempty"`
	TUKID                *string    `json:"tuk_id,omitempty"`
	TUKName              *string    `json:"tuk_name,omitempty"`
	Province             *string    `json:"province,omitempty"`
	City                 *string    `json:"city,omitempty"`
	DailyAssignmentLimit int        `json:"daily_assignment_limit"`
	CreatedAt            *time.Time `json:"created_at,omitempty"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty"`
}

func toProfileDTO(profile domain.Profile) ProfileDTO {
	createdAt := profile.CreatedAt
	updatedAt := profile.UpdatedAt
	return ProfileDTO{
		UserID:               profile.UserID,
		TUKID:                profile.TUKID,
		Province:             profile.Province,
		City:                 profile.City,
		DailyAssignmentLimit: profile.DailyAssignmentLimit,
		CreatedAt:            &createdAt,
		UpdatedAt:            &updatedAt,
	}
}

func toListItemDTO(item domain.ListItem) ProfileDTO {
	return ProfileDTO{
		UserID:               item.UserID,
		Name:                 item.Name,
		Email:                item.Email,
		TUKID:                item.TUKID,
		TUKName:              item.TUKName,
		Province:             item.Province,
		City:                 item.City,
		DailyAssignmentLimit: item.DailyAssignmentLimit,
	}
}

type ListCertificationsInput struct {
	UserID string `path:"user_id"`
}

type CreateCertificationInput struct {
	UserID string            `path:"user_id"`
	Body   CertificationBody `json:"body"`
}

type UpdateCertificationInput struct {
	CertID string                 `path:"cert_id"`
	Body   CertificationPatchBody `json:"body"`
}

type DeleteCertificationInput struct {
	UserID string `path:"user_id"`
	CertID string `path:"cert_id"`
}

type CertificationBody struct {
	Name      string  `json:"name" minLength:"1" maxLength:"255" doc:"Nama sertifikasi/skema"`
	Issuer    string  `json:"issuer" minLength:"1" maxLength:"255" doc:"Lembaga penerbit sertifikat"`
	Number    *string `json:"number,omitempty" maxLength:"120" doc:"Nomor sertifikat"`
	IssuedAt  *string `json:"issued_at,omitempty" pattern:"^\\d{4}-\\d{2}-\\d{2}$" doc:"Tanggal terbit (YYYY-MM-DD)"`
	ExpiredAt *string `json:"expired_at,omitempty" pattern:"^\\d{4}-\\d{2}-\\d{2}$" doc:"Tanggal kadaluarsa (YYYY-MM-DD)"`
}

type CertificationPatchBody struct {
	Name      *string `json:"name,omitempty" minLength:"1" maxLength:"255" doc:"Nama sertifikasi/skema"`
	Issuer    *string `json:"issuer,omitempty" minLength:"1" maxLength:"255" doc:"Lembaga penerbit sertifikat"`
	Number    *string `json:"number,omitempty" maxLength:"120" doc:"Nomor sertifikat"`
	IssuedAt  *string `json:"issued_at,omitempty" pattern:"^\\d{4}-\\d{2}-\\d{2}$" doc:"Tanggal terbit (YYYY-MM-DD)"`
	ExpiredAt *string `json:"expired_at,omitempty" pattern:"^\\d{4}-\\d{2}-\\d{2}$" doc:"Tanggal kadaluarsa (YYYY-MM-DD)"`
}

type CertificationOutput struct {
	Body CertificationDTO `json:"body"`
}

type CertificationListOutput struct {
	Body struct {
		Data []CertificationDTO `json:"data"`
	} `json:"body"`
}

type DeleteCertificationOutput struct{}

type CertificationDTO struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Issuer    string    `json:"issuer"`
	Number    *string   `json:"number,omitempty"`
	IssuedAt  *string   `json:"issued_at,omitempty"`
	ExpiredAt *string   `json:"expired_at,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toCertificationDTO(cert domain.Certification) CertificationDTO {
	return CertificationDTO{
		ID:        cert.ID,
		UserID:    cert.UserID,
		Name:      cert.Name,
		Issuer:    cert.Issuer,
		Number:    cert.Number,
		IssuedAt:  formatDate(cert.IssuedAt),
		ExpiredAt: formatDate(cert.ExpiredAt),
		CreatedAt: cert.CreatedAt,
		UpdatedAt: cert.UpdatedAt,
	}
}

// parseDate mengubah tanggal "YYYY-MM-DD" menjadi time.Time (nil bila kosong).
func parseDate(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	t, err := time.Parse(time.DateOnly, *value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func formatDate(value *time.Time) *string {
	if value == nil {
		return nil
	}
	s := value.Format(time.DateOnly)
	return &s
}
