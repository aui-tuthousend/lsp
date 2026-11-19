package domain

import "time"

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleAsesi     Role = "asesi"
	RoleAsesor    Role = "asesor"
	RoleValidator Role = "validator"
)

type User struct {
	ID              string
	Role            Role
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
