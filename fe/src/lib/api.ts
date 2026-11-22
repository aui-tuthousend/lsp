import type { Session } from './auth/types'
import { API_BASE_URL } from './api-base'

export type PagedResponse<T> = {
  data: T[]
  meta: {
    page: number
    limit: number
    total: number
  }
}

export async function apiFetch<T>(path: string, session?: Session | null, init?: RequestInit): Promise<T> {
  const headers = new Headers(init?.headers)
  headers.set('Accept', 'application/json')
  if (init?.body && !headers.has('Content-Type')) headers.set('Content-Type', 'application/json')
  if (session?.token) headers.set('Authorization', `Bearer ${session.token}`)

  const res = await fetch(`${API_BASE_URL}${path}`, { ...init, headers })
  if (!res.ok) {
    const message = await res.text().catch(() => '')
    throw new Error(message || `Request failed with status ${res.status}`)
  }
  if (res.status === 204) return undefined as T
  return res.json() as Promise<T>
}
