package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"lsp/internal/features/users/domain"
	"lsp/internal/platform/database"
)

// RolesRepository adalah adapter Postgres untuk master data roles.
type RolesRepository struct {
	pool *pgxpool.Pool
}

func NewRolesRepository(pool *pgxpool.Pool) *RolesRepository {
	return &RolesRepository{pool: pool}
}

func (r *RolesRepository) FindByCode(ctx context.Context, code domain.Role) (domain.RoleRecord, error) {
	db := database.GetExecutor(ctx, r.pool)
	var role domain.RoleRecord
	var roleCode string
	err := db.QueryRow(ctx, `
		SELECT id, code::text, name, description
		FROM roles
		WHERE code = $1
	`, string(code)).Scan(&role.ID, &roleCode, &role.Name, &role.Description)
	if err != nil {
		return domain.RoleRecord{}, mapError(err)
	}
	role.Code = domain.Role(roleCode)
	return role, nil
}

func (r *RolesRepository) List(ctx context.Context) ([]domain.RoleRecord, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT id, code::text, name, description
		FROM roles
		ORDER BY code
	`)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	roles := make([]domain.RoleRecord, 0)
	for rows.Next() {
		var role domain.RoleRecord
		var roleCode string
		if err := rows.Scan(&role.ID, &roleCode, &role.Name, &role.Description); err != nil {
			return nil, mapError(err)
		}
		role.Code = domain.Role(roleCode)
		roles = append(roles, role)
	}
	return roles, mapError(rows.Err())
}

// Create dipakai oleh seeder untuk mengisi master data role.
func (r *RolesRepository) Create(ctx context.Context, role domain.RoleRecord) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO roles (id, code, name, description)
		VALUES ($1, $2, $3, $4)
	`, role.ID, string(role.Code), role.Name, role.Description)
	return mapError(err)
}

func (r *RolesRepository) FindByID(ctx context.Context, id string) (domain.RoleRecord, error) {
	db := database.GetExecutor(ctx, r.pool)
	var role domain.RoleRecord
	var roleCode string
	err := db.QueryRow(ctx, `
		SELECT id, code::text, name, description
		FROM roles
		WHERE id = $1
	`, id).Scan(&role.ID, &roleCode, &role.Name, &role.Description)
	if err != nil {
		return domain.RoleRecord{}, mapError(err)
	}
	role.Code = domain.Role(roleCode)
	return role, nil
}

func (r *RolesRepository) Patch(ctx context.Context, patch domain.RolePatch) (domain.RoleRecord, error) {
	db := database.GetExecutor(ctx, r.pool)
	var code *string
	if patch.Code != nil {
		value := string(*patch.Code)
		code = &value
	}
	var role domain.RoleRecord
	var roleCode string
	err := db.QueryRow(ctx, `
		UPDATE roles
		SET code = COALESCE($2::varchar, code),
		    name = COALESCE($3::varchar, name),
		    description = COALESCE($4::text, description)
		WHERE id = $1
		RETURNING id, code::text, name, description
	`, patch.ID, code, patch.Name, patch.Description).Scan(&role.ID, &roleCode, &role.Name, &role.Description)
	if err != nil {
		return domain.RoleRecord{}, mapError(err)
	}
	role.Code = domain.Role(roleCode)
	return role, nil
}

func (r *RolesRepository) Delete(ctx context.Context, id string) error {
	db := database.GetExecutor(ctx, r.pool)
	commandTag, err := db.Exec(ctx, `DELETE FROM roles WHERE id = $1`, id)
	if err != nil {
		return mapError(err)
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}
