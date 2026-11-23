import type { Session } from '@/lib/auth/types'
import { apiFetch, type PagedResponse } from '@/lib/api'
import type {
  AsesorAccount,
  AsesorProfile,
  CreateAsesorBody,
  ListAsesorParams,
  TukOption,
  UpdateAsesorAccountBody,
  UpdateAsesorProfileBody,
} from './-types'

const BASE = '/asesor-profiles'
const USERS = '/users'

function queryString(params: ListAsesorParams): string {
  const search = new URLSearchParams()
  if (params.page) search.set('page', String(params.page))
  if (params.limit) search.set('limit', String(params.limit))
  const value = search.toString()
  return value ? `?${value}` : ''
}

/** Akses data asesor (akun user + profil asesor). Semua request memakai token sesi (admin). */
export const asesorService = {
  list(params: ListAsesorParams, session: Session | null) {
    return apiFetch<PagedResponse<AsesorProfile>>(`${BASE}${queryString(params)}`, session)
  },

  getProfile(userId: string, session: Session | null) {
    return apiFetch<AsesorProfile>(`${BASE}/${userId}`, session)
  },

  /** Akun user asesor (tabel `users`) untuk prefill form edit. */
  getAccount(id: string, session: Session | null) {
    return apiFetch<AsesorAccount>(`${USERS}/${id}`, session)
  },

  /** Opsi TUK untuk dropdown afiliasi. */
  listTuks(session: Session | null) {
    return apiFetch<PagedResponse<TukOption>>('/tuks?limit=100', session)
  },

  /** Buat akun user ber-role asesor (profil default dibuat otomatis oleh backend). */
  createAccount(body: CreateAsesorBody, session: Session | null) {
    return apiFetch<AsesorAccount>(USERS, session, {
      method: 'POST',
      body: JSON.stringify({ ...body, role: 'asesor' }),
    })
  },

  /** Update parsial akun user. */
  updateAccount(id: string, body: UpdateAsesorAccountBody, session: Session | null) {
    return apiFetch<AsesorAccount>(`${USERS}/${id}`, session, { method: 'PATCH', body: JSON.stringify(body) })
  },

  /** Update parsial profil asesor. */
  updateProfile(userId: string, body: UpdateAsesorProfileBody, session: Session | null) {
    return apiFetch<AsesorProfile>(`${BASE}/${userId}`, session, { method: 'PATCH', body: JSON.stringify(body) })
  },

  /** Hapus (soft delete) akun asesor. */
  remove(id: string, session: Session | null) {
    return apiFetch<void>(`${USERS}/${id}`, session, { method: 'DELETE' })
  },
}
