package domain

import "context"

type ListParams struct {
	Limit  int
	Offset int
}

type ListResult struct {
	Items []TUK
	Total int
}

type Repository interface {
	Create(ctx context.Context, tuk TUK) error
	Update(ctx context.Context, tuk TUK) (TUK, error)
	FindByID(ctx context.Context, id string) (TUK, error)
	List(ctx context.Context, params ListParams) (ListResult, error)
	Delete(ctx context.Context, id string) error
}
