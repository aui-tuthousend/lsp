import type { ApplicationStatus, AssessmentMode, AssessmentPurpose, DocumentStatus } from './-types'

export const PURPOSE_OPTIONS: { value: AssessmentPurpose; label: string }[] = [
  { value: 'SERTIFIKASI', label: 'Sertifikasi' },
  { value: 'SERTIFIKASI_ULANG', label: 'Sertifikasi Ulang' },
  { value: 'REKOGNISI_PEMBELAJARAN_LAMPAU', label: 'Rekognisi Pembelajaran Lampau' },
  { value: 'PENGAKUAN_KOMPETENSI_TERKINI', label: 'Pengakuan Kompetensi Terkini' },
]

export const MODE_OPTIONS: { value: AssessmentMode; label: string }[] = [
  { value: 'ONLINE', label: 'Online' },
  { value: 'OFFLINE', label: 'Offline' },
]

export const STATUS_LABELS: Record<ApplicationStatus, string> = {
  DRAFT: 'Draft',
  SUBMITTED: 'Diajukan',
  APPROVED: 'Disetujui',
  REJECTED: 'Ditolak',
}

export const STATUS_CLASSES: Record<ApplicationStatus, string> = {
  DRAFT: 'bg-slate-100 text-slate-600',
  SUBMITTED: 'bg-amber-100 text-amber-700',
  APPROVED: 'bg-emerald-100 text-emerald-700',
  REJECTED: 'bg-rose-100 text-rose-700',
}

export const DOCUMENT_STATUS_LABELS: Record<DocumentStatus, string> = {
  BELUM_LENGKAP: 'Belum lengkap',
  MEMENUHI_SYARAT: 'Memenuhi syarat',
  TIDAK_MEMENUHI_SYARAT: 'Tidak memenuhi syarat',
}

export const DOCUMENT_STATUS_CLASSES: Record<DocumentStatus, string> = {
  BELUM_LENGKAP: 'bg-slate-100 text-slate-600',
  MEMENUHI_SYARAT: 'bg-emerald-100 text-emerald-700',
  TIDAK_MEMENUHI_SYARAT: 'bg-rose-100 text-rose-700',
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

/** Kelas kontrol bersama untuk textarea (input/select memakai komponen UI). */
export const textareaClass =
  'flex min-h-20 w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-slate-400 focus:ring-2 focus:ring-slate-200'
