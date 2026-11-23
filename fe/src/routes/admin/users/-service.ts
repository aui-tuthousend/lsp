import type { Session } from '@/lib/auth/types'
import { apiFetch, type PagedResponse } from '@/lib/api'
import type { CreateUserBody, ListUsersParams, UpdateUserBody, User } from './-types'

const BASE = '/users'

function queryString(params: ListUsersParams): string {
  const search = new URLSearchParams()
  if (params.page) search.set('page', String(params.page))
  if (params.limit) search.set('limit', String(params.limit))
  for (const role of params.roles ?? []) search.append('roles', role)
  const value = search.toString()
  return value ? `?${value}` : ''
}

/** Akses data user. Semua request memakai token sesi (admin). */
export const userService = {
  list(params: ListUsersParams, session: Session | null) {
    return apiFetch<PagedResponse<User>>(`${BASE}${queryString(params)}`, session)
  },

  get(id: string, session: Session | null) {
    return apiFetch<User>(`${BASE}/${id}`, session)
  },

  create(body: CreateUserBody, session: Session | null) {
    return apiFetch<User>(BASE, session, { method: 'POST', body: JSON.stringify(body) })
  },

  update(id: string, body: UpdateUserBody, session: Session | null) {
    return apiFetch<User>(`${BASE}/${id}`, session, { method: 'PATCH', body: JSON.stringify(body) })
  },

  remove(id: string, session: Session | null) {
    return apiFetch<void>(`${BASE}/${id}`, session, { method: 'DELETE' })
  },
}
