import type { Session } from '@/lib/auth/types'
import { apiFetch, type PagedResponse } from '@/lib/api'
import type { ApplicationDetail } from '../asesi/permohonan/-types'
import type {
  AsesorCandidatesResponse,
  AssignAsesorBody,
  ListQueueParams,
  QueueItem,
} from './-types'

const APL01 = '/apl01'
const APL02 = '/apl02'

function queryString(params: ListQueueParams): string {
  const search = new URLSearchParams()
  if (params.page) search.set('page', String(params.page))
  if (params.limit) search.set('limit', String(params.limit))
  if (params.application_status) search.set('application_status', params.application_status)
  if (params.assigned !== undefined) search.set('assigned', String(params.assigned))
  const value = search.toString()
  return value ? `?${value}` : ''
}

/**
 * Akses data untuk laman manager sertifikasi: antrian permohonan APPROVED,
 * daftar asesor (untuk plotting), dan plotting asesor itu sendiri.
 */
export const managerService = {
  /** Antrian permohonan (asesmen mandiri) seluruh asesi (paged + filter status). */
  listQueue(params: ListQueueParams, session: Session | null) {
    return apiFetch<PagedResponse<QueueItem>>(`${APL02}${queryString(params)}`, session)
  },

  /** Kandidat asesor untuk plotting (dari TUK jadwal, + kapasitas harian). */
  listAsesorCandidates(applicationId: string, session: Session | null) {
    return apiFetch<AsesorCandidatesResponse>(`${APL01}/asesor-candidates/${applicationId}`, session)
  },

  /** Tugaskan asesor untuk sebuah permohonan (status awal PENDING). */
  assignAsesor(applicationId: string, body: AssignAsesorBody, session: Session | null) {
    return apiFetch<ApplicationDetail>(`${APL01}/${applicationId}/assign-asesor`, session, {
      method: 'POST',
      body: JSON.stringify(body),
    })
  },
}
