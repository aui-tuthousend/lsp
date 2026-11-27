package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"lsp/internal/features/apl02/domain"
	"lsp/internal/platform/database"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const assessmentColumns = `id, application_id, code, status::text, asesor_note, asesor_decision::text,
	submitted_at, reviewed_by, reviewed_at, zoom_link, zoom_date, zoom_time, created_at, updated_at`

// assessmentColumnsA sama seperti assessmentColumns namun ber-alias `a.`,
// dipakai saat query JOIN agar kolom tidak ambigu.
const assessmentColumnsA = `a.id, a.application_id, a.code, a.status::text, a.asesor_note, a.asesor_decision::text,
	a.submitted_at, a.reviewed_by, a.reviewed_at, a.zoom_link, a.zoom_date, a.zoom_time, a.created_at, a.updated_at`

const unitColumns = `u.id, u.assessment_id, u.competency_unit_id, u.is_competent,
	u.created_at, u.updated_at, cu.code, cu.title`

const unitFileColumns = `id, assessment_unit_id, file_url, created_at`

const signatureColumns = `id, application_id, entity, role, token_hash, status,
	expires_at, used_at, revoked_at, created_at, updated_at`

// --- Konteks permohonan APL-01 ---

func (r *Repository) FindApplicationContext(ctx context.Context, applicationID string) (domain.ApplicationContext, error) {
	db := database.GetExecutor(ctx, r.pool)
	var app domain.ApplicationContext
	err := db.QueryRow(ctx, `
		SELECT a.user_id, a.status::text, d.certification_schema_id, a.asesor_id,
		       cs.name, COALESCE(cs.schema_code, cs.code), d.purpose::text, d.created_at
		FROM certification_applications a
		JOIN certification_application_data d ON d.application_id = a.id
		JOIN certification_schemas cs ON cs.id = d.certification_schema_id
		WHERE a.id = $1 AND a.deleted_at IS NULL
	`, applicationID).Scan(&app.UserID, &app.Status, &app.SchemaID, &app.AsesorID,
		&app.SchemaName, &app.SchemaCode, &app.Purpose, &app.DataCreatedAt)
	if err != nil {
		return domain.ApplicationContext{}, mapError(err)
	}
	return app, nil
}

// --- Header asesmen mandiri ---

func (r *Repository) CreateAssessment(ctx context.Context, assessment domain.Assessment) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO independent_assessments (id, application_id, code, status)
		VALUES ($1, $2, $3, $4)
	`, assessment.ID, assessment.ApplicationID, assessment.Code, assessment.Status)
	return mapError(err)
}

func (r *Repository) FindAssessmentByID(ctx context.Context, id string) (domain.Assessment, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanAssessment(db.QueryRow(ctx, `
		SELECT `+assessmentColumns+`
		FROM independent_assessments
		WHERE id = $1
	`, id))
}

func (r *Repository) FindAssessmentByApplicationID(ctx context.Context, applicationID string) (domain.Assessment, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanAssessment(db.QueryRow(ctx, `
		SELECT `+assessmentColumns+`
		FROM independent_assessments
		WHERE application_id = $1
	`, applicationID))
}

func (r *Repository) ListAssessments(ctx context.Context, params domain.ListParams) (domain.ListResult, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT `+assessmentColumnsA+`, COALESCE(cs.name, ''), COALESCE(cs.schema_code, cs.code, ''), COUNT(*) OVER() AS total_count
		FROM independent_assessments a
		JOIN certification_applications c ON c.id = a.application_id
		LEFT JOIN certification_application_data d ON d.application_id = a.application_id
		LEFT JOIN certification_schemas cs ON cs.id = d.certification_schema_id
		WHERE ($1::text = '' OR c.user_id = $1::text)
		ORDER BY a.created_at DESC
		LIMIT $2 OFFSET $3
	`, params.UserID, params.Limit, params.Offset)
	if err != nil {
		return domain.ListResult{}, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.Assessment, 0)
	total := 0
	for rows.Next() {
		item, rowTotal, err := scanAssessmentWithTotal(rows)
		if err != nil {
			return domain.ListResult{}, err
		}
		total = rowTotal
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.ListResult{}, mapError(err)
	}
	return domain.ListResult{Items: items, Total: total}, nil
}

func (r *Repository) MarkSubmitted(ctx context.Context, id string, submittedAt time.Time) (domain.Assessment, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanAssessment(db.QueryRow(ctx, `
		UPDATE independent_assessments
		SET status = 'SUBMITTED', submitted_at = $2
		WHERE id = $1 AND status = 'DRAFT'
		RETURNING `+assessmentColumns,
		id, submittedAt))
}

