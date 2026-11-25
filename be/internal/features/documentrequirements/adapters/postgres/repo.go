package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"lsp/internal/features/documentrequirements/domain"
	"lsp/internal/platform/database"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const documentTypeColumns = `id, code, label, allowed_mime_types, created_at, updated_at`

// --- Master jenis dokumen ---

func (r *Repository) CreateDocumentType(ctx context.Context, docType domain.DocumentType) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO document_types (id, code, label, allowed_mime_types)
		VALUES ($1, $2, $3, $4)
	`, docType.ID, docType.Code, docType.Label, docType.AllowedMimeTypes)
	return mapError(err)
}

func (r *Repository) PatchDocumentType(ctx context.Context, patch domain.DocumentTypePatch) (domain.DocumentType, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanDocumentType(db.QueryRow(ctx, `
		UPDATE document_types
		SET code = COALESCE($2::varchar, code),
		    label = COALESCE($3::varchar, label),
		    allowed_mime_types = COALESCE($4::varchar[], allowed_mime_types)
		WHERE id = $1
		RETURNING `+documentTypeColumns,
		patch.ID, patch.Code, patch.Label, patch.AllowedMimeTypes))
}

func (r *Repository) FindDocumentTypeByID(ctx context.Context, id string) (domain.DocumentType, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanDocumentType(db.QueryRow(ctx, `SELECT `+documentTypeColumns+` FROM document_types WHERE id = $1`, id))
}

// AllowedMimeTypes mengembalikan daftar MIME yang diizinkan untuk sebuah jenis
// dokumen. Dipakai modul uploads saat menegakkan aturan unggah (purpose 'document').
func (r *Repository) AllowedMimeTypes(ctx context.Context, id string) ([]string, error) {
	db := database.GetExecutor(ctx, r.pool)
	var types []string
	if err := db.QueryRow(ctx, `SELECT allowed_mime_types FROM document_types WHERE id = $1`, id).Scan(&types); err != nil {
		return nil, mapError(err)
	}
	return types, nil
}

func (r *Repository) ListDocumentTypes(ctx context.Context, params domain.ListParams) (domain.DocumentTypeListResult, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT `+documentTypeColumns+`, COUNT(*) OVER() AS total_count
		FROM document_types
		ORDER BY label, created_at
		LIMIT $1 OFFSET $2
	`, params.Limit, params.Offset)
	if err != nil {
		return domain.DocumentTypeListResult{}, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.DocumentType, 0)
	total := 0
	for rows.Next() {
		item, rowTotal, err := scanDocumentTypeWithTotal(rows)
		if err != nil {
			return domain.DocumentTypeListResult{}, err
		}
		total = rowTotal
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.DocumentTypeListResult{}, mapError(err)
	}
	if len(items) == 0 && params.Offset > 0 {
		if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM document_types`).Scan(&total); err != nil {
			return domain.DocumentTypeListResult{}, mapError(err)
		}
	}
	return domain.DocumentTypeListResult{Items: items, Total: total}, nil
}

func (r *Repository) DeleteDocumentType(ctx context.Context, id string) error {
	db := database.GetExecutor(ctx, r.pool)
	commandTag, err := db.Exec(ctx, `DELETE FROM document_types WHERE id = $1`, id)
	if err != nil {
		return mapError(err)
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// --- Persyaratan dokumen per skema ---

func (r *Repository) ListRequirementsBySchema(ctx context.Context, schemaID string) ([]domain.RequirementDetail, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT r.id, r.certification_schema_id, r.document_type_id, r.is_required, r.created_at, r.updated_at,
		       dt.id, dt.code, dt.label, dt.allowed_mime_types, dt.created_at, dt.updated_at
		FROM schema_document_requirements r
		JOIN document_types dt ON dt.id = r.document_type_id
		WHERE r.certification_schema_id = $1
		ORDER BY dt.label, dt.code
	`, schemaID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.RequirementDetail, 0)
	for rows.Next() {
		item, err := scanRequirementDetail(rows)
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

// DeleteRequirementsNotIn menghapus persyaratan skema yang document_type_id-nya
// tidak termasuk dalam daftar. Daftar kosong menghapus seluruh persyaratan skema.
func (r *Repository) DeleteRequirementsNotIn(ctx context.Context, schemaID string, documentTypeIDs []string) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		DELETE FROM schema_document_requirements
		WHERE certification_schema_id = $1
		  AND NOT (document_type_id = ANY($2::varchar[]))
	`, schemaID, documentTypeIDs)
	return mapError(err)
}

func (r *Repository) UpsertRequirement(ctx context.Context, requirement domain.SchemaDocumentRequirement) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO schema_document_requirements (id, certification_schema_id, document_type_id, is_required)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (certification_schema_id, document_type_id)
		DO UPDATE SET is_required = EXCLUDED.is_required
	`, requirement.ID, requirement.CertificationSchemaID, requirement.DocumentTypeID, requirement.IsRequired)
	return mapError(err)
}

// --- scan helpers ---

func scanDocumentType(row pgx.Row) (domain.DocumentType, error) {
	var dt domain.DocumentType
	if err := row.Scan(&dt.ID, &dt.Code, &dt.Label, &dt.AllowedMimeTypes, &dt.CreatedAt, &dt.UpdatedAt); err != nil {
		return domain.DocumentType{}, mapError(err)
	}
	return dt, nil
}

func scanDocumentTypeWithTotal(rows pgx.Rows) (domain.DocumentType, int, error) {
	var dt domain.DocumentType
	var total int
	if err := rows.Scan(&dt.ID, &dt.Code, &dt.Label, &dt.AllowedMimeTypes, &dt.CreatedAt, &dt.UpdatedAt, &total); err != nil {
		return domain.DocumentType{}, 0, mapError(err)
	}
	return dt, total, nil
}

func scanRequirementDetail(row pgx.Row) (domain.RequirementDetail, error) {
	var d domain.RequirementDetail
	if err := row.Scan(
		&d.ID, &d.CertificationSchemaID, &d.DocumentTypeID, &d.IsRequired, &d.CreatedAt, &d.UpdatedAt,
		&d.DocumentType.ID, &d.DocumentType.Code, &d.DocumentType.Label, &d.DocumentType.AllowedMimeTypes, &d.DocumentType.CreatedAt, &d.DocumentType.UpdatedAt,
	); err != nil {
		return domain.RequirementDetail{}, mapError(err)
	}
	return d, nil
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
		case "23505":
			return domain.ErrConflict
		case "23503":
			// FK violation: skema / jenis dokumen tidak ada.
			return domain.ErrNotFound
		}
	}
	return err
}
