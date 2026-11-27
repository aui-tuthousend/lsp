// Tipe untuk laman verifikasi admin LSP (/admin_lsp).
//
// Dashboard menampilkan antrian pengajuan: asesmen mandiri (APL-02) yang sudah
// diajukan asesi (status SUBMITTED/REVIEWED) beserta ringkasan permohonan
// APL-01. Detail menampilkan data APL-01 & APL-02 secara read-only beserta
// panel keputusan admin (APPROVED/REJECTED).

import type { Assessment, AssessmentStatus } from '../asesi/asesmen-mandiri/-types'
import type { ApplicationStatus } from '../asesi/permohonan/-types'

export type { AssessmentStatus, ApplicationStatus }

/**
 * Satu baris antrian verifikasi. Field `asesi_name`, `application_code`, dan
 * `application_status` hanya terisi pada daftar staff (hasil join backend).
 */
export type QueueItem = Assessment & {
  asesi_name?: string
  application_code?: string
  application_status?: ApplicationStatus
}

export type ListQueueParams = {
  page?: number
  limit?: number
  /** Filter status permohonan APL-01 (keputusan admin: SUBMITTED/APPROVED/REJECTED). */
  application_status?: ApplicationStatus
}

/** Body POST /apl01/{id}/decision. */
export type DecideBody = {
  decision: 'APPROVED' | 'REJECTED'
  note?: string
}
