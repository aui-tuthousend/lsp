// Tipe untuk asesmen mandiri APL-02 (/asesi/asesmen-mandiri).
//
// Field mengikuti DTO backend modul apl02. Jawaban per unit kompetensi berupa
// boolean (kompeten/tidak) + daftar berkas bukti (object key hasil upload).

import type { ElementNode, SchemaTree } from '../../admin/skema/-types'
import type { SignatureStatus } from '../permohonan/-types'

export type { ElementNode, SchemaTree, SignatureStatus }

export type AssessmentStatus = 'DRAFT' | 'SUBMITTED' | 'REVIEWED'

export type AsesorDecision = 'RECOMMENDED' | 'NOT_RECOMMENDED'

/** Header asesmen mandiri (1:1 dengan permohonan APL-01). */
export type Assessment = {
  id: string
  application_id: string
  code: string
  status: AssessmentStatus
  asesor_note?: string | null
  asesor_decision?: AsesorDecision | null
  submitted_at?: string | null
  reviewed_by?: string | null
  reviewed_at?: string | null
  zoom_link?: string | null
  zoom_date?: string | null
  zoom_time?: string | null
  /** Skema yang dipilih (hanya terisi pada daftar). */
  schema_name?: string
  schema_code?: string
  created_at: string
  updated_at: string
}

/** Berkas bukti untuk sebuah jawaban unit. */
export type AssessmentUnitFile = {
  id: string
  /** Object key berkas (dipakai saat menyimpan ulang). */
  file_url: string
  /** URL presigned (bucket privat) untuk pratinjau berkas. */
  file_view_url?: string
  created_at: string
}

/** Jawaban asesi untuk satu unit kompetensi. */
export type AssessmentUnit = {
  id: string
  competency_unit_id: string
  unit_code: string
  unit_title: string
  is_competent: boolean | null
  files: AssessmentUnitFile[]
  created_at: string
  updated_at: string
}

/** Info permohonan/skema yang ditampilkan pada laman asesmen mandiri (read-only). */
export type ApplicationInfo = {
  certification_schema_id: string
  schema_name: string
  schema_code: string
  purpose: string
  data_created_at: string
}

/** Detail asesmen mandiri lengkap seluruh unit (GET /apl02/{id}). */
export type AssessmentDetail = {
  assessment: Assessment
  application: ApplicationInfo
  units: AssessmentUnit[]
}

export type ListAssessmentsParams = {
  page?: number
  limit?: number
}

/** Satu baris jawaban unit untuk PUT /apl02/{id}/units. */
export type UnitAnswerItemBody = {
  competency_unit_id: string
  is_competent?: boolean | null
  files?: string[]
}

/** Body PUT /apl02/{id}/units (set-semantik: daftar lengkap jawaban). */
export type PutUnitsBody = {
  items: UnitAnswerItemBody[]
}

/** Hasil POST /apl02/{id}/submit — token mentah untuk QR (muncul sekali). */
export type SubmitAssessmentResult = {
  assessment: AssessmentDetail
  verify_token: string
}

/** Status token verifikasi (GET /apl02/verify/{token}). */
export type AssessmentVerifyStatus = {
  status: SignatureStatus
  code?: string
  expires_at: string
  verified_at?: string | null
}
