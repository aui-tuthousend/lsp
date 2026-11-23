import { useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { Loader2, Pencil, Plus, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'
import { AsesorFormDialog } from './-components/asesor-form-dialog'
import { useAsesorList, useDeleteAsesor } from './-hooks'
import type { AsesorProfile } from './-types'

export const Route = createFileRoute('/admin/asesor/')({
  component: AsesorPage,
})

const PAGE_SIZE = 10

function AsesorPage() {
  const { auth } = Route.useRouteContext()
  const [page, setPage] = useState(1)
  const { data, isLoading, isError } = useAsesorList({ page, limit: PAGE_SIZE }, auth)

  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<AsesorProfile | null>(null)
  const [deleting, setDeleting] = useState<AsesorProfile | null>(null)
  const remove = useDeleteAsesor(auth)

  function openCreate() {
    setEditing(null)
    setFormOpen(true)
  }

  function openEdit(asesor: AsesorProfile) {
    setEditing(asesor)
    setFormOpen(true)
  }

  const items = data?.data ?? []
  const meta = data?.meta
  const totalPages = meta ? Math.max(1, Math.ceil(meta.total / meta.limit)) : 1

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-lg font-semibold text-slate-900">Master Asesor</h1>
          <p className="text-sm text-slate-500">Kelola akun dan profil asesor.</p>
        </div>
        <Button onClick={openCreate} className="gap-2">
          <Plus className="h-4 w-4" />
          Tambah Asesor
        </Button>
      </div>

      <Card>
        <CardContent className="p-0">
          {isLoading ? (
            <div className="flex items-center gap-2 p-6 text-sm text-slate-500">
              <Loader2 className="h-4 w-4 animate-spin" />
              Memuat...
            </div>
          ) : isError ? (
            <p className="p-6 text-sm text-rose-600">Gagal memuat data asesor.</p>
          ) : items.length === 0 ? (
            <p className="p-6 text-sm text-slate-500">Belum ada data asesor.</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full min-w-[760px] text-sm">
                <thead className="border-b border-slate-200 bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500">
                  <tr>
                    <th className="px-4 py-3 font-medium">Nama</th>
                    <th className="px-4 py-3 font-medium">TUK</th>
                    <th className="px-4 py-3 font-medium">Domisili</th>
                    <th className="px-4 py-3 font-medium">Batas Harian</th>
                    <th className="px-4 py-3 text-right font-medium">Aksi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {items.map((asesor) => (
                    <tr key={asesor.user_id} className="hover:bg-slate-50/70">
                      <td className="px-4 py-3">
                        <div className="font-medium text-slate-900">{asesor.name}</div>
                        <div className="truncate text-xs text-slate-400">{asesor.email}</div>
                      </td>
                      <td className="px-4 py-3 text-slate-600">{asesor.tuk_name ?? <span className="text-slate-400">—</span>}</td>
                      <td className="px-4 py-3 text-slate-600">
                        {asesor.city || asesor.province ? `${asesor.city ?? ''}${asesor.city && asesor.province ? ', ' : ''}${asesor.province ?? ''}` : <span className="text-slate-400">—</span>}
                      </td>
                      <td className="px-4 py-3 text-slate-600">
                        {asesor.daily_assignment_limit === 0 ? 'Tanpa batas' : asesor.daily_assignment_limit}
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center justify-end gap-1">
                          <Button variant="ghost" size="icon" aria-label={`Ubah ${asesor.name}`} onClick={() => openEdit(asesor)}>
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            aria-label={`Hapus ${asesor.name}`}
                            className="text-rose-600 hover:bg-rose-50"
                            onClick={() => setDeleting(asesor)}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>

      {meta && totalPages > 1 ? (
        <div className="flex flex-wrap items-center justify-between gap-3">
          <p className="text-xs text-slate-500">
            Halaman {meta.page} dari {totalPages} · {meta.total} data
          </p>
          <div className="flex gap-2">
            <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage((prev) => Math.max(1, prev - 1))}>
              Sebelumnya
            </Button>
            <Button variant="outline" size="sm" disabled={page >= totalPages} onClick={() => setPage((prev) => prev + 1)}>
              Berikutnya
            </Button>
          </div>
        </div>
      ) : null}

      <AsesorFormDialog open={formOpen} onOpenChange={setFormOpen} asesor={editing} session={auth} />

      <ConfirmDialog
        open={Boolean(deleting)}
        onOpenChange={(open) => (!open ? setDeleting(null) : undefined)}
        title="Hapus Asesor"
        description={`Yakin ingin menghapus akun "${deleting?.name}"? Tindakan ini tidak bisa dibatalkan.`}
        confirmLabel="Hapus"
        pendingLabel="Menghapus..."
        destructive
        pending={remove.isPending}
        onConfirm={() => {
          if (!deleting) return
          remove.mutate(deleting.user_id, { onSuccess: () => setDeleting(null) })
        }}
      />
    </div>
  )
}
