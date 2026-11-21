package http

import (
	"context"

	"lsp/internal/features/tuk/app"
	httpplatform "lsp/internal/platform/http"
)

type Handler struct {
	service *app.Service
}

func NewHandler(service *app.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(ctx context.Context, input *CreateInput) (*TUKOutput, error) {
	cmd := toCreateCommand(input.Body)
	tuk, err := h.service.Create(ctx, app.CreateCommand{
		Name:       cmd.Name,
		Logo:       cmd.Logo,
		PostalCode: cmd.PostalCode,
		Phone:      cmd.Phone,
		Email:      cmd.Email,
		Province:   cmd.Province,
		City:       cmd.City,
		Address:    cmd.Address,
		MapsURL:    cmd.MapsURL,
		Letterhead: cmd.Letterhead,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &TUKOutput{Body: toDTO(tuk)}, nil
}

func (h *Handler) Update(ctx context.Context, input *UpdateInput) (*TUKOutput, error) {
	tuk, err := h.service.Patch(ctx, app.PatchCommand{
		ID:         input.ID,
		Name:       input.Body.Name,
		Logo:       input.Body.Logo,
		PostalCode: input.Body.PostalCode,
		Phone:      input.Body.Phone,
		Email:      input.Body.Email,
		Province:   input.Body.Province,
		City:       input.Body.City,
		Address:    input.Body.Address,
		MapsURL:    input.Body.MapsURL,
		Letterhead: input.Body.Letterhead,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &TUKOutput{Body: toDTO(tuk)}, nil
}

func (h *Handler) Get(ctx context.Context, input *GetInput) (*TUKOutput, error) {
	tuk, err := h.service.Get(ctx, input.ID)
	if err != nil {
		return nil, mapError(err)
	}
	return &TUKOutput{Body: toDTO(tuk)}, nil
}

func (h *Handler) List(ctx context.Context, input *ListInput) (*httpplatform.PagedResponse[TUKDTO], error) {
	limit := input.SafeLimit()
	result, err := h.service.List(ctx, limit, input.Offset())
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]TUKDTO, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toDTO(item))
	}
	return httpplatform.NewPagedResponse(items, input.Page, limit, result.Total), nil
}

func (h *Handler) Delete(ctx context.Context, input *DeleteInput) (*DeleteOutput, error) {
	if err := h.service.Delete(ctx, input.ID); err != nil {
		return nil, mapError(err)
	}
	return &DeleteOutput{}, nil
}
