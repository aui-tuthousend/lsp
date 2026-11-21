package http

import (
	"context"

	"lsp/internal/features/users/app"
	"lsp/internal/features/users/domain"
	"lsp/internal/platform/auth"
	httpplatform "lsp/internal/platform/http"
)

type Handler struct {
	service *app.Service
}

func NewHandler(service *app.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(ctx context.Context, input *RegisterInput) (*AuthOutput, error) {
	result, err := h.service.Register(ctx, app.RegisterCommand{Email: input.Body.Email, Password: input.Body.Password, Name: input.Body.Name, Phone: input.Body.Phone, NomorRegistrasi: input.Body.NomorRegistrasi})
	if err != nil {
		return nil, mapError(err)
	}
	return &AuthOutput{Body: AuthResponse{Token: result.Token, User: toUserDTO(result.User)}}, nil
}

// Create dipakai admin untuk membuat user dengan role apa pun.
func (h *Handler) Create(ctx context.Context, input *CreateUserInput) (*UserOutput, error) {
	user, err := h.service.CreateUser(ctx, app.CreateUserCommand{Email: input.Body.Email, Password: input.Body.Password, Name: input.Body.Name, Role: domain.Role(input.Body.Role), Phone: input.Body.Phone, NomorRegistrasi: input.Body.NomorRegistrasi})
	if err != nil {
		return nil, mapError(err)
	}
	return &UserOutput{Body: toUserDTO(user)}, nil
}

func (h *Handler) Login(ctx context.Context, input *LoginInput) (*AuthOutput, error) {
	result, err := h.service.Login(ctx, app.LoginCommand{Email: input.Body.Email, Password: input.Body.Password})
	if err != nil {
		return nil, mapError(err)
	}
	return &AuthOutput{Body: AuthResponse{Token: result.Token, User: toUserDTO(result.User)}}, nil
}

func (h *Handler) GetMe(ctx context.Context, input *EmptyInput) (*UserOutput, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok {
		return nil, mapError(domain.ErrInvalidCredentials)
	}
	user, err := h.service.GetMe(ctx, claims)
	if err != nil {
		return nil, mapError(err)
	}
	return &UserOutput{Body: toUserDTO(user)}, nil
}

func (h *Handler) List(ctx context.Context, input *ListUsersInput) (*httpplatform.PagedResponse[UserDTO], error) {
	limit := input.SafeLimit()
	result, err := h.service.List(ctx, limit, input.Offset())
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]UserDTO, 0, len(result.Users))
	for _, user := range result.Users {
		items = append(items, toUserDTO(user))
	}
	return httpplatform.NewPagedResponse(items, input.Page, limit, result.Total), nil
}
