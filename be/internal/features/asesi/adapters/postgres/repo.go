package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"lsp/internal/features/asesi/domain"
	"lsp/internal/platform/database"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, profile domain.Profile) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO asesi_profiles (
			id, user_id, full_name, nik, phone, gender, birth_place, birth_date,
			nationality, occupation, education_level, ktp_address, domicile_same_as_ktp,
			domicile_address, province, city, district, village, postal_code, photo,
			ktp_file, cv_file, npwp_file
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11, $12, $13,
			$14, $15, $16, $17, $18, $19, $20,
			$21, $22, $23
		)
	`, profile.ID, profile.UserID, profile.FullName, profile.NIK, profile.Phone, profile.Gender, profile.BirthPlace, profile.BirthDate,
		profile.Nationality, profile.Occupation, profile.EducationLevel, profile.KTPAddress, profile.DomicileSameAsKTP,
		profile.DomicileAddress, profile.Province, profile.City, profile.District, profile.Village, profile.PostalCode, profile.Photo,
		profile.KTPFile, profile.CVFile, profile.NPWPFile)
	return mapError(err)
}

func (r *Repository) Patch(ctx context.Context, patch domain.Patch) (domain.Profile, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanProfile(db.QueryRow(ctx, `
		UPDATE asesi_profiles
		SET full_name = COALESCE($2::varchar, full_name),
		    nik = COALESCE($3::varchar, nik),
		    phone = COALESCE($4::varchar, phone),
		    gender = COALESCE($5::gender, gender),
		    birth_place = COALESCE($6::varchar, birth_place),
		    birth_date = COALESCE($7::date, birth_date),
		    nationality = COALESCE($8::varchar, nationality),
		    occupation = COALESCE($9::varchar, occupation),
		    education_level = COALESCE($10::education_level, education_level),
		    ktp_address = COALESCE($11::text, ktp_address),
		    domicile_same_as_ktp = COALESCE($12::boolean, domicile_same_as_ktp),
		    domicile_address = COALESCE($13::text, domicile_address),
		    province = COALESCE($14::varchar, province),
		    city = COALESCE($15::varchar, city),
		    district = COALESCE($16::varchar, district),
		    village = COALESCE($17::varchar, village),
		    postal_code = COALESCE($18::varchar, postal_code),
		    photo = COALESCE($19::varchar, photo),
		    ktp_file = COALESCE($20::varchar, ktp_file),
		    cv_file = COALESCE($21::varchar, cv_file),
		    npwp_file = COALESCE($22::varchar, npwp_file)
		WHERE id = $1
		RETURNING id, user_id, full_name, nik, phone, gender::text, birth_place, birth_date,
		          nationality, occupation, education_level::text, ktp_address, domicile_same_as_ktp,
		          domicile_address, province, city, district, village, postal_code, photo,
		          ktp_file, cv_file, npwp_file, created_at, updated_at
	`, patch.ID, patch.FullName, patch.NIK, patch.Phone, patch.Gender, patch.BirthPlace, patch.BirthDate,
		patch.Nationality, patch.Occupation, patch.EducationLevel, patch.KTPAddress, patch.DomicileSameAsKTP,
		patch.DomicileAddress, patch.Province, patch.City, patch.District, patch.Village, patch.PostalCode, patch.Photo,
		patch.KTPFile, patch.CVFile, patch.NPWPFile))
}

func (r *Repository) FindByID(ctx context.Context, id string) (domain.Profile, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanProfile(db.QueryRow(ctx, selectProfileSQL()+` WHERE id = $1`, id))
}

func (r *Repository) FindByUserID(ctx context.Context, userID string) (domain.Profile, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanProfile(db.QueryRow(ctx, selectProfileSQL()+` WHERE user_id = $1`, userID))
}

func (r *Repository) List(ctx context.Context, params domain.ListParams) (domain.ListResult, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT id, user_id, full_name, nik, phone, gender::text, birth_place, birth_date,
		       nationality, occupation, education_level::text, ktp_address, domicile_same_as_ktp,
		       domicile_address, province, city, district, village, postal_code, photo,
		       ktp_file, cv_file, npwp_file, created_at, updated_at,
		       COUNT(*) OVER() AS total_count
		FROM asesi_profiles
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, params.Limit, params.Offset)
	if err != nil {
		return domain.ListResult{}, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.Profile, 0)
	total := 0
	for rows.Next() {
		item, rowTotal, err := scanProfileWithTotal(rows)
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
		if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM asesi_profiles`).Scan(&total); err != nil {
			return domain.ListResult{}, mapError(err)
		}
	}
	return domain.ListResult{Items: items, Total: total}, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	db := database.GetExecutor(ctx, r.pool)
	commandTag, err := db.Exec(ctx, `DELETE FROM asesi_profiles WHERE id = $1`, id)
	if err != nil {
		return mapError(err)
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// --- Data pekerjaan asesi (asesi_occupations) ---

func (r *Repository) CreateOccupation(ctx context.Context, occupation domain.Occupation) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO asesi_occupations (
			id, asesi_profile_id, company_name, position, tenure_months, subordinates_count,
			address, province, city, district, postal_code, office_email, office_phone
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`, occupation.ID, occupation.AsesiProfileID, occupation.CompanyName, occupation.Position, occupation.TenureMonths, occupation.SubordinatesCount,
		occupation.Address, occupation.Province, occupation.City, occupation.District, occupation.PostalCode, occupation.OfficeEmail, occupation.OfficePhone)
	return mapError(err)
}

func (r *Repository) PatchOccupation(ctx context.Context, patch domain.OccupationPatch) (domain.Occupation, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanOccupation(db.QueryRow(ctx, `
		UPDATE asesi_occupations AS o
		SET company_name = COALESCE($2::varchar, company_name),
		    position = COALESCE($3::varchar, position),
		    tenure_months = COALESCE($4::integer, tenure_months),
		    subordinates_count = COALESCE($5::integer, subordinates_count),
		    address = COALESCE($6::text, address),
		    province = COALESCE($7::varchar, province),
		    city = COALESCE($8::varchar, city),
		    district = COALESCE($9::varchar, district),
		    postal_code = COALESCE($10::varchar, postal_code),
		    office_email = COALESCE($11::citext, office_email),
		    office_phone = COALESCE($12::varchar, office_phone)
		WHERE o.id = $1
		RETURNING `+occupationColumns,
		patch.ID, patch.CompanyName, patch.Position, patch.TenureMonths, patch.SubordinatesCount,
		patch.Address, patch.Province, patch.City, patch.District, patch.PostalCode, patch.OfficeEmail, patch.OfficePhone))
}

func (r *Repository) FindOccupationByID(ctx context.Context, id string) (domain.Occupation, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanOccupation(db.QueryRow(ctx, `SELECT `+occupationColumns+` FROM asesi_occupations o WHERE o.id = $1`, id))
}

func (r *Repository) FindOccupationByUserID(ctx context.Context, userID string) (domain.Occupation, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanOccupation(db.QueryRow(ctx, `
		SELECT `+occupationColumns+`
		FROM asesi_occupations o
		JOIN asesi_profiles p ON p.id = o.asesi_profile_id
		WHERE p.user_id = $1
	`, userID))
}

// UpsertOccupationByUser membuat atau memperbarui data pekerjaan milik user
// (asesi). asesi_profile_id di-resolve dari user_id; bila profil tidak ada,
// NOT NULL/FK akan gagal dan dipetakan ke ErrNotFound.
func (r *Repository) UpsertOccupationByUser(ctx context.Context, userID string, occupation domain.Occupation) (domain.Occupation, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanOccupation(db.QueryRow(ctx, `
		INSERT INTO asesi_occupations (
			id, asesi_profile_id, company_name, position, tenure_months, subordinates_count,
			address, province, city, district, postal_code, office_email, office_phone
		)
		VALUES (
			$1, (SELECT id FROM asesi_profiles WHERE user_id = $2), $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13
		)
		ON CONFLICT (asesi_profile_id) DO UPDATE SET
		    company_name = EXCLUDED.company_name,
		    position = EXCLUDED.position,
		    tenure_months = EXCLUDED.tenure_months,
		    subordinates_count = EXCLUDED.subordinates_count,
		    address = EXCLUDED.address,
		    province = EXCLUDED.province,
		    city = EXCLUDED.city,
		    district = EXCLUDED.district,
		    postal_code = EXCLUDED.postal_code,
		    office_email = EXCLUDED.office_email,
		    office_phone = EXCLUDED.office_phone
		RETURNING `+occupationColumns,
		occupation.ID, userID, occupation.CompanyName, occupation.Position, occupation.TenureMonths, occupation.SubordinatesCount,
		occupation.Address, occupation.Province, occupation.City, occupation.District, occupation.PostalCode, occupation.OfficeEmail, occupation.OfficePhone))
}

func (r *Repository) ListOccupations(ctx context.Context, params domain.ListParams) (domain.OccupationListResult, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT `+occupationColumns+`, COUNT(*) OVER() AS total_count
		FROM asesi_occupations o
		ORDER BY o.created_at DESC
		LIMIT $1 OFFSET $2
	`, params.Limit, params.Offset)
	if err != nil {
		return domain.OccupationListResult{}, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.Occupation, 0)
	total := 0
	for rows.Next() {
		item, rowTotal, err := scanOccupationWithTotal(rows)
		if err != nil {
			return domain.OccupationListResult{}, err
		}
		total = rowTotal
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.OccupationListResult{}, mapError(err)
	}
	if len(items) == 0 && params.Offset > 0 {
		if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM asesi_occupations`).Scan(&total); err != nil {
			return domain.OccupationListResult{}, mapError(err)
		}
	}
	return domain.OccupationListResult{Items: items, Total: total}, nil
}

func (r *Repository) DeleteOccupation(ctx context.Context, id string) error {
	db := database.GetExecutor(ctx, r.pool)
	commandTag, err := db.Exec(ctx, `DELETE FROM asesi_occupations WHERE id = $1`, id)
	if err != nil {
		return mapError(err)
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

const occupationColumns = `o.id, o.asesi_profile_id, o.company_name, o.position, o.tenure_months, o.subordinates_count,
	o.address, o.province, o.city, o.district, o.postal_code, o.office_email, o.office_phone,
	o.created_at, o.updated_at`

func scanOccupation(row pgx.Row) (domain.Occupation, error) {
	var o domain.Occupation
	if err := row.Scan(
		&o.ID, &o.AsesiProfileID, &o.CompanyName, &o.Position, &o.TenureMonths, &o.SubordinatesCount,
		&o.Address, &o.Province, &o.City, &o.District, &o.PostalCode, &o.OfficeEmail, &o.OfficePhone,
		&o.CreatedAt, &o.UpdatedAt,
	); err != nil {
		return domain.Occupation{}, mapError(err)
	}
	return o, nil
}

func scanOccupationWithTotal(rows pgx.Rows) (domain.Occupation, int, error) {
	var o domain.Occupation
	var total int
	if err := rows.Scan(
		&o.ID, &o.AsesiProfileID, &o.CompanyName, &o.Position, &o.TenureMonths, &o.SubordinatesCount,
		&o.Address, &o.Province, &o.City, &o.District, &o.PostalCode, &o.OfficeEmail, &o.OfficePhone,
		&o.CreatedAt, &o.UpdatedAt, &total,
	); err != nil {
		return domain.Occupation{}, 0, mapError(err)
	}
	return o, total, nil
}

func selectProfileSQL() string {
	return `
		SELECT id, user_id, full_name, nik, phone, gender::text, birth_place, birth_date,
		       nationality, occupation, education_level::text, ktp_address, domicile_same_as_ktp,
		       domicile_address, province, city, district, village, postal_code, photo,
		       ktp_file, cv_file, npwp_file, created_at, updated_at
		FROM asesi_profiles`
}

func scanProfile(row pgx.Row) (domain.Profile, error) {
	var profile domain.Profile
	var gender string
	var educationLevel string
	if err := row.Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FullName,
		&profile.NIK,
		&profile.Phone,
		&gender,
		&profile.BirthPlace,
		&profile.BirthDate,
		&profile.Nationality,
		&profile.Occupation,
		&educationLevel,
		&profile.KTPAddress,
		&profile.DomicileSameAsKTP,
		&profile.DomicileAddress,
		&profile.Province,
		&profile.City,
		&profile.District,
		&profile.Village,
		&profile.PostalCode,
		&profile.Photo,
		&profile.KTPFile,
		&profile.CVFile,
		&profile.NPWPFile,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	); err != nil {
		return domain.Profile{}, mapError(err)
	}
	profile.Gender = domain.Gender(gender)
	profile.EducationLevel = domain.EducationLevel(educationLevel)
	return profile, nil
}

func scanProfileWithTotal(rows pgx.Rows) (domain.Profile, int, error) {
	var profile domain.Profile
	var gender string
	var educationLevel string
	var total int
	if err := rows.Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FullName,
		&profile.NIK,
		&profile.Phone,
		&gender,
		&profile.BirthPlace,
		&profile.BirthDate,
		&profile.Nationality,
		&profile.Occupation,
		&educationLevel,
		&profile.KTPAddress,
		&profile.DomicileSameAsKTP,
		&profile.DomicileAddress,
		&profile.Province,
		&profile.City,
		&profile.District,
		&profile.Village,
		&profile.PostalCode,
		&profile.Photo,
		&profile.KTPFile,
		&profile.CVFile,
		&profile.NPWPFile,
		&profile.CreatedAt,
		&profile.UpdatedAt,
		&total,
	); err != nil {
		return domain.Profile{}, 0, mapError(err)
	}
	profile.Gender = domain.Gender(gender)
	profile.EducationLevel = domain.EducationLevel(educationLevel)
	return profile, total, nil
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
			// NOT NULL violation dari subquery asesi_profile_id yang kosong
			// (upsert data pekerjaan saat profil asesi tidak ada).
			return domain.ErrNotFound
		case "23514":
			if pgErr.ConstraintName == "ck_asesi_profiles_domicile" {
				return domain.ErrInvalidDomicile
			}
		}
	}
	return err
}
