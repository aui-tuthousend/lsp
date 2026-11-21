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

func (r *Repository) PatchProfile(ctx context.Context, patch domain.ProfilePatch) (domain.Profile, error) {
	db := database.GetExecutor(ctx, r.pool)
	var out domain.Profile
	err := db.QueryRow(ctx, `
		UPDATE asesor_profiles
		SET tuk_id = COALESCE($2::varchar, tuk_id),
		    province = COALESCE($3::varchar, province),
		    city = COALESCE($4::varchar, city),
		    daily_assignment_limit = COALESCE($5::integer, daily_assignment_limit)
		WHERE user_id = $1
		RETURNING user_id, tuk_id, province, city, daily_assignment_limit, created_at, updated_at
	`, patch.UserID, patch.TUKID, patch.Province, patch.City, patch.DailyAssignmentLimit).
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

func (r *Repository) ListCertifications(ctx context.Context, userID string) ([]domain.Certification, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT id, user_id, name, issuer, number, issued_at, expired_at, created_at, updated_at
		FROM asesor_certifications
		WHERE user_id = $1
		ORDER BY issued_at DESC NULLS LAST, created_at DESC
	`, userID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.Certification, 0)
	for rows.Next() {
		item, err := scanCertification(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, mapError(err)
	}
	return items, nil
}

func (r *Repository) FindCertificationByID(ctx context.Context, userID, certID string) (domain.Certification, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanCertification(db.QueryRow(ctx, `
		SELECT id, user_id, name, issuer, number, issued_at, expired_at, created_at, updated_at
		FROM asesor_certifications
		WHERE id = $1 AND user_id = $2
	`, certID, userID))
}

func (r *Repository) CreateCertification(ctx context.Context, cert domain.Certification) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO asesor_certifications (id, user_id, name, issuer, number, issued_at, expired_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, cert.ID, cert.UserID, cert.Name, cert.Issuer, cert.Number, cert.IssuedAt, cert.ExpiredAt)
	return mapError(err)
}

func (r *Repository) PatchCertification(ctx context.Context, patch domain.CertificationPatch) (domain.Certification, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanCertification(db.QueryRow(ctx, `
		UPDATE asesor_certifications
		SET name = COALESCE($2::varchar, name),
		    issuer = COALESCE($3::varchar, issuer),
		    number = COALESCE($4::varchar, number),
		    issued_at = COALESCE($5::date, issued_at),
		    expired_at = COALESCE($6::date, expired_at)
		WHERE id = $1
		RETURNING id, user_id, name, issuer, number, issued_at, expired_at, created_at, updated_at
	`, patch.ID, patch.Name, patch.Issuer, patch.Number, patch.IssuedAt, patch.ExpiredAt))
}

func (r *Repository) DeleteCertification(ctx context.Context, userID, certID string) error {
	db := database.GetExecutor(ctx, r.pool)
	commandTag, err := db.Exec(ctx, `
		DELETE FROM asesor_certifications
		WHERE id = $1 AND user_id = $2
	`, certID, userID)
	if err != nil {
		return mapError(err)
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrCertificationNotFound
	}
	return nil
}

func scanCertification(row pgx.Row) (domain.Certification, error) {
	var c domain.Certification
	if err := row.Scan(&c.ID, &c.UserID, &c.Name, &c.Issuer, &c.Number, &c.IssuedAt, &c.ExpiredAt, &c.CreatedAt, &c.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Certification{}, domain.ErrCertificationNotFound
		}
		return domain.Certification{}, mapError(err)
	}
	return c, nil
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
			switch pgErr.ConstraintName {
			case "fk_asesor_profiles_tuk":
				return domain.ErrTUKNotFound
			case "fk_asesor_certifications_profile":
				return domain.ErrNotFound // profil asesor tidak ada
			default:
				return domain.ErrNotAsesor // user tidak ada
			}
		case "23514": // check constraint
			if pgErr.ConstraintName == "ck_asesor_certifications_period" {
				return domain.ErrInvalidCertificationPeriod
			}
			return domain.ErrInvalidLimit
		}
	}
	return err
}
