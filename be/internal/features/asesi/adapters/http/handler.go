package http

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/asesi/app"
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
	return &AsesiProfileOutput{Body: h.toDTO(ctx, profile)}, nil
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
	return &AsesiProfileOutput{Body: h.toDTO(ctx, profile)}, nil
}

func (h *Handler) Get(ctx context.Context, input *GetInput) (*AsesiProfileOutput, error) {
	profile, err := h.service.Get(ctx, input.ID)
	if err != nil {
		return nil, mapError(err)
	}
	return &AsesiProfileOutput{Body: h.toDTO(ctx, profile)}, nil
}

func (h *Handler) GetByUser(ctx context.Context, input *GetByUserInput) (*AsesiProfileOutput, error) {
	profile, err := h.service.GetByUserID(ctx, input.UserID)
	if err != nil {
		return nil, mapError(err)
	}
	return &AsesiProfileOutput{Body: h.toDTO(ctx, profile)}, nil
}

func (h *Handler) List(ctx context.Context, input *ListInput) (*httpplatform.PagedResponse[AsesiProfileDTO], error) {
	limit := input.SafeLimit()
	result, err := h.service.List(ctx, limit, input.Offset())
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]AsesiProfileDTO, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, h.toDTO(ctx, item))
	}
	return httpplatform.NewPagedResponse(items, input.Page, limit, result.Total), nil
}

func (h *Handler) Delete(ctx context.Context, input *DeleteInput) (*DeleteOutput, error) {
	if err := h.service.Delete(ctx, input.ID); err != nil {
		return nil, mapError(err)
	}
	return &DeleteOutput{}, nil
}

// --- Data pekerjaan asesi (asesi_occupations) ---

func (h *Handler) CreateOccupation(ctx context.Context, input *CreateOccupationInput) (*OccupationOutput, error) {
	occupation, err := h.service.CreateOccupation(ctx, app.CreateOccupationCommand{
		AsesiProfileID:    input.Body.AsesiProfileID,
		CompanyName:       input.Body.CompanyName,
		Position:          input.Body.Position,
		TenureMonths:      input.Body.TenureMonths,
		SubordinatesCount: input.Body.SubordinatesCount,
		Address:           input.Body.Address,
		Province:          input.Body.Province,
		City:              input.Body.City,
		District:          input.Body.District,
		PostalCode:        input.Body.PostalCode,
		OfficeEmail:       input.Body.OfficeEmail,
		OfficePhone:       input.Body.OfficePhone,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &OccupationOutput{Body: toOccupationDTO(occupation)}, nil
}

// UpsertOccupationByUser membuat/memperbarui data pekerjaan milik asesi sendiri.
func (h *Handler) UpsertOccupationByUser(ctx context.Context, input *UpsertOccupationByUserInput) (*OccupationOutput, error) {
	occupation, err := h.service.UpsertOccupationByUser(ctx, app.UpsertOccupationCommand{
		UserID:            input.UserID,
		CompanyName:       input.Body.CompanyName,
		Position:          input.Body.Position,
		TenureMonths:      input.Body.TenureMonths,
		SubordinatesCount: input.Body.SubordinatesCount,
		Address:           input.Body.Address,
		Province:          input.Body.Province,
		City:              input.Body.City,
		District:          input.Body.District,
		PostalCode:        input.Body.PostalCode,
		OfficeEmail:       input.Body.OfficeEmail,
		OfficePhone:       input.Body.OfficePhone,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &OccupationOutput{Body: toOccupationDTO(occupation)}, nil
}

func (h *Handler) UpdateOccupation(ctx context.Context, input *UpdateOccupationInput) (*OccupationOutput, error) {
	occupation, err := h.service.PatchOccupation(ctx, app.PatchOccupationCommand{
		ID:                input.ID,
		CompanyName:       input.Body.CompanyName,
		Position:          input.Body.Position,
		TenureMonths:      input.Body.TenureMonths,
		SubordinatesCount: input.Body.SubordinatesCount,
		Address:           input.Body.Address,
		Province:          input.Body.Province,
		City:              input.Body.City,
		District:          input.Body.District,
		PostalCode:        input.Body.PostalCode,
		OfficeEmail:       input.Body.OfficeEmail,
		OfficePhone:       input.Body.OfficePhone,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &OccupationOutput{Body: toOccupationDTO(occupation)}, nil
}

func (h *Handler) GetOccupation(ctx context.Context, input *GetOccupationInput) (*OccupationOutput, error) {
	occupation, err := h.service.GetOccupation(ctx, input.ID)
	if err != nil {
		return nil, mapError(err)
	}
	return &OccupationOutput{Body: toOccupationDTO(occupation)}, nil
}

func (h *Handler) GetOccupationByUser(ctx context.Context, input *GetOccupationByUserInput) (*OccupationOutput, error) {
	occupation, err := h.service.GetOccupationByUserID(ctx, input.UserID)
	if err != nil {
		return nil, mapError(err)
	}
	return &OccupationOutput{Body: toOccupationDTO(occupation)}, nil
}

func (h *Handler) ListOccupations(ctx context.Context, input *ListOccupationInput) (*httpplatform.PagedResponse[OccupationDTO], error) {
	limit := input.SafeLimit()
	result, err := h.service.ListOccupations(ctx, limit, input.Offset())
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]OccupationDTO, 0, len(result.Items))
	for _, item := range result.Items {
		items = append(items, toOccupationDTO(item))
	}
	return httpplatform.NewPagedResponse(items, input.Page, limit, result.Total), nil
}

func (h *Handler) DeleteOccupation(ctx context.Context, input *DeleteOccupationInput) (*DeleteOccupationOutput, error) {
	if err := h.service.DeleteOccupation(ctx, input.ID); err != nil {
		return nil, mapError(err)
	}
	return &DeleteOccupationOutput{}, nil
}
