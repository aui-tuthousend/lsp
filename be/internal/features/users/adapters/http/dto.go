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
	Role            string  `json:"role" enum:"admin,admin_lsp,asesi,asesor,validator,manager_sertifikasi"`
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

// UpdateMeInput/UpdateMeBody dipakai user yang login untuk mengubah akunnya
// sendiri (hanya nama & telepon).
type UpdateMeInput struct {
	Body UpdateMeBody `json:"body"`
}

type UpdateMeBody struct {
	Name  *string `json:"name,omitempty" minLength:"2" maxLength:"255"`
	Phone *string `json:"phone,omitempty" maxLength:"20"`
}

type EmptyInput struct{}

type ListUsersInput struct {
	httpplatform.PaginationQuery
	Roles []string `query:"roles,explode" doc:"Filter kode role (bisa berulang)"`
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

// --- Master data roles ---

type CreateRoleInput struct {
	Body RoleBody `json:"body"`
}

type UpdateRoleInput struct {
	ID   string        `path:"id"`
	Body RolePatchBody `json:"body"`
}

type GetRoleInput struct {
	ID string `path:"id"`
}

type DeleteRoleInput struct {
	ID string `path:"id"`
}

type ListRolesInput struct{}

type RoleBody struct {
	Code        string  `json:"code" minLength:"1" maxLength:"50"`
	Name        string  `json:"name" minLength:"1" maxLength:"100"`
	Description *string `json:"description,omitempty"`
}

type RolePatchBody struct {
	Code        *string `json:"code,omitempty" minLength:"1" maxLength:"50"`
	Name        *string `json:"name,omitempty" minLength:"1" maxLength:"100"`
	Description *string `json:"description,omitempty"`
}

type RoleOutput struct {
	Body RoleDTO `json:"body"`
}

type RoleListOutput struct {
	Body struct {
		Data []RoleDTO `json:"data"`
	} `json:"body"`
}

type DeleteRoleOutput struct{}

type RoleDTO struct {
	ID          string  `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

func toRoleDTO(role domain.RoleRecord) RoleDTO {
	return RoleDTO{ID: role.ID, Code: string(role.Code), Name: role.Name, Description: role.Description}
}
