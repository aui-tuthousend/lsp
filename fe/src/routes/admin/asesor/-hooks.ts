import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { Session } from '@/lib/auth/types'
import { asesorKeys } from './-queryKeys'
import { asesorService } from './-service'
import type {
  CreateAsesorBody,
  ListAsesorParams,
  UpdateAsesorAccountBody,
  UpdateAsesorProfileBody,
} from './-types'

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : 'Silakan coba lagi'
}

/** Daftar asesor (paged, memuat profil + data akun). */
export function useAsesorList(params: ListAsesorParams, session: Session | null) {
  return useQuery({
    queryKey: asesorKeys.list(params),
    queryFn: () => asesorService.list(params, session),
    enabled: Boolean(session),
  })
}

/** Opsi TUK untuk dropdown afiliasi. */
export function useTukOptions(session: Session | null) {
  return useQuery({
    queryKey: asesorKeys.tukOptions(),
    queryFn: () => asesorService.listTuks(session),
    enabled: Boolean(session),
  })
}

/** Akun user asesor (untuk prefill form edit). */
export function useAsesorAccount(userId: string, session: Session | null) {
  return useQuery({
    queryKey: asesorKeys.detail(userId),
    queryFn: () => asesorService.getAccount(userId, session),
    enabled: Boolean(userId && session),
  })
}

/** Buat akun asesor lalu terapkan profil (afiliasi/domisili/batas harian). */
export function useCreateAsesor(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ account, profile }: { account: CreateAsesorBody; profile: UpdateAsesorProfileBody }) => {
      const created = await asesorService.createAccount(account, session)
      if (Object.keys(profile).length > 0) {
        await asesorService.updateProfile(created.id, profile, session)
      }
      return created
    },
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: asesorKeys.lists() })
      toast.success('Asesor berhasil ditambahkan')
    },
    onError: (error) => {
      toast.error('Gagal menambahkan asesor', { description: errorMessage(error) })
    },
  })
}

/** Update parsial akun dan/atau profil asesor. */
export function useUpdateAsesor(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({
      userId,
      account,
      profile,
    }: {
      userId: string
      account: UpdateAsesorAccountBody
      profile: UpdateAsesorProfileBody
    }) => {
      if (Object.keys(account).length > 0) {
        await asesorService.updateAccount(userId, account, session)
      }
      if (Object.keys(profile).length > 0) {
        await asesorService.updateProfile(userId, profile, session)
      }
    },
    onSuccess: (_result, variables) => {
      void queryClient.invalidateQueries({ queryKey: asesorKeys.lists() })
      queryClient.removeQueries({ queryKey: asesorKeys.detail(variables.userId) })
      toast.success('Asesor berhasil diperbarui')
    },
    onError: (error) => {
      toast.error('Gagal memperbarui asesor', { description: errorMessage(error) })
    },
  })
}

/** Hapus (soft delete) akun asesor. */
export function useDeleteAsesor(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (userId: string) => asesorService.remove(userId, session),
    onSuccess: (_result, userId) => {
      void queryClient.invalidateQueries({ queryKey: asesorKeys.lists() })
      queryClient.removeQueries({ queryKey: asesorKeys.detail(userId) })
      toast.success('Asesor berhasil dihapus')
    },
    onError: (error) => {
      toast.error('Gagal menghapus asesor', { description: errorMessage(error) })
    },
  })
}
