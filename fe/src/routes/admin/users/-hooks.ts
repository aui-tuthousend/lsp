import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { Session } from '@/lib/auth/types'
import { userKeys } from './-queryKeys'
import { userService } from './-service'
import type { CreateUserBody, ListUsersParams, UpdateUserBody } from './-types'

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : 'Silakan coba lagi'
}

/** Daftar user (paged, difilter per role). */
export function useUsers(params: ListUsersParams, session: Session | null) {
  return useQuery({
    queryKey: userKeys.list(params),
    queryFn: () => userService.list(params, session),
    enabled: Boolean(session),
  })
}

export function useCreateUser(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (body: CreateUserBody) => userService.create(body, session),
    onSuccess: (user) => {
      void queryClient.invalidateQueries({ queryKey: userKeys.lists() })
      queryClient.setQueryData(userKeys.detail(user.id), user)
      toast.success('User berhasil ditambahkan')
    },
    onError: (error) => {
      toast.error('Gagal menambahkan user', { description: errorMessage(error) })
    },
  })
}

export function useUpdateUser(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, body }: { id: string; body: UpdateUserBody }) => userService.update(id, body, session),
    onSuccess: (user) => {
      void queryClient.invalidateQueries({ queryKey: userKeys.lists() })
      queryClient.setQueryData(userKeys.detail(user.id), user)
      toast.success('User berhasil diperbarui')
    },
    onError: (error) => {
      toast.error('Gagal memperbarui user', { description: errorMessage(error) })
    },
  })
}

export function useDeleteUser(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => userService.remove(id, session),
    onSuccess: (_result, id) => {
      void queryClient.invalidateQueries({ queryKey: userKeys.lists() })
      queryClient.removeQueries({ queryKey: userKeys.detail(id) })
      toast.success('User berhasil dihapus')
    },
    onError: (error) => {
      toast.error('Gagal menghapus user', { description: errorMessage(error) })
    },
  })
}
