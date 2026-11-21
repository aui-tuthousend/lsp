package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"lsp/internal/features/tuk/domain"
	"lsp/internal/platform/database"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, tuk domain.TUK) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO tuks (id, name, logo, postal_code, phone, email, province, city, address, maps_url, letterhead)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, tuk.ID, tuk.Name, tuk.Logo, tuk.PostalCode, tuk.Phone, tuk.Email, tuk.Province, tuk.City, tuk.Address, tuk.MapsURL, tuk.Letterhead)
	return mapError(err)
}

func (r *Repository) Patch(ctx context.Context, patch domain.Patch) (domain.TUK, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanTUK(db.QueryRow(ctx, `
		UPDATE tuks
		SET name = COALESCE($2::varchar, name),
		    logo = COALESCE($3::text, logo),
		    postal_code = COALESCE($4::varchar, postal_code),
		    phone = COALESCE($5::varchar, phone),
		    email = COALESCE($6::citext, email),
		    province = COALESCE($7::varchar, province),
		    city = COALESCE($8::varchar, city),
		    address = COALESCE($9::text, address),
		    maps_url = COALESCE($10::text, maps_url),
		    letterhead = COALESCE($11::text, letterhead)
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, name, logo, postal_code, phone, email::text, province, city, address, maps_url, letterhead, created_at, updated_at, deleted_at
	`, patch.ID, patch.Name, patch.Logo, patch.PostalCode, patch.Phone, patch.Email, patch.Province, patch.City, patch.Address, patch.MapsURL, patch.Letterhead))
}

func (r *Repository) FindByID(ctx context.Context, id string) (domain.TUK, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanTUK(db.QueryRow(ctx, `
		SELECT id, name, logo, postal_code, phone, email::text, province, city, address, maps_url, letterhead, created_at, updated_at, deleted_at
		FROM tuks
		WHERE id = $1 AND deleted_at IS NULL
	`, id))
}

func (r *Repository) List(ctx context.Context, params domain.ListParams) (domain.ListResult, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT id, name, logo, postal_code, phone, email::text, province, city, address, maps_url, letterhead, created_at, updated_at, deleted_at,
		       COUNT(*) OVER() AS total_count
		FROM tuks
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, params.Limit, params.Offset)
	if err != nil {
		return domain.ListResult{}, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.TUK, 0)
	total := 0
	for rows.Next() {
		item, rowTotal, err := scanTUKWithTotal(rows)
		if err != nil {
			return domain.ListResult{}, err
		}
		total = rowTotal
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.ListResult{}, mapError(err)
	}
	if len(items) == 0 && params.Offset > 0 {
		if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM tuks WHERE deleted_at IS NULL`).Scan(&total); err != nil {
			return domain.ListResult{}, mapError(err)
		}
	}
	return domain.ListResult{Items: items, Total: total}, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	db := database.GetExecutor(ctx, r.pool)
	commandTag, err := db.Exec(ctx, `
		UPDATE tuks
		SET deleted_at = now()
		WHERE id = $1 AND deleted_at IS NULL
	`, id)
	if err != nil {
		return mapError(err)
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func scanTUK(row pgx.Row) (domain.TUK, error) {
	var tuk domain.TUK
	if err := row.Scan(&tuk.ID, &tuk.Name, &tuk.Logo, &tuk.PostalCode, &tuk.Phone, &tuk.Email, &tuk.Province, &tuk.City, &tuk.Address, &tuk.MapsURL, &tuk.Letterhead, &tuk.CreatedAt, &tuk.UpdatedAt, &tuk.DeletedAt); err != nil {
		return domain.TUK{}, mapError(err)
	}
	return tuk, nil
}

func scanTUKWithTotal(rows pgx.Rows) (domain.TUK, int, error) {
	var tuk domain.TUK
	var total int
	if err := rows.Scan(&tuk.ID, &tuk.Name, &tuk.Logo, &tuk.PostalCode, &tuk.Phone, &tuk.Email, &tuk.Province, &tuk.City, &tuk.Address, &tuk.MapsURL, &tuk.Letterhead, &tuk.CreatedAt, &tuk.UpdatedAt, &tuk.DeletedAt, &total); err != nil {
		return domain.TUK{}, 0, mapError(err)
	}
	return tuk, total, nil
}

func mapError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.ErrNotFound
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return domain.ErrConflict
	}
	return err
}
