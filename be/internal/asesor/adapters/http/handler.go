package http

import (
	"context"

	"lsp/internal/asesor/app"
	httpplatform "lsp/internal/platform/http"
)

type Handler struct {
	service *app.Service
}

func NewHandler(service *app.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Upsert(ctx context.Context, input *UpsertProfileInput) (*ProfileOutput, error) {
	profile, err := h.service.Upsert(ctx, app.UpsertCommand{
		UserID:               input.UserID,
		TUKID:                input.Body.TUKID,
		Province:             input.Body.Province,
		City:                 input.Body.City,
		DailyAssignmentLimit: input.Body.DailyAssignmentLimit,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &ProfileOutput{Body: toProfileDTO(profile)}, nil
}

func (h *Handler) Get(ctx context.Context, input *GetProfileInput) (*ProfileOutput, error) {
	profile, err := h.service.Get(ctx, input.UserID)
	if err != nil {
		return nil, mapError(err)
	}
	return &ProfileOutput{Body: toProfileDTO(profile)}, nil
}

func (h *Handler) List(ctx context.Context, input *ListProfilesInput) (*httpplatform.PagedResponse[ProfileDTO], error) {
	limit := input.SafeLimit()
	result, err := h.service.List(ctx, limit, input.Offset())
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]ProfileDTO, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toListItemDTO(item))
	}
	return httpplatform.NewPagedResponse(items, input.Page, limit, result.Total), nil
}