func (r *Repository) MarkReviewed(ctx context.Context, review domain.Review) (domain.Assessment, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanAssessment(db.QueryRow(ctx, `
		UPDATE independent_assessments
		SET status = 'REVIEWED',
		    asesor_decision = $2,
		    asesor_note = $3,
		    reviewed_by = $4,
		    reviewed_at = $5,
		    zoom_link = $6,
		    zoom_date = $7,
		    zoom_time = $8
		WHERE id = $1 AND status = 'SUBMITTED'
		RETURNING `+assessmentColumns,
		review.AssessmentID, review.Decision, review.Note, review.ReviewedBy, review.ReviewedAt,
		review.ZoomLink, review.ZoomDate, review.ZoomTime))
}

// --- Jawaban per unit kompetensi ---

func (r *Repository) ListUnitAnswers(ctx context.Context, assessmentID string) ([]domain.UnitAnswer, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT `+unitColumns+`
		FROM independent_assessment_units u
		JOIN competency_units cu ON cu.id = u.competency_unit_id
		WHERE u.assessment_id = $1
		ORDER BY cu.sort_order, cu.code
	`, assessmentID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.UnitAnswer, 0)
	index := make(map[string]int)
	for rows.Next() {
		var u domain.UnitAnswer
		if err := rows.Scan(&u.ID, &u.AssessmentID, &u.CompetencyUnitID, &u.IsCompetent,
			&u.CreatedAt, &u.UpdatedAt, &u.UnitCode, &u.UnitTitle); err != nil {
			return nil, mapError(err)
		}
		u.Files = make([]domain.UnitFile, 0)
		index[u.ID] = len(items)
		items = append(items, u)
	}
	if err := rows.Err(); err != nil {
		return nil, mapError(err)
	}

	fileRows, err := db.Query(ctx, `
		SELECT f.id, f.assessment_unit_id, f.file_url, f.created_at
		FROM independent_assessment_unit_files f
		JOIN independent_assessment_units u ON u.id = f.assessment_unit_id
		WHERE u.assessment_id = $1
		ORDER BY f.created_at
	`, assessmentID)
	if err != nil {
		return nil, mapError(err)
	}
	defer fileRows.Close()
	for fileRows.Next() {
		var f domain.UnitFile
		if err := fileRows.Scan(&f.ID, &f.AssessmentUnitID, &f.FileURL, &f.CreatedAt); err != nil {
			return nil, mapError(err)
		}
		if pos, ok := index[f.AssessmentUnitID]; ok {
			items[pos].Files = append(items[pos].Files, f)
		}
	}
	if err := fileRows.Err(); err != nil {
		return nil, mapError(err)
	}
	return items, nil
}

func (r *Repository) CreateUnitAnswers(ctx context.Context, answers []domain.UnitAnswer) error {
	db := database.GetExecutor(ctx, r.pool)
	for _, answer := range answers {
		if _, err := db.Exec(ctx, `
			INSERT INTO independent_assessment_units (id, assessment_id, competency_unit_id)
			VALUES ($1, $2, $3)
			ON CONFLICT (assessment_id, competency_unit_id) DO NOTHING
		`, answer.ID, answer.AssessmentID, answer.CompetencyUnitID); err != nil {
			return mapError(err)
		}
	}
	return nil
}

// ReplaceUnitAnswers menyimpan jawaban (set-semantik): memperbarui is_competent
// lalu mengganti seluruh berkas bukti untuk tiap unit yang dikirim.
func (r *Repository) ReplaceUnitAnswers(ctx context.Context, assessmentID string, answers []domain.UnitAnswer) error {
	db := database.GetExecutor(ctx, r.pool)
	for _, answer := range answers {
		var unitID string
		if err := db.QueryRow(ctx, `
			SELECT id FROM independent_assessment_units
			WHERE assessment_id = $1 AND competency_unit_id = $2
		`, assessmentID, answer.CompetencyUnitID).Scan(&unitID); err != nil {
			return mapError(err)
		}
		if _, err := db.Exec(ctx, `
			UPDATE independent_assessment_units SET is_competent = $2 WHERE id = $1
		`, unitID, answer.IsCompetent); err != nil {
			return mapError(err)
		}
		if _, err := db.Exec(ctx, `
			DELETE FROM independent_assessment_unit_files WHERE assessment_unit_id = $1
		`, unitID); err != nil {
			return mapError(err)
		}
		for _, file := range answer.Files {
			if _, err := db.Exec(ctx, `
				INSERT INTO independent_assessment_unit_files (id, assessment_unit_id, file_url)
				VALUES ($1, $2, $3)
			`, file.ID, unitID, file.FileURL); err != nil {
				return mapError(err)
			}
		}
	}
	return nil
}

func (r *Repository) CountUnansweredUnits(ctx context.Context, assessmentID string) (int, error) {
	db := database.GetExecutor(ctx, r.pool)
	var count int
	err := db.QueryRow(ctx, `
		SELECT COUNT(*) FROM independent_assessment_units
		WHERE assessment_id = $1 AND is_competent IS NULL
	`, assessmentID).Scan(&count)
	if err != nil {
		return 0, mapError(err)
	}
	return count, nil
}

func (r *Repository) CountCompetentWithoutFiles(ctx context.Context, assessmentID string) (int, error) {
	db := database.GetExecutor(ctx, r.pool)
	var count int
	err := db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM independent_assessment_units u
		WHERE u.assessment_id = $1
		  AND u.is_competent = true
		  AND NOT EXISTS (
		      SELECT 1 FROM independent_assessment_unit_files f
		      WHERE f.assessment_unit_id = u.id
		  )
	`, assessmentID).Scan(&count)
	if err != nil {
		return 0, mapError(err)
	}
	return count, nil
}

