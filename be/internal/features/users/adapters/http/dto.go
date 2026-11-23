package http

import (
	"time"

	"lsp/internal/features/users/domain"
	httpplatform "lsp/internal/platform/http"
)

type CreateUserInput struct {
	Body CreateUserBody `json:"body"`
}

type CreateUserBody struct {
	Email           string  `json:"email" format:"email" maxLength:"255"`
	Password        string  `json:"password" minLength:"8" maxLength:"72"`
	Name            string  `json:"name" minLength:"2" maxLength:"255"`
	Role            string  `json:"role" enum:"admin,asesi,asesor,validator,manager_sertifikasi"`
	Phone           *string `json:"phone,omitempty" maxLength:"20"`
	NomorRegistrasi *string `json:"nomor_registrasi,omitempty" maxLength:"50"`
}

type LoginInput struct {
	Body LoginBody `json:"body"`
}

type LoginBody struct {
	Email    string `json:"email" format:"email" maxLength:"255"`
	Password string `json:"password" minLength:"8" maxLength:"72"`
}

type GetUserInput struct {
	ID string `path:"id"`
}

type UpdateUserInput struct {
	ID   string         `path:"id"`
	Body UpdateUserBody `json:"body"`
}

type UpdateUserBody struct {
	Name            *string `json:"name,omitempty" minLength:"2" maxLength:"255"`
	Phone           *string `json:"phone,omitempty" maxLength:"20"`
	NomorRegistrasi *string `json:"nomor_registrasi,omitempty" maxLength:"50"`
	Password        *string `json:"password,omitempty" minLength:"8" maxLength:"72"`
	IsValid         *bool   `json:"isvalid,omitempty"`
}

type DeleteUserInput struct {
	ID string `path:"id"`
}

type DeleteUserOutput struct{}

type EmptyInput struct{}

type ListUsersInput struct {
	httpplatform.PaginationQuery
	Roles []string `query:"roles" doc:"Filter kode role (bisa berulang)"`
}

type AuthOutput struct {
	Body AuthResponse `json:"body"`
}

type AuthResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

type UserOutput struct {
	Body UserDTO `json:"body"`
}

type UserDTO struct {
	ID              string     `json:"id"`
	Email           string     `json:"email"`
	Name            string     `json:"name"`
	Role            string     `json:"role"`
	IsValid         bool       `json:"isvalid"`
	Phone           *string    `json:"phone,omitempty"`
	NomorRegistrasi *string    `json:"nomor_registrasi,omitempty"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func toUserDTO(user domain.User) UserDTO {
	return UserDTO{ID: user.ID, Email: user.Email, Name: user.Name, Role: string(user.Role), IsValid: user.IsValid, Phone: user.Phone, NomorRegistrasi: user.NomorRegistrasi, EmailVerifiedAt: user.EmailVerifiedAt, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}
}
