package http

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"lsp/internal/features/registration/app"
)

type Handler struct {
	service *app.Service
}

func NewHandler(service *app.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(ctx context.Context, input *RegisterInput) (*RegisterOutput, error) {
	birthDate, err := time.Parse(time.DateOnly, input.Body.BirthDate)
	if err != nil {
		return nil, huma.Error400BadRequest("birth_date must be a valid date (YYYY-MM-DD)")
	}
	result, err := h.service.Register(ctx, app.RegisterCommand{
		Email:             input.Body.Email,
		Password:          input.Body.Password,
		FullName:          input.Body.FullName,
		Phone:             input.Body.Phone,
		NIK:               input.Body.NIK,
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
	return &RegisterOutput{Body: RegisterResponse{
		Token: result.Token,
		User:  RegisterUserDTO{ID: result.UserID, Email: result.Email, Name: result.Name, Role: result.Role},
	}}, nil
}
