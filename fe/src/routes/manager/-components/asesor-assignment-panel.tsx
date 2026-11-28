import { useState } from 'react'
import { toast } from 'sonner'
import { UserCheck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import type { Session } from '@/lib/auth/types'
import { formatDateTime } from '../../asesi/permohonan/-helpers'
import type { Application, AsesorAssignment } from '../../asesi/permohonan/-types'
import { useAsesorCandidates, useAssignAsesor } from '../-hooks'
import { AssignmentStatusBadge } from './assignment-status-badge'

/**
 * Panel plotting asesor: manager memilih asesor yang menangani permohonan yang
 * sudah disetujui admin LSP. Kandidat dibatasi pada asesor yang terafiliasi TUK
 * jadwal asesi; asesor yang kapasitas hariannya penuh ditampilkan namun tidak
 * dapat dipilih. Penugasan disimpan dengan status awal PENDING.
 */
export function AsesorAssignmentPanel({
  application,
  assignment,
  session,
}: {
  application: Application
  assignment: AsesorAssignment | null
  session: Session
}) {
  const candidates = useAsesorCandidates(application.id, session)
  const assign = useAssignAsesor(application.id, session)
  const [selected, setSelected] = useState('')

  if (application.status !== 'APPROVED') {
    return (
      <Card>
        <CardContent className="p-4 text-sm text-slate-500">
          Permohonan belum disetujui admin LSP sehingga asesor belum dapat ditugaskan.
        </CardContent>
      </Card>
    )
  }

  const rows = candidates.data?.items ?? []
  const currentName =
    assignment?.asesor_name ||
    rows.find((row) => row.user_id === assignment?.asesor_id)?.name ||
    assignment?.asesor_id
  // Penugasan aktif (PENDING/ACCEPTED) mengunci plotting: manager tidak dapat
  // mengganti asesor sampai asesor sebelumnya menolak (status REJECTED).
  const hasActiveAssignment = assignment?.status === 'PENDING' || assignment?.status === 'ACCEPTED'

  function submit() {
    if (!selected) {
      toast.error('Pilih asesor terlebih dahulu')
      return
    }
    assign.mutate({ asesor_id: selected }, { onSuccess: () => setSelected('') })
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Penugasan Asesor</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {assignment ? (
          <div className="space-y-1 rounded-lg border border-slate-200 bg-slate-50 p-3">
            <p className="text-sm text-slate-700">
              Asesor saat ini: <span className="font-medium text-slate-900">{currentName}</span>
            </p>
            <div className="flex flex-wrap items-center gap-2">
              <AssignmentStatusBadge status={assignment.status} />
              <span className="text-xs text-slate-400">Ditugaskan {formatDateTime(assignment.assigned_at)}</span>
            </div>
            {assignment.status === 'REJECTED' && assignment.note ? (
              <p className="text-xs text-rose-700">Alasan: {assignment.note}</p>
            ) : null}
          </div>
        ) : (
          <p className="text-sm text-slate-500">Belum ada asesor yang ditugaskan untuk permohonan ini.</p>
        )}

        {hasActiveAssignment ? (
          <p className="text-sm text-slate-500">
            Asesor sudah ditugaskan. Penggantian asesor hanya dapat dilakukan setelah asesor
            sebelumnya menolak penugasan.
          </p>
        ) : (
          <>
            <div className="space-y-2">
              <p className="text-sm font-medium text-slate-700">Pilih asesor</p>

              {candidates.isLoading ? (
                <p className="text-sm text-slate-500">Memuat asesor...</p>
              ) : candidates.isError ? (
                <p className="text-sm text-rose-600">Gagal memuat daftar asesor.</p>
              ) : rows.length === 0 ? (
                <p className="text-sm text-slate-500">
                  Belum ada asesor terdaftar pada TUK jadwal permohonan ini.
                </p>
              ) : (
                <div className="overflow-x-auto rounded-lg border border-slate-200">
                  <table className="w-full text-sm">
                    <thead className="bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500">
                      <tr>
                        <th className="w-10 px-3 py-2" />
                        <th className="px-3 py-2 font-medium">Nama</th>
                        <th className="px-3 py-2 font-medium">Domisili</th>
                        <th className="px-3 py-2 font-medium">Penugasan</th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-slate-100">
                      {rows.map((row) => {
                        const domisili = [row.city, row.province].filter(Boolean).join(', ')
                        const countLabel =
                          row.daily_assignment_limit > 0
                            ? `${row.assigned_count}/${row.daily_assignment_limit}`
                            : `${row.assigned_count}`
                        return (
                          <tr
                            key={row.user_id}
                            className={row.disabled ? 'opacity-60' : 'cursor-pointer hover:bg-slate-50'}
                            onClick={() => {
                              if (!row.disabled) setSelected(row.user_id)
                            }}
                          >
                            <td className="px-3 py-2">
                              <input
                                type="radio"
                                name="asesor"
                                value={row.user_id}
                                checked={selected === row.user_id}
                                disabled={row.disabled}
                                onChange={() => setSelected(row.user_id)}
                              />
                            </td>
                            <td className="px-3 py-2 font-medium text-slate-900">{row.name}</td>
                            <td className="px-3 py-2 text-slate-600">{domisili || '—'}</td>
                            <td className="px-3 py-2 text-slate-700">
                              {countLabel}
                              {row.disabled ? (
                                <span className="ml-2 rounded bg-rose-50 px-1.5 py-0.5 text-xs text-rose-700">
                                  Penuh
                                </span>
                              ) : null}
                            </td>
                          </tr>
                        )
                      })}
                    </tbody>
                  </table>
                </div>
              )}
            </div>

            <Button className="gap-1.5" disabled={assign.isPending || !selected} onClick={submit}>
              <UserCheck className="h-4 w-4" />
              {assign.isPending ? 'Memproses...' : 'Tugaskan Asesor'}
            </Button>
          </>
        )}
      </CardContent>
    </Card>
  )
}
