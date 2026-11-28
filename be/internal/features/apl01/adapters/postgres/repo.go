package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"lsp/internal/features/apl01/domain"
	"lsp/internal/platform/database"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const applicationColumns = `id, user_id, code, status::text, current_step, submitted_at,
	verified_by, verified_at, verification_note, asesor_id, asesor_assigned_at,
	created_at, updated_at, deleted_at`

const dataColumns = `id, application_id, certification_schema_id, purpose::text, created_at, updated_at`

const documentColumns = `id, application_id, document_type_id, file_url, status::text, note, verified_at, created_at, updated_at`

const scheduleColumns = `id, application_id, mode::text, scheduled_date, tuk_id, created_at, updated_at`

const signatureColumns = `id, application_id, entity::text, role::text, token_hash, status::text,
	expires_at, used_at, used_ip::text, used_user_agent, revoked_at, created_at, updated_at`

// --- Header permohonan ---

func (r *Repository) CreateApplication(ctx context.Context, application domain.Application) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO certification_applications (id, user_id, code, status, current_step)
		VALUES ($1, $2, $3, $4, $5)
	`, application.ID, application.UserID, application.Code, application.Status, application.CurrentStep)
	return mapError(err)
}

func (r *Repository) FindApplicationByID(ctx context.Context, id string) (domain.Application, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanApplication(db.QueryRow(ctx, `
		SELECT `+applicationColumns+`
		FROM certification_applications
		WHERE id = $1 AND deleted_at IS NULL
	`, id))
}

func (r *Repository) FindActiveDraft(ctx context.Context, userID string) (domain.Application, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanApplication(db.QueryRow(ctx, `
		SELECT `+applicationColumns+`
		FROM certification_applications
		WHERE user_id = $1 AND status = 'DRAFT' AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT 1
	`, userID))
}

func (r *Repository) ListApplications(ctx context.Context, params domain.ListParams) (domain.ListResult, error) {
	db := database.GetExecutor(ctx, r.pool)
	var status *string
	if params.Status != nil {
		value := string(*params.Status)
		status = &value
	}
	rows, err := db.Query(ctx, `
		SELECT `+applicationColumns+`, COUNT(*) OVER() AS total_count
		FROM certification_applications
		WHERE deleted_at IS NULL
		  AND ($1::text = '' OR user_id = $1::text)
		  AND ($2::text IS NULL OR status::text = $2::text)
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4
	`, params.UserID, status, params.Limit, params.Offset)
	if err != nil {
		return domain.ListResult{}, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.Application, 0)
	total := 0
	for rows.Next() {
		item, rowTotal, err := scanApplicationWithTotal(rows)
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
		if err := db.QueryRow(ctx, `
			SELECT COUNT(*) FROM certification_applications
			WHERE deleted_at IS NULL
			  AND ($1::text = '' OR user_id = $1::text)
			  AND ($2::text IS NULL OR status::text = $2::text)
		`, params.UserID, status).Scan(&total); err != nil {
			return domain.ListResult{}, mapError(err)
		}
	}
	return domain.ListResult{Items: items, Total: total}, nil
}

func (r *Repository) UpdateCurrentStep(ctx context.Context, id string, step int) (domain.Application, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanApplication(db.QueryRow(ctx, `
		UPDATE certification_applications
		SET current_step = $2
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING `+applicationColumns,
		id, step))
}

// MarkSubmitted mengubah status menjadi SUBMITTED. Dipanggil saat asesi
// memverifikasi permohonan via QR.
func (r *Repository) MarkSubmitted(ctx context.Context, id string, submittedAt time.Time) (domain.Application, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanApplication(db.QueryRow(ctx, `
		UPDATE certification_applications
		SET status = 'SUBMITTED', submitted_at = $2
		WHERE id = $1 AND status = 'DRAFT' AND deleted_at IS NULL
		RETURNING `+applicationColumns,
		id, submittedAt))
}

func (r *Repository) ApplyDecision(ctx context.Context, decision domain.Decision) (domain.Application, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanApplication(db.QueryRow(ctx, `
		UPDATE certification_applications
		SET status = $2,
		    verified_by = $3,
		    verified_at = $4,
		    verification_note = $5,
		    asesor_id = COALESCE($6::varchar, asesor_id),
		    asesor_assigned_at = COALESCE($7::timestamptz, asesor_assigned_at)
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING `+applicationColumns,
		decision.ApplicationID, decision.Status, decision.VerifiedBy, decision.VerifiedAt,
		decision.Note, decision.AsesorID, decision.AsesorAssignedAt))
}

// --- Penugasan asesor (plotting) ---

const assignmentColumns = `a.id, a.application_id, a.asesor_id, COALESCE(u.name, ''), a.status::text,
	a.assigned_by, a.assigned_at, a.responded_at, a.note, a.created_at, a.updated_at`

func (r *Repository) FindActiveAssignment(ctx context.Context, applicationID string) (*domain.AsesorAssignment, error) {
	db := database.GetExecutor(ctx, r.pool)
	assignment, err := scanAssignment(db.QueryRow(ctx, `
		SELECT `+assignmentColumns+`
		FROM asesor_assignments a
		LEFT JOIN users u ON u.id = a.asesor_id
		WHERE a.application_id = $1 AND a.status IN ('PENDING', 'ACCEPTED')
		LIMIT 1
	`, applicationID))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &assignment, nil
}

func (r *Repository) AssignAsesor(ctx context.Context, assignment domain.AsesorAssignment) error {
	db := database.GetExecutor(ctx, r.pool)
	// Batalkan penugasan aktif sebelumnya (bila ada) agar hanya tersisa satu.
	if _, err := db.Exec(ctx, `
		UPDATE asesor_assignments
		SET status = 'CANCELLED', responded_at = $2
		WHERE application_id = $1 AND status IN ('PENDING', 'ACCEPTED')
	`, assignment.ApplicationID, assignment.AssignedAt); err != nil {
		return mapError(err)
	}
	// Simpan penugasan baru.
	if _, err := db.Exec(ctx, `
		INSERT INTO asesor_assignments (id, application_id, asesor_id, status, assigned_by, assigned_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, assignment.ID, assignment.ApplicationID, assignment.AsesorID, assignment.Status,
		assignment.AssignedBy, assignment.AssignedAt); err != nil {
		return mapError(err)
	}
	// Sinkronkan asesor aktif di header permohonan.
	if _, err := db.Exec(ctx, `
		UPDATE certification_applications
		SET asesor_id = $2, asesor_assigned_at = $3
		WHERE id = $1 AND deleted_at IS NULL
	`, assignment.ApplicationID, assignment.AsesorID, assignment.AssignedAt); err != nil {
		return mapError(err)
	}
	return nil
}

