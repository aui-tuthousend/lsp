import type { ListApplicationsParams } from './-types'

export const apl01Keys = {
  all: ['apl01'] as const,
  lists: () => [...apl01Keys.all, 'list'] as const,
  list: (params: ListApplicationsParams = {}) => [...apl01Keys.lists(), params] as const,
  details: () => [...apl01Keys.all, 'detail'] as const,
  detail: (id: string) => [...apl01Keys.details(), id] as const,
  schemaOptions: () => [...apl01Keys.all, 'schema-options'] as const,
  requirements: (schemaId: string) => [...apl01Keys.all, 'requirements', schemaId] as const,
  tukOptions: () => [...apl01Keys.all, 'tuk-options'] as const,
}
