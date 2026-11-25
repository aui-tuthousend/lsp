import { cn } from '@/lib/utils'
import { STATUS_CLASSES, STATUS_LABELS } from '../-helpers'
import type { ApplicationStatus } from '../-types'

/** Badge status permohonan APL-01. */
export function ApplicationStatusBadge({ status, className }: { status: ApplicationStatus; className?: string }) {
  return (
    <span
      className={cn(
        'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
        STATUS_CLASSES[status],
        className,
      )}
    >
      {STATUS_LABELS[status]}
    </span>
  )
}
