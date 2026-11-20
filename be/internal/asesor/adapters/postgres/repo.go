package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"lsp/internal/asesor/domain"
	"lsp/internal/platform/database"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Upsert(ctx context.Context, profile domain.Profile) (domain.Profile, error) {
	db := database.GetExecutor(ctx, r.pool)
	var out domain.Profile
	err := db.QueryRow(ctx, `
		INSERT INTO asesor_profiles (user_id, tuk_id, province, city, daily_assignment_limit)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id) DO UPDATE SET
			tuk_id                 = EXCLUDED.tuk_id,
			province               = EXCLUDED.province,
			city                   = EXCLUDED.city,
			daily_assignment_limit = EXCLUDED.daily_assignment_limit
		RETURNING user_id, tuk_id, province, city, daily_assignment_limit, created_at, updated_at
	`, profile.UserID, profile.TUKID, profile.Province, profile.City, profile.DailyAssignmentLimit).
		Scan(&out.UserID, &out.TUKID, &out.Province, &out.City, &out.DailyAssignmentLimit, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		return domain.Profile{}, mapError(err)
	}
	return out, nil
}

func (r *Repository) FindByUserID(ctx context.Context, userID string) (domain.Profile, error) {
	db := database.GetExecutor(ctx, r.pool)
	var out domain.Profile
	err := db.QueryRow(ctx, `
		SELECT user_id, tuk_id, province, city, daily_assignment_limit, created_at, updated_at
		FROM asesor_profiles
		WHERE user_id = $1
	`, userID).Scan(&out.UserID, &out.TUKID, &out.Province, &out.City, &out.DailyAssignmentLimit, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		return domain.Profile{}, mapError(err)
	}
	return out, nil
}

func (r *Repository) List(ctx context.Context, params domain.ListParams) (domain.ListResult, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT ap.user_id, u.name, u.email::text, ap.tuk_id, t.name, ap.province, ap.city, ap.daily_assignment_limit,
		       COUNT(*) OVER() AS total_count
		FROM asesor_profiles ap
		JOIN users u ON u.id = ap.user_id
		LEFT JOIN tuks t ON t.id = ap.tuk_id
		WHERE u.deleted_at IS NULL
		ORDER BY u.name
		LIMIT $1 OFFSET $2
	`, params.Limit, params.Offset)
	if err != nil {
		return domain.ListResult{}, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.ListItem, 0)
	total := 0
	for rows.Next() {
		var item domain.ListItem
		if err := rows.Scan(&item.UserID, &item.Name, &item.Email, &item.TUKID, &item.TUKName, &item.Province, &item.City, &item.DailyAssignmentLimit, &total); err != nil {
			return domain.ListResult{}, mapError(err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.ListResult{}, mapError(err)
	}
	if len(items) == 0 && params.Offset > 0 {
		if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM asesor_profiles`).Scan(&total); err != nil {
			return domain.ListResult{}, mapError(err)
		}
	}
	return domain.ListResult{Items: items, Total: total}, nil
}

// EnsureProfile membuat profil asesor default (limit 0) bila belum ada.
// Implementasi port AsesorProfileInitializer milik modul users.
func (r *Repository) EnsureProfile(ctx context.Context, userID string) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO asesor_profiles (user_id, daily_assignment_limit)
		VALUES ($1, 0)
		ON CONFLICT (user_id) DO NOTHING
	`, userID)
	return mapError(err)
}

func mapError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503": // foreign key violation
			if pgErr.ConstraintName == "fk_asesor_profiles_tuk" {
				return domain.ErrTUKNotFound
			}
			return domain.ErrNotAsesor // user tidak ada
		case "23514": // check constraint: limit negatif
			return domain.ErrInvalidLimit
		}
	}
	return err
}
