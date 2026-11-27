import type { Session } from '@/lib/auth/types'
import { apiFetch, type PagedResponse } from '@/lib/api'
import type { ApplicationDetail } from '../asesi/permohonan/-types'
import type { DecideBody, ListQueueParams, QueueItem } from './-types'

const APL01 = '/apl01'
const APL02 = '/apl02'

function queryString(params: ListQueueParams): string {
  const search = new URLSearchParams()
  if (params.page) search.set('page', String(params.page))
  if (params.limit) search.set('limit', String(params.limit))
  if (params.application_status) search.set('application_status', params.application_status)
  const value = search.toString()
  return value ? `?${value}` : ''
}

/**
 * Akses antrian verifikasi & keputusan admin LSP. Detail APL-01/APL-02 dan
 * data pendukung (skema, profil asesi) memakai ulang hook modul terkait.
 */
export const adminLspService = {
  /** Antrian asesmen mandiri (APL-02) milik seluruh asesi (paged + filter status). */
  listQueue(params: ListQueueParams, session: Session | null) {
    return apiFetch<PagedResponse<QueueItem>>(`${APL02}${queryString(params)}`, session)
  },

  /** Tetapkan keputusan admin atas sebuah permohonan (APPROVED/REJECTED + catatan). */
  decide(applicationId: string, body: DecideBody, session: Session | null) {
    return apiFetch<ApplicationDetail>(`${APL01}/${applicationId}/decision`, session, {
      method: 'POST',
      body: JSON.stringify(body),
    })
  },
}
