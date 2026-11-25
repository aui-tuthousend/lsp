import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { ApiError } from '@/lib/api'
import type { Session } from '@/lib/auth/types'
import { qk } from '@/lib/query-keys'
import { uploadFile } from '@/lib/upload'
import { asesiProfileKeys } from './-queryKeys'
import { asesiProfileService } from './-service'
import type { CreateAsesiProfileBody, ListAsesiProfilesParams, UpdateAccountBody, UpdateAsesiProfileBody, UpsertOccupationBody } from './-types'

/** Field berkas pendukung pada profil asesi (tabel `asesi_profiles`). */
export type ProfileFileField = 'photo' | 'ktp_file' | 'cv_file' | 'npwp_file'

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

/**
 * Ganti satu berkas pendukung: unggah berkas baru ke object storage lalu PATCH
 * profil. Backend otomatis menghapus objek lama saat nilainya berubah.
 */
export function useReplaceProfileFile(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: async ({ id, field, file }: { id: string; field: ProfileFileField; file: File }) => {
      const objectKey = await uploadFile(file, field)
      const body: UpdateAsesiProfileBody = {}
      body[field] = objectKey
      return asesiProfileService.update(id, body, session)
    },
    onSuccess: (profile) => {
      void queryClient.invalidateQueries({ queryKey: asesiProfileKeys.lists() })
      queryClient.setQueryData(asesiProfileKeys.detail(profile.id), profile)
      queryClient.setQueryData(asesiProfileKeys.byUser(profile.user_id), profile)
      toast.success('Berkas diperbarui')
    },
    onError: (error) => {
      toast.error('Gagal memperbarui berkas', { description: errorMessage(error) })
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

/** Data pekerjaan asesi milik user; null bila belum pernah diisi (404). */
export function useOccupationByUser(userId: string, session: Session | null) {
  return useQuery({
    queryKey: asesiProfileKeys.occupationByUser(userId),
    queryFn: async () => {
      try {
        return await asesiProfileService.getOccupationByUser(userId, session)
      } catch (error) {
        if (error instanceof ApiError && error.status === 404) return null
        throw error
      }
    },
    enabled: Boolean(userId && session),
  })
}

/** Isi/perbarui data pekerjaan milik user sendiri (upsert). */
export function useUpsertOccupation(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ userId, body }: { userId: string; body: UpsertOccupationBody }) =>
      asesiProfileService.upsertOccupationByUser(userId, body, session),
    onSuccess: (_occupation, variables) => {
      void queryClient.invalidateQueries({ queryKey: asesiProfileKeys.occupationByUser(variables.userId) })
      toast.success('Data pekerjaan disimpan')
    },
    onError: (error) => {
      toast.error('Gagal menyimpan data pekerjaan', { description: errorMessage(error) })
    },
  })
}

/** Ubah data akun sendiri (nama & telepon). */
export function useUpdateMe(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (body: UpdateAccountBody) => asesiProfileService.updateMe(body, session),
    onSuccess: (account) => {
      queryClient.setQueryData(qk.me, account)
      toast.success('Data akun diperbarui')
    },
    onError: (error) => {
      toast.error('Gagal memperbarui data akun', { description: errorMessage(error) })
    },
  })
}
