import type { ListAsesiProfilesParams } from './-types'

export const asesiProfileKeys = {
  all: ['asesi-profiles'] as const,
  lists: () => [...asesiProfileKeys.all, 'list'] as const,
  list: (params: ListAsesiProfilesParams = {}) => [...asesiProfileKeys.lists(), params] as const,
  details: () => [...asesiProfileKeys.all, 'detail'] as const,
  detail: (id: string) => [...asesiProfileKeys.details(), id] as const,
  byUser: (userId: string) => [...asesiProfileKeys.details(), 'user', userId] as const,
}
