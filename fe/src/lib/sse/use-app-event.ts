import { useEffect, useRef } from 'react'
import { subscribeAppEvent, type AppEvent } from './events'

/**
 * Berlangganan satu tipe event SSE untuk aksi imperatif (toast, redirect, dll).
 * Handler terbaru selalu dipakai tanpa perlu melepas/memasang ulang listener.
 */
export function useAppEvent(type: string, handler: (event: AppEvent) => void) {
  const handlerRef = useRef(handler)
  handlerRef.current = handler

  useEffect(() => {
    return subscribeAppEvent(type, (event) => handlerRef.current(event))
  }, [type])
}
