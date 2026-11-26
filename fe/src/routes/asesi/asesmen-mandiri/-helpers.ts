import { PURPOSE_OPTIONS } from '../permohonan/-helpers'
import type { AsesorDecision, AssessmentStatus } from './-types'

export const ASSESSMENT_STATUS_LABELS: Record<AssessmentStatus, string> = {
  DRAFT: 'Draft',
  SUBMITTED: 'Diajukan',
  REVIEWED: 'Ditinjau',
}

export const ASSESSMENT_STATUS_CLASSES: Record<AssessmentStatus, string> = {
  DRAFT: 'bg-slate-100 text-slate-600',
  SUBMITTED: 'bg-amber-100 text-amber-700',
  REVIEWED: 'bg-emerald-100 text-emerald-700',
}

export const DECISION_LABELS: Record<AsesorDecision, string> = {
  RECOMMENDED: 'Direkomendasikan',
  NOT_RECOMMENDED: 'Tidak direkomendasikan',
}

const PURPOSE_LABELS = new Map(PURPOSE_OPTIONS.map((option) => [option.value, option.label]))

/** Label tujuan asesmen (mis. 'SERTIFIKASI' -> 'Sertifikasi'). */
export function purposeLabel(value?: string | null): string {
  if (!value) return '—'
  return PURPOSE_LABELS.get(value as (typeof PURPOSE_OPTIONS)[number]['value']) ?? value
}

/** Format tanggal (YYYY-MM-DD) menjadi tanggal Indonesia. */
export function formatDate(value?: string | null): string {
  if (!value) return '—'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '—'
  return new Intl.DateTimeFormat('id-ID', { dateStyle: 'long' }).format(date)
}

/** Format timestamp menjadi tanggal + jam Indonesia. */
export function formatDateTime(value?: string | null): string {
  if (!value) return '—'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '—'
  return new Intl.DateTimeFormat('id-ID', { dateStyle: 'medium', timeStyle: 'short' }).format(date)
}

// Token QR APL-02 disimpan di browser karena server hanya menyimpan hash-nya,
// sehingga token yang sama bisa ditampilkan lagi setelah refresh/re-submit.
const VERIFY_TOKEN_PREFIX = 'apl02:verify-token:'

/** Token QR terakhir yang diketahui browser untuk sebuah asesmen mandiri. */
export function readVerifyToken(assessmentId: string): string | null {
  if (typeof window === 'undefined') return null
  return window.localStorage.getItem(VERIFY_TOKEN_PREFIX + assessmentId)
}

export function writeVerifyToken(assessmentId: string, token: string): void {
  if (typeof window === 'undefined') return
  window.localStorage.setItem(VERIFY_TOKEN_PREFIX + assessmentId, token)
}
