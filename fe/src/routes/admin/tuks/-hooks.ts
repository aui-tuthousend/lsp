import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { Session } from '@/lib/auth/types'
import { tukKeys } from './-queryKeys'
import { tukService } from './-service'
import type { CreateTUKBody, ListTUKParams, UpdateTUKBody } from './-types'

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : 'Silakan coba lagi'
}

/** Daftar TUK (paged). */
export function useTuks(params: ListTUKParams, session: Session | null) {
  return useQuery({
    queryKey: tukKeys.list(params),
    queryFn: () => tukService.list(params, session),
    enabled: Boolean(session),
    // staleTime: 1000 * 60 * 5,
    // gcTime: 1000 * 60 * 5,
  })
}

export function useCreateTuk(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (body: CreateTUKBody) => tukService.create(body, session),
    onSuccess: (tuk) => {
      void queryClient.invalidateQueries({ queryKey: tukKeys.lists() })
      queryClient.setQueryData(tukKeys.detail(tuk.id), tuk)
      toast.success('TUK berhasil ditambahkan')
    },
    onError: (error) => {
      toast.error('Gagal menambahkan TUK', { description: errorMessage(error) })
    },
  })
}

export function useUpdateTuk(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, body }: { id: string; body: UpdateTUKBody }) => tukService.update(id, body, session),
    onSuccess: (tuk) => {
      void queryClient.invalidateQueries({ queryKey: tukKeys.lists() })
      queryClient.setQueryData(tukKeys.detail(tuk.id), tuk)
      toast.success('TUK berhasil diperbarui')
    },
    onError: (error) => {
      toast.error('Gagal memperbarui TUK', { description: errorMessage(error) })
    },
  })
}

export function useDeleteTuk(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => tukService.remove(id, session),
    onSuccess: (_result, id) => {
      void queryClient.invalidateQueries({ queryKey: tukKeys.lists() })
      queryClient.removeQueries({ queryKey: tukKeys.detail(id) })
      toast.success('TUK berhasil dihapus')
    },
    onError: (error) => {
      toast.error('Gagal menghapus TUK', { description: errorMessage(error) })
    },
  })
}
