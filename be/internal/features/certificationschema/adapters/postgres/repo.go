package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"lsp/internal/features/certificationschema/domain"
	"lsp/internal/platform/database"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const (
	schemaColumns   = `id, code, name, schema_code, sort_order, is_active, created_at, updated_at`
	unitColumns     = `id, certification_schema_id, code, title, sort_order, created_at, updated_at`
	elementColumns  = `id, competency_unit_id, title, sort_order, created_at, updated_at`
	criteriaColumns = `id, competency_unit_element_id, content, sort_order, created_at, updated_at`
)

// --- Skema sertifikasi ---

func (r *Repository) CreateSchema(ctx context.Context, schema domain.CertificationSchema) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO certification_schemas (id, code, name, schema_code, sort_order, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, schema.ID, schema.Code, schema.Name, schema.SchemaCode, schema.SortOrder, schema.IsActive)
	return mapError(err)
}

func (r *Repository) PatchSchema(ctx context.Context, patch domain.CertificationSchemaPatch) (domain.CertificationSchema, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanSchema(db.QueryRow(ctx, `
		UPDATE certification_schemas
		SET code = COALESCE($2::varchar, code),
		    name = COALESCE($3::varchar, name),
		    schema_code = COALESCE($4::varchar, schema_code),
		    sort_order = COALESCE($5::smallint, sort_order),
		    is_active = COALESCE($6::boolean, is_active)
		WHERE id = $1
		RETURNING `+schemaColumns,
		patch.ID, patch.Code, patch.Name, patch.SchemaCode, patch.SortOrder, patch.IsActive))
}

func (r *Repository) FindSchemaByID(ctx context.Context, id string) (domain.CertificationSchema, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanSchema(db.QueryRow(ctx, `SELECT `+schemaColumns+` FROM certification_schemas WHERE id = $1`, id))
}

func (r *Repository) ListSchemas(ctx context.Context, params domain.ListParams) (domain.SchemaListResult, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT `+schemaColumns+`, COUNT(*) OVER() AS total_count
		FROM certification_schemas
		ORDER BY sort_order, created_at
		LIMIT $1 OFFSET $2
	`, params.Limit, params.Offset)
	if err != nil {
		return domain.SchemaListResult{}, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.CertificationSchema, 0)
	total := 0
	for rows.Next() {
		item, rowTotal, err := scanSchemaWithTotal(rows)
		if err != nil {
			return domain.SchemaListResult{}, err
		}
		total = rowTotal
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.SchemaListResult{}, mapError(err)
	}
	if len(items) == 0 && params.Offset > 0 {
		if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM certification_schemas`).Scan(&total); err != nil {
			return domain.SchemaListResult{}, mapError(err)
		}
	}
	return domain.SchemaListResult{Items: items, Total: total}, nil
}

