package http

import (
	"time"

	"lsp/internal/features/apl01/domain"
	httpplatform "lsp/internal/platform/http"
)

// --- Input ---

type CreateDraftInput struct{}

type ListMineInput struct {
	httpplatform.PaginationQuery
	Status string `query:"status" enum:"DRAFT,SUBMITTED,APPROVED,REJECTED" doc:"Filter status"`
}

type ListAllInput struct {
	httpplatform.PaginationQuery
	UserID string `query:"user_id" maxLength:"32" doc:"Filter berdasarkan id asesi"`
	Status string `query:"status" enum:"DRAFT,SUBMITTED,APPROVED,REJECTED" doc:"Filter status"`
}

type GetInput struct {
	ID string `path:"id"`
}

type UpdateStepInput struct {
	ID   string         `path:"id"`
	Body UpdateStepBody `json:"body"`
}

type PutDataInput struct {
	ID   string   `path:"id"`
	Body DataBody `json:"body"`
}

type PutDocumentsInput struct {
	ID   string        `path:"id"`
	Body DocumentsBody `json:"body"`
}

type PutScheduleInput struct {
	ID   string       `path:"id"`
	Body ScheduleBody `json:"body"`
}

type SubmitInput struct {
	ID string `path:"id"`
}

type VerifyInput struct {
	UserAgent    string     `header:"User-Agent"`
	ForwardedFor string     `header:"X-Forwarded-For"`
	Body         VerifyBody `json:"body"`
}

type DecideInput struct {
	ID   string     `path:"id"`
	Body DecideBody `json:"body"`
}

// --- Body ---

type UpdateStepBody struct {
	CurrentStep int `json:"current_step" minimum:"1" maximum:"4" doc:"Step terakhir yang dikerjakan (1..4)"`
}

type DataBody struct {
	CertificationSchemaID string `json:"certification_schema_id" minLength:"1" maxLength:"32" doc:"ID skema sertifikasi"`
	Purpose               string `json:"purpose" enum:"SERTIFIKASI,SERTIFIKASI_ULANG,REKOGNISI_PEMBELAJARAN_LAMPAU,PENGAKUAN_KOMPETENSI_TERKINI" doc:"Tujuan asesmen"`
}

type DocumentsBody struct {
	Items []DocumentItemBody `json:"items" maxItems:"200" doc:"Daftar lengkap bukti kelengkapan yang diunggah"`
}

type DocumentItemBody struct {
	DocumentTypeID string  `json:"document_type_id" minLength:"1" maxLength:"32" doc:"ID master jenis dokumen"`
	FileURL        *string `json:"file_url,omitempty" maxLength:"255" doc:"URL hasil upload; kosong bila belum diunggah"`
}

type ScheduleBody struct {
	Mode          string  `json:"mode" enum:"ONLINE,OFFLINE" doc:"Metode pelaksanaan"`
	ScheduledDate string  `json:"scheduled_date" pattern:"^\\d{4}-\\d{2}-\\d{2}$" doc:"Tanggal jadwal (YYYY-MM-DD)"`
	TukID         *string `json:"tuk_id,omitempty" maxLength:"32" doc:"ID TUK (opsional, untuk offline)"`
}

type VerifyBody struct {
	Token string `json:"token" minLength:"1" doc:"Token hasil scan QR"`
}

type DecideBody struct {
	Decision string  `json:"decision" enum:"APPROVED,REJECTED" doc:"Hasil verifikasi admin"`
	Note     *string `json:"note,omitempty" doc:"Catatan verifikasi"`
	AsesorID *string `json:"asesor_id,omitempty" maxLength:"32" doc:"Asesor yang di-plotting (bila disetujui)"`
}

// --- Output ---

type ApplicationOutput struct {
	Body ApplicationDetailDTO `json:"body"`
}

type SubmitOutput struct {
	Body SubmitDTO `json:"body"`
}

type SubmitDTO struct {
	Application ApplicationDetailDTO `json:"application"`
	VerifyToken string               `json:"verify_token" doc:"Token mentah untuk ditampilkan sebagai QR (hanya muncul sekali)"`
}

// --- DTO ---

type ApplicationDetailDTO struct {
	Application ApplicationDTO           `json:"application"`
	Data        *ApplicationDataDTO      `json:"data,omitempty"`
	Documents   []ApplicationDocumentDTO `json:"documents"`
	Schedule    *AssessmentScheduleDTO   `json:"schedule,omitempty"`
	Signatures  []FormSignatureDTO       `json:"signatures"`
}

