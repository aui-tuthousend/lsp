package app

import (
	"context"
	"time"

	"lsp/internal/features/asesor/domain"
	"lsp/internal/platform/id"
)

type Service struct {
	repo    domain.Repository
	checker domain.AsesorChecker
}

func NewService(repo domain.Repository, checker domain.AsesorChecker) *Service {
	return &Service{repo: repo, checker: checker}
}

type UpsertCommand struct {
	UserID               string
	TUKID                *string
	Province             *string
	City                 *string
	DailyAssignmentLimit int
}

type PatchProfileCommand struct {
	UserID               string
	TUKID                *string
	Province             *string
	City                 *string
	DailyAssignmentLimit *int
}

// Upsert menetapkan profil asesor: afiliasi TUK, domisili, dan batas
// penugasan harian.
func (s *Service) Upsert(ctx context.Context, cmd UpsertCommand) (domain.Profile, error) {
	if cmd.DailyAssignmentLimit < 0 {
		return domain.Profile{}, domain.ErrInvalidLimit
	}
	isAsesor, err := s.checker.IsAsesor(ctx, cmd.UserID)
	if err != nil {
		return domain.Profile{}, err
	}
	if !isAsesor {
		return domain.Profile{}, domain.ErrNotAsesor
	}
	return s.repo.Upsert(ctx, domain.Profile{
		UserID:               cmd.UserID,
		TUKID:                cmd.TUKID,
		Province:             cmd.Province,
		City:                 cmd.City,
		DailyAssignmentLimit: cmd.DailyAssignmentLimit,
	})
}

func (s *Service) PatchProfile(ctx context.Context, cmd PatchProfileCommand) (domain.Profile, error) {
	if cmd.DailyAssignmentLimit != nil && *cmd.DailyAssignmentLimit < 0 {
		return domain.Profile{}, domain.ErrInvalidLimit
	}
	return s.repo.PatchProfile(ctx, domain.ProfilePatch{
		UserID:               cmd.UserID,
		TUKID:                cmd.TUKID,
		Province:             cmd.Province,
		City:                 cmd.City,
		DailyAssignmentLimit: cmd.DailyAssignmentLimit,
	})
}

func (s *Service) Get(ctx context.Context, userID string) (domain.Profile, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *Service) List(ctx context.Context, limit, offset int) (domain.ListResult, error) {
	return s.repo.List(ctx, domain.ListParams{Limit: limit, Offset: offset})
}

type CreateCertificationCommand struct {
	UserID    string
	Name      string
	Issuer    string
	Number    *string
	IssuedAt  *time.Time
	ExpiredAt *time.Time
}

// CreateCertification menambah satu sertifikasi untuk seorang asesor.
func (s *Service) CreateCertification(ctx context.Context, cmd CreateCertificationCommand) (domain.Certification, error) {
	isAsesor, err := s.checker.IsAsesor(ctx, cmd.UserID)
	if err != nil {
		return domain.Certification{}, err
	}
	if !isAsesor {
		return domain.Certification{}, domain.ErrNotAsesor
	}
	cert := domain.Certification{
		ID:        id.New(),
		UserID:    cmd.UserID,
		Name:      cmd.Name,
		Issuer:    cmd.Issuer,
		Number:    cmd.Number,
		IssuedAt:  cmd.IssuedAt,
		ExpiredAt: cmd.ExpiredAt,
	}
	if err := s.repo.CreateCertification(ctx, cert); err != nil {
		return domain.Certification{}, err
	}
	return s.repo.FindCertificationByID(ctx, cmd.UserID, cert.ID)
}

type PatchCertificationCommand struct {
	ID        string
	Name      *string
	Issuer    *string
	Number    *string
	IssuedAt  *time.Time
	ExpiredAt *time.Time
}

func (s *Service) PatchCertification(ctx context.Context, cmd PatchCertificationCommand) (domain.Certification, error) {
	return s.repo.PatchCertification(ctx, domain.CertificationPatch{
		ID:        cmd.ID,
		Name:      cmd.Name,
		Issuer:    cmd.Issuer,
		Number:    cmd.Number,
		IssuedAt:  cmd.IssuedAt,
		ExpiredAt: cmd.ExpiredAt,
	})
}

func (s *Service) ListCertifications(ctx context.Context, userID string) ([]domain.Certification, error) {
	return s.repo.ListCertifications(ctx, userID)
}

func (s *Service) DeleteCertification(ctx context.Context, userID, certID string) error {
	return s.repo.DeleteCertification(ctx, userID, certID)
}
