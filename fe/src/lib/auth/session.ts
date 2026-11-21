import { createServerFn } from '@tanstack/react-start'
import { deleteCookie, getCookie, setCookie } from '@tanstack/react-start/server'
import type { Role, Session } from './types'

const COOKIE = 'lsp_session'
const API_BASE_URL = process.env.API_BASE_URL ?? 'http://localhost:8080'

type BackendUser = {
  id: string
  role: Role
  name: string
}

type LoginResponse = {
  token: string
  user: BackendUser
}

export const getSessionFn = createServerFn({ method: 'GET' }).handler(async (): Promise<Session | null> => {
  const raw = getCookie(COOKIE)
  if (!raw) return null
  try {
    return JSON.parse(raw) as Session
  } catch {
    return null
  }
})

export const loginFn = createServerFn({ method: 'POST' })
  .validator((data: { email: string; password: string }) => data)
  .handler(async ({ data }): Promise<Session> => {
    const res = await fetch(`${API_BASE_URL}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    })
    if (!res.ok) throw new Error('login failed')
    const { token, user } = (await res.json()) as LoginResponse
    const session: Session = { id: user.id, role: user.role, name: user.name, token }
    setCookie(COOKIE, JSON.stringify(session), {
      httpOnly: true,
      sameSite: 'lax',
      path: '/',
      secure: process.env.NODE_ENV === 'production',
    })
    return session
  })

export const logoutFn = createServerFn({ method: 'POST' }).handler(async () => {
  deleteCookie(COOKIE, { path: '/' })
})
