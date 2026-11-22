import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { Session } from '@/lib/auth/types'
import { qk } from '@/lib/query-keys'
import { asesiProfileKeys } from './-queryKeys'
import { asesiProfileService } from './-service'
import type { CreateAsesiProfileBody, ListAsesiProfilesParams, UpdateAsesiProfileBody } from './-types'

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : 'Silakan coba lagi'
}

/** Daftar profil asesi (paged). */
export function useAsesiProfiles(params: ListAsesiProfilesParams, session: Session | null) {
  return useQuery({
    queryKey: asesiProfileKeys.list(params),
    queryFn: () => asesiProfileService.list(params, session),
    enabled: Boolean(session),
  })
}

/** Detail profil asesi berdasarkan id profil. */
export function useAsesiProfileDetail(id: string, session: Session | null) {
  return useQuery({
    queryKey: asesiProfileKeys.detail(id),
    queryFn: () => asesiProfileService.getProfileById(id, session),
    enabled: Boolean(id && session),
  })
}

/** Profil asesi milik user tertentu — dipakai halaman profil sendiri (id = session.id). */
export function useAsesiProfileByUser(userId: string, session: Session | null) {
  return useQuery({
    queryKey: asesiProfileKeys.byUser(userId),
    queryFn: () => asesiProfileService.getProfileByUser(userId, session),
    enabled: Boolean(userId && session),
  })
}

/** Akun user saat ini (tabel `users`) — GET /auth/me. */
export function useMe(session: Session | null) {
  return useQuery({
    queryKey: qk.me,
    queryFn: () => asesiProfileService.getMe(session),
    enabled: Boolean(session),
  })
}

// export function useCreateAsesiProfile(session: Session | null) {
//   const queryClient = useQueryClient()

//   return useMutation({
//     mutationFn: (body: CreateAsesiProfileBody) => asesiProfileService.create(body, session),
//     onSuccess: (profile) => {
//       void queryClient.invalidateQueries({ queryKey: asesiProfileKeys.lists() })
//       queryClient.setQueryData(asesiProfileKeys.detail(profile.id), profile)
//       queryClient.setQueryData(asesiProfileKeys.byUser(profile.user_id), profile)
//       toast.success('Profil asesi dibuat')
//     },
//     onError: (error) => {
//       toast.error('Gagal membuat profil', { description: errorMessage(error) })
//     },
//   })
// }

export function useUpdateAsesiProfile(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, body }: { id: string; body: UpdateAsesiProfileBody }) =>
      asesiProfileService.update(id, body, session),
    onSuccess: (profile) => {
      void queryClient.invalidateQueries({ queryKey: asesiProfileKeys.lists() })
      queryClient.setQueryData(asesiProfileKeys.detail(profile.id), profile)
      queryClient.setQueryData(asesiProfileKeys.byUser(profile.user_id), profile)
      toast.success('Profil asesi diperbarui')
    },
    onError: (error) => {
      toast.error('Gagal memperbarui profil', { description: errorMessage(error) })
    },
  })
}

// export function useDeleteAsesiProfile(session: Session | null) {
//   const queryClient = useQueryClient()

//   return useMutation({
//     mutationFn: (id: string) => asesiProfileService.remove(id, session),
//     onSuccess: (_result, id) => {
//       void queryClient.invalidateQueries({ queryKey: asesiProfileKeys.lists() })
//       queryClient.removeQueries({ queryKey: asesiProfileKeys.detail(id) })
//       toast.success('Profil asesi dihapus')
//     },
//     onError: (error) => {
//       toast.error('Gagal menghapus profil', { description: errorMessage(error) })
//     },
//   })
// }
