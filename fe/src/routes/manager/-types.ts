// Tipe untuk laman manager sertifikasi (/manager).
//
// Dashboard menampilkan permohonan APL-01 yang sudah DISETUJUI admin LSP
// (read-only). Manager memilih asesor yang akan menangani permohonan tersebut.

import type { Assessment, AssessmentStatus } from '../asesi/asesmen-mandiri/-types'
import type { AsesorAssignmentStatus, ApplicationStatus } from '../asesi/permohonan/-types'

export type { AssessmentStatus, ApplicationStatus, AsesorAssignmentStatus }

/**
 * Satu baris antrian permohonan (APPROVED). Field `asesi_name`,
 * `application_code`, `application_status`, dan info asesor hanya terisi pada
 * daftar staff (hasil join backend).
 */
export type QueueItem = Assessment & {
  asesi_name?: string
  application_code?: string
  application_status?: ApplicationStatus
  asesor_id?: string | null
  asesor_name?: string
  asesor_assignment_status?: AsesorAssignmentStatus | null
}

export type ListQueueParams = {
  page?: number
  limit?: number
  application_status?: ApplicationStatus
  /** Filter penugasan asesor: true = sudah ditugaskan, false = belum. */
  assigned?: boolean
}

/** Body POST /apl01/{id}/assign-asesor. */
export type AssignAsesorBody = {
  asesor_id: string
}

/**
 * Kandidat asesor untuk plotting (GET /apl01/{id}/asesor-candidates): asesor
 * yang terafiliasi TUK jadwal asesi, beserta jumlah penugasan aktif pada
 * tanggal jadwal. `disabled` true bila kapasitas hariannya sudah penuh.
 */
export type AsesorCandidate = {
  user_id: string
  name: string
  province?: string
  city?: string
  /** 0 = tanpa batas */
  daily_assignment_limit: number
  /** Jumlah penugasan aktif (PENDING/ACCEPTED) pada tanggal jadwal */
  assigned_count: number
  disabled: boolean
}

export type AsesorCandidatesResponse = {
  items: AsesorCandidate[]
}
