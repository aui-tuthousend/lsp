import type { Session } from '@/lib/auth/types'
import { apiFetch, type PagedResponse } from '@/lib/api'
import type { CreateDocumentTypeBody, DocumentType, ListDocumentTypeParams, UpdateDocumentTypeBody } from './-types'

const BASE = '/document-types'

function queryString(params: ListDocumentTypeParams): string {
  const search = new URLSearchParams()
  if (params.page) search.set('page', String(params.page))
  if (params.limit) search.set('limit', String(params.limit))
  const value = search.toString()
  return value ? `?${value}` : ''
}

/** Akses data master jenis dokumen. Semua request memakai token sesi (admin). */
export const documentTypeService = {
  list(params: ListDocumentTypeParams, session: Session | null) {
    return apiFetch<PagedResponse<DocumentType>>(`${BASE}${queryString(params)}`, session)
  },

  get(id: string, session: Session | null) {
    return apiFetch<DocumentType>(`${BASE}/${id}`, session)
  },

  create(body: CreateDocumentTypeBody, session: Session | null) {
    return apiFetch<DocumentType>(BASE, session, { method: 'POST', body: JSON.stringify(body) })
  },

  update(id: string, body: UpdateDocumentTypeBody, session: Session | null) {
    return apiFetch<DocumentType>(`${BASE}/${id}`, session, { method: 'PATCH', body: JSON.stringify(body) })
  },

  remove(id: string, session: Session | null) {
    return apiFetch<void>(`${BASE}/${id}`, session, { method: 'DELETE' })
  },
}
