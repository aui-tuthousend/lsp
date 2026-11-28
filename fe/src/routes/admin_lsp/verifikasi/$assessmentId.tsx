import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { toast } from 'sonner'
import { Check, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import type { Session } from '@/lib/auth/types'
import { ApplicationStatusBadge } from '../../asesi/permohonan/-components/application-status-badge'
import { formatDateTime, textareaClass } from '../../asesi/permohonan/-helpers'
import type { Application } from '../../asesi/permohonan/-types'
import { ApplicationReview } from '../-components/application-review'
import { useDecide } from '../-hooks'

export const Route = createFileRoute('/admin_lsp/verifikasi/$assessmentId')({ component: VerifyPage })

function VerifyPage() {
  const { assessmentId } = Route.useParams()
  const { auth } = Route.useRouteContext()
  const session = auth!

  return (
    <ApplicationReview
      assessmentId={assessmentId}
      session={session}
      backTo="/admin_lsp"
      backLabel="Kembali ke daftar"
      title="Verifikasi Pengajuan"
      footer={({ applicationDetail }) => (
        <DecisionPanel application={applicationDetail.application} session={session} />
      )}
    />
  )
}

function DecisionPanel({ application, session }: { application: Application; session: Session }) {
  const decide = useDecide(application.id, session)
  const [approveOpen, setApproveOpen] = useState(false)
  const [rejectOpen, setRejectOpen] = useState(false)
  const [note, setNote] = useState('')

  if (application.status === 'APPROVED' || application.status === 'REJECTED') {
    return (
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0">
          <CardTitle>Keputusan Verifikasi</CardTitle>
          <ApplicationStatusBadge status={application.status} />
        </CardHeader>
        <CardContent className="space-y-2">
          {application.verification_note ? (
            <p className="whitespace-pre-line text-sm text-slate-600">{application.verification_note}</p>
          ) : (
            <p className="text-sm text-slate-400">Tanpa catatan.</p>
          )}
          {application.verified_at ? (
            <p className="text-xs text-slate-400">Diputuskan {formatDateTime(application.verified_at)}</p>
          ) : null}
        </CardContent>
      </Card>
    )
  }

  if (application.status !== 'SUBMITTED') {
    return (
      <Card>
        <CardContent className="p-4 text-sm text-slate-500">
          Permohonan belum diajukan sehingga belum dapat diputuskan.
        </CardContent>
      </Card>
    )
  }

  function submit(decision: 'APPROVED' | 'REJECTED') {
    const trimmed = note.trim()
    if (decision === 'REJECTED' && trimmed === '') {
      toast.error('Catatan wajib diisi untuk penolakan')
      return
    }
    decide.mutate(
      { decision, note: trimmed || undefined },
      {
        onSuccess: () => {
          setApproveOpen(false)
          setRejectOpen(false)
          setNote('')
        },
      },
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Keputusan Verifikasi</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <p className="text-sm text-slate-500">
          Pastikan seluruh data APL-01 &amp; APL-02 sudah sesuai sebelum mengambil keputusan.
        </p>
        <div className="flex flex-wrap gap-2">
          <Button
            className="gap-1.5"
            disabled={decide.isPending}
            onClick={() => {
              setNote('')
              setApproveOpen(true)
            }}
          >
            <Check className="h-4 w-4" />
            Setujui
          </Button>
          <Button
            variant="outline"
            className="gap-1.5 text-rose-600 hover:bg-rose-50"
            disabled={decide.isPending}
            onClick={() => {
              setNote('')
              setRejectOpen(true)
            }}
          >
            <X className="h-4 w-4" />
            Tolak
          </Button>
        </div>
      </CardContent>

      <Dialog open={approveOpen} onOpenChange={setApproveOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Setujui Pengajuan</DialogTitle>
            <DialogDescription>
              Permohonan akan disetujui. Anda dapat menambahkan catatan (opsional).
            </DialogDescription>
          </DialogHeader>
          <textarea
            className={textareaClass}
            value={note}
            onChange={(event) => setNote(event.target.value)}
            placeholder="Catatan (opsional)"
          />
          <DialogFooter>
            <Button variant="outline" onClick={() => setApproveOpen(false)}>
              Batal
            </Button>
            <Button className="gap-1.5" disabled={decide.isPending} onClick={() => submit('APPROVED')}>
              <Check className="h-4 w-4" />
              {decide.isPending ? 'Memproses...' : 'Setujui'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={rejectOpen} onOpenChange={setRejectOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Tolak Pengajuan</DialogTitle>
            <DialogDescription>
              Permohonan akan ditolak. Catatan wajib diisi sebagai alasan penolakan.
            </DialogDescription>
          </DialogHeader>
          <textarea
            className={textareaClass}
            value={note}
            onChange={(event) => setNote(event.target.value)}
            placeholder="Alasan penolakan (wajib)"
          />
          <DialogFooter>
            <Button variant="outline" onClick={() => setRejectOpen(false)}>
              Batal
            </Button>
            <Button
              className="gap-1.5 bg-rose-600 hover:bg-rose-700"
              disabled={decide.isPending}
              onClick={() => submit('REJECTED')}
            >
              <X className="h-4 w-4" />
              {decide.isPending ? 'Memproses...' : 'Tolak'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Card>
  )
}
