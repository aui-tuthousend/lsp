import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { Session } from '@/lib/auth/types'
import { documentTypeKeys } from './-queryKeys'
import { documentTypeService } from './-service'
import type { CreateDocumentTypeBody, ListDocumentTypeParams, UpdateDocumentTypeBody } from './-types'

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : 'Silakan coba lagi'
}

/** Daftar jenis dokumen (paged). */
export function useDocumentTypes(params: ListDocumentTypeParams, session: Session | null) {
  return useQuery({
    queryKey: documentTypeKeys.list(params),
    queryFn: () => documentTypeService.list(params, session),
    enabled: Boolean(session),
  })
}

export function useCreateDocumentType(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (body: CreateDocumentTypeBody) => documentTypeService.create(body, session),
    onSuccess: (docType) => {
      void queryClient.invalidateQueries({ queryKey: documentTypeKeys.lists() })
      queryClient.setQueryData(documentTypeKeys.detail(docType.id), docType)
      toast.success('Jenis dokumen berhasil ditambahkan')
    },
    onError: (error) => {
      toast.error('Gagal menambahkan jenis dokumen', { description: errorMessage(error) })
    },
  })
}

export function useUpdateDocumentType(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, body }: { id: string; body: UpdateDocumentTypeBody }) => documentTypeService.update(id, body, session),
    onSuccess: (docType) => {
      void queryClient.invalidateQueries({ queryKey: documentTypeKeys.lists() })
      queryClient.setQueryData(documentTypeKeys.detail(docType.id), docType)
      toast.success('Jenis dokumen berhasil diperbarui')
    },
    onError: (error) => {
      toast.error('Gagal memperbarui jenis dokumen', { description: errorMessage(error) })
    },
  })
}

export function useDeleteDocumentType(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => documentTypeService.remove(id, session),
    onSuccess: (_result, id) => {
      void queryClient.invalidateQueries({ queryKey: documentTypeKeys.lists() })
      queryClient.removeQueries({ queryKey: documentTypeKeys.detail(id) })
      toast.success('Jenis dokumen berhasil dihapus')
    },
    onError: (error) => {
      toast.error('Gagal menghapus jenis dokumen', { description: errorMessage(error) })
    },
  })
}
