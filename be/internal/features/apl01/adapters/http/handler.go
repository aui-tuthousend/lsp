package http

import (
	"context"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/apl01/app"
	"lsp/internal/features/apl01/domain"
	"lsp/internal/platform/auth"
	httpplatform "lsp/internal/platform/http"
	"lsp/internal/platform/storage"
)

type Handler struct {
	service *app.Service
	files   storage.ObjectStorage
}

func NewHandler(service *app.Service, files storage.ObjectStorage) *Handler {
	return &Handler{service: service, files: files}
}

// CreateDraft membuat draft APL-01 baru atau mengembalikan draft aktif asesi.
func (h *Handler) CreateDraft(ctx context.Context, _ *CreateDraftInput) (*ApplicationOutput, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("unauthenticated")
	}
	detail, err := h.service.CreateOrGetDraft(ctx, claims.UserID)
	if err != nil {
		return nil, mapError(err)
	}
	return &ApplicationOutput{Body: h.toApplicationDetailDTO(ctx, detail)}, nil
}

// ListMine menampilkan daftar permohonan milik asesi yang sedang login.
func (h *Handler) ListMine(ctx context.Context, input *ListMineInput) (*httpplatform.PagedResponse[ApplicationDTO], error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("unauthenticated")
	}
	return h.list(ctx, input.PaginationQuery, claims.UserID, input.Status)
}

// ListAll menampilkan seluruh permohonan (khusus admin/staff).
func (h *Handler) ListAll(ctx context.Context, input *ListAllInput) (*httpplatform.PagedResponse[ApplicationDTO], error) {
	return h.list(ctx, input.PaginationQuery, input.UserID, input.Status)
}

func (h *Handler) list(ctx context.Context, page httpplatform.PaginationQuery, userID, status string) (*httpplatform.PagedResponse[ApplicationDTO], error) {
	var filter *domain.Status
	if status != "" {
		value := domain.Status(status)
		filter = &value
	}
	limit := page.SafeLimit()
	result, err := h.service.List(ctx, userID, filter, limit, page.Offset())
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]ApplicationDTO, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toApplicationDTO(item))
	}
	return httpplatform.NewPagedResponse(items, page.Page, limit, result.Total), nil
}

func (h *Handler) Get(ctx context.Context, input *GetInput) (*ApplicationOutput, error) {
	detail, err := h.service.Detail(ctx, input.ID)
	if err != nil {
		return nil, mapError(err)
	}
	return &ApplicationOutput{Body: h.toApplicationDetailDTO(ctx, detail)}, nil
}

func (h *Handler) UpdateStep(ctx context.Context, input *UpdateStepInput) (*ApplicationOutput, error) {
	detail, err := h.service.UpdateStep(ctx, input.ID, input.Body.CurrentStep)
	if err != nil {
		return nil, mapError(err)
	}
	return &ApplicationOutput{Body: h.toApplicationDetailDTO(ctx, detail)}, nil
}

