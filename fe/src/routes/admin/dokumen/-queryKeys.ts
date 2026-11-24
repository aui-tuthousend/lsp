import type { ListDocumentTypeParams } from './-types'

export const documentTypeKeys = {
  all: ['document-types'] as const,
  lists: () => [...documentTypeKeys.all, 'list'] as const,
  list: (params: ListDocumentTypeParams = {}) => [...documentTypeKeys.lists(), params] as const,
  details: () => [...documentTypeKeys.all, 'detail'] as const,
  detail: (id: string) => [...documentTypeKeys.details(), id] as const,
}
