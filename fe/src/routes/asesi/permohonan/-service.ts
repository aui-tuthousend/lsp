import type { Session } from '@/lib/auth/types'
import { apiFetch, type PagedResponse } from '@/lib/api'
import type {
  Application,
  ApplicationDetail,
  CertificationSchema,
  ListApplicationsParams,
  PutDataBody,
  PutDocumentsBody,
  PutScheduleBody,
  SchemaDocumentRequirement,
  SubmitResult,
  TUK,
  VerifyTokenStatus,
} from './-types'

const BASE = '/apl01'
const SCHEMAS = '/certification-schemas'

/** Respons daftar persyaratan dokumen: { data: SchemaDocumentRequirement[] }. */
export type RequirementsResponse = {
  data: SchemaDocumentRequirement[]
}

type SchemaListParams = {
  page?: number
  limit?: number
}

function queryString(params: ListApplicationsParams | SchemaListParams): string {
  const search = new URLSearchParams()
  if (params.page) search.set('page', String(params.page))
  if (params.limit) search.set('limit', String(params.limit))
  if ('status' in params && params.status) search.set('status', params.status)
  const value = search.toString()
  return value ? `?${value}` : ''
}

/** Akses data permohonan APL-01 (asesi) + master skema/dokumen yang dipakai form. */
export const apl01Service = {
  /** Buat draft baru atau lanjutkan draft aktif milik asesi. */
  createDraft(session: Session | null) {
    return apiFetch<ApplicationDetail>(BASE, session, { method: 'POST' })
  },

  /** Daftar permohonan milik asesi yang login (paged). */
  listMine(params: ListApplicationsParams, session: Session | null) {
    return apiFetch<PagedResponse<Application>>(`${BASE}/mine${queryString(params)}`, session)
  },

  /** Detail permohonan lengkap seluruh step. */
  get(id: string, session: Session | null) {
    return apiFetch<ApplicationDetail>(`${BASE}/${id}`, session)
  },

  /** Simpan step terakhir yang dikerjakan (resume draft). */
  updateStep(id: string, currentStep: number, session: Session | null) {
    return apiFetch<ApplicationDetail>(`${BASE}/${id}`, session, {
      method: 'PATCH',
      body: JSON.stringify({ current_step: currentStep }),
    })
  },

  /** Step 2 — simpan skema + tujuan asesmen. */
  putData(id: string, body: PutDataBody, session: Session | null) {
    return apiFetch<ApplicationDetail>(`${BASE}/${id}/data`, session, {
      method: 'PUT',
      body: JSON.stringify(body),
    })
  },

  /** Step 3 — ganti seluruh bukti kelengkapan. */
  putDocuments(id: string, body: PutDocumentsBody, session: Session | null) {
    return apiFetch<ApplicationDetail>(`${BASE}/${id}/documents`, session, {
      method: 'PUT',
      body: JSON.stringify(body),
    })
  },

  /** Step 4 — simpan jadwal asesmen. */
  putSchedule(id: string, body: PutScheduleBody, session: Session | null) {
    return apiFetch<ApplicationDetail>(`${BASE}/${id}/schedule`, session, {
      method: 'PUT',
      body: JSON.stringify(body),
    })
  },

  /** Submit permohonan + terbitkan token QR (token aktif yang dikirim dipakai ulang). */
  submit(id: string, verifyToken: string | null, session: Session | null) {
    return apiFetch<SubmitResult>(`${BASE}/${id}/submit`, session, {
      method: 'POST',
      body: JSON.stringify({ verify_token: verifyToken ?? undefined }),
    })
  },

  /** Konfirmasi hasil scan QR (token menentukan permohonan). */
  verify(token: string, session: Session | null) {
    return apiFetch<ApplicationDetail>(`${BASE}/verify`, session, {
      method: 'POST',
      body: JSON.stringify({ token }),
    })
  },

  /** Status token verifikasi tanpa memakainya (untuk laman verifikasi). */
  verifyStatus(token: string, session: Session | null) {
    return apiFetch<VerifyTokenStatus>(`${BASE}/verify/${token}`, session)
  },

  /** Master skema sertifikasi aktif (untuk step 2). */
  listSchemas(params: SchemaListParams, session: Session | null) {
    return apiFetch<PagedResponse<CertificationSchema>>(`${SCHEMAS}${queryString(params)}`, session)
  },

  /** Persyaratan dokumen sebuah skema (untuk step 3). */
  listRequirements(schemaId: string, session: Session | null) {
    return apiFetch<RequirementsResponse>(`${SCHEMAS}/${schemaId}/document-requirements`, session)
  },

  /** Master TUK (untuk pilihan lokasi asesmen offline di step 4). */
  listTuks(params: SchemaListParams, session: Session | null) {
    return apiFetch<PagedResponse<TUK>>(`/tuks${queryString(params)}`, session)
  },
}
