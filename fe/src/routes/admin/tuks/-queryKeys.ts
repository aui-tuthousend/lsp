import type { ListTUKParams } from './-types'

export const tukKeys = {
  all: ['tuks'] as const,
  lists: () => [...tukKeys.all, 'list'] as const,
  list: (params: ListTUKParams = {}) => [...tukKeys.lists(), params] as const,
  details: () => [...tukKeys.all, 'detail'] as const,
  detail: (id: string) => [...tukKeys.details(), id] as const,
}
