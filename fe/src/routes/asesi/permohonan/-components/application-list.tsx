import { useNavigate } from '@tanstack/react-router'
import { AlertCircle, ArrowRight, ClipboardCheck, FileText, Loader2 } from 'lucide-react'
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
  /**
   * Bila diisi, permohonan yang sudah diajukan (SUBMITTED/APPROVED) menampilkan
   * tombol menuju asesmen mandiri APL-02.
   */
  onOpenAssessment?: (application: Application) => void
}

/** Daftar permohonan APL-01 milik asesi (dipakai di beranda). */
export function ApplicationList({ items, isLoading, isError, emptyMessage, onOpenAssessment }: ApplicationListProps) {
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
              {item.status === 'REJECTED' && item.verification_note ? (
                <p className="mt-2 flex items-start gap-1.5 rounded-md bg-rose-50 px-2.5 py-1.5 text-xs text-rose-700">
                  <AlertCircle className="mt-0.5 h-3.5 w-3.5 shrink-0" />
                  <span className="whitespace-pre-line">{item.verification_note}</span>
                </p>
              ) : null}
            </div>
            <div className="flex flex-wrap items-center gap-2">
              {onOpenAssessment && (item.status === 'SUBMITTED' || item.status === 'APPROVED') ? (
                <Button size="sm" className="gap-1.5" onClick={() => onOpenAssessment(item)}>
                  <ClipboardCheck className="h-4 w-4" />
                  Asesmen Mandiri
                </Button>
              ) : null}
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
            </div>
          </CardContent>
        </Card>
      ))}
    </div>
  )
}
