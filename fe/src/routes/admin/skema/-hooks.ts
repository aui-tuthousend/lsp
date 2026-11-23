import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import type { Session } from '@/lib/auth/types'
import { schemaKeys } from './-queryKeys'
import { schemaService } from './-service'
import type {
  CreateCriteriaBody,
  CreateElementBody,
  CreateSchemaBody,
  CreateUnitBody,
  ListSchemaParams,
  UpdateCriteriaBody,
  UpdateElementBody,
  UpdateSchemaBody,
  UpdateUnitBody,
} from './-types'

function errorMessage(error: unknown): string {
  return error instanceof Error ? error.message : 'Silakan coba lagi'
}

// --- Skema ---

/** Daftar skema sertifikasi (paged). */
export function useSchemas(params: ListSchemaParams, session: Session | null) {
  return useQuery({
    queryKey: schemaKeys.list(params),
    queryFn: () => schemaService.list(params, session),
    enabled: Boolean(session),
  })
}

/** Skema beserta struktur lengkapnya (unit -> elemen -> KUK). */
export function useSchemaTree(id: string, session: Session | null) {
  return useQuery({
    queryKey: schemaKeys.tree(id),
    queryFn: () => schemaService.tree(id, session),
    enabled: Boolean(session) && Boolean(id),
  })
}

export function useCreateSchema(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (body: CreateSchemaBody) => schemaService.create(body, session),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: schemaKeys.lists() })
      toast.success('Skema sertifikasi berhasil ditambahkan')
    },
    onError: (error) => {
      toast.error('Gagal menambahkan skema', { description: errorMessage(error) })
    },
  })
}

export function useUpdateSchema(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, body }: { id: string; body: UpdateSchemaBody }) => schemaService.update(id, body, session),
    onSuccess: (schema) => {
      void queryClient.invalidateQueries({ queryKey: schemaKeys.lists() })
      void queryClient.invalidateQueries({ queryKey: schemaKeys.tree(schema.id) })
      toast.success('Skema sertifikasi berhasil diperbarui')
    },
    onError: (error) => {
      toast.error('Gagal memperbarui skema', { description: errorMessage(error) })
    },
  })
}

export function useDeleteSchema(session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => schemaService.remove(id, session),
    onSuccess: (_result, id) => {
      void queryClient.invalidateQueries({ queryKey: schemaKeys.lists() })
      queryClient.removeQueries({ queryKey: schemaKeys.tree(id) })
      toast.success('Skema sertifikasi berhasil dihapus')
    },
    onError: (error) => {
      toast.error('Gagal menghapus skema', { description: errorMessage(error) })
    },
  })
}

// --- Unit kompetensi ---

export function useCreateUnit(schemaId: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (body: CreateUnitBody) => schemaService.createUnit(schemaId, body, session),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: schemaKeys.tree(schemaId) })
      toast.success('Unit kompetensi berhasil ditambahkan')
    },
    onError: (error) => {
      toast.error('Gagal menambahkan unit', { description: errorMessage(error) })
    },
  })
}

export function useUpdateUnit(schemaId: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, body }: { id: string; body: UpdateUnitBody }) => schemaService.updateUnit(id, body, session),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: schemaKeys.tree(schemaId) })
      toast.success('Unit kompetensi berhasil diperbarui')
    },
    onError: (error) => {
      toast.error('Gagal memperbarui unit', { description: errorMessage(error) })
    },
  })
}

export function useDeleteUnit(schemaId: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => schemaService.removeUnit(id, session),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: schemaKeys.tree(schemaId) })
      toast.success('Unit kompetensi berhasil dihapus')
    },
    onError: (error) => {
      toast.error('Gagal menghapus unit', { description: errorMessage(error) })
    },
  })
}

// --- Elemen unit kompetensi ---

export function useCreateElement(schemaId: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ unitId, body }: { unitId: string; body: CreateElementBody }) => schemaService.createElement(unitId, body, session),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: schemaKeys.tree(schemaId) })
      toast.success('Elemen berhasil ditambahkan')
    },
    onError: (error) => {
      toast.error('Gagal menambahkan elemen', { description: errorMessage(error) })
    },
  })
}

export function useUpdateElement(schemaId: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, body }: { id: string; body: UpdateElementBody }) => schemaService.updateElement(id, body, session),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: schemaKeys.tree(schemaId) })
      toast.success('Elemen berhasil diperbarui')
    },
    onError: (error) => {
      toast.error('Gagal memperbarui elemen', { description: errorMessage(error) })
    },
  })
}

export function useDeleteElement(schemaId: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => schemaService.removeElement(id, session),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: schemaKeys.tree(schemaId) })
      toast.success('Elemen berhasil dihapus')
    },
    onError: (error) => {
      toast.error('Gagal menghapus elemen', { description: errorMessage(error) })
    },
  })
}

// --- Kriteria unjuk kerja (KUK) ---

export function useCreateCriteria(schemaId: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ elementId, body }: { elementId: string; body: CreateCriteriaBody }) => schemaService.createCriteria(elementId, body, session),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: schemaKeys.tree(schemaId) })
      toast.success('KUK berhasil ditambahkan')
    },
    onError: (error) => {
      toast.error('Gagal menambahkan KUK', { description: errorMessage(error) })
    },
  })
}

export function useUpdateCriteria(schemaId: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, body }: { id: string; body: UpdateCriteriaBody }) => schemaService.updateCriteria(id, body, session),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: schemaKeys.tree(schemaId) })
      toast.success('KUK berhasil diperbarui')
    },
    onError: (error) => {
      toast.error('Gagal memperbarui KUK', { description: errorMessage(error) })
    },
  })
}

export function useDeleteCriteria(schemaId: string, session: Session | null) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: string) => schemaService.removeCriteria(id, session),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: schemaKeys.tree(schemaId) })
      toast.success('KUK berhasil dihapus')
    },
    onError: (error) => {
      toast.error('Gagal menghapus KUK', { description: errorMessage(error) })
    },
  })
}
