import { useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { Loader2, Pencil, Plus, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'
import { RoleFormDialog } from './-components/role-form-dialog'
import { useDeleteRole, useRoles } from './-hooks'
import type { Role } from './-types'

export const Route = createFileRoute('/admin/roles/')({
  component: RolesPage,
})

function RolesPage() {
  const { auth } = Route.useRouteContext()
  const { data, isLoading, isError } = useRoles(auth)

  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<Role | null>(null)
  const [deleting, setDeleting] = useState<Role | null>(null)
  const remove = useDeleteRole(auth)

  function openCreate() {
    setEditing(null)
    setFormOpen(true)
  }

  function openEdit(role: Role) {
    setEditing(role)
    setFormOpen(true)
  }

  const items = data?.data ?? []

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-lg font-semibold text-slate-900">Master Role</h1>
          <p className="text-sm text-slate-500">Kelola daftar role pengguna.</p>
        </div>
        <Button onClick={openCreate} className="gap-2">
          <Plus className="h-4 w-4" />
          Tambah Role
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
            <p className="p-6 text-sm text-rose-600">Gagal memuat data role.</p>
          ) : items.length === 0 ? (
            <p className="p-6 text-sm text-slate-500">Belum ada data role.</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full min-w-[640px] text-sm">
                <thead className="border-b border-slate-200 bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500">
                  <tr>
                    <th className="px-4 py-3 font-medium">Kode</th>
                    <th className="px-4 py-3 font-medium">Nama</th>
                    <th className="px-4 py-3 font-medium">Deskripsi</th>
                    <th className="px-4 py-3 text-right font-medium">Aksi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {items.map((role) => (
                    <tr key={role.id} className="hover:bg-slate-50/70">
                      <td className="px-4 py-3">
                        <code className="rounded bg-slate-100 px-1.5 py-0.5 text-xs font-medium text-slate-700">{role.code}</code>
                      </td>
                      <td className="px-4 py-3 font-medium text-slate-900">{role.name}</td>
                      <td className="px-4 py-3 text-slate-600">{role.description || <span className="text-slate-400">—</span>}</td>
                      <td className="px-4 py-3">
                        <div className="flex items-center justify-end gap-1">
                          <Button variant="ghost" size="icon" aria-label={`Ubah ${role.name}`} onClick={() => openEdit(role)}>
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            aria-label={`Hapus ${role.name}`}
                            className="text-rose-600 hover:bg-rose-50"
                            onClick={() => setDeleting(role)}
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

      <RoleFormDialog open={formOpen} onOpenChange={setFormOpen} role={editing} session={auth} />

      <ConfirmDialog
        open={Boolean(deleting)}
        onOpenChange={(open) => (!open ? setDeleting(null) : undefined)}
        title="Hapus Role"
        description={`Yakin ingin menghapus role "${deleting?.name}" (${deleting?.code})? Role yang masih dipakai user tidak bisa dihapus.`}
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
