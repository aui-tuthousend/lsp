import { useMutation } from '@tanstack/react-query'
import { useRouter } from '@tanstack/react-router'
import { toast } from 'sonner'
import { homeForRole } from '@/lib/auth/types'
import { registerAsesi } from './-service'
import { useRegisterDraft } from './-store'
import type { CreateAsesiProfiles } from './-types'

/** Mutation pendaftaran asesi: kirim data, lalu login otomatis & arahkan ke beranda. */
export function useRegisterAsesi() {
  const router = useRouter()

  return useMutation({
    mutationFn: (values: CreateAsesiProfiles) => registerAsesi(values),
    onSuccess: async (session) => {
      useRegisterDraft.getState().reset()
      await router.invalidate()
      toast.success('Pendaftaran berhasil', { description: `Selamat datang, ${session.name}` })
      await router.navigate({ to: homeForRole(session.role), replace: true })
    },
    onError: (error) => {
      toast.error('Pendaftaran gagal', {
        description: error instanceof Error ? error.message : 'Silakan coba lagi',
      })
    },
  })
}
