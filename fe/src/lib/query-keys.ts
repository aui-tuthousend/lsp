export const qk = {
  me: ['me'] as const,
  tuks: (params?: object) => ['tuks', params] as const,
  asesorProfiles: (params?: object) => ['asesor-profiles', params] as const,
  assignments: (params?: object) => ['assignments', params] as const,
  assignment: (id: string) => ['assignments', id] as const,
  verification: (id: string) => ['verification', id] as const,
  notifications: ['notifications'] as const,
}