// ListAsesorCandidates mengembalikan asesor yang terafiliasi dengan TUK jadwal
// permohonan, beserta jumlah penugasan aktif (PENDING/ACCEPTED) pada tanggal
// jadwal tersebut. TUK dan tanggal diambil dari assessment_schedules.
func (r *Repository) ListAsesorCandidates(ctx context.Context, applicationID string) ([]domain.AsesorCandidate, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT ap.user_id, COALESCE(u.name, ''), COALESCE(ap.province, ''), COALESCE(ap.city, ''),
		       ap.daily_assignment_limit,
		       (SELECT COUNT(*)
		        FROM asesor_assignments aa
		        JOIN assessment_schedules sc ON sc.application_id = aa.application_id
		        WHERE aa.status IN ('PENDING', 'ACCEPTED')
		          AND aa.asesor_id = ap.user_id
		          AND sc.scheduled_date = s.scheduled_date) AS assigned_count
		FROM assessment_schedules s
		JOIN asesor_profiles ap ON ap.tuk_id = s.tuk_id
		JOIN users u ON u.id = ap.user_id
		WHERE s.application_id = $1 AND u.deleted_at IS NULL
		ORDER BY u.name
	`, applicationID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.AsesorCandidate, 0)
	for rows.Next() {
		var c domain.AsesorCandidate
		if err := rows.Scan(&c.UserID, &c.Name, &c.Province, &c.City,
			&c.DailyAssignmentLimit, &c.AssignedCount); err != nil {
			return nil, mapError(err)
		}
		c.Disabled = c.DailyAssignmentLimit > 0 && c.AssignedCount >= c.DailyAssignmentLimit
		items = append(items, c)
	}
	if err := rows.Err(); err != nil {
		return nil, mapError(err)
	}
	return items, nil
}

// --- Data sertifikasi (step 2) ---

func (r *Repository) UpsertApplicationData(ctx context.Context, data domain.ApplicationData) (domain.ApplicationData, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanApplicationData(db.QueryRow(ctx, `
		INSERT INTO certification_application_data (id, application_id, certification_schema_id, purpose)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (application_id) DO UPDATE SET
		    certification_schema_id = EXCLUDED.certification_schema_id,
		    purpose = EXCLUDED.purpose
		RETURNING `+dataColumns,
		data.ID, data.ApplicationID, data.CertificationSchemaID, data.Purpose))
}

func (r *Repository) FindApplicationData(ctx context.Context, applicationID string) (domain.ApplicationData, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanApplicationData(db.QueryRow(ctx, `
		SELECT `+dataColumns+`
		FROM certification_application_data
		WHERE application_id = $1
	`, applicationID))
}

// --- Bukti kelengkapan (step 3) ---

func (r *Repository) ListApplicationDocuments(ctx context.Context, applicationID string) ([]domain.ApplicationDocument, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT `+documentColumns+`
		FROM application_documents
		WHERE application_id = $1
		ORDER BY created_at, document_type_id
	`, applicationID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.ApplicationDocument, 0)
	for rows.Next() {
		item, err := scanDocument(rows)
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

func (r *Repository) ReplaceApplicationDocuments(ctx context.Context, applicationID string, documents []domain.ApplicationDocument) error {
	db := database.GetExecutor(ctx, r.pool)

	documentTypeIDs := make([]string, 0, len(documents))
	for _, doc := range documents {
		documentTypeIDs = append(documentTypeIDs, doc.DocumentTypeID)
	}
	if _, err := db.Exec(ctx, `
		DELETE FROM application_documents
		WHERE application_id = $1
		  AND NOT (document_type_id = ANY($2::varchar[]))
	`, applicationID, documentTypeIDs); err != nil {
		return mapError(err)
	}

	for _, doc := range documents {
		if _, err := db.Exec(ctx, `
			INSERT INTO application_documents (id, application_id, document_type_id, file_url, status)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (application_id, document_type_id) DO UPDATE SET
			    file_url = EXCLUDED.file_url,
			    status = EXCLUDED.status,
			    note = NULL,
			    verified_at = NULL
		`, doc.ID, doc.ApplicationID, doc.DocumentTypeID, doc.FileURL, doc.Status); err != nil {
			return mapError(err)
		}
	}
	return nil
}

func (r *Repository) CountMissingRequiredDocuments(ctx context.Context, applicationID string) (int, error) {
	db := database.GetExecutor(ctx, r.pool)
	var count int
	err := db.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM schema_document_requirements r
		JOIN certification_application_data d ON d.application_id = $1
		WHERE r.certification_schema_id = d.certification_schema_id
		  AND r.is_required = true
		  AND NOT EXISTS (
		      SELECT 1 FROM application_documents ad
		      WHERE ad.application_id = $1
		        AND ad.document_type_id = r.document_type_id
		        AND ad.file_url IS NOT NULL AND ad.file_url <> ''
		  )
	`, applicationID).Scan(&count)
	if err != nil {
		return 0, mapError(err)
	}
	return count, nil
}

// --- Jadwal asesmen (step 4) ---

func (r *Repository) UpsertAssessmentSchedule(ctx context.Context, schedule domain.AssessmentSchedule) (domain.AssessmentSchedule, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanSchedule(db.QueryRow(ctx, `
		INSERT INTO assessment_schedules (id, application_id, mode, scheduled_date, tuk_id)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (application_id) DO UPDATE SET
		    mode = EXCLUDED.mode,
		    scheduled_date = EXCLUDED.scheduled_date,
		    tuk_id = EXCLUDED.tuk_id
		RETURNING `+scheduleColumns,
		schedule.ID, schedule.ApplicationID, schedule.Mode, schedule.ScheduledDate, schedule.TukID))
}

func (r *Repository) FindAssessmentSchedule(ctx context.Context, applicationID string) (domain.AssessmentSchedule, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanSchedule(db.QueryRow(ctx, `
		SELECT `+scheduleColumns+`
		FROM assessment_schedules
		WHERE application_id = $1
	`, applicationID))
}

// HasScheduleConflict mengecek apakah asesi pemilik permohonan punya permohonan
// lain (status selain REJECTED) dengan jadwal pada tanggal yang sama.
func (r *Repository) HasScheduleConflict(ctx context.Context, applicationID string, date time.Time) (bool, error) {
	db := database.GetExecutor(ctx, r.pool)
	var exists bool
	err := db.QueryRow(ctx, `
		SELECT EXISTS (
		    SELECT 1
		    FROM assessment_schedules s
		    JOIN certification_applications a ON a.id = s.application_id
		    WHERE s.scheduled_date = $2::date
		      AND a.id <> $1
		      AND a.status <> 'REJECTED'
		      AND a.deleted_at IS NULL
		      AND a.user_id = (SELECT user_id FROM certification_applications WHERE id = $1)
		)
	`, applicationID, date.Format(time.DateOnly)).Scan(&exists)
	if err != nil {
		return false, mapError(err)
	}
	return exists, nil
}

// --- Token verifikasi QR ---

func (r *Repository) CreateSignature(ctx context.Context, signature domain.FormSignature) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO form_signatures (id, application_id, entity, role, token_hash, status, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, signature.ID, signature.ApplicationID, signature.Entity, signature.Role,
		signature.TokenHash, signature.Status, signature.ExpiresAt)
	return mapError(err)
}

func (r *Repository) RevokeActiveSignatures(ctx context.Context, applicationID string, entity domain.FormEntity, role domain.SignerRole, revokedAt time.Time) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		UPDATE form_signatures
		SET status = 'dibatalkan', revoked_at = $4
		WHERE application_id = $1 AND entity = $2 AND role = $3 AND status = 'menunggu'
	`, applicationID, entity, role, revokedAt)
	return mapError(err)
}

