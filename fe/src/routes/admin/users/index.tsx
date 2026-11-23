import { useState } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { Loader2, Pencil, Plus, Trash2 } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'
import { UserFormDialog } from './-components/user-form-dialog'
import { useDeleteUser, useUsers } from './-hooks'
import { MANAGED_ROLES, roleLabel } from './-schemas'
import type { User } from './-types'

export const Route = createFileRoute('/admin/users/')({
  component: UsersPage,
})

const PAGE_SIZE = 10
const MANAGED_ROLE_LIST = [...MANAGED_ROLES]

function UsersPage() {
  const { auth } = Route.useRouteContext()
  const [page, setPage] = useState(1)
  const { data, isLoading, isError } = useUsers({ page, limit: PAGE_SIZE, roles: MANAGED_ROLE_LIST }, auth)

  const [formOpen, setFormOpen] = useState(false)
  const [editing, setEditing] = useState<User | null>(null)
  const [deleting, setDeleting] = useState<User | null>(null)
  const remove = useDeleteUser(auth)

  function openCreate() {
    setEditing(null)
    setFormOpen(true)
  }

  function openEdit(user: User) {
    setEditing(user)
    setFormOpen(true)
  }

  const items = data?.data ?? []
  const meta = data?.meta
  const totalPages = meta ? Math.max(1, Math.ceil(meta.total / meta.limit)) : 1

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-lg font-semibold text-slate-900">Master User</h1>
          <p className="text-sm text-slate-500">Kelola akun validator &amp; manager sertifikasi.</p>
        </div>
        <Button onClick={openCreate} className="gap-2">
          <Plus className="h-4 w-4" />
          Tambah User
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
            <p className="p-6 text-sm text-rose-600">Gagal memuat data user.</p>
          ) : items.length === 0 ? (
            <p className="p-6 text-sm text-slate-500">Belum ada data user.</p>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full min-w-[720px] text-sm">
                <thead className="border-b border-slate-200 bg-slate-50 text-left text-xs uppercase tracking-wide text-slate-500">
                  <tr>
                    <th className="px-4 py-3 font-medium">Nama</th>
                    <th className="px-4 py-3 font-medium">Role</th>
                    <th className="px-4 py-3 font-medium">Kontak</th>
                    <th className="px-4 py-3 text-right font-medium">Aksi</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-100">
                  {items.map((user) => (
                    <tr key={user.id} className="hover:bg-slate-50/70">
                      <td className="px-4 py-3">
                        <div className="font-medium text-slate-900">{user.name}</div>
                        <div className="truncate text-xs text-slate-400">{user.email}</div>
                      </td>
                      <td className="px-4 py-3 text-slate-600">
                        <span className="inline-flex rounded-full bg-slate-100 px-2 py-0.5 text-xs font-medium text-slate-700">
                          {roleLabel(user.role)}
                        </span>
                      </td>
                      <td className="px-4 py-3 text-slate-600">
                        <div>{user.phone ?? <span className="text-slate-400">—</span>}</div>
                        {user.nomor_registrasi ? <div className="text-xs text-slate-400">{user.nomor_registrasi}</div> : null}
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center justify-end gap-1">
                          <Button variant="ghost" size="icon" aria-label={`Ubah ${user.name}`} onClick={() => openEdit(user)}>
                            <Pencil className="h-4 w-4" />
                          </Button>
                          <Button
                            variant="ghost"
                            size="icon"
                            aria-label={`Hapus ${user.name}`}
                            className="text-rose-600 hover:bg-rose-50"
                            onClick={() => setDeleting(user)}
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

      <UserFormDialog open={formOpen} onOpenChange={setFormOpen} user={editing} session={auth} />

      <ConfirmDialog
        open={Boolean(deleting)}
        onOpenChange={(open) => (!open ? setDeleting(null) : undefined)}
        title="Hapus User"
        description={`Yakin ingin menghapus akun "${deleting?.name}"? Tindakan ini tidak bisa dibatalkan.`}
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
