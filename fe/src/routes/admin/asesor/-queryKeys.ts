import type { ListAsesorParams } from './-types'

export const asesorKeys = {
  all: ['asesor'] as const,
  lists: () => [...asesorKeys.all, 'list'] as const,
  list: (params: ListAsesorParams = {}) => [...asesorKeys.lists(), params] as const,
  details: () => [...asesorKeys.all, 'detail'] as const,
  detail: (userId: string) => [...asesorKeys.details(), userId] as const,
  tukOptions: () => [...asesorKeys.all, 'tuk-options'] as const,
}
