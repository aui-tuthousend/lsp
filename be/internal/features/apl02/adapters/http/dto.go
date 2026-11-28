package http

import (
	"context"
	"strings"
	"time"

	"lsp/internal/features/apl02/domain"
	httpplatform "lsp/internal/platform/http"
)

// --- Input ---

type CreateInput struct {
	Body CreateBody `json:"body"`
}

type ListMineInput struct {
	httpplatform.PaginationQuery
}

type ListAllInput struct {
	httpplatform.PaginationQuery
	Status string `query:"status" enum:"DRAFT,SUBMITTED,REVIEWED" doc:"Filter status asesmen mandiri"`
	// ApplicationStatus memfilter berdasarkan status permohonan APL-01
	// (keputusan admin): SUBMITTED = menunggu, APPROVED/REJECTED = sudah diputuskan.
	ApplicationStatus string `query:"application_status" enum:"DRAFT,SUBMITTED,APPROVED,REJECTED" doc:"Filter status permohonan APL-01"`
	// Assigned memfilter ada/tidaknya penugasan asesor aktif (PENDING/ACCEPTED):
	// "true" = sudah ditugaskan, "false" = belum, kosong = tanpa filter.
	Assigned string `query:"assigned" enum:"true,false" doc:"Filter penugasan asesor: true=sudah, false=belum"`
}

type GetInput struct {
	ID string `path:"id"`
}

type PutUnitsInput struct {
	ID   string       `path:"id"`
	Body PutUnitsBody `json:"body"`
}

type SubmitInput struct {
	ID   string               `path:"id"`
	Body AssessmentSubmitBody `json:"body"`
}

type VerifyInput struct {
	UserAgent    string               `header:"User-Agent"`
	ForwardedFor string               `header:"X-Forwarded-For"`
	Body         AssessmentVerifyBody `json:"body"`
}

type VerifyStatusInput struct {
	Token string `path:"token" minLength:"1" maxLength:"200" doc:"Token hasil scan QR"`
}

type ReviewInput struct {
	ID   string     `path:"id"`
	Body ReviewBody `json:"body"`
}

// --- Body ---

type CreateBody struct {
	ApplicationID string `json:"application_id" minLength:"1" maxLength:"32" doc:"ID permohonan APL-01"`
}

type UnitAnswerBody struct {
	CompetencyUnitID string   `json:"competency_unit_id" minLength:"1" maxLength:"32"`
	IsCompetent      *bool    `json:"is_competent,omitempty" doc:"Pernyataan mandiri asesi; kosong = belum dijawab"`
	Files            []string `json:"files,omitempty" maxItems:"20" doc:"URL/object key berkas bukti (wajib bila kompeten)"`
}

type PutUnitsBody struct {
	Items []UnitAnswerBody `json:"items" maxItems:"500" doc:"Daftar lengkap jawaban unit (set-semantik)"`
}

type AssessmentSubmitBody struct {
	VerifyToken *string `json:"verify_token,omitempty" maxLength:"200" doc:"Token QR yang sudah dimiliki; bila masih aktif dipakai ulang"`
}

type AssessmentVerifyBody struct {
	Token string `json:"token" minLength:"1" doc:"Token hasil scan QR"`
}

type ReviewBody struct {
	Decision string  `json:"decision" enum:"RECOMMENDED,NOT_RECOMMENDED" doc:"Rekomendasi asesor"`
	Note     *string `json:"note,omitempty" doc:"Catatan tinjauan asesor"`
	ZoomLink *string `json:"zoom_link,omitempty" maxLength:"255" doc:"Link pertemuan zoom"`
	ZoomDate *string `json:"zoom_date,omitempty" pattern:"^\\d{4}-\\d{2}-\\d{2}$" doc:"Tanggal pertemuan (YYYY-MM-DD)"`
	ZoomTime *string `json:"zoom_time,omitempty" pattern:"^\\d{2}:\\d{2}$" doc:"Jam pertemuan (HH:MM)"`
}

// --- Output ---

