package app

import (
	"context"
	"time"

	"lsp/internal/features/asesi/domain"
	"lsp/internal/platform/id"
	"lsp/internal/platform/storage"
)

type Service struct {
	repo    domain.Repository
	checker domain.AsesiChecker
	files   storage.ObjectStorage
}

func NewService(repo domain.Repository, checker domain.AsesiChecker, files storage.ObjectStorage) *Service {
	return &Service{repo: repo, checker: checker, files: files}
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
	previous, err := s.repo.FindByID(ctx, cmd.ID)
	if err != nil {
		return domain.Profile{}, err
	}
	updated, err := s.repo.Patch(ctx, domain.Patch{
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
	if err != nil {
		return domain.Profile{}, err
	}
	s.removeReplacedFiles(ctx, previous, updated)
	return updated, nil
}

// removeReplacedFiles menghapus objek lama di object storage saat sebuah berkas
// diganti melalui Patch. Best-effort: kegagalan hapus tidak menggagalkan
// permintaan karena nilai baru sudah tersimpan di database.
func (s *Service) removeReplacedFiles(ctx context.Context, previous, updated domain.Profile) {
	if s.files == nil {
		return
	}
	pairs := [][2]*string{
		{previous.Photo, updated.Photo},
		{previous.KTPFile, updated.KTPFile},
		{previous.CVFile, updated.CVFile},
		{previous.NPWPFile, updated.NPWPFile},
	}
	for _, pair := range pairs {
		oldKey, newKey := pair[0], pair[1]
		if oldKey == nil || *oldKey == "" {
			continue
		}
		if newKey != nil && *newKey == *oldKey {
			continue
		}
		_ = s.files.Delete(ctx, *oldKey)
	}
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

// --- Data pekerjaan asesi (asesi_occupations) ---

type CreateOccupationCommand struct {
	AsesiProfileID    string
	CompanyName       string
	Position          string
	TenureMonths      *int
	SubordinatesCount *int
	Address           string
	Province          string
	City              string
	District          string
	PostalCode        *string
	OfficeEmail       *string
	OfficePhone       *string
}

func (s *Service) CreateOccupation(ctx context.Context, cmd CreateOccupationCommand) (domain.Occupation, error) {
	occupation := domain.Occupation{
		ID:                id.New(),
		AsesiProfileID:    cmd.AsesiProfileID,
		CompanyName:       cmd.CompanyName,
		Position:          cmd.Position,
		TenureMonths:      cmd.TenureMonths,
		SubordinatesCount: cmd.SubordinatesCount,
		Address:           cmd.Address,
		Province:          cmd.Province,
		City:              cmd.City,
		District:          cmd.District,
		PostalCode:        cmd.PostalCode,
		OfficeEmail:       cmd.OfficeEmail,
		OfficePhone:       cmd.OfficePhone,
	}
	if err := s.repo.CreateOccupation(ctx, occupation); err != nil {
		return domain.Occupation{}, err
	}
	return s.repo.FindOccupationByID(ctx, occupation.ID)
}

// UpsertOccupationCommand membuat/memperbarui data pekerjaan milik user (asesi) sendiri.
type UpsertOccupationCommand struct {
	UserID            string
	CompanyName       string
	Position          string
	TenureMonths      *int
	SubordinatesCount *int
	Address           string
	Province          string
	City              string
	District          string
	PostalCode        *string
	OfficeEmail       *string
	OfficePhone       *string
}

func (s *Service) UpsertOccupationByUser(ctx context.Context, cmd UpsertOccupationCommand) (domain.Occupation, error) {
	return s.repo.UpsertOccupationByUser(ctx, cmd.UserID, domain.Occupation{
		ID:                id.New(),
		CompanyName:       cmd.CompanyName,
		Position:          cmd.Position,
		TenureMonths:      cmd.TenureMonths,
		SubordinatesCount: cmd.SubordinatesCount,
		Address:           cmd.Address,
		Province:          cmd.Province,
		City:              cmd.City,
		District:          cmd.District,
		PostalCode:        cmd.PostalCode,
		OfficeEmail:       cmd.OfficeEmail,
		OfficePhone:       cmd.OfficePhone,
	})
}

type PatchOccupationCommand struct {
	ID                string
	CompanyName       *string
	Position          *string
	TenureMonths      *int
	SubordinatesCount *int
	Address           *string
	Province          *string
	City              *string
	District          *string
	PostalCode        *string
	OfficeEmail       *string
	OfficePhone       *string
}

func (s *Service) PatchOccupation(ctx context.Context, cmd PatchOccupationCommand) (domain.Occupation, error) {
	return s.repo.PatchOccupation(ctx, domain.OccupationPatch{
		ID:                cmd.ID,
		CompanyName:       cmd.CompanyName,
		Position:          cmd.Position,
		TenureMonths:      cmd.TenureMonths,
		SubordinatesCount: cmd.SubordinatesCount,
		Address:           cmd.Address,
		Province:          cmd.Province,
		City:              cmd.City,
		District:          cmd.District,
		PostalCode:        cmd.PostalCode,
		OfficeEmail:       cmd.OfficeEmail,
		OfficePhone:       cmd.OfficePhone,
	})
}

func (s *Service) GetOccupation(ctx context.Context, id string) (domain.Occupation, error) {
	return s.repo.FindOccupationByID(ctx, id)
}

func (s *Service) GetOccupationByUserID(ctx context.Context, userID string) (domain.Occupation, error) {
	return s.repo.FindOccupationByUserID(ctx, userID)
}

func (s *Service) ListOccupations(ctx context.Context, limit, offset int) (domain.OccupationListResult, error) {
	return s.repo.ListOccupations(ctx, domain.ListParams{Limit: limit, Offset: offset})
}

func (s *Service) DeleteOccupation(ctx context.Context, id string) error {
	return s.repo.DeleteOccupation(ctx, id)
}
