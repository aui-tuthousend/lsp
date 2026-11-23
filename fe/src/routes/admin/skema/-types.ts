// Tipe untuk master data skema sertifikasi (/admin/skema).
//
// Field mengikuti DTO backend modul certificationschema:
// `SchemaDTO`, `UnitDTO`, `ElementDTO`, `CriteriaDTO` dan tree `SchemaTreeDTO`.
// Struktur: skema -> unit kompetensi -> elemen -> kriteria unjuk kerja (KUK).

/** Skema sertifikasi seperti dikembalikan backend. */
export type CertificationSchema = {
  id: string
  code: string
  name: string
  schema_code?: string | null
  sort_order: number
  is_active: boolean
  created_at: string
  updated_at: string
}

/** Body pembuatan skema (POST /certification-schemas). */
export type CreateSchemaBody = {
  code: string
  name: string
  schema_code?: string
  sort_order?: number
  is_active?: boolean
}

/** Body update parsial skema (PATCH /certification-schemas/{id}). */
export type UpdateSchemaBody = Partial<CreateSchemaBody>

export type ListSchemaParams = {
  page?: number
  limit?: number
}

/** Unit kompetensi milik sebuah skema. */
export type CompetencyUnit = {
  id: string
  certification_schema_id: string
  code: string
  title: string
  sort_order: number
  created_at: string
  updated_at: string
}

export type CreateUnitBody = {
  code: string
  title: string
  sort_order?: number
}

export type UpdateUnitBody = Partial<CreateUnitBody>

/** Elemen milik sebuah unit kompetensi. */
export type CompetencyUnitElement = {
  id: string
  competency_unit_id: string
  title?: string | null
  sort_order: number
  created_at: string
  updated_at: string
}

export type CreateElementBody = {
  title?: string
  sort_order?: number
}

export type UpdateElementBody = Partial<CreateElementBody>

/** Kriteria unjuk kerja (KUK) milik sebuah elemen. */
export type CompetencyUnitCriteria = {
  id: string
  competency_unit_element_id: string
  content: string
  sort_order: number
  created_at: string
  updated_at: string
}

export type CreateCriteriaBody = {
  content: string
  sort_order?: number
}

export type UpdateCriteriaBody = Partial<CreateCriteriaBody>

/** Node elemen beserta kriterianya (untuk tree). */
export type ElementNode = CompetencyUnitElement & {
  criteria: CompetencyUnitCriteria[]
}

/** Node unit beserta elemennya (untuk tree). */
export type UnitNode = CompetencyUnit & {
  elements: ElementNode[]
}

/** Skema beserta struktur lengkapnya (GET /certification-schemas/{id}). */
export type SchemaTree = CertificationSchema & {
  units: UnitNode[]
}
