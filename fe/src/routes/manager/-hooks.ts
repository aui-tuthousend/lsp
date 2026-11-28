import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { Session } from '@/lib/auth/types'
import { apl01Keys } from '../asesi/permohonan/-queryKeys'
import { managerKeys } from './-queryKeys'
import { managerService } from './-service'
import type { AssignAsesorBody, ListQueueParams } from './-types'

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : 'Silakan coba lagi'
}

/** Antrian permohonan (APPROVED) untuk manager sertifikasi. */
export function useManagerQueue(params: ListQueueParams, session: Session | null) {
  return useQuery({
    queryKey: managerKeys.queue(params),
    queryFn: () => managerService.listQueue(params, session),
    enabled: Boolean(session),
  })
}

/** Kandidat asesor (matriks plotting) untuk sebuah permohonan. */
export function useAsesorCandidates(applicationId: string, session: Session | null) {
  return useQuery({
    queryKey: managerKeys.asesorCandidates(applicationId),
    queryFn: () => managerService.listAsesorCandidates(applicationId, session),
    enabled: Boolean(session) && Boolean(applicationId),
  })
}

/** Plotting asesor untuk sebuah permohonan. */
export function useAssignAsesor(applicationId: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (body: AssignAsesorBody) => managerService.assignAsesor(applicationId, body, session),
    onSuccess: (detail) => {
      queryClient.setQueryData(apl01Keys.detail(applicationId), detail)
      void queryClient.invalidateQueries({ queryKey: managerKeys.queues() })
      void queryClient.invalidateQueries({ queryKey: managerKeys.asesorCandidates(applicationId) })
      toast.success('Asesor berhasil ditugaskan')
    },
    onError: (error) => {
      toast.error('Gagal menugaskan asesor', { description: errorMessage(error) })
    },
  })
}
