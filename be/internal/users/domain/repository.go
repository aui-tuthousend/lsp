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

// Repository adalah port penyimpanan user.
type Repository interface {
	Create(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id string) (User, error)
	List(ctx context.Context, params ListParams) (ListResult, error)
}

// RolesRepository adalah port untuk master data roles.
type RolesRepository interface {
	FindByCode(ctx context.Context, code Role) (RoleRecord, error)
	List(ctx context.Context) ([]RoleRecord, error)
}

// AsesorProfileInitializer membuat profil asesor default saat user ber-role
// asesor dibuat. Port didefinisikan di sini dan diimplementasikan modul asesor
// (internal/asesor/adapters/postgres), lalu disambungkan di bootstrap.
type AsesorProfileInitializer interface {
	EnsureProfile(ctx context.Context, userID string) error
}
