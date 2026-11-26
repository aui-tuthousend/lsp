import type { ListAssessmentsParams } from './-types'

export const apl02Keys = {
  all: ['apl02'] as const,
  lists: () => [...apl02Keys.all, 'list'] as const,
  list: (params: ListAssessmentsParams = {}) => [...apl02Keys.lists(), params] as const,
  details: () => [...apl02Keys.all, 'detail'] as const,
  detail: (id: string) => [...apl02Keys.details(), id] as const,
  schemaTree: (schemaId: string) => [...apl02Keys.all, 'schema-tree', schemaId] as const,
  verifyStatuses: () => [...apl02Keys.all, 'verify-status'] as const,
  verifyStatus: (token: string) => [...apl02Keys.verifyStatuses(), token] as const,
}
