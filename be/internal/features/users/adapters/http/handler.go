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

// Create dipakai admin untuk membuat user dengan role apa pun.
func (h *Handler) Create(ctx context.Context, input *CreateUserInput) (*UserOutput, error) {
	user, err := h.service.CreateUser(ctx, app.CreateUserCommand{Email: input.Body.Email, Password: input.Body.Password, Name: input.Body.Name, Role: domain.Role(input.Body.Role), Phone: input.Body.Phone, NomorRegistrasi: input.Body.NomorRegistrasi})
	if err != nil {
		return nil, mapError(err)
	}
	return &UserOutput{Body: toUserDTO(user)}, nil
}

// Get dipakai admin untuk mengambil satu user.
func (h *Handler) Get(ctx context.Context, input *GetUserInput) (*UserOutput, error) {
	user, err := h.service.GetByID(ctx, input.ID)
	if err != nil {
		return nil, mapError(err)
	}
	return &UserOutput{Body: toUserDTO(user)}, nil
}

// Update dipakai admin untuk mengubah sebagian data user.
func (h *Handler) Update(ctx context.Context, input *UpdateUserInput) (*UserOutput, error) {
	user, err := h.service.PatchUser(ctx, app.PatchUserCommand{
		ID:              input.ID,
		Name:            input.Body.Name,
		Phone:           input.Body.Phone,
		NomorRegistrasi: input.Body.NomorRegistrasi,
		Password:        input.Body.Password,
		IsValid:         input.Body.IsValid,
	})
	if err != nil {
		return nil, mapError(err)
	}
	return &UserOutput{Body: toUserDTO(user)}, nil
}

// Delete dipakai admin untuk menghapus (soft delete) user.
func (h *Handler) Delete(ctx context.Context, input *DeleteUserInput) (*DeleteUserOutput, error) {
	if err := h.service.DeleteUser(ctx, input.ID); err != nil {
		return nil, mapError(err)
	}
	return &DeleteUserOutput{}, nil
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
	result, err := h.service.List(ctx, limit, input.Offset(), input.Roles)
	if err != nil {
		return nil, mapError(err)
	}
	items := make([]UserDTO, 0, len(result.Users))
	for _, user := range result.Users {
		items = append(items, toUserDTO(user))
	}
	return httpplatform.NewPagedResponse(items, input.Page, limit, result.Total), nil
}
