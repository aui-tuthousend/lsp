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
		INSERT INTO users (id, role_id, isvalid, name, email, phone, nomor_registrasi, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, user.ID, user.RoleID, user.IsValid, user.Name, user.Email, user.Phone, user.NomorRegistrasi, user.PasswordHash)
	return mapError(err)
}

func (r *Repository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanUser(db.QueryRow(ctx, `
		SELECT u.id, u.role_id, r.code::text, u.isvalid, u.name, u.email::text, u.phone, u.nomor_registrasi, u.password_hash,
		       u.email_verified_at, u.remember_token, u.created_at, u.updated_at, u.deleted_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.email = $1 AND u.deleted_at IS NULL
	`, email))
}

func (r *Repository) FindByID(ctx context.Context, id string) (domain.User, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanUser(db.QueryRow(ctx, `
		SELECT u.id, u.role_id, r.code::text, u.isvalid, u.name, u.email::text, u.phone, u.nomor_registrasi, u.password_hash,
		       u.email_verified_at, u.remember_token, u.created_at, u.updated_at, u.deleted_at
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.id = $1 AND u.deleted_at IS NULL
	`, id))
}

func (r *Repository) List(ctx context.Context, params domain.ListParams) (domain.ListResult, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT u.id, u.role_id, r.name::text, u.isvalid, u.name, u.email::text, u.phone, u.nomor_registrasi, u.password_hash,
		       u.email_verified_at, u.remember_token, u.created_at, u.updated_at, u.deleted_at,
		       COUNT(*) OVER() AS total_count
		FROM users u
		JOIN roles r ON r.id = u.role_id
		WHERE u.deleted_at IS NULL
		ORDER BY u.created_at DESC
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

// IsAsesor mengembalikan true bila user dengan id tersebut ber-role asesor.
// Dipakai modul asesor sebagai implementasi port AsesorChecker.
func (r *Repository) IsAsesor(ctx context.Context, userID string) (bool, error) {
	db := database.GetExecutor(ctx, r.pool)
	var exists bool
	if err := db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM users u
			JOIN roles r ON r.id = u.role_id
			WHERE u.id = $1 AND r.code = 'asesor' AND u.deleted_at IS NULL
		)
	`, userID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func scanUser(row pgx.Row) (domain.User, error) {
	var user domain.User
	var roleCode string
	if err := row.Scan(&user.ID, &user.RoleID, &roleCode, &user.IsValid, &user.Name, &user.Email, &user.Phone, &user.NomorRegistrasi, &user.PasswordHash, &user.EmailVerifiedAt, &user.RememberToken, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
		return domain.User{}, mapError(err)
	}
	user.Role = domain.Role(roleCode)
	return user, nil
}

func scanUserWithTotal(rows pgx.Rows) (domain.User, int, error) {
	var user domain.User
	var roleCode string
	var total int
	if err := rows.Scan(&user.ID, &user.RoleID, &roleCode, &user.IsValid, &user.Name, &user.Email, &user.Phone, &user.NomorRegistrasi, &user.PasswordHash, &user.EmailVerifiedAt, &user.RememberToken, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt, &total); err != nil {
		return domain.User{}, 0, mapError(err)
	}
	user.Role = domain.Role(roleCode)
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
