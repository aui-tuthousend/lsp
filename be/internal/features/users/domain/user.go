package domain

import "time"

// Role adalah kode role pada master data `roles`.
type Role string

const (
	RoleAdmin              Role = "admin"
	RoleAsesi              Role = "asesi"
	RoleAsesor             Role = "asesor"
	RoleValidator          Role = "validator"
	RoleManagerSertifikasi Role = "manager_sertifikasi"
)

// RoleRecord adalah satu baris master data roles.
type RoleRecord struct {
	ID          string
	Code        Role
	Name        string
	Description *string
}

type User struct {
	ID              string
	RoleID          string
	Role            Role // kode role (hasil join ke tabel roles)
	IsValid         bool
	Name            string
	Email           string
	Phone           *string
	NomorRegistrasi *string
	PasswordHash    string
	EmailVerifiedAt *time.Time
	RememberToken   *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}

func IsValidRole(role Role) bool {
	switch role {
	case RoleAdmin, RoleAsesi, RoleAsesor, RoleValidator, RoleManagerSertifikasi:
		return true
	default:
		return false
	}
}
