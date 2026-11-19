package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"lsp/internal/platform/database"
	"lsp/internal/users/domain"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, user domain.User) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO users (id, role, isvalid, name, email, phone, nomor_registrasi, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, user.ID, string(user.Role), user.IsValid, user.Name, user.Email, user.Phone, user.NomorRegistrasi, user.PasswordHash)
	return mapError(err)
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanUser(db.QueryRow(ctx, `
		SELECT id, role::text, isvalid, name, email::text, phone, nomor_registrasi, password_hash,
		       email_verified_at, remember_token, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`, email))
}

func (r *Repository) FindByID(ctx context.Context, id string) (domain.User, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanUser(db.QueryRow(ctx, `
		SELECT id, role::text, isvalid, name, email::text, phone, nomor_registrasi, password_hash,
		       email_verified_at, remember_token, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`, id))
}

func (r *Repository) List(ctx context.Context, params domain.ListParams) (domain.ListResult, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT id, role::text, isvalid, name, email::text, phone, nomor_registrasi, password_hash,
		       email_verified_at, remember_token, created_at, updated_at, deleted_at,
		       COUNT(*) OVER() AS total_count
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, params.Limit, params.Offset)
	if err != nil {
		return domain.ListResult{}, mapError(err)
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	total := 0
	for rows.Next() {
		user, rowTotal, err := scanUserWithTotal(rows)
		if err != nil {
			return domain.ListResult{}, err
		}
		total = rowTotal
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return domain.ListResult{}, mapError(err)
	}
	if len(users) == 0 && params.Offset > 0 {
		if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`).Scan(&total); err != nil {
			return domain.ListResult{}, mapError(err)
		}
	}
	return domain.ListResult{Users: users, Total: total}, nil
}

func scanUser(row pgx.Row) (domain.User, error) {
	var user domain.User
	var role string
	if err := row.Scan(&user.ID, &role, &user.IsValid, &user.Name, &user.Email, &user.Phone, &user.NomorRegistrasi, &user.PasswordHash, &user.EmailVerifiedAt, &user.RememberToken, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
		return domain.User{}, mapError(err)
	}
	user.Role = domain.Role(role)
	return user, nil
}

func scanUserWithTotal(rows pgx.Rows) (domain.User, int, error) {
	var user domain.User
	var role string
	var total int
	if err := rows.Scan(&user.ID, &role, &user.IsValid, &user.Name, &user.Email, &user.Phone, &user.NomorRegistrasi, &user.PasswordHash, &user.EmailVerifiedAt, &user.RememberToken, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &total); err != nil {
		return domain.User{}, 0, mapError(err)
	}
	user.Role = domain.Role(role)
	return user, total, nil
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
		case "23505", "23503":
			return domain.ErrConflict
		}
	}
	return err
}
