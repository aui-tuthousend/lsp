package app

import (
	"context"

	"lsp/internal/platform/id"
	"lsp/internal/tuk/domain"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{repo: repo}
}

type CreateCommand struct {
	Name       string
	Logo       *string
	PostalCode string
	Phone      string
	Email      string
	Province   string
	City       string
	Address    string
	MapsURL    string
	Letterhead *string
}

type UpdateCommand struct {
	ID         string
	Name       string
	Logo       *string
	PostalCode string
	Phone      string
	Email      string
	Province   string
	City       string
	Address    string
	MapsURL    string
	Letterhead *string
}

func (s *Service) Create(ctx context.Context, cmd CreateCommand) (domain.TUK, error) {
	tuk := domain.TUK{
		ID:         id.New(),
		Name:       cmd.Name,
		Logo:       cmd.Logo,
		PostalCode: cmd.PostalCode,
		Phone:      cmd.Phone,
		Email:      cmd.Email,
		Province:   cmd.Province,
		City:       cmd.City,
		Address:    cmd.Address,
		MapsURL:    cmd.MapsURL,
		Letterhead: cmd.Letterhead,
	}
	if err := s.repo.Create(ctx, tuk); err != nil {
		return domain.TUK{}, err
	}
	return s.repo.FindByID(ctx, tuk.ID)
}

func (s *Service) Update(ctx context.Context, cmd UpdateCommand) (domain.TUK, error) {
	return s.repo.Update(ctx, domain.TUK{
		ID:         cmd.ID,
		Name:       cmd.Name,
		Logo:       cmd.Logo,
		PostalCode: cmd.PostalCode,
		Phone:      cmd.Phone,
		Email:      cmd.Email,
		Province:   cmd.Province,
		City:       cmd.City,
		Address:    cmd.Address,
		MapsURL:    cmd.MapsURL,
		Letterhead: cmd.Letterhead,
	})
}

func (s *Service) Get(ctx context.Context, id string) (domain.TUK, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) List(ctx context.Context, limit, offset int) (domain.ListResult, error) {
	return s.repo.List(ctx, domain.ListParams{Limit: limit, Offset: offset})
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
