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
		case "23514":
			if pgErr.ConstraintName == "ck_asesi_profiles_domicile" {
				return domain.ErrInvalidDomicile
			}
		}
	}
	return err
}