type ApplicationDTO struct {
	ID               string     `json:"id"`
	UserID           string     `json:"user_id"`
	Code             string     `json:"code"`
	Status           string     `json:"status"`
	CurrentStep      int        `json:"current_step"`
	SubmittedAt      *time.Time `json:"submitted_at,omitempty"`
	VerifiedBy       *string    `json:"verified_by,omitempty"`
	VerifiedAt       *time.Time `json:"verified_at,omitempty"`
	VerificationNote *string    `json:"verification_note,omitempty"`
	AsesorID         *string    `json:"asesor_id,omitempty"`
	AsesorAssignedAt *time.Time `json:"asesor_assigned_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type ApplicationDataDTO struct {
	ID                    string    `json:"id"`
	ApplicationID         string    `json:"application_id"`
	CertificationSchemaID string    `json:"certification_schema_id"`
	Purpose               string    `json:"purpose"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type ApplicationDocumentDTO struct {
	ID             string     `json:"id"`
	ApplicationID  string     `json:"application_id"`
	DocumentTypeID string     `json:"document_type_id"`
	FileURL        *string    `json:"file_url,omitempty"`
	Status         string     `json:"status"`
	Note           *string    `json:"note,omitempty"`
	VerifiedAt     *time.Time `json:"verified_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type AssessmentScheduleDTO struct {
	ID            string    `json:"id"`
	ApplicationID string    `json:"application_id"`
	Mode          string    `json:"mode"`
	ScheduledDate string    `json:"scheduled_date"`
	TukID         *string   `json:"tuk_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type FormSignatureDTO struct {
	ID            string     `json:"id"`
	ApplicationID string     `json:"application_id"`
	Entity        string     `json:"entity"`
	Role          string     `json:"role"`
	Status        string     `json:"status"`
	ExpiresAt     time.Time  `json:"expires_at"`
	UsedAt        *time.Time `json:"used_at,omitempty"`
	RevokedAt     *time.Time `json:"revoked_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// --- Konversi ---

func toApplicationDTO(a domain.Application) ApplicationDTO {
	return ApplicationDTO{
		ID:               a.ID,
		UserID:           a.UserID,
		Code:             a.Code,
		Status:           string(a.Status),
		CurrentStep:      a.CurrentStep,
		SubmittedAt:      a.SubmittedAt,
		VerifiedBy:       a.VerifiedBy,
		VerifiedAt:       a.VerifiedAt,
		VerificationNote: a.VerificationNote,
		AsesorID:         a.AsesorID,
		AsesorAssignedAt: a.AsesorAssignedAt,
		CreatedAt:        a.CreatedAt,
		UpdatedAt:        a.UpdatedAt,
	}
}

func toApplicationDetailDTO(d domain.ApplicationDetail) ApplicationDetailDTO {
	detail := ApplicationDetailDTO{
		Application: toApplicationDTO(d.Application),
		Documents:   toApplicationDocumentDTOs(d.Documents),
		Signatures:  toFormSignatureDTOs(d.Signatures),
	}
	if d.Data != nil {
		data := toApplicationDataDTO(*d.Data)
		detail.Data = &data
	}
	if d.Schedule != nil {
		schedule := toAssessmentScheduleDTO(*d.Schedule)
		detail.Schedule = &schedule
	}
	return detail
}

func toApplicationDataDTO(d domain.ApplicationData) ApplicationDataDTO {
	return ApplicationDataDTO{
		ID:                    d.ID,
		ApplicationID:         d.ApplicationID,
		CertificationSchemaID: d.CertificationSchemaID,
		Purpose:               string(d.Purpose),
		CreatedAt:             d.CreatedAt,
		UpdatedAt:             d.UpdatedAt,
	}
}

func toApplicationDocumentDTO(d domain.ApplicationDocument) ApplicationDocumentDTO {
	return ApplicationDocumentDTO{
		ID:             d.ID,
		ApplicationID:  d.ApplicationID,
		DocumentTypeID: d.DocumentTypeID,
		FileURL:        d.FileURL,
		Status:         string(d.Status),
		Note:           d.Note,
		VerifiedAt:     d.VerifiedAt,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}
}

func toApplicationDocumentDTOs(items []domain.ApplicationDocument) []ApplicationDocumentDTO {
	result := make([]ApplicationDocumentDTO, 0, len(items))
	for _, item := range items {
		result = append(result, toApplicationDocumentDTO(item))
	}
	return result
}

func toAssessmentScheduleDTO(s domain.AssessmentSchedule) AssessmentScheduleDTO {
	return AssessmentScheduleDTO{
		ID:            s.ID,
		ApplicationID: s.ApplicationID,
		Mode:          string(s.Mode),
		ScheduledDate: s.ScheduledDate.Format(time.DateOnly),
		TukID:         s.TukID,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}

func toFormSignatureDTO(s domain.FormSignature) FormSignatureDTO {
	return FormSignatureDTO{
		ID:            s.ID,
		ApplicationID: s.ApplicationID,
		Entity:        string(s.Entity),
		Role:          string(s.Role),
		Status:        string(s.Status),
		ExpiresAt:     s.ExpiresAt,
		UsedAt:        s.UsedAt,
		RevokedAt:     s.RevokedAt,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}

func toFormSignatureDTOs(items []domain.FormSignature) []FormSignatureDTO {
	result := make([]FormSignatureDTO, 0, len(items))
	for _, item := range items {
		result = append(result, toFormSignatureDTO(item))
	}
	return result
}
