import type { Session } from '@/lib/auth/types'
import { apiFetch, type PagedResponse } from '@/lib/api'
import type {
  CertificationSchema,
  CompetencyUnit,
  CompetencyUnitCriteria,
  CompetencyUnitElement,
  CreateCriteriaBody,
  CreateElementBody,
  CreateSchemaBody,
  CreateUnitBody,
  ListSchemaParams,
  SchemaTree,
  UpdateCriteriaBody,
  UpdateElementBody,
  UpdateSchemaBody,
  UpdateUnitBody,
} from './-types'

const SCHEMAS = '/certification-schemas'
const UNITS = '/competency-units'
const ELEMENTS = '/competency-unit-elements'
const CRITERIA = '/competency-unit-criteria'

function queryString(params: ListSchemaParams): string {
  const search = new URLSearchParams()
  if (params.page) search.set('page', String(params.page))
  if (params.limit) search.set('limit', String(params.limit))
  const value = search.toString()
  return value ? `?${value}` : ''
}

/** Akses data master skema sertifikasi (skema, unit, elemen, KUK). Semua request admin. */
export const schemaService = {
  list(params: ListSchemaParams, session: Session | null) {
    return apiFetch<PagedResponse<CertificationSchema>>(`${SCHEMAS}${queryString(params)}`, session)
  },

  /** Ambil skema beserta struktur unit -> elemen -> kriteria. */
  tree(id: string, session: Session | null) {
    return apiFetch<SchemaTree>(`${SCHEMAS}/${id}`, session)
  },

  create(body: CreateSchemaBody, session: Session | null) {
    return apiFetch<CertificationSchema>(SCHEMAS, session, { method: 'POST', body: JSON.stringify(body) })
  },

  update(id: string, body: UpdateSchemaBody, session: Session | null) {
    return apiFetch<CertificationSchema>(`${SCHEMAS}/${id}`, session, { method: 'PATCH', body: JSON.stringify(body) })
  },

  remove(id: string, session: Session | null) {
    return apiFetch<void>(`${SCHEMAS}/${id}`, session, { method: 'DELETE' })
  },

  createUnit(schemaId: string, body: CreateUnitBody, session: Session | null) {
    return apiFetch<CompetencyUnit>(`${SCHEMAS}/${schemaId}/units`, session, { method: 'POST', body: JSON.stringify(body) })
  },

  updateUnit(id: string, body: UpdateUnitBody, session: Session | null) {
    return apiFetch<CompetencyUnit>(`${UNITS}/${id}`, session, { method: 'PATCH', body: JSON.stringify(body) })
  },

  removeUnit(id: string, session: Session | null) {
    return apiFetch<void>(`${UNITS}/${id}`, session, { method: 'DELETE' })
  },

  createElement(unitId: string, body: CreateElementBody, session: Session | null) {
    return apiFetch<CompetencyUnitElement>(`${UNITS}/${unitId}/elements`, session, { method: 'POST', body: JSON.stringify(body) })
  },

  updateElement(id: string, body: UpdateElementBody, session: Session | null) {
    return apiFetch<CompetencyUnitElement>(`${ELEMENTS}/${id}`, session, { method: 'PATCH', body: JSON.stringify(body) })
  },

  removeElement(id: string, session: Session | null) {
    return apiFetch<void>(`${ELEMENTS}/${id}`, session, { method: 'DELETE' })
  },

  createCriteria(elementId: string, body: CreateCriteriaBody, session: Session | null) {
    return apiFetch<CompetencyUnitCriteria>(`${ELEMENTS}/${elementId}/criteria`, session, { method: 'POST', body: JSON.stringify(body) })
  },

  updateCriteria(id: string, body: UpdateCriteriaBody, session: Session | null) {
    return apiFetch<CompetencyUnitCriteria>(`${CRITERIA}/${id}`, session, { method: 'PATCH', body: JSON.stringify(body) })
  },

  removeCriteria(id: string, session: Session | null) {
    return apiFetch<void>(`${CRITERIA}/${id}`, session, { method: 'DELETE' })
  },
}
