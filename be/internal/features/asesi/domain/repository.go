package domain

import "context"

type ListParams struct {
	Limit  int
	Offset int
}

type ListResult struct {
	Items []Profile
	Total int
}

type OccupationListResult struct {
	Items []Occupation
	Total int
}

type Repository interface {
	Create(ctx context.Context, profile Profile) error
	Patch(ctx context.Context, patch Patch) (Profile, error)
	FindByID(ctx context.Context, id string) (Profile, error)
	FindByUserID(ctx context.Context, userID string) (Profile, error)
	List(ctx context.Context, params ListParams) (ListResult, error)
	Delete(ctx context.Context, id string) error

	// Occupation (data pekerjaan asesi, relasi 1:1 ke profil).
	CreateOccupation(ctx context.Context, occupation Occupation) error
	PatchOccupation(ctx context.Context, patch OccupationPatch) (Occupation, error)
	FindOccupationByID(ctx context.Context, id string) (Occupation, error)
	FindOccupationByUserID(ctx context.Context, userID string) (Occupation, error)
	// UpsertOccupationByUser membuat/memperbarui data pekerjaan milik user (asesi).
	UpsertOccupationByUser(ctx context.Context, userID string, occupation Occupation) (Occupation, error)
	ListOccupations(ctx context.Context, params ListParams) (OccupationListResult, error)
	DeleteOccupation(ctx context.Context, id string) error
}

type AsesiChecker interface {
	IsAsesi(ctx context.Context, userID string) (bool, error)
}
