import type { Session } from '@/lib/auth/types'
import { apiFetch, type PagedResponse } from '@/lib/api'
import type {
  AsesiOccupation,
  AsesiProfile,
  CreateAsesiProfileBody,
  ListAsesiProfilesParams,
  UpdateAccountBody,
  UpdateAsesiProfileBody,
  UpsertOccupationBody,
  UserAccount,
} from './-types'

const BASE = '/asesi-profiles'

function queryString(params: ListAsesiProfilesParams): string {
  const search = new URLSearchParams()
  if (params.page) search.set('page', String(params.page))
  if (params.limit) search.set('limit', String(params.limit))
  const value = search.toString()
  return value ? `?${value}` : ''
}

/**
 * Akses data profil asesi ke backend. Semua request memakai token sesi.
 * `session` boleh null; pemanggil (hooks) yang memastikan hanya jalan saat ada sesi.
 */
export const asesiProfileService = {
  list(params: ListAsesiProfilesParams, session: Session | null) {
    return apiFetch<PagedResponse<AsesiProfile>>(`${BASE}${queryString(params)}`, session)
  },

  getProfileById(id: string, session: Session | null) {
    return apiFetch<AsesiProfile>(`${BASE}/${id}`, session)
  },

  getProfileByUser(userId: string, session: Session | null) {
    return apiFetch<AsesiProfile>(`/users/${userId}/asesi-profile`, session)
  },

  /** Akun user saat ini (tabel `users`). */
  getMe(session: Session | null) {
    return apiFetch<UserAccount>('/auth/me', session)
  },

  // create(body: CreateAsesiProfileBody, session: Session | null) {
  //   return apiFetch<AsesiProfile>(BASE, session, { method: 'POST', body: JSON.stringify(body) })
  // },

  update(id: string, body: UpdateAsesiProfileBody, session: Session | null) {
    return apiFetch<AsesiProfile>(`${BASE}/${id}`, session, { method: 'PATCH', body: JSON.stringify(body) })
  },

  // remove(id: string, session: Session | null) {
  //   return apiFetch<void>(`${BASE}/${id}`, session, { method: 'DELETE' })
  // },

  /** Data pekerjaan asesi milik user (GET /users/{id}/asesi-occupation). */
  getOccupationByUser(userId: string, session: Session | null) {
    return apiFetch<AsesiOccupation>(`/users/${userId}/asesi-occupation`, session)
  },

  /** Isi/perbarui data pekerjaan milik user sendiri (PUT upsert). */
  upsertOccupationByUser(userId: string, body: UpsertOccupationBody, session: Session | null) {
    return apiFetch<AsesiOccupation>(`/users/${userId}/asesi-occupation`, session, { method: 'PUT', body: JSON.stringify(body) })
  },

  /** Ubah data akun sendiri (PATCH /auth/me). */
  updateMe(body: UpdateAccountBody, session: Session | null) {
    return apiFetch<UserAccount>('/auth/me', session, { method: 'PATCH', body: JSON.stringify(body) })
  },
}
