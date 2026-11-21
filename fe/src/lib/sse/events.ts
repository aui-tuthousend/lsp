export type AppEvent = {
  type: string
  topic: string
  data?: Record<string, unknown>
}

export const eventInvalidation: Record<string, (event: AppEvent) => readonly unknown[][]> = {
  'verification.scanned': (event) => [
    ['assignments'],
    ['verification', event.data?.session_id],
  ],
  'apl01.approved': () => [['apl01']],
  'apl02.submitted': () => [['apl02']],
}
