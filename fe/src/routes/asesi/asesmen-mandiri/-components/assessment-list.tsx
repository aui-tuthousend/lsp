import { useNavigate } from '@tanstack/react-router'
import { ArrowRight, ClipboardCheck, Layers, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { formatDateTime } from '../-helpers'
import type { Assessment } from '../-types'
import { AssessmentStatusBadge } from './assessment-status-badge'

export type AssessmentListProps = {
  items: Assessment[]
  isLoading: boolean
  isError: boolean
  /** Pesan saat daftar kosong. */
  emptyMessage?: string
}

/** Daftar asesmen mandiri APL-02 milik asesi (dipakai di laman daftar). */
export function AssessmentList({ items, isLoading, isError, emptyMessage }: AssessmentListProps) {
  const navigate = useNavigate()

  if (isLoading) {
    return (
      <div className="flex items-center gap-2 py-6 text-sm text-slate-500">
        <Loader2 className="h-4 w-4 animate-spin" />
        Memuat asesmen mandiri...
      </div>
    )
  }

  if (isError) {
    return <p className="py-6 text-sm text-rose-600">Gagal memuat asesmen mandiri.</p>
  }

  if (items.length === 0) {
    return (
      <div className="flex flex-col items-center gap-2 py-8 text-center">
        <ClipboardCheck className="h-8 w-8 text-slate-300" />
        <p className="text-sm text-slate-500">
          {emptyMessage ?? 'Belum ada asesmen mandiri. Buka dari daftar permohonan yang sudah diajukan.'}
        </p>
      </div>
    )
  }

  return (
    <div className="grid gap-3">
      {items.map((item) => (
        <Card key={item.id}>
          <CardContent className="flex flex-wrap items-center justify-between gap-3 p-4">
            <div className="min-w-0">
              <div className="flex flex-wrap items-center gap-2">
                <span className="font-medium text-slate-900">{item.code}</span>
                <AssessmentStatusBadge status={item.status} />
              </div>
              {item.schema_name ? (
                <p className="mt-0.5 flex flex-wrap items-center gap-1.5 text-xs text-slate-600">
                  <Layers className="h-3.5 w-3.5 text-slate-400" />
                  <span className="font-medium">{item.schema_name}</span>
                  {item.schema_code ? <span className="text-slate-400">· {item.schema_code}</span> : null}
                </p>
              ) : null}
              <p className="mt-0.5 text-xs text-slate-500">Dibuat {formatDateTime(item.created_at)}</p>
            </div>
            <Button
              variant="outline"
              size="sm"
              className="gap-1.5"
              onClick={() =>
                navigate({ to: '/asesi/asesmen-mandiri/$assessmentId', params: { assessmentId: item.id } })
              }
            >
              {item.status === 'DRAFT' ? 'Lanjutkan' : 'Lihat'}
              <ArrowRight className="h-4 w-4" />
            </Button>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}
