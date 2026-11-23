import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { Session } from '@/lib/auth/types'
import { roleKeys } from './-queryKeys'
import { roleService } from './-service'
import type { CreateRoleBody, UpdateRoleBody } from './-types'

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : 'Silakan coba lagi'
}

/** Daftar master data role. */
export function useRoles(session: Session | null) {
  return useQuery({
    queryKey: roleKeys.list(),
    queryFn: () => roleService.list(session),
    enabled: Boolean(session),
  })
}

export function useCreateRole(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (body: CreateRoleBody) => roleService.create(body, session),
    onSuccess: (role) => {
      void queryClient.invalidateQueries({ queryKey: roleKeys.lists() })
      queryClient.setQueryData(roleKeys.detail(role.id), role)
      toast.success('Role berhasil ditambahkan')
    },
    onError: (error) => {
      toast.error('Gagal menambahkan role', { description: errorMessage(error) })
    },
  })
}

export function useUpdateRole(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, body }: { id: string; body: UpdateRoleBody }) => roleService.update(id, body, session),
    onSuccess: (role) => {
      void queryClient.invalidateQueries({ queryKey: roleKeys.lists() })
      queryClient.setQueryData(roleKeys.detail(role.id), role)
      toast.success('Role berhasil diperbarui')
    },
    onError: (error) => {
      toast.error('Gagal memperbarui role', { description: errorMessage(error) })
    },
  })
}

export function useDeleteRole(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => roleService.remove(id, session),
    onSuccess: (_result, id) => {
      void queryClient.invalidateQueries({ queryKey: roleKeys.lists() })
      queryClient.removeQueries({ queryKey: roleKeys.detail(id) })
      toast.success('Role berhasil dihapus')
    },
    onError: (error) => {
      toast.error('Gagal menghapus role', { description: errorMessage(error) })
    },
  })
}