// HasPendingSignature melaporkan apakah permohonan masih punya token QR asesi
// yang menunggu (belum dipakai & belum kedaluwarsa). Selama token ini ada,
// formulir APL-01 terkunci dari perubahan.
func (r *Repository) HasPendingSignature(ctx context.Context, applicationID string) (bool, error) {
	db := database.GetExecutor(ctx, r.pool)
	var exists bool
	err := db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM form_signatures
			WHERE application_id = $1 AND entity = 'permohonan' AND role = 'asesi'
			  AND status = 'menunggu' AND expires_at > now()
		)
	`, applicationID).Scan(&exists)
	if err != nil {
		return false, mapError(err)
	}
	return exists, nil
}

func (r *Repository) FindSignatureByTokenHash(ctx context.Context, tokenHash string) (domain.FormSignature, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanSignature(db.QueryRow(ctx, `
		SELECT `+signatureColumns+`
		FROM form_signatures
		WHERE token_hash = $1
	`, tokenHash))
}

func (r *Repository) MarkSignatureUsed(ctx context.Context, id string, usedAt time.Time, ip, userAgent *string) (domain.FormSignature, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanSignature(db.QueryRow(ctx, `
		UPDATE form_signatures
		SET status = 'dipakai', used_at = $2, used_ip = $3::inet, used_user_agent = $4
		WHERE id = $1 AND status = 'menunggu'
		RETURNING `+signatureColumns,
		id, usedAt, ip, userAgent))
}

func (r *Repository) ListSignatures(ctx context.Context, applicationID string) ([]domain.FormSignature, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT `+signatureColumns+`
		FROM form_signatures
		WHERE application_id = $1
		ORDER BY created_at DESC
	`, applicationID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.FormSignature, 0)
	for rows.Next() {
		item, err := scanSignature(rows)
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

// --- scan helpers ---

func scanApplication(row pgx.Row) (domain.Application, error) {
	var a domain.Application
	var status string
	if err := row.Scan(&a.ID, &a.UserID, &a.Code, &status, &a.CurrentStep, &a.SubmittedAt,
		&a.VerifiedBy, &a.VerifiedAt, &a.VerificationNote, &a.AsesorID, &a.AsesorAssignedAt,
		&a.CreatedAt, &a.UpdatedAt, &a.DeletedAt); err != nil {
		return domain.Application{}, mapError(err)
	}
	a.Status = domain.Status(status)
	return a, nil
}

func scanApplicationWithTotal(rows pgx.Rows) (domain.Application, int, error) {
	var a domain.Application
	var status string
	var total int
	if err := rows.Scan(&a.ID, &a.UserID, &a.Code, &status, &a.CurrentStep, &a.SubmittedAt,
		&a.VerifiedBy, &a.VerifiedAt, &a.VerificationNote, &a.AsesorID, &a.AsesorAssignedAt,
		&a.CreatedAt, &a.UpdatedAt, &a.DeletedAt, &total); err != nil {
		return domain.Application{}, 0, mapError(err)
	}
	a.Status = domain.Status(status)
	return a, total, nil
}

func scanAssignment(row pgx.Row) (domain.AsesorAssignment, error) {
	var a domain.AsesorAssignment
	if err := row.Scan(&a.ID, &a.ApplicationID, &a.AsesorID, &a.AsesorName, &a.Status,
		&a.AssignedBy, &a.AssignedAt, &a.RespondedAt, &a.Note, &a.CreatedAt, &a.UpdatedAt); err != nil {
		return domain.AsesorAssignment{}, mapError(err)
	}
	return a, nil
}

func scanApplicationData(row pgx.Row) (domain.ApplicationData, error) {
	var d domain.ApplicationData
	var purpose string
	if err := row.Scan(&d.ID, &d.ApplicationID, &d.CertificationSchemaID, &purpose, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return domain.ApplicationData{}, mapError(err)
	}
	d.Purpose = domain.Purpose(purpose)
	return d, nil
}

func scanDocument(row pgx.Row) (domain.ApplicationDocument, error) {
	var d domain.ApplicationDocument
	var status string
	if err := row.Scan(&d.ID, &d.ApplicationID, &d.DocumentTypeID, &d.FileURL, &status, &d.Note,
		&d.VerifiedAt, &d.CreatedAt, &d.UpdatedAt); err != nil {
		return domain.ApplicationDocument{}, mapError(err)
	}
	d.Status = domain.DocumentStatus(status)
	return d, nil
}

func scanSchedule(row pgx.Row) (domain.AssessmentSchedule, error) {
	var s domain.AssessmentSchedule
	var mode string
	if err := row.Scan(&s.ID, &s.ApplicationID, &mode, &s.ScheduledDate, &s.TukID, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return domain.AssessmentSchedule{}, mapError(err)
	}
	s.Mode = domain.Mode(mode)
	return s, nil
}

func scanSignature(row pgx.Row) (domain.FormSignature, error) {
	var s domain.FormSignature
	var entity string
	var role string
	var status string
	if err := row.Scan(&s.ID, &s.ApplicationID, &entity, &role, &s.TokenHash, &status,
		&s.ExpiresAt, &s.UsedAt, &s.UsedIP, &s.UsedUserAgent, &s.RevokedAt, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return domain.FormSignature{}, mapError(err)
	}
	s.Entity = domain.FormEntity(entity)
	s.Role = domain.SignerRole(role)
	s.Status = domain.SignatureStatus(status)
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
