import type { ListQueueParams } from './-types'

export const managerKeys = {
  all: ['manager'] as const,
  queues: () => [...managerKeys.all, 'queue'] as const,
  queue: (params: ListQueueParams = {}) => [...managerKeys.queues(), params] as const,
  asesorCandidates: (applicationId: string) => [...managerKeys.all, 'asesor-candidates', applicationId] as const,
}