func (h *Handler) PutData(ctx context.Context, input *PutDataInput) (*ApplicationOutput, error) {
	detail, err := h.service.PutData(ctx, app.PutDataCommand{
		ApplicationID:         input.ID,
		CertificationSchemaID: input.Body.CertificationSchemaID,
		Purpose:               domain.Purpose(input.Body.Purpose),
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &ApplicationOutput{Body: h.toApplicationDetailDTO(ctx, detail)}, nil
}

func (h *Handler) PutDocuments(ctx context.Context, input *PutDocumentsInput) (*ApplicationOutput, error) {
	items := make([]app.DocumentItem, 0, len(input.Body.Items))
	for _, item := range input.Body.Items {
		items = append(items, app.DocumentItem{
			DocumentTypeID: item.DocumentTypeID,
			FileURL:        item.FileURL,
		})
	}
	detail, err := h.service.PutDocuments(ctx, app.PutDocumentsCommand{
		ApplicationID: input.ID,
		Items:         items,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &ApplicationOutput{Body: h.toApplicationDetailDTO(ctx, detail)}, nil
}

func (h *Handler) PutSchedule(ctx context.Context, input *PutScheduleInput) (*ApplicationOutput, error) {
	scheduledDate, err := time.Parse(time.DateOnly, input.Body.ScheduledDate)
	if err != nil {
		return nil, huma.Error400BadRequest("scheduled_date must be a valid date (YYYY-MM-DD)")
	}
	detail, err := h.service.PutSchedule(ctx, app.PutScheduleCommand{
		ApplicationID: input.ID,
		Mode:          domain.Mode(input.Body.Mode),
		ScheduledDate: scheduledDate,
		TukID:         input.Body.TukID,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &ApplicationOutput{Body: h.toApplicationDetailDTO(ctx, detail)}, nil
}

func (h *Handler) Submit(ctx context.Context, input *SubmitInput) (*SubmitOutput, error) {
	result, err := h.service.Submit(ctx, input.ID, derefString(input.Body.VerifyToken))
	if err != nil {
		return nil, mapError(err)
	}
	return &SubmitOutput{Body: SubmitDTO{
		Application: h.toApplicationDetailDTO(ctx, result.Application),
		VerifyToken: result.VerifyToken,
	}}, nil
}

func (h *Handler) Verify(ctx context.Context, input *VerifyInput) (*ApplicationOutput, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("unauthenticated")
	}
	detail, err := h.service.Verify(ctx, claims.UserID, input.Body.Token, clientIP(input.ForwardedFor), headerValue(input.UserAgent))
	if err != nil {
		return nil, mapError(err)
	}
	return &ApplicationOutput{Body: h.toApplicationDetailDTO(ctx, detail)}, nil
}

// VerifyStatus mengembalikan status token verifikasi tanpa memakainya.
func (h *Handler) VerifyStatus(ctx context.Context, input *VerifyStatusInput) (*VerifyStatusOutput, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("unauthenticated")
	}
	signature, application, err := h.service.VerifyStatus(ctx, claims.UserID, input.Token)
	if err != nil {
		return nil, mapError(err)
	}
	return &VerifyStatusOutput{Body: VerifyStatusDTO{
		Status:          string(signature.Status),
		ApplicationCode: application.Code,
		ExpiresAt:       signature.ExpiresAt,
		VerifiedAt:      signature.UsedAt,
	}}, nil
}

func (h *Handler) Decide(ctx context.Context, input *DecideInput) (*ApplicationOutput, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("unauthenticated")
	}
	detail, err := h.service.Decide(ctx, app.DecideCommand{
		ApplicationID: input.ID,
		Decision:      domain.Status(input.Body.Decision),
		Note:          input.Body.Note,
		AsesorID:      input.Body.AsesorID,
		VerifiedBy:    claims.UserID,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &ApplicationOutput{Body: h.toApplicationDetailDTO(ctx, detail)}, nil
}

// AssignAsesor menugaskan asesor untuk permohonan yang sudah disetujui admin.
func (h *Handler) AssignAsesor(ctx context.Context, input *AssignAsesorInput) (*ApplicationOutput, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("unauthenticated")
	}
	detail, err := h.service.AssignAsesor(ctx, app.AssignAsesorCommand{
		ApplicationID: input.ID,
		AsesorID:      input.Body.AsesorID,
		AssignedBy:    claims.UserID,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &ApplicationOutput{Body: h.toApplicationDetailDTO(ctx, detail)}, nil
}

// AsesorCandidates menampilkan matriks kandidat asesor untuk plotting.
func (h *Handler) AsesorCandidates(ctx context.Context, input *AsesorCandidatesInput) (*AsesorCandidatesOutput, error) {
	candidates, err := h.service.AsesorCandidates(ctx, input.ID)
	if err != nil {
		return nil, mapError(err)
	}
	return &AsesorCandidatesOutput{Body: AsesorCandidatesDTO{Items: toAsesorCandidateDTOs(candidates)}}, nil
}

func headerValue(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

// clientIP mengambil IP klien dari header X-Forwarded-For (nilai pertama).
func clientIP(forwardedFor string) *string {
	if forwardedFor == "" {
		return nil
	}
	ip := strings.TrimSpace(strings.Split(forwardedFor, ",")[0])
	if ip == "" {
		return nil
	}
	return &ip
}
