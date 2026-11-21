import { useEffect, useRef } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import type { Session } from '@/lib/auth/types'
import { eventInvalidation, type AppEvent } from './events'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

export function useNotifications(session: Session | null) {
  const queryClient = useQueryClient()
  const eventSourceRef = useRef<EventSource | null>(null)

  useEffect(() => {
    if (!session) return

    const url = `${API_BASE_URL}/events?token=${encodeURIComponent(session.token)}`
    const eventSource = new EventSource(url)
    eventSourceRef.current = eventSource

    eventSource.addEventListener('notification', (raw) => {
      const event = JSON.parse((raw as MessageEvent).data) as AppEvent
      if (event.type === 'connected') return

      const keys = eventInvalidation[event.type]?.(event) ?? []
      keys.forEach((queryKey) => queryClient.invalidateQueries({ queryKey }))
    })

    return () => {
      eventSource.close()
      eventSourceRef.current = null
    }
  }, [queryClient, session])
}
