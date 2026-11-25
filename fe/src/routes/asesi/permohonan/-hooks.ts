import { useMutation, useQuery, useQueryClient, type QueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { Session } from '@/lib/auth/types'
import { apl01Keys } from './-queryKeys'
import { apl01Service } from './-service'
import type { ApplicationDetail, ListApplicationsParams, PutDataBody, PutDocumentsBody, PutScheduleBody } from './-types'

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : 'Silakan coba lagi'
}

/** Simpan detail ke cache + segarkan daftar (dipakai tiap mutasi). */
function cacheDetail(queryClient: QueryClient, detail: ApplicationDetail): void {
  queryClient.setQueryData(apl01Keys.detail(detail.application.id), detail)
  void queryClient.invalidateQueries({ queryKey: apl01Keys.lists() })
}

// --- Query ---

/** Daftar permohonan milik asesi yang login (paged). */
export function useMyApplications(params: ListApplicationsParams, session: Session | null) {
  return useQuery({
    queryKey: apl01Keys.list(params),
    queryFn: () => apl01Service.listMine(params, session),
    enabled: Boolean(session),
  })
}

/** Detail permohonan lengkap seluruh step. */
export function useApplication(id: string, session: Session | null) {
  return useQuery({
    queryKey: apl01Keys.detail(id),
    queryFn: () => apl01Service.get(id, session),
    enabled: Boolean(session) && Boolean(id),
  })
}

/** Master skema sertifikasi (untuk pilihan step 2). */
export function useSchemaOptions(session: Session | null) {
  return useQuery({
    queryKey: apl01Keys.schemaOptions(),
    queryFn: () => apl01Service.listSchemas({ page: 1, limit: 100 }, session),
    enabled: Boolean(session),
  })
}

/** Persyaratan dokumen sebuah skema (untuk step 3). */
export function useSchemaRequirements(schemaId: string, session: Session | null) {
  return useQuery({
    queryKey: apl01Keys.requirements(schemaId),
    queryFn: () => apl01Service.listRequirements(schemaId, session),
    enabled: Boolean(session) && Boolean(schemaId),
  })
}

/** Master TUK (untuk pilihan lokasi asesmen offline di step 4). */
export function useTukOptions(session: Session | null) {
  return useQuery({
    queryKey: apl01Keys.tukOptions(),
    queryFn: () => apl01Service.listTuks({ page: 1, limit: 100 }, session),
    enabled: Boolean(session),
  })
}

// --- Mutation ---

/** Buat draft baru atau lanjutkan draft aktif. */
export function useCreateDraft(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => apl01Service.createDraft(session),
    onSuccess: (detail) => {
      cacheDetail(queryClient, detail)
      toast.success('Permohonan baru dibuat')
    },
    onError: (error) => {
      toast.error('Gagal membuat permohonan', { description: errorMessage(error) })
    },
  })
}

/** Simpan step terakhir yang dikerjakan (resume draft). */
export function useUpdateStep(id: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (currentStep: number) => apl01Service.updateStep(id, currentStep, session),
    onSuccess: (detail) => cacheDetail(queryClient, detail),
    onError: (error) => {
      toast.error('Gagal menyimpan langkah', { description: errorMessage(error) })
    },
  })
}

/** Step 2 — simpan skema + tujuan asesmen. */
export function usePutData(id: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (body: PutDataBody) => apl01Service.putData(id, body, session),
    onSuccess: (detail) => cacheDetail(queryClient, detail),
    onError: (error) => {
      toast.error('Gagal menyimpan data sertifikasi', { description: errorMessage(error) })
    },
  })
}

/** Step 3 — ganti seluruh bukti kelengkapan. */
export function usePutDocuments(id: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (body: PutDocumentsBody) => apl01Service.putDocuments(id, body, session),
    onSuccess: (detail) => cacheDetail(queryClient, detail),
    onError: (error) => {
      toast.error('Gagal menyimpan dokumen', { description: errorMessage(error) })
    },
  })
}

/** Step 4 — simpan jadwal asesmen. */
export function usePutSchedule(id: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (body: PutScheduleBody) => apl01Service.putSchedule(id, body, session),
    onSuccess: (detail) => cacheDetail(queryClient, detail),
    onError: (error) => {
      toast.error('Gagal menyimpan jadwal', { description: errorMessage(error) })
    },
  })
}

/** Submit permohonan + terbitkan token QR. */
export function useSubmit(id: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => apl01Service.submit(id, session),
    onSuccess: (result) => cacheDetail(queryClient, result.application),
    onError: (error) => {
      toast.error('Gagal mengajukan permohonan', { description: errorMessage(error) })
    },
  })
}
