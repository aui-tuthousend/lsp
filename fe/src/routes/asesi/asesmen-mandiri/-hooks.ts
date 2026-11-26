import { useMutation, useQuery, useQueryClient, type QueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { Session } from '@/lib/auth/types'
import { apl02Keys } from './-queryKeys'
import { apl02Service } from './-service'
import type { AssessmentDetail, ListAssessmentsParams, PutUnitsBody } from './-types'

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : 'Silakan coba lagi'
}

/** Simpan detail ke cache + segarkan daftar (dipakai tiap mutasi). */
function cacheDetail(queryClient: QueryClient, detail: AssessmentDetail): void {
  queryClient.setQueryData(apl02Keys.detail(detail.assessment.id), detail)
  void queryClient.invalidateQueries({ queryKey: apl02Keys.lists() })
}

// --- Query ---

/** Daftar asesmen mandiri milik asesi yang login (paged). */
export function useMyAssessments(params: ListAssessmentsParams, session: Session | null) {
  return useQuery({
    queryKey: apl02Keys.list(params),
    queryFn: () => apl02Service.listMine(params, session),
    enabled: Boolean(session),
  })
}

/** Detail asesmen mandiri lengkap seluruh unit. */
export function useAssessment(id: string, session: Session | null) {
  return useQuery({
    queryKey: apl02Keys.detail(id),
    queryFn: () => apl02Service.get(id, session),
    enabled: Boolean(session) && Boolean(id),
  })
}

/** Struktur skema (unit -> elemen -> KUK) untuk menampilkan elemen & KUK. */
export function useSchemaTree(schemaId: string, session: Session | null) {
  return useQuery({
    queryKey: apl02Keys.schemaTree(schemaId),
    queryFn: () => apl02Service.schemaTree(schemaId, session),
    enabled: Boolean(session) && Boolean(schemaId),
  })
}

/** Status token QR asesmen mandiri (menunggu/dipakai/kadaluarsa/dibatalkan). */
export function useAssessmentVerifyStatus(token: string | null, session: Session | null) {
  return useQuery({
    queryKey: apl02Keys.verifyStatus(token ?? ''),
    queryFn: () => apl02Service.verifyStatus(token as string, session),
    enabled: Boolean(session) && Boolean(token),
  })
}

// --- Mutation ---

/** Buat atau lanjutkan asesmen mandiri untuk sebuah permohonan. */
export function useCreateOrGetAssessment(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (applicationId: string) => apl02Service.createOrGet(applicationId, session),
    onSuccess: (detail) => cacheDetail(queryClient, detail),
    onError: (error) => {
      toast.error('Gagal membuka asesmen mandiri', { description: errorMessage(error) })
    },
  })
}

/** Simpan seluruh jawaban unit (set-semantik). */
export function usePutUnits(id: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (body: PutUnitsBody) => apl02Service.putUnits(id, body, session),
    onSuccess: (detail) => {
      cacheDetail(queryClient, detail)
      toast.success('Jawaban disimpan')
    },
    onError: (error) => {
      toast.error('Gagal menyimpan jawaban', { description: errorMessage(error) })
    },
  })
}

/** Submit asesmen mandiri + terbitkan token QR (token aktif dipakai ulang). */
export function useSubmitAssessment(id: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (verifyToken: string | null) => apl02Service.submit(id, verifyToken, session),
    onSuccess: (result) => cacheDetail(queryClient, result.assessment),
    onError: (error) => {
      toast.error('Gagal mengajukan asesmen mandiri', { description: errorMessage(error) })
    },
  })
}

/** Konfirmasi hasil scan QR APL-02 (dipakai laman verifikasi). */
export function useVerifyAssessment(session: Session | null) {
  return useMutation({
    mutationFn: (token: string) => apl02Service.verify(token, session),
  })
}
