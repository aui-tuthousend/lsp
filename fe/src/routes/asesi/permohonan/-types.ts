// Tipe untuk form permohonan sertifikasi APL-01 (/asesi/permohonan).
//
// Field mengikuti DTO backend modul apl01 (migrasi 000009 & 000010). Tipe
// master yang dipakai saat memilih skema/dokumen diimpor dari modul admin
// terkait agar tidak duplikasi.

import type { DocumentType } from '../../admin/dokumen/-types'
import type { CertificationSchema, SchemaDocumentRequirement } from '../../admin/skema/-types'
import type { TUK } from '../../admin/tuks/-types'

export type { CertificationSchema, DocumentType, SchemaDocumentRequirement, TUK }

export type ApplicationStatus = 'DRAFT' | 'SUBMITTED' | 'APPROVED' | 'REJECTED'

export type AssessmentPurpose =
  | 'SERTIFIKASI'
  | 'SERTIFIKASI_ULANG'
  | 'REKOGNISI_PEMBELAJARAN_LAMPAU'
  | 'PENGAKUAN_KOMPETENSI_TERKINI'

export type AssessmentMode = 'ONLINE' | 'OFFLINE'

export type DocumentStatus = 'BELUM_LENGKAP' | 'MEMENUHI_SYARAT' | 'TIDAK_MEMENUHI_SYARAT'

export type SignerRole = 'asesi' | 'asesor' | 'asesor_ulang'

export type SignatureStatus = 'menunggu' | 'dipakai' | 'kadaluarsa' | 'dibatalkan'

/** Header permohonan APL-01. */
export type Application = {
  id: string
  user_id: string
  code: string
  status: ApplicationStatus
  current_step: number
  submitted_at?: string | null
  verified_by?: string | null
  verified_at?: string | null
  verification_note?: string | null
  asesor_id?: string | null
  asesor_assigned_at?: string | null
  created_at: string
  updated_at: string
}

/** Data sertifikasi (step 2): skema + tujuan asesmen. */
export type ApplicationData = {
  id: string
  application_id: string
  certification_schema_id: string
  purpose: AssessmentPurpose
  created_at: string
  updated_at: string
}

/** Bukti kelengkapan (step 3). */
export type ApplicationDocument = {
  id: string
  application_id: string
  document_type_id: string
  file_url?: string | null
  /** URL presigned (bucket privat) untuk pratinjau berkas. */
  file_view_url?: string | null
  status: DocumentStatus
  note?: string | null
  verified_at?: string | null
  created_at: string
  updated_at: string
}

/** Jadwal asesmen (step 4). */
export type AssessmentSchedule = {
  id: string
  application_id: string
  mode: AssessmentMode
  scheduled_date: string
  tuk_id?: string | null
  created_at: string
  updated_at: string
}

/** Token verifikasi QR. */
export type FormSignature = {
  id: string
  application_id: string
  entity: string
  role: SignerRole
  status: SignatureStatus
  expires_at: string
  used_at?: string | null
  revoked_at?: string | null
  created_at: string
  updated_at: string
}

/** Detail permohonan lengkap seluruh step (GET /apl01/{id}). */
export type ApplicationDetail = {
  application: Application
  data?: ApplicationData | null
  documents: ApplicationDocument[]
  schedule?: AssessmentSchedule | null
  signatures: FormSignature[]
}

export type ListApplicationsParams = {
  page?: number
  limit?: number
  status?: ApplicationStatus
}

/** Body PUT /apl01/{id}/data. */
export type PutDataBody = {
  certification_schema_id: string
  purpose: AssessmentPurpose
}

/** Satu baris bukti kelengkapan untuk PUT /apl01/{id}/documents. */
export type DocumentItemBody = {
  document_type_id: string
  file_url?: string | null
}

/** Body PUT /apl01/{id}/documents (set-semantik: daftar lengkap). */
export type PutDocumentsBody = {
  items: DocumentItemBody[]
}

/** Body PUT /apl01/{id}/schedule. */
export type PutScheduleBody = {
  mode: AssessmentMode
  scheduled_date: string
  tuk_id?: string | null
}

/** Hasil POST /apl01/{id}/submit — token mentah untuk QR (muncul sekali). */
export type SubmitResult = {
  application: ApplicationDetail
  verify_token: string
}

/** Status token verifikasi (GET /apl01/verify/{token}). */
export type VerifyTokenStatus = {
  status: SignatureStatus
  application_code?: string
  expires_at: string
  verified_at?: string | null
}