func (r *Repository) ListSchemaUnitIDs(ctx context.Context, schemaID string) ([]string, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT id FROM competency_units
		WHERE certification_schema_id = $1
		ORDER BY sort_order, code
	`, schemaID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	ids := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, mapError(err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, mapError(err)
	}
	return ids, nil
}

// --- Token verifikasi QR (form_signatures) ---

func (r *Repository) CreateSignature(ctx context.Context, signature domain.Signature) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO form_signatures (id, application_id, entity, role, token_hash, status, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, signature.ID, signature.ApplicationID, signature.Entity, signature.Role,
		signature.TokenHash, signature.Status, signature.ExpiresAt)
	return mapError(err)
}

func (r *Repository) RevokeActiveSignatures(ctx context.Context, applicationID, entity, role string, revokedAt time.Time) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		UPDATE form_signatures
		SET status = 'dibatalkan', revoked_at = $4
		WHERE application_id = $1 AND entity = $2 AND role = $3 AND status = 'menunggu'
	`, applicationID, entity, role, revokedAt)
	return mapError(err)
}

// HasPendingSignature melaporkan apakah asesmen mandiri masih punya token QR
// asesi yang menunggu (belum dipakai & belum kedaluwarsa). Selama token ini
// ada, formulir APL-02 terkunci dari perubahan.
func (r *Repository) HasPendingSignature(ctx context.Context, applicationID string) (bool, error) {
	db := database.GetExecutor(ctx, r.pool)
	var exists bool
	err := db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM form_signatures
			WHERE application_id = $1 AND entity = 'asesmen_mandiri' AND role = 'asesi'
			  AND status = 'menunggu' AND expires_at > now()
		)
	`, applicationID).Scan(&exists)
	if err != nil {
		return false, mapError(err)
	}
	return exists, nil
}

func (r *Repository) FindSignatureByTokenHash(ctx context.Context, tokenHash string) (domain.Signature, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanSignature(db.QueryRow(ctx, `
		SELECT `+signatureColumns+`
		FROM form_signatures
		WHERE token_hash = $1
	`, tokenHash))
}

func (r *Repository) MarkSignatureUsed(ctx context.Context, id string, usedAt time.Time, ip, userAgent *string) (domain.Signature, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanSignature(db.QueryRow(ctx, `
		UPDATE form_signatures
		SET status = 'dipakai', used_at = $2, used_ip = $3::inet, used_user_agent = $4
		WHERE id = $1 AND status = 'menunggu'
		RETURNING `+signatureColumns,
		id, usedAt, ip, userAgent))
}

// --- scan helpers ---

func scanAssessment(row pgx.Row) (domain.Assessment, error) {
	var a domain.Assessment
	var status string
	var decision *string
	if err := row.Scan(&a.ID, &a.ApplicationID, &a.Code, &status, &a.AsesorNote, &decision,
		&a.SubmittedAt, &a.ReviewedBy, &a.ReviewedAt, &a.ZoomLink, &a.ZoomDate, &a.ZoomTime,
		&a.CreatedAt, &a.UpdatedAt); err != nil {
		return domain.Assessment{}, mapError(err)
	}
	a.Status = domain.Status(status)
	if decision != nil {
		value := domain.Decision(*decision)
		a.AsesorDecision = &value
	}
	return a, nil
}

func scanAssessmentWithTotal(rows pgx.Rows) (domain.Assessment, int, error) {
	var a domain.Assessment
	var status string
	var decision *string
	var total int
	if err := rows.Scan(&a.ID, &a.ApplicationID, &a.Code, &status, &a.AsesorNote, &decision,
		&a.SubmittedAt, &a.ReviewedBy, &a.ReviewedAt, &a.ZoomLink, &a.ZoomDate, &a.ZoomTime,
		&a.CreatedAt, &a.UpdatedAt, &a.SchemaName, &a.SchemaCode, &total); err != nil {
		return domain.Assessment{}, 0, mapError(err)
	}
	a.Status = domain.Status(status)
	if decision != nil {
		value := domain.Decision(*decision)
		a.AsesorDecision = &value
	}
	return a, total, nil
}

func scanSignature(row pgx.Row) (domain.Signature, error) {
	var s domain.Signature
	if err := row.Scan(&s.ID, &s.ApplicationID, &s.Entity, &s.Role, &s.TokenHash, &s.Status,
		&s.ExpiresAt, &s.UsedAt, &s.RevokedAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return domain.Signature{}, mapError(err)
	}
	return s, nil
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
			return domain.ErrNotFound
		case "23502":
			return domain.ErrIncomplete
		}
	}
	return err
}
