import type { Session } from '@/lib/auth/types'
import { apiFetch } from '@/lib/api'
import type { CreateRoleBody, Role, RoleListResponse, UpdateRoleBody } from './-types'

const BASE = '/roles'

/** Akses data master role. Semua request memakai token sesi (admin). */
export const roleService = {
  list(session: Session | null) {
    return apiFetch<RoleListResponse>(BASE, session)
  },

  get(id: string, session: Session | null) {
    return apiFetch<Role>(`${BASE}/${id}`, session)
  },

  create(body: CreateRoleBody, session: Session | null) {
    return apiFetch<Role>(BASE, session, { method: 'POST', body: JSON.stringify(body) })
  },

  update(id: string, body: UpdateRoleBody, session: Session | null) {
    return apiFetch<Role>(`${BASE}/${id}`, session, { method: 'PATCH', body: JSON.stringify(body) })
  },

  remove(id: string, session: Session | null) {
    return apiFetch<void>(`${BASE}/${id}`, session, { method: 'DELETE' })
  },
}
