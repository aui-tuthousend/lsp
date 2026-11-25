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
  'apl01.verified': () => [['apl01']],
  'apl02.submitted': () => [['apl02']],
}

// ── Registry listener imperatif ──────────────────────────────────────────────
//
// Sebagian event butuh aksi di luar invalidasi query (mis. toast + redirect).
// Komponen berlangganan lewat registry ini, sementara `useNotifications`
// (satu koneksi SSE global) yang meneruskan event ke sini — jadi tidak perlu
// membuka koneksi SSE tambahan.

export type AppEventListener = (event: AppEvent) => void

const listeners = new Map<string, Set<AppEventListener>>()

/** Berlangganan satu tipe event. Mengembalikan fungsi unsubscribe. */
export function subscribeAppEvent(type: string, listener: AppEventListener): () => void {
  const set = listeners.get(type) ?? new Set<AppEventListener>()
  set.add(listener)
  listeners.set(type, set)
  return () => {
    set.delete(listener)
    if (set.size === 0) listeners.delete(type)
  }
}

/** Teruskan event ke seluruh listener yang berlangganan tipe tersebut. */
export function emitAppEvent(event: AppEvent): void {
  listeners.get(event.type)?.forEach((listener) => listener(event))
}
