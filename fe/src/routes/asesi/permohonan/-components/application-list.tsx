import { useNavigate } from '@tanstack/react-router'
import { ArrowRight, FileText, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { formatDateTime } from '../-helpers'
import type { Application } from '../-types'
import { ApplicationStatusBadge } from './application-status-badge'

export type ApplicationListProps = {
  items: Application[]
  isLoading: boolean
  isError: boolean
  /** Pesan saat daftar kosong (mis. saat filter status aktif). */
  emptyMessage?: string
}

/** Daftar permohonan APL-01 milik asesi (dipakai di beranda). */
export function ApplicationList({ items, isLoading, isError, emptyMessage }: ApplicationListProps) {
  const navigate = useNavigate()

  if (isLoading) {
    return (
      <div className="flex items-center gap-2 py-6 text-sm text-slate-500">
        <Loader2 className="h-4 w-4 animate-spin" />
        Memuat permohonan...
      </div>
    )
  }

  if (isError) {
    return <p className="py-6 text-sm text-rose-600">Gagal memuat permohonan.</p>
  }

  if (items.length === 0) {
    return (
      <div className="flex flex-col items-center gap-2 py-8 text-center">
        <FileText className="h-8 w-8 text-slate-300" />
        <p className="text-sm text-slate-500">
          {emptyMessage ?? 'Belum ada permohonan. Mulai dengan menekan tombol di atas.'}
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
                <ApplicationStatusBadge status={item.status} />
              </div>
              <p className="mt-0.5 text-xs text-slate-500">Dibuat {formatDateTime(item.created_at)}</p>
            </div>
            <Button
              variant="outline"
              size="sm"
              className="gap-1.5"
              onClick={() =>
                navigate({ to: '/asesi/permohonan/$applicationId', params: { applicationId: item.id } })
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
