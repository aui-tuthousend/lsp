import type { Session } from '@/lib/auth/types'
import { apiFetch, type PagedResponse } from '@/lib/api'
import type { CreateTUKBody, ListTUKParams, TUK, UpdateTUKBody } from './-types'

const BASE = '/tuks'

function queryString(params: ListTUKParams): string {
  const search = new URLSearchParams()
  if (params.page) search.set('page', String(params.page))
  if (params.limit) search.set('limit', String(params.limit))
  const value = search.toString()
  return value ? `?${value}` : ''
}

/** Akses data master TUK. Semua request memakai token sesi (admin). */
export const tukService = {
  list(params: ListTUKParams, session: Session | null) {
    return apiFetch<PagedResponse<TUK>>(`${BASE}${queryString(params)}`, session)
  },

  get(id: string, session: Session | null) {
    return apiFetch<TUK>(`${BASE}/${id}`, session)
  },

  create(body: CreateTUKBody, session: Session | null) {
    return apiFetch<TUK>(BASE, session, { method: 'POST', body: JSON.stringify(body) })
  },

  update(id: string, body: UpdateTUKBody, session: Session | null) {
    return apiFetch<TUK>(`${BASE}/${id}`, session, { method: 'PATCH', body: JSON.stringify(body) })
  },

  remove(id: string, session: Session | null) {
    return apiFetch<void>(`${BASE}/${id}`, session, { method: 'DELETE' })
  },
}
