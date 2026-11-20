package app

import (
	"context"

	"lsp/internal/asesor/domain"
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

func (s *Service) Get(ctx context.Context, userID string) (domain.Profile, error) {
	return s.repo.FindByUserID(ctx, userID)
}

func (s *Service) List(ctx context.Context, limit, offset int) (domain.ListResult, error) {
	return s.repo.List(ctx, domain.ListParams{Limit: limit, Offset: offset})
}
