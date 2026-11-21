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

// Repository adalah port penyimpanan profil asesor beserta sertifikasinya.
type Repository interface {
	Upsert(ctx context.Context, profile Profile) (Profile, error)
	PatchProfile(ctx context.Context, patch ProfilePatch) (Profile, error)
	FindByUserID(ctx context.Context, userID string) (Profile, error)
	List(ctx context.Context, params ListParams) (ListResult, error)

	ListCertifications(ctx context.Context, userID string) ([]Certification, error)
	FindCertificationByID(ctx context.Context, userID, certID string) (Certification, error)
	CreateCertification(ctx context.Context, cert Certification) error
	PatchCertification(ctx context.Context, patch CertificationPatch) (Certification, error)
	DeleteCertification(ctx context.Context, userID, certID string) error
}

// AsesorChecker memastikan sebuah user benar-benar ber-role asesor.
// Port didefinisikan di sini dan diimplementasikan oleh modul users
// (lihat internal/features/users/adapters/postgres), lalu disambungkan di bootstrap.
type AsesorChecker interface {
	IsAsesor(ctx context.Context, userID string) (bool, error)
}
