package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/asesi/app"
	httpplatform "lsp/internal/platform/http"
)

type Handler struct {
	service *app.Service
}

func NewHandler(service *app.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(ctx context.Context, input *CreateInput) (*AsesiProfileOutput, error) {
	birthDate, err := parseDate(input.Body.BirthDate)
	if err != nil {
		return nil, huma.Error400BadRequest("birth_date must be a valid date (YYYY-MM-DD)")
	}
	profile, err := h.service.Create(ctx, app.CreateCommand{
		UserID:            input.Body.UserID,
		FullName:          input.Body.FullName,
		NIK:               input.Body.NIK,
		Phone:             input.Body.Phone,
		Gender:            input.Body.Gender,
		BirthPlace:        input.Body.BirthPlace,
		BirthDate:         birthDate,
		Nationality:       input.Body.Nationality,
		Occupation:        input.Body.Occupation,
		EducationLevel:    input.Body.EducationLevel,
		KTPAddress:        input.Body.KTPAddress,
		DomicileSameAsKTP: input.Body.DomicileSameAsKTP,
		DomicileAddress:   input.Body.DomicileAddress,
		Province:          input.Body.Province,
		City:              input.Body.City,
		District:          input.Body.District,
		Village:           input.Body.Village,
		PostalCode:        input.Body.PostalCode,
		Photo:             input.Body.Photo,
		KTPFile:           input.Body.KTPFile,
		CVFile:            input.Body.CVFile,
		NPWPFile:          input.Body.NPWPFile,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &AsesiProfileOutput{Body: toDTO(profile)}, nil
}

func (h *Handler) Update(ctx context.Context, input *UpdateInput) (*AsesiProfileOutput, error) {
	birthDate, err := parseOptionalDate(input.Body.BirthDate)
	if err != nil {
		return nil, huma.Error400BadRequest("birth_date must be a valid date (YYYY-MM-DD)")
	}
	profile, err := h.service.Patch(ctx, app.PatchCommand{
		ID:                input.ID,
		FullName:          input.Body.FullName,
		NIK:               input.Body.NIK,
		Phone:             input.Body.Phone,
		Gender:            input.Body.Gender,
		BirthPlace:        input.Body.BirthPlace,
		BirthDate:         birthDate,
		Nationality:       input.Body.Nationality,
		Occupation:        input.Body.Occupation,
		EducationLevel:    input.Body.EducationLevel,
		KTPAddress:        input.Body.KTPAddress,
		DomicileSameAsKTP: input.Body.DomicileSameAsKTP,
		DomicileAddress:   input.Body.DomicileAddress,
		Province:          input.Body.Province,
		City:              input.Body.City,
		District:          input.Body.District,
		Village:           input.Body.Village,
		PostalCode:        input.Body.PostalCode,
		Photo:             input.Body.Photo,
		KTPFile:           input.Body.KTPFile,
		CVFile:            input.Body.CVFile,
		NPWPFile:          input.Body.NPWPFile,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &AsesiProfileOutput{Body: toDTO(profile)}, nil
}

func (h *Handler) Get(ctx context.Context, input *GetInput) (*AsesiProfileOutput, error) {
	profile, err := h.service.Get(ctx, input.ID)
	if err != nil {
		return nil, mapError(err)
	}
	return &AsesiProfileOutput{Body: toDTO(profile)}, nil
}

func (h *Handler) GetByUser(ctx context.Context, input *GetByUserInput) (*AsesiProfileOutput, error) {
	profile, err := h.service.GetByUserID(ctx, input.UserID)
	if err != nil {
		return nil, mapError(err)
	}
	return &AsesiProfileOutput{Body: toDTO(profile)}, nil
}

func (h *Handler) List(ctx context.Context, input *ListInput) (*httpplatform.PagedResponse[AsesiProfileDTO], error) {
	limit := input.SafeLimit()
	result, err := h.service.List(ctx, limit, input.Offset())
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]AsesiProfileDTO, 0, len(result.Items))
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
