import type { ListQueueParams } from './-types'

export const adminLspKeys = {
  all: ['admin-lsp'] as const,
  queues: () => [...adminLspKeys.all, 'queue'] as const,
  queue: (params: ListQueueParams = {}) => [...adminLspKeys.queues(), params] as const,
}