type AssessmentOutput struct {
	Body AssessmentDetailDTO `json:"body"`
}

type AssessmentSubmitOutput struct {
	Body AssessmentSubmitDTO `json:"body"`
}

type AssessmentSubmitDTO struct {
	Assessment  AssessmentDetailDTO `json:"assessment"`
	VerifyToken string              `json:"verify_token" doc:"Token mentah untuk QR (hanya muncul sekali)"`
}

type AssessmentVerifyStatusOutput struct {
	Body AssessmentVerifyStatusDTO `json:"body"`
}

type AssessmentVerifyStatusDTO struct {
	Status     string     `json:"status" doc:"menunggu | dipakai | kadaluarsa | dibatalkan"`
	Code       string     `json:"code,omitempty"`
	ExpiresAt  time.Time  `json:"expires_at"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
}

// --- DTO ---

type AssessmentDetailDTO struct {
	Assessment  AssessmentDTO       `json:"assessment"`
	Application ApplicationInfoDTO  `json:"application"`
	Units       []AssessmentUnitDTO `json:"units"`
	// Locked: formulir tidak boleh diubah (sudah diajukan / QR sudah terbit).
	Locked bool `json:"locked"`
}

// ApplicationInfoDTO — info permohonan/skema untuk laman asesmen mandiri (read-only).
type ApplicationInfoDTO struct {
	CertificationSchemaID string    `json:"certification_schema_id"`
	SchemaName            string    `json:"schema_name"`
	SchemaCode            string    `json:"schema_code"`
	Purpose               string    `json:"purpose"`
	DataCreatedAt         time.Time `json:"data_created_at"`
}

type AssessmentDTO struct {
	ID             string     `json:"id"`
	ApplicationID  string     `json:"application_id"`
	Code           string     `json:"code"`
	Status         string     `json:"status"`
	AsesorNote     *string    `json:"asesor_note,omitempty"`
	AsesorDecision *string    `json:"asesor_decision,omitempty"`
	SubmittedAt    *time.Time `json:"submitted_at,omitempty"`
	ReviewedBy     *string    `json:"reviewed_by,omitempty"`
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty"`
	ZoomLink       *string    `json:"zoom_link,omitempty"`
	ZoomDate       *string    `json:"zoom_date,omitempty"`
	ZoomTime       *string    `json:"zoom_time,omitempty"`
	SchemaName     string     `json:"schema_name,omitempty"`
	SchemaCode     string     `json:"schema_code,omitempty"`
	// AsesiName/ApplicationCode terisi pada daftar staff (antrian verifikasi).
	AsesiName       string `json:"asesi_name,omitempty"`
	ApplicationCode string `json:"application_code,omitempty"`
	// ApplicationStatus adalah status permohonan APL-01 (DRAFT/SUBMITTED/
	// APPROVED/REJECTED); terisi pada daftar staff.
	ApplicationStatus string `json:"application_status,omitempty"`
	// ApplicationNote adalah catatan verifikasi admin (mis. alasan penolakan).
	ApplicationNote *string `json:"application_note,omitempty"`
	// AsesorID/AsesorName/AsesorAssignmentStatus terisi pada daftar staff bila
	// permohonan sudah di-plotting ke seorang asesor.
	AsesorID               *string   `json:"asesor_id,omitempty"`
	AsesorName             string    `json:"asesor_name,omitempty"`
	AsesorAssignmentStatus *string   `json:"asesor_assignment_status,omitempty"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type AssessmentUnitDTO struct {
	ID               string                  `json:"id"`
	CompetencyUnitID string                  `json:"competency_unit_id"`
	UnitCode         string                  `json:"unit_code"`
	UnitTitle        string                  `json:"unit_title"`
	IsCompetent      *bool                   `json:"is_competent"`
	Files            []AssessmentUnitFileDTO `json:"files"`
	CreatedAt        time.Time               `json:"created_at"`
	UpdatedAt        time.Time               `json:"updated_at"`
}

type AssessmentUnitFileDTO struct {
	ID          string    `json:"id"`
	FileURL     string    `json:"file_url" doc:"Object key berkas (dipakai saat menyimpan ulang)"`
	FileViewURL string    `json:"file_view_url,omitempty" doc:"URL presigned untuk pratinjau"`
	CreatedAt   time.Time `json:"created_at"`
}

// --- Konversi ---

// fileURLExpiry adalah masa berlaku presigned GET URL untuk file privat.
const fileURLExpiry = 15 * time.Minute

// fileURL menandatangani object key menjadi presigned GET URL. Bila storage
// tidak dikonfigurasi, nilai sudah berupa URL penuh, atau penandatanganan
// gagal, nilai asli dikembalikan apa adanya.
func (h *Handler) fileURL(ctx context.Context, key string) string {
	if key == "" {
		return key
	}
	if h.files == nil || strings.HasPrefix(key, "http://") || strings.HasPrefix(key, "https://") {
		return key
	}
	url, err := h.files.PresignedGetURL(ctx, key, fileURLExpiry)
	if err != nil {
		return key
	}
	return url
}

func (h *Handler) toAssessmentDetailDTO(ctx context.Context, d domain.AssessmentDetail) AssessmentDetailDTO {
	units := make([]AssessmentUnitDTO, 0, len(d.Units))
	for _, unit := range d.Units {
		units = append(units, h.toAssessmentUnitDTO(ctx, unit))
	}
	return AssessmentDetailDTO{
		Assessment: toAssessmentDTO(d.Assessment),
		Application: ApplicationInfoDTO{
			CertificationSchemaID: d.Application.SchemaID,
			SchemaName:            d.Application.SchemaName,
			SchemaCode:            d.Application.SchemaCode,
			Purpose:               d.Application.Purpose,
			DataCreatedAt:         d.Application.DataCreatedAt,
		},
		Units:  units,
		Locked: d.Locked,
	}
}

func toAssessmentDTO(a domain.Assessment) AssessmentDTO {
	dto := AssessmentDTO{
		ID:                     a.ID,
		ApplicationID:          a.ApplicationID,
		Code:                   a.Code,
		Status:                 string(a.Status),
		AsesorNote:             a.AsesorNote,
		SubmittedAt:            a.SubmittedAt,
		ReviewedBy:             a.ReviewedBy,
		ReviewedAt:             a.ReviewedAt,
		ZoomLink:               a.ZoomLink,
		ZoomTime:               a.ZoomTime,
		SchemaName:             a.SchemaName,
		SchemaCode:             a.SchemaCode,
		AsesiName:              a.AsesiName,
		ApplicationCode:        a.ApplicationCode,
		ApplicationStatus:      a.ApplicationStatus,
		ApplicationNote:        a.ApplicationNote,
		AsesorID:               a.AsesorID,
		AsesorName:             a.AsesorName,
		AsesorAssignmentStatus: a.AsesorAssignmentStatus,
		CreatedAt:              a.CreatedAt,
		UpdatedAt:              a.UpdatedAt,
	}
	if a.AsesorDecision != nil {
		value := string(*a.AsesorDecision)
		dto.AsesorDecision = &value
	}
	if a.ZoomDate != nil {
		value := a.ZoomDate.Format(time.DateOnly)
		dto.ZoomDate = &value
	}
	return dto
}

func (h *Handler) toAssessmentUnitDTO(ctx context.Context, u domain.UnitAnswer) AssessmentUnitDTO {
	files := make([]AssessmentUnitFileDTO, 0, len(u.Files))
	for _, file := range u.Files {
		files = append(files, AssessmentUnitFileDTO{
			ID:          file.ID,
			FileURL:     file.FileURL,
			FileViewURL: h.fileURL(ctx, file.FileURL),
			CreatedAt:   file.CreatedAt,
		})
	}
	return AssessmentUnitDTO{
		ID:               u.ID,
		CompetencyUnitID: u.CompetencyUnitID,
		UnitCode:         u.UnitCode,
		UnitTitle:        u.UnitTitle,
		IsCompetent:      u.IsCompetent,
		Files:            files,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
	}
}
