import { cn } from '@/lib/utils'
import { ASSESSMENT_STATUS_CLASSES, ASSESSMENT_STATUS_LABELS } from '../-helpers'
import type { AssessmentStatus } from '../-types'

/** Badge status asesmen mandiri APL-02. */
export function AssessmentStatusBadge({ status, className }: { status: AssessmentStatus; className?: string }) {
  return (
    <span
      className={cn(
        'inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium',
        ASSESSMENT_STATUS_CLASSES[status],
        className,
      )}
    >
      {ASSESSMENT_STATUS_LABELS[status]}
    </span>
  )
}
