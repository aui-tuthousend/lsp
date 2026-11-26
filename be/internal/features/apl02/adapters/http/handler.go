package http

import (
	"context"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/apl02/app"
	"lsp/internal/features/apl02/domain"
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

// Create membuat (atau mengembalikan) asesmen mandiri untuk sebuah permohonan.
func (h *Handler) Create(ctx context.Context, input *CreateInput) (*AssessmentOutput, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("unauthenticated")
	}
	detail, err := h.service.CreateOrGet(ctx, claims.UserID, input.Body.ApplicationID)
	if err != nil {
		return nil, mapError(err)
	}
	return &AssessmentOutput{Body: h.toAssessmentDetailDTO(ctx, detail)}, nil
}

// ListMine menampilkan daftar asesmen mandiri milik asesi yang login.
func (h *Handler) ListMine(ctx context.Context, input *ListMineInput) (*httpplatform.PagedResponse[AssessmentDTO], error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("unauthenticated")
	}
	limit := input.SafeLimit()
	result, err := h.service.List(ctx, claims.UserID, limit, input.Offset())
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]AssessmentDTO, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toAssessmentDTO(item))
	}
	return httpplatform.NewPagedResponse(items, input.Page, limit, result.Total), nil
}

func (h *Handler) Get(ctx context.Context, input *GetInput) (*AssessmentOutput, error) {
	detail, err := h.service.Detail(ctx, input.ID)
	if err != nil {
		return nil, mapError(err)
	}
	return &AssessmentOutput{Body: h.toAssessmentDetailDTO(ctx, detail)}, nil
}

func (h *Handler) PutUnits(ctx context.Context, input *PutUnitsInput) (*AssessmentOutput, error) {
	items := make([]app.UnitAnswerItem, 0, len(input.Body.Items))
	for _, item := range input.Body.Items {
		items = append(items, app.UnitAnswerItem{
			CompetencyUnitID: item.CompetencyUnitID,
			IsCompetent:      item.IsCompetent,
			FileURLs:         item.Files,
		})
	}
	detail, err := h.service.PutAnswers(ctx, app.PutAnswersCommand{
		AssessmentID: input.ID,
		Items:        items,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &AssessmentOutput{Body: h.toAssessmentDetailDTO(ctx, detail)}, nil
}

func (h *Handler) Submit(ctx context.Context, input *SubmitInput) (*AssessmentSubmitOutput, error) {
	result, err := h.service.Submit(ctx, input.ID, derefString(input.Body.VerifyToken))
	if err != nil {
		return nil, mapError(err)
	}
	return &AssessmentSubmitOutput{Body: AssessmentSubmitDTO{
		Assessment:  h.toAssessmentDetailDTO(ctx, result.Assessment),
		VerifyToken: result.VerifyToken,
	}}, nil
}

func (h *Handler) Verify(ctx context.Context, input *VerifyInput) (*AssessmentOutput, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("unauthenticated")
	}
	detail, err := h.service.Verify(ctx, claims.UserID, input.Body.Token, clientIP(input.ForwardedFor), headerValue(input.UserAgent))
	if err != nil {
		return nil, mapError(err)
	}
	return &AssessmentOutput{Body: h.toAssessmentDetailDTO(ctx, detail)}, nil
}

func (h *Handler) VerifyStatus(ctx context.Context, input *VerifyStatusInput) (*AssessmentVerifyStatusOutput, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("unauthenticated")
	}
	signature, assessment, err := h.service.VerifyStatus(ctx, claims.UserID, input.Token)
	if err != nil {
		return nil, mapError(err)
	}
	return &AssessmentVerifyStatusOutput{Body: AssessmentVerifyStatusDTO{
		Status:     signature.Status,
		Code:       assessment.Code,
		ExpiresAt:  signature.ExpiresAt,
		VerifiedAt: signature.UsedAt,
	}}, nil
}

func (h *Handler) Review(ctx context.Context, input *ReviewInput) (*AssessmentOutput, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, huma.Error401Unauthorized("unauthenticated")
	}
	var zoomDate *time.Time
	if input.Body.ZoomDate != nil && *input.Body.ZoomDate != "" {
		parsed, err := time.Parse(time.DateOnly, *input.Body.ZoomDate)
		if err != nil {
			return nil, huma.Error400BadRequest("zoom_date must be a valid date (YYYY-MM-DD)")
		}
		zoomDate = &parsed
	}
	decision := domain.Decision(input.Body.Decision)
	detail, err := h.service.Review(ctx, app.ReviewCommand{
		AssessmentID: input.ID,
		Decision:     &decision,
		Note:         input.Body.Note,
		ZoomLink:     input.Body.ZoomLink,
		ZoomDate:     zoomDate,
		ZoomTime:     input.Body.ZoomTime,
		ReviewedBy:   claims.UserID,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &AssessmentOutput{Body: h.toAssessmentDetailDTO(ctx, detail)}, nil
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
