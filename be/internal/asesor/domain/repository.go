package domain

import "context"

type ListParams struct {
	Limit  int
	Offset int
}

type ListResult struct {
	Items []ListItem
	Total int
}

// Repository adalah port penyimpanan profil asesor.
type Repository interface {
	Upsert(ctx context.Context, profile Profile) (Profile, error)
	FindByUserID(ctx context.Context, userID string) (Profile, error)
	List(ctx context.Context, params ListParams) (ListResult, error)
}

// AsesorChecker memastikan sebuah user benar-benar ber-role asesor.
// Port didefinisikan di sini dan diimplementasikan oleh modul users
// (lihat internal/users/adapters/postgres), lalu disambungkan di bootstrap.
type AsesorChecker interface {
	IsAsesor(ctx context.Context, userID string) (bool, error)
}