func (r *Repository) DeleteSchema(ctx context.Context, id string) error {
	db := database.GetExecutor(ctx, r.pool)
	commandTag, err := db.Exec(ctx, `DELETE FROM certification_schemas WHERE id = $1`, id)
	if err != nil {
		return mapError(err)
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// SchemaTree mengambil skema beserta unit -> elemen -> kriteria dalam 3 query
// lalu dirakit di memori.
func (r *Repository) SchemaTree(ctx context.Context, id string) (domain.SchemaTree, error) {
	schema, err := r.FindSchemaByID(ctx, id)
	if err != nil {
		return domain.SchemaTree{}, err
	}

	units, err := r.listUnits(ctx, id)
	if err != nil {
		return domain.SchemaTree{}, err
	}
	elements, err := r.listElements(ctx, id)
	if err != nil {
		return domain.SchemaTree{}, err
	}
	criteria, err := r.listCriteria(ctx, id)
	if err != nil {
		return domain.SchemaTree{}, err
	}

	criteriaByElement := make(map[string][]domain.CompetencyUnitCriteria)
	for _, c := range criteria {
		criteriaByElement[c.CompetencyUnitElementID] = append(criteriaByElement[c.CompetencyUnitElementID], c)
	}

	elementsByUnit := make(map[string][]domain.ElementNode)
	for _, e := range elements {
		node := domain.ElementNode{Element: e, Criteria: criteriaByElement[e.ID]}
		if node.Criteria == nil {
			node.Criteria = make([]domain.CompetencyUnitCriteria, 0)
		}
		elementsByUnit[e.CompetencyUnitID] = append(elementsByUnit[e.CompetencyUnitID], node)
	}

	unitNodes := make([]domain.UnitNode, 0, len(units))
	for _, u := range units {
		node := domain.UnitNode{Unit: u, Elements: elementsByUnit[u.ID]}
		if node.Elements == nil {
			node.Elements = make([]domain.ElementNode, 0)
		}
		unitNodes = append(unitNodes, node)
	}

	return domain.SchemaTree{Schema: schema, Units: unitNodes}, nil
}

func (r *Repository) listUnits(ctx context.Context, schemaID string) ([]domain.CompetencyUnit, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT `+unitColumns+`
		FROM competency_units
		WHERE certification_schema_id = $1
		ORDER BY sort_order, created_at
	`, schemaID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.CompetencyUnit, 0)
	for rows.Next() {
		item, err := scanUnit(rows)
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

func (r *Repository) listElements(ctx context.Context, schemaID string) ([]domain.CompetencyUnitElement, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT e.id, e.competency_unit_id, e.title, e.sort_order, e.created_at, e.updated_at
		FROM competency_unit_elements e
		JOIN competency_units u ON u.id = e.competency_unit_id
		WHERE u.certification_schema_id = $1
		ORDER BY e.sort_order, e.created_at
	`, schemaID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.CompetencyUnitElement, 0)
	for rows.Next() {
		item, err := scanElement(rows)
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

func (r *Repository) listCriteria(ctx context.Context, schemaID string) ([]domain.CompetencyUnitCriteria, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT c.id, c.competency_unit_element_id, c.content, c.sort_order, c.created_at, c.updated_at
		FROM competency_unit_criteria c
		JOIN competency_unit_elements e ON e.id = c.competency_unit_element_id
		JOIN competency_units u ON u.id = e.competency_unit_id
		WHERE u.certification_schema_id = $1
		ORDER BY c.sort_order, c.created_at
	`, schemaID)
	if err != nil {
		return nil, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.CompetencyUnitCriteria, 0)
	for rows.Next() {
		item, err := scanCriteria(rows)
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

// --- Unit kompetensi ---

func (r *Repository) CreateUnit(ctx context.Context, unit domain.CompetencyUnit) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO competency_units (id, certification_schema_id, code, title, sort_order)
		VALUES ($1, $2, $3, $4, $5)
	`, unit.ID, unit.CertificationSchemaID, unit.Code, unit.Title, unit.SortOrder)
	return mapError(err)
}

func (r *Repository) PatchUnit(ctx context.Context, patch domain.CompetencyUnitPatch) (domain.CompetencyUnit, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanUnit(db.QueryRow(ctx, `
		UPDATE competency_units
		SET code = COALESCE($2::varchar, code),
		    title = COALESCE($3::varchar, title),
		    sort_order = COALESCE($4::smallint, sort_order)
		WHERE id = $1
		RETURNING `+unitColumns,
		patch.ID, patch.Code, patch.Title, patch.SortOrder))
}

func (r *Repository) FindUnitByID(ctx context.Context, id string) (domain.CompetencyUnit, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanUnit(db.QueryRow(ctx, `SELECT `+unitColumns+` FROM competency_units WHERE id = $1`, id))
}

func (r *Repository) ListUnitsBySchema(ctx context.Context, schemaID string, params domain.ListParams) (domain.UnitListResult, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT `+unitColumns+`, COUNT(*) OVER() AS total_count
		FROM competency_units
		WHERE certification_schema_id = $1
		ORDER BY sort_order, created_at
		LIMIT $2 OFFSET $3
	`, schemaID, params.Limit, params.Offset)
	if err != nil {
		return domain.UnitListResult{}, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.CompetencyUnit, 0)
	total := 0
	for rows.Next() {
		item, rowTotal, err := scanUnitWithTotal(rows)
		if err != nil {
			return domain.UnitListResult{}, err
		}
		total = rowTotal
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.UnitListResult{}, mapError(err)
	}
	if len(items) == 0 && params.Offset > 0 {
		if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM competency_units WHERE certification_schema_id = $1`, schemaID).Scan(&total); err != nil {
			return domain.UnitListResult{}, mapError(err)
		}
	}
	return domain.UnitListResult{Items: items, Total: total}, nil
}

