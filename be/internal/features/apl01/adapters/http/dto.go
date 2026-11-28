package http

import (
	"context"
	"strings"
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
	ID   string     `path:"id"`
	Body SubmitBody `json:"body"`
}

type VerifyInput struct {
	UserAgent    string     `header:"User-Agent"`
	ForwardedFor string     `header:"X-Forwarded-For"`
	Body         VerifyBody `json:"body"`
}

type VerifyStatusInput struct {
	Token string `path:"token" minLength:"1" maxLength:"200" doc:"Token hasil scan QR"`
}

type DecideInput struct {
	ID   string     `path:"id"`
	Body DecideBody `json:"body"`
}

type AssignAsesorInput struct {
	ID   string           `path:"id"`
	Body AssignAsesorBody `json:"body"`
}

type AsesorCandidatesInput struct {
	ID string `path:"id"`
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

type SubmitBody struct {
	VerifyToken *string `json:"verify_token,omitempty" maxLength:"200" doc:"Token QR yang sudah dimiliki; bila masih aktif, token yang sama dipakai ulang (tidak diterbitkan baru)"`
}

type DecideBody struct {
	Decision string  `json:"decision" enum:"APPROVED,REJECTED" doc:"Hasil verifikasi admin"`
	Note     *string `json:"note,omitempty" doc:"Catatan verifikasi"`
	AsesorID *string `json:"asesor_id,omitempty" maxLength:"32" doc:"Asesor yang di-plotting (bila disetujui)"`
}

type AssignAsesorBody struct {
	AsesorID string `json:"asesor_id" minLength:"1" maxLength:"32" doc:"ID asesor yang ditugaskan"`
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

type VerifyStatusOutput struct {
	Body VerifyStatusDTO `json:"body"`
}

type AsesorCandidatesOutput struct {
	Body AsesorCandidatesDTO `json:"body"`
}

type VerifyStatusDTO struct {
	Status          string     `json:"status" doc:"menunggu | dipakai | kadaluarsa | dibatalkan"`
	ApplicationCode string     `json:"application_code,omitempty"`
	ExpiresAt       time.Time  `json:"expires_at"`
	VerifiedAt      *time.Time `json:"verified_at,omitempty"`
}

// --- DTO ---

type ApplicationDetailDTO struct {
	Application ApplicationDTO           `json:"application"`
	Data        *ApplicationDataDTO      `json:"data,omitempty"`
	Documents   []ApplicationDocumentDTO `json:"documents"`
	Schedule    *AssessmentScheduleDTO   `json:"schedule,omitempty"`
	Signatures  []FormSignatureDTO       `json:"signatures"`
	// Assignment: penugasan asesor aktif (bila sudah di-plotting).
	Assignment *AsesorAssignmentDTO `json:"assignment,omitempty"`
	// Locked: formulir tidak boleh diubah (sudah diajukan / QR sudah terbit).
	Locked bool `json:"locked"`
}

type AsesorAssignmentDTO struct {
	ID            string     `json:"id"`
	ApplicationID string     `json:"application_id"`
	AsesorID      string     `json:"asesor_id"`
	AsesorName    string     `json:"asesor_name,omitempty"`
	Status        string     `json:"status"`
	AssignedBy    *string    `json:"assigned_by,omitempty"`
	AssignedAt    time.Time  `json:"assigned_at"`
	RespondedAt   *time.Time `json:"responded_at,omitempty"`
	Note          *string    `json:"note,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type AsesorCandidateDTO struct {
	UserID               string `json:"user_id"`
	Name                 string `json:"name"`
	Province             string `json:"province,omitempty"`
	City                 string `json:"city,omitempty"`
	DailyAssignmentLimit int    `json:"daily_assignment_limit" doc:"0 = tanpa batas"`
	AssignedCount        int    `json:"assigned_count" doc:"Jumlah penugasan aktif (PENDING/ACCEPTED) pada tanggal jadwal"`
	Disabled             bool   `json:"disabled" doc:"True bila kapasitas harian sudah penuh"`
}

type AsesorCandidatesDTO struct {
	Items []AsesorCandidateDTO `json:"items"`
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
	FileURL        *string    `json:"file_url,omitempty" doc:"Object key berkas (dipakai saat menyimpan ulang)"`
	FileViewURL    *string    `json:"file_view_url,omitempty" doc:"URL presigned untuk pratinjau"`
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

func (h *Handler) toApplicationDetailDTO(ctx context.Context, d domain.ApplicationDetail) ApplicationDetailDTO {
	detail := ApplicationDetailDTO{
		Application: toApplicationDTO(d.Application),
		Documents:   h.toApplicationDocumentDTOs(ctx, d.Documents),
		Signatures:  toFormSignatureDTOs(d.Signatures),
		Locked:      d.Locked,
	}
	if d.Data != nil {
		data := toApplicationDataDTO(*d.Data)
		detail.Data = &data
	}
	if d.Schedule != nil {
		schedule := toAssessmentScheduleDTO(*d.Schedule)
		detail.Schedule = &schedule
	}
	if d.Assignment != nil {
		assignment := toAsesorAssignmentDTO(*d.Assignment)
		detail.Assignment = &assignment
	}
	return detail
}

func toAsesorCandidateDTOs(items []domain.AsesorCandidate) []AsesorCandidateDTO {
	result := make([]AsesorCandidateDTO, 0, len(items))
	for _, item := range items {
		result = append(result, AsesorCandidateDTO{
			UserID:               item.UserID,
			Name:                 item.Name,
			Province:             item.Province,
			City:                 item.City,
			DailyAssignmentLimit: item.DailyAssignmentLimit,
			AssignedCount:        item.AssignedCount,
			Disabled:             item.Disabled,
		})
	}
	return result
}

func toAsesorAssignmentDTO(a domain.AsesorAssignment) AsesorAssignmentDTO {
	return AsesorAssignmentDTO{
		ID:            a.ID,
		ApplicationID: a.ApplicationID,
		AsesorID:      a.AsesorID,
		AsesorName:    a.AsesorName,
		Status:        a.Status,
		AssignedBy:    a.AssignedBy,
		AssignedAt:    a.AssignedAt,
		RespondedAt:   a.RespondedAt,
		Note:          a.Note,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
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

// fileURLExpiry adalah masa berlaku presigned GET URL untuk file privat.
const fileURLExpiry = 15 * time.Minute

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

func (h *Handler) toApplicationDocumentDTO(ctx context.Context, d domain.ApplicationDocument) ApplicationDocumentDTO {
	return ApplicationDocumentDTO{
		ID:             d.ID,
		ApplicationID:  d.ApplicationID,
		DocumentTypeID: d.DocumentTypeID,
		FileURL:        d.FileURL,
		FileViewURL:    h.fileURL(ctx, d.FileURL),
		Status:         string(d.Status),
		Note:           d.Note,
		VerifiedAt:     d.VerifiedAt,
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
	}
}

func (h *Handler) toApplicationDocumentDTOs(ctx context.Context, items []domain.ApplicationDocument) []ApplicationDocumentDTO {
	result := make([]ApplicationDocumentDTO, 0, len(items))
	for _, item := range items {
		result = append(result, h.toApplicationDocumentDTO(ctx, item))
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
