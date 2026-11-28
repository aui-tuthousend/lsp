import type { AsesorAssignmentStatus } from '../../asesi/permohonan/-types'

const LABELS: Record<AsesorAssignmentStatus, { label: string; className: string }> = {
  PENDING: { label: 'Menunggu konfirmasi asesor', className: 'border-amber-200 bg-amber-50 text-amber-700' },
  ACCEPTED: { label: 'Diterima asesor', className: 'border-emerald-200 bg-emerald-50 text-emerald-700' },
  REJECTED: { label: 'Ditolak asesor', className: 'border-rose-200 bg-rose-50 text-rose-700' },
  CANCELLED: { label: 'Dibatalkan', className: 'border-slate-200 bg-slate-100 text-slate-600' },
}

/** Lencana status penugasan asesor. */
export function AssignmentStatusBadge({ status }: { status: AsesorAssignmentStatus }) {
  const item = LABELS[status]
  return (
    <span className={`inline-flex items-center rounded-full border px-2 py-0.5 text-xs font-medium ${item.className}`}>
      {item.label}
    </span>
  )
}