func (r *Repository) DeleteUnit(ctx context.Context, id string) error {
	db := database.GetExecutor(ctx, r.pool)
	commandTag, err := db.Exec(ctx, `DELETE FROM competency_units WHERE id = $1`, id)
	if err != nil {
		return mapError(err)
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// --- Elemen unit kompetensi ---

func (r *Repository) CreateElement(ctx context.Context, element domain.CompetencyUnitElement) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO competency_unit_elements (id, competency_unit_id, title, sort_order)
		VALUES ($1, $2, $3, $4)
	`, element.ID, element.CompetencyUnitID, element.Title, element.SortOrder)
	return mapError(err)
}

func (r *Repository) PatchElement(ctx context.Context, patch domain.CompetencyUnitElementPatch) (domain.CompetencyUnitElement, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanElement(db.QueryRow(ctx, `
		UPDATE competency_unit_elements
		SET title = COALESCE($2::varchar, title),
		    sort_order = COALESCE($3::smallint, sort_order)
		WHERE id = $1
		RETURNING `+elementColumns,
		patch.ID, patch.Title, patch.SortOrder))
}

func (r *Repository) FindElementByID(ctx context.Context, id string) (domain.CompetencyUnitElement, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanElement(db.QueryRow(ctx, `SELECT `+elementColumns+` FROM competency_unit_elements WHERE id = $1`, id))
}

func (r *Repository) ListElementsByUnit(ctx context.Context, unitID string, params domain.ListParams) (domain.ElementListResult, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT `+elementColumns+`, COUNT(*) OVER() AS total_count
		FROM competency_unit_elements
		WHERE competency_unit_id = $1
		ORDER BY sort_order, created_at
		LIMIT $2 OFFSET $3
	`, unitID, params.Limit, params.Offset)
	if err != nil {
		return domain.ElementListResult{}, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.CompetencyUnitElement, 0)
	total := 0
	for rows.Next() {
		item, rowTotal, err := scanElementWithTotal(rows)
		if err != nil {
			return domain.ElementListResult{}, err
		}
		total = rowTotal
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.ElementListResult{}, mapError(err)
	}
	if len(items) == 0 && params.Offset > 0 {
		if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM competency_unit_elements WHERE competency_unit_id = $1`, unitID).Scan(&total); err != nil {
			return domain.ElementListResult{}, mapError(err)
		}
	}
	return domain.ElementListResult{Items: items, Total: total}, nil
}

func (r *Repository) DeleteElement(ctx context.Context, id string) error {
	db := database.GetExecutor(ctx, r.pool)
	commandTag, err := db.Exec(ctx, `DELETE FROM competency_unit_elements WHERE id = $1`, id)
	if err != nil {
		return mapError(err)
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// --- Kriteria unjuk kerja (KUK) ---

func (r *Repository) CreateCriteria(ctx context.Context, criteria domain.CompetencyUnitCriteria) error {
	db := database.GetExecutor(ctx, r.pool)
	_, err := db.Exec(ctx, `
		INSERT INTO competency_unit_criteria (id, competency_unit_element_id, content, sort_order)
		VALUES ($1, $2, $3, $4)
	`, criteria.ID, criteria.CompetencyUnitElementID, criteria.Content, criteria.SortOrder)
	return mapError(err)
}

func (r *Repository) PatchCriteria(ctx context.Context, patch domain.CompetencyUnitCriteriaPatch) (domain.CompetencyUnitCriteria, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanCriteria(db.QueryRow(ctx, `
		UPDATE competency_unit_criteria
		SET content = COALESCE($2::text, content),
		    sort_order = COALESCE($3::smallint, sort_order)
		WHERE id = $1
		RETURNING `+criteriaColumns,
		patch.ID, patch.Content, patch.SortOrder))
}

func (r *Repository) FindCriteriaByID(ctx context.Context, id string) (domain.CompetencyUnitCriteria, error) {
	db := database.GetExecutor(ctx, r.pool)
	return scanCriteria(db.QueryRow(ctx, `SELECT `+criteriaColumns+` FROM competency_unit_criteria WHERE id = $1`, id))
}

func (r *Repository) ListCriteriaByElement(ctx context.Context, elementID string, params domain.ListParams) (domain.CriteriaListResult, error) {
	db := database.GetExecutor(ctx, r.pool)
	rows, err := db.Query(ctx, `
		SELECT `+criteriaColumns+`, COUNT(*) OVER() AS total_count
		FROM competency_unit_criteria
		WHERE competency_unit_element_id = $1
		ORDER BY sort_order, created_at
		LIMIT $2 OFFSET $3
	`, elementID, params.Limit, params.Offset)
	if err != nil {
		return domain.CriteriaListResult{}, mapError(err)
	}
	defer rows.Close()

	items := make([]domain.CompetencyUnitCriteria, 0)
	total := 0
	for rows.Next() {
		item, rowTotal, err := scanCriteriaWithTotal(rows)
		if err != nil {
			return domain.CriteriaListResult{}, err
		}
		total = rowTotal
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return domain.CriteriaListResult{}, mapError(err)
	}
	if len(items) == 0 && params.Offset > 0 {
		if err := db.QueryRow(ctx, `SELECT COUNT(*) FROM competency_unit_criteria WHERE competency_unit_element_id = $1`, elementID).Scan(&total); err != nil {
			return domain.CriteriaListResult{}, mapError(err)
		}
	}
	return domain.CriteriaListResult{Items: items, Total: total}, nil
}

func (r *Repository) DeleteCriteria(ctx context.Context, id string) error {
	db := database.GetExecutor(ctx, r.pool)
	commandTag, err := db.Exec(ctx, `DELETE FROM competency_unit_criteria WHERE id = $1`, id)
	if err != nil {
		return mapError(err)
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// --- scan helpers ---

func scanSchema(row pgx.Row) (domain.CertificationSchema, error) {
	var s domain.CertificationSchema
	if err := row.Scan(&s.ID, &s.Code, &s.Name, &s.SchemaCode, &s.SortOrder, &s.IsActive, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return domain.CertificationSchema{}, mapError(err)
	}
	return s, nil
}

func scanSchemaWithTotal(rows pgx.Rows) (domain.CertificationSchema, int, error) {
	var s domain.CertificationSchema
	var total int
	if err := rows.Scan(&s.ID, &s.Code, &s.Name, &s.SchemaCode, &s.SortOrder, &s.IsActive, &s.CreatedAt, &s.UpdatedAt, &total); err != nil {
		return domain.CertificationSchema{}, 0, mapError(err)
	}
	return s, total, nil
}

func scanUnit(row pgx.Row) (domain.CompetencyUnit, error) {
	var u domain.CompetencyUnit
	if err := row.Scan(&u.ID, &u.CertificationSchemaID, &u.Code, &u.Title, &u.SortOrder, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return domain.CompetencyUnit{}, mapError(err)
	}
	return u, nil
}

func scanUnitWithTotal(rows pgx.Rows) (domain.CompetencyUnit, int, error) {
	var u domain.CompetencyUnit
	var total int
	if err := rows.Scan(&u.ID, &u.CertificationSchemaID, &u.Code, &u.Title, &u.SortOrder, &u.CreatedAt, &u.UpdatedAt, &total); err != nil {
		return domain.CompetencyUnit{}, 0, mapError(err)
	}
	return u, total, nil
}

func scanElement(row pgx.Row) (domain.CompetencyUnitElement, error) {
	var e domain.CompetencyUnitElement
	if err := row.Scan(&e.ID, &e.CompetencyUnitID, &e.Title, &e.SortOrder, &e.CreatedAt, &e.UpdatedAt); err != nil {
		return domain.CompetencyUnitElement{}, mapError(err)
	}
	return e, nil
}

func scanElementWithTotal(rows pgx.Rows) (domain.CompetencyUnitElement, int, error) {
	var e domain.CompetencyUnitElement
	var total int
	if err := rows.Scan(&e.ID, &e.CompetencyUnitID, &e.Title, &e.SortOrder, &e.CreatedAt, &e.UpdatedAt, &total); err != nil {
		return domain.CompetencyUnitElement{}, 0, mapError(err)
	}
	return e, total, nil
}

func scanCriteria(row pgx.Row) (domain.CompetencyUnitCriteria, error) {
	var c domain.CompetencyUnitCriteria
	if err := row.Scan(&c.ID, &c.CompetencyUnitElementID, &c.Content, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return domain.CompetencyUnitCriteria{}, mapError(err)
	}
	return c, nil
}

func scanCriteriaWithTotal(rows pgx.Rows) (domain.CompetencyUnitCriteria, int, error) {
	var c domain.CompetencyUnitCriteria
	var total int
	if err := rows.Scan(&c.ID, &c.CompetencyUnitElementID, &c.Content, &c.SortOrder, &c.CreatedAt, &c.UpdatedAt, &total); err != nil {
		return domain.CompetencyUnitCriteria{}, 0, mapError(err)
	}
	return c, total, nil
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
			// FK violation: parent (schema/unit/elemen) tidak ada.
			return domain.ErrNotFound
		}
	}
	return err
}
