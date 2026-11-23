import { useState } from 'react'
import { createFileRoute, Link } from '@tanstack/react-router'
import { Layers, Loader2, Pencil, Plus, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'
import { SchemaFormDialog } from './-components/schema-form-dialog'
import { useDeleteSchema, useSchemas } from './-hooks'
import type { CertificationSchema } from './-types'

export const Route = createFileRoute('/admin/skema/')({
  component: SkemaPage,
})

const PAGE_SIZE = 10

function SkemaPage() {
  const { auth } = Route.useRouteContext()
  const [page, setPage] = useState(1)
  const { data, isLoading, isError } = useSchemas({ page, limit: PAGE_SIZE }, auth)

  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<CertificationSchema | null>(null)
  const [deleting, setDeleting] = useState<CertificationSchema | null>(null)
  const remove = useDeleteSchema(auth)

  function openCreate() {
    setEditing(null)
    setFormOpen(true)
  }

  function openEdit(schema: CertificationSchema) {
    setEditing(schema)
    setFormOpen(true)
  }

  const items = data?.data ?? []
  const meta = data?.meta
  const totalPages = meta ? Math.max(1, Math.ceil(meta.total / meta.limit)) : 1

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-lg font-semibold text-slate-900">Master Skema Sertifikasi</h1>
          <p className="text-sm text-slate-500">Kelola skema beserta unit, elemen, dan kriteria unjuk kerja (KUK).</p>
        </div>
        <Button onClick={openCreate} className="gap-2">
          <Plus className="h-4 w-4" />
          Tambah Skema
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
            <p className="p-6 text-sm text-rose-600">Gagal memuat data skema sertifikasi.</p>
          ) : items.length === 0 ? (
            <p className="p-6 text-sm text-slate-500">Belum ada data skema sertifikasi.</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full min-w-[760px] text-sm">
                <thead className="border-b border-slate-200 bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500">
                  <tr>
                    <th className="px-4 py-3 font-medium">Kode</th>
                    <th className="px-4 py-3 font-medium">Nama</th>
                    <th className="px-4 py-3 font-medium">Nomor Skema</th>
                    <th className="px-4 py-3 font-medium">Urutan</th>
                    <th className="px-4 py-3 font-medium">Status</th>
                    <th className="px-4 py-3 text-right font-medium">Aksi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {items.map((schema) => (
                    <tr key={schema.id} className="hover:bg-slate-50/70">
                      <td className="px-4 py-3">
                        <code className="rounded bg-slate-100 px-1.5 py-0.5 text-xs font-medium text-slate-700">{schema.code}</code>
                      </td>
                      <td className="px-4 py-3 font-medium text-slate-900">{schema.name}</td>
                      <td className="px-4 py-3 text-slate-600">{schema.schema_code || <span className="text-slate-400">—</span>}</td>
                      <td className="px-4 py-3 text-slate-600">{schema.sort_order}</td>
                      <td className="px-4 py-3">
                        <span
                          className={
                            schema.is_active
                              ? 'rounded-full bg-emerald-50 px-2 py-0.5 text-xs font-medium text-emerald-700'
                              : 'rounded-full bg-slate-100 px-2 py-0.5 text-xs font-medium text-slate-500'
                          }
                        >
                          {schema.is_active ? 'Aktif' : 'Tidak aktif'}
                        </span>
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center justify-end gap-1">
                          <Link
                            to="/admin/skema/$schemaId"
                            params={{ schemaId: schema.id }}
                            className="inline-flex h-9 items-center gap-1.5 rounded-md px-3 text-sm font-medium text-slate-700 transition-colors hover:bg-slate-100"
                          >
                            <Layers className="h-4 w-4" />
                            Struktur
                          </Link>
                          <Button variant="ghost" size="icon" aria-label={`Ubah ${schema.name}`} onClick={() => openEdit(schema)}>
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            aria-label={`Hapus ${schema.name}`}
                            className="text-rose-600 hover:bg-rose-50"
                            onClick={() => setDeleting(schema)}
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

      <SchemaFormDialog open={formOpen} onOpenChange={setFormOpen} schema={editing} session={auth} />

      <ConfirmDialog
        open={Boolean(deleting)}
        onOpenChange={(open) => (!open ? setDeleting(null) : undefined)}
        title="Hapus Skema Sertifikasi"
        description={`Yakin ingin menghapus "${deleting?.name}"? Seluruh unit, elemen, dan KUK di dalamnya akan ikut terhapus.`}
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
