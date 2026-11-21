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

type Repository interface {
	Create(ctx context.Context, profile Profile) error
	Patch(ctx context.Context, patch Patch) (Profile, error)
	FindByID(ctx context.Context, id string) (Profile, error)
	FindByUserID(ctx context.Context, userID string) (Profile, error)
	List(ctx context.Context, params ListParams) (ListResult, error)
	Delete(ctx context.Context, id string) error
}

type AsesiChecker interface {
	IsAsesi(ctx context.Context, userID string) (bool, error)
}
