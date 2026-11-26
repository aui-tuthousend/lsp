import type { Session } from '@/lib/auth/types'
import { apiFetch, type PagedResponse } from '@/lib/api'
import type {
  Assessment,
  AssessmentDetail,
  AssessmentVerifyStatus,
  ListAssessmentsParams,
  PutUnitsBody,
  SchemaTree,
  SubmitAssessmentResult,
} from './-types'

const BASE = '/apl02'
const SCHEMAS = '/certification-schemas'

function queryString(params: ListAssessmentsParams): string {
  const search = new URLSearchParams()
  if (params.page) search.set('page', String(params.page))
  if (params.limit) search.set('limit', String(params.limit))
  const value = search.toString()
  return value ? `?${value}` : ''
}

/** Akses data asesmen mandiri APL-02 (asesi). */
export const apl02Service = {
  /** Buat atau lanjutkan asesmen mandiri untuk sebuah permohonan yang diajukan. */
  createOrGet(applicationId: string, session: Session | null) {
    return apiFetch<AssessmentDetail>(BASE, session, {
      method: 'POST',
      body: JSON.stringify({ application_id: applicationId }),
    })
  },

  /** Daftar asesmen mandiri milik asesi yang login (paged). */
  listMine(params: ListAssessmentsParams, session: Session | null) {
    return apiFetch<PagedResponse<Assessment>>(`${BASE}/mine${queryString(params)}`, session)
  },

  /** Detail asesmen mandiri lengkap seluruh unit. */
  get(id: string, session: Session | null) {
    return apiFetch<AssessmentDetail>(`${BASE}/${id}`, session)
  },

  /** Simpan seluruh jawaban unit (set-semantik) beserta berkas bukti. */
  putUnits(id: string, body: PutUnitsBody, session: Session | null) {
    return apiFetch<AssessmentDetail>(`${BASE}/${id}/units`, session, {
      method: 'PUT',
      body: JSON.stringify(body),
    })
  },

  /** Submit asesmen mandiri + terbitkan token QR (token aktif dipakai ulang). */
  submit(id: string, verifyToken: string | null, session: Session | null) {
    return apiFetch<SubmitAssessmentResult>(`${BASE}/${id}/submit`, session, {
      method: 'POST',
      body: JSON.stringify({ verify_token: verifyToken ?? undefined }),
    })
  },

  /** Konfirmasi hasil scan QR (token menentukan asesmen). */
  verify(token: string, session: Session | null) {
    return apiFetch<AssessmentDetail>(`${BASE}/verify`, session, {
      method: 'POST',
      body: JSON.stringify({ token }),
    })
  },

  /** Status token verifikasi tanpa memakainya (untuk laman verifikasi). */
  verifyStatus(token: string, session: Session | null) {
    return apiFetch<AssessmentVerifyStatus>(`${BASE}/verify/${token}`, session)
  },

  /** Struktur skema (unit -> elemen -> KUK) untuk menampilkan elemen & KUK. */
  schemaTree(schemaId: string, session: Session | null) {
    return apiFetch<SchemaTree>(`${SCHEMAS}/${schemaId}`, session)
  },
}
