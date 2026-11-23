import { useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { Loader2, Pencil, Plus, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'
import { TUKFormDialog } from './-components/tuk-form-dialog'
import { useDeleteTuk, useTuks } from './-hooks'
import type { TUK } from './-types'

export const Route = createFileRoute('/admin/tuks/')({
  component: TuksPage,
})

const PAGE_SIZE = 10

function TuksPage() {
  const { auth } = Route.useRouteContext()
  const [page, setPage] = useState(1)
  const { data, isLoading, isError } = useTuks({ page, limit: PAGE_SIZE }, auth)

  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<TUK | null>(null)
  const [deleting, setDeleting] = useState<TUK | null>(null)
  const remove = useDeleteTuk(auth)

  function openCreate() {
    setEditing(null)
    setFormOpen(true)
  }

  function openEdit(tuk: TUK) {
    setEditing(tuk)
    setFormOpen(true)
  }

  const items = data?.data ?? []
  const meta = data?.meta
  const totalPages = meta ? Math.max(1, Math.ceil(meta.total / meta.limit)) : 1

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-lg font-semibold text-slate-900">Master TUK</h1>
          <p className="text-sm text-slate-500">Kelola Tempat Uji Kompetensi.</p>
        </div>
        <Button onClick={openCreate} className="gap-2">
          <Plus className="h-4 w-4" />
          Tambah TUK
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
            <p className="p-6 text-sm text-rose-600">Gagal memuat data TUK.</p>
          ) : items.length === 0 ? (
            <p className="p-6 text-sm text-slate-500">Belum ada data TUK.</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full min-w-[720px] text-sm">
                <thead className="border-b border-slate-200 bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500">
                  <tr>
                    <th className="px-4 py-3 font-medium">Nama</th>
                    <th className="px-4 py-3 font-medium">Kontak</th>
                    <th className="px-4 py-3 font-medium">Wilayah</th>
                    <th className="px-4 py-3 text-right font-medium">Aksi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {items.map((tuk) => (
                    <tr key={tuk.id} className="hover:bg-slate-50/70">
                      <td className="px-4 py-3 font-medium text-slate-900">{tuk.name}</td>
                      <td className="px-4 py-3 text-slate-600">
                        <div className="truncate">{tuk.email}</div>
                        <div className="text-xs text-slate-400">{tuk.phone}</div>
                      </td>
                      <td className="px-4 py-3 text-slate-600">
                        {tuk.city}, {tuk.province}
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center justify-end gap-1">
                          <Button variant="ghost" size="icon" aria-label={`Ubah ${tuk.name}`} onClick={() => openEdit(tuk)}>
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            aria-label={`Hapus ${tuk.name}`}
                            className="text-rose-600 hover:bg-rose-50"
                            onClick={() => setDeleting(tuk)}
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

      <TUKFormDialog open={formOpen} onOpenChange={setFormOpen} tuk={editing} session={auth} />

      <ConfirmDialog
        open={Boolean(deleting)}
        onOpenChange={(open) => (!open ? setDeleting(null) : undefined)}
        title="Hapus TUK"
        description={`Yakin ingin menghapus "${deleting?.name}"? Tindakan ini tidak bisa dibatalkan.`}
        confirmLabel="Hapus"
        pendingLabel="Menghapus..."
        destructive
        pending={remove.isPending}
        onConfirm={() => {
          if (!deleting) return
          remove.mutate(deleting.id, { onSuccess: () => setDeleting(null) })
        }}
      />
    </div>
  )
}
