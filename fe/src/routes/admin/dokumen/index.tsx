import { useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { Loader2, Pencil, Plus, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'
import { DocumentTypeFormDialog } from './-components/document-type-form-dialog'
import { useDeleteDocumentType, useDocumentTypes } from './-hooks'
import type { DocumentType } from './-types'

export const Route = createFileRoute('/admin/dokumen/')({
  component: DocumentTypesPage,
})

const PAGE_SIZE = 10

function DocumentTypesPage() {
  const { auth } = Route.useRouteContext()
  const [page, setPage] = useState(1)
  const { data, isLoading, isError } = useDocumentTypes({ page, limit: PAGE_SIZE }, auth)

  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<DocumentType | null>(null)
  const [deleting, setDeleting] = useState<DocumentType | null>(null)
  const remove = useDeleteDocumentType(auth)

  function openCreate() {
    setEditing(null)
    setFormOpen(true)
  }

  function openEdit(docType: DocumentType) {
    setEditing(docType)
    setFormOpen(true)
  }

  const items = data?.data ?? []
  const meta = data?.meta
  const totalPages = meta ? Math.max(1, Math.ceil(meta.total / meta.limit)) : 1

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-lg font-semibold text-slate-900">Master Jenis Dokumen</h1>
          <p className="text-sm text-slate-500">Kelola daftar jenis dokumen yang dapat dipersyaratkan pada skema.</p>
        </div>
        <Button onClick={openCreate} className="gap-2">
          <Plus className="h-4 w-4" />
          Tambah Jenis Dokumen
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
            <p className="p-6 text-sm text-rose-600">Gagal memuat data jenis dokumen.</p>
          ) : items.length === 0 ? (
            <p className="p-6 text-sm text-slate-500">Belum ada data jenis dokumen.</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full min-w-[640px] text-sm">
                <thead className="border-b border-slate-200 bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500">
                  <tr>
                    <th className="px-4 py-3 font-medium">Kode</th>
                    <th className="px-4 py-3 font-medium">Label</th>
                    <th className="px-4 py-3 text-right font-medium">Aksi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {items.map((docType) => (
                    <tr key={docType.id} className="hover:bg-slate-50/70">
                      <td className="px-4 py-3">
                        <code className="rounded bg-slate-100 px-1.5 py-0.5 text-xs font-medium text-slate-700">{docType.code}</code>
                      </td>
                      <td className="px-4 py-3 font-medium text-slate-900">{docType.label}</td>
                      <td className="px-4 py-3">
                        <div className="flex items-center justify-end gap-1">
                          <Button variant="ghost" size="icon" aria-label={`Ubah ${docType.label}`} onClick={() => openEdit(docType)}>
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            aria-label={`Hapus ${docType.label}`}
                            className="text-rose-600 hover:bg-rose-50"
                            onClick={() => setDeleting(docType)}
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

      <DocumentTypeFormDialog open={formOpen} onOpenChange={setFormOpen} docType={editing} session={auth} />

      <ConfirmDialog
        open={Boolean(deleting)}
        onOpenChange={(open) => (!open ? setDeleting(null) : undefined)}
        title="Hapus Jenis Dokumen"
        description={`Yakin ingin menghapus "${deleting?.label}"? Jenis dokumen yang masih dipakai skema tidak bisa dihapus.`}
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
