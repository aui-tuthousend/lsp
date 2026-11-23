import type { ListSchemaParams } from './-types'

export const schemaKeys = {
  all: ['certification-schemas'] as const,
  lists: () => [...schemaKeys.all, 'list'] as const,
  list: (params: ListSchemaParams = {}) => [...schemaKeys.lists(), params] as const,
  trees: () => [...schemaKeys.all, 'tree'] as const,
  tree: (id: string) => [...schemaKeys.trees(), id] as const,
}
