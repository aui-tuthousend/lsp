package app

import (
	"context"
	"time"

	"lsp/internal/features/asesi/domain"
	"lsp/internal/platform/id"
)

type Service struct {
	repo    domain.Repository
	checker domain.AsesiChecker
}

func NewService(repo domain.Repository, checker domain.AsesiChecker) *Service {
	return &Service{repo: repo, checker: checker}
}

type CreateCommand struct {
	UserID            string
	FullName          string
	NIK               string
	Phone             string
	Gender            domain.Gender
	BirthPlace        string
	BirthDate         time.Time
	Nationality       string
	Occupation        string
	EducationLevel    domain.EducationLevel
	KTPAddress        string
	DomicileSameAsKTP bool
	DomicileAddress   *string
	Province          string
	City              string
	District          string
	Village           string
	PostalCode        *string
	Photo             *string
	KTPFile           *string
	CVFile            *string
	NPWPFile          *string
}

func (s *Service) Create(ctx context.Context, cmd CreateCommand) (domain.Profile, error) {
	isAsesi, err := s.checker.IsAsesi(ctx, cmd.UserID)
	if err != nil {
		return domain.Profile{}, err
	}
	if !isAsesi {
		return domain.Profile{}, domain.ErrNotAsesi
	}
	profile := domain.Profile{
		ID:                id.New(),
		UserID:            cmd.UserID,
		FullName:          cmd.FullName,
		NIK:               cmd.NIK,
		Phone:             cmd.Phone,
		Gender:            cmd.Gender,
		BirthPlace:        cmd.BirthPlace,
		BirthDate:         cmd.BirthDate,
		Nationality:       cmd.Nationality,
		Occupation:        cmd.Occupation,
		EducationLevel:    cmd.EducationLevel,
		KTPAddress:        cmd.KTPAddress,
		DomicileSameAsKTP: cmd.DomicileSameAsKTP,
		DomicileAddress:   cmd.DomicileAddress,
		Province:          cmd.Province,
		City:              cmd.City,
		District:          cmd.District,
		Village:           cmd.Village,
		PostalCode:        cmd.PostalCode,
		Photo:             cmd.Photo,
		KTPFile:           cmd.KTPFile,
		CVFile:            cmd.CVFile,
		NPWPFile:          cmd.NPWPFile,
	}
	if err := s.repo.Create(ctx, profile); err != nil {
		return domain.Profile{}, err
	}
	return s.repo.FindByID(ctx, profile.ID)
}

type PatchCommand struct {
	ID                string
	FullName          *string
	NIK               *string
	Phone             *string
	Gender            *domain.Gender
	BirthPlace        *string
	BirthDate         *time.Time
	Nationality       *string
	Occupation        *string
	EducationLevel    *domain.EducationLevel
	KTPAddress        *string
	DomicileSameAsKTP *bool
	DomicileAddress   *string
	Province          *string
	City              *string
	District          *string
	Village           *string
	PostalCode        *string
	Photo             *string
	KTPFile           *string
	CVFile            *string
	NPWPFile          *string
}

func (s *Service) Patch(ctx context.Context, cmd PatchCommand) (domain.Profile, error) {
	return s.repo.Patch(ctx, domain.Patch{
		ID:                cmd.ID,
		FullName:          cmd.FullName,
		NIK:               cmd.NIK,
		Phone:             cmd.Phone,
		Gender:            cmd.Gender,
		BirthPlace:        cmd.BirthPlace,
		BirthDate:         cmd.BirthDate,
		Nationality:       cmd.Nationality,
		Occupation:        cmd.Occupation,
		EducationLevel:    cmd.EducationLevel,
		KTPAddress:        cmd.KTPAddress,
		DomicileSameAsKTP: cmd.DomicileSameAsKTP,
		DomicileAddress:   cmd.DomicileAddress,
		Province:          cmd.Province,
		City:              cmd.City,
		District:          cmd.District,
		Village:           cmd.Village,
		PostalCode:        cmd.PostalCode,
		Photo:             cmd.Photo,
		KTPFile:           cmd.KTPFile,
		CVFile:            cmd.CVFile,
		NPWPFile:          cmd.NPWPFile,
	})
}

func (s *Service) Get(ctx context.Context, id string) (domain.Profile, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) GetByUserID(ctx context.Context, userID string) (domain.Profile, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *Service) List(ctx context.Context, limit, offset int) (domain.ListResult, error) {
	return s.repo.List(ctx, domain.ListParams{Limit: limit, Offset: offset})
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
