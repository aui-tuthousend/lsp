package domain

import "context"

type ListParams struct {
	Limit  int
	Offset int
}

type ListResult struct {
	Users []User
	Total int
}

type Repository interface {
	Create(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id string) (User, error)
	List(ctx context.Context, params ListParams) (ListResult, error)
}

func IsValidRole(role Role) bool {
	switch role {
	case RoleAdmin, RoleAsesi, RoleAsesor, RoleValidator:
		return true
	default:
		return false
	}
}
