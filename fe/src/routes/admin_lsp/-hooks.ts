import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { Session } from '@/lib/auth/types'
import { apl01Keys } from '../asesi/permohonan/-queryKeys'
import { adminLspKeys } from './-queryKeys'
import { adminLspService } from './-service'
import type { DecideBody, ListQueueParams } from './-types'

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : 'Silakan coba lagi'
}

/** Antrian verifikasi (asesmen mandiri APL-02) untuk admin LSP (paged + filter status). */
export function useQueue(params: ListQueueParams, session: Session | null) {
  return useQuery({
    queryKey: adminLspKeys.queue(params),
    queryFn: () => adminLspService.listQueue(params, session),
    enabled: Boolean(session),
  })
}

/** Keputusan admin atas sebuah permohonan (APPROVED/REJECTED). */
export function useDecide(applicationId: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (body: DecideBody) => adminLspService.decide(applicationId, body, session),
    onSuccess: (detail) => {
      queryClient.setQueryData(apl01Keys.detail(applicationId), detail)
      void queryClient.invalidateQueries({ queryKey: adminLspKeys.queues() })
      toast.success(detail.application.status === 'APPROVED' ? 'Permohonan disetujui' : 'Permohonan ditolak')
    },
    onError: (error) => {
      toast.error('Gagal menyimpan keputusan', { description: errorMessage(error) })
    },
  })
}
