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

type AuthResponse = {
  token: string
  user: BackendUser
}

/** Payload pendaftaran asesi (data akun + profil), sesuai body `POST /auth/register`. */
export type RegisterPayload = {
  email: string
  password: string
  full_name: string
  phone: string
  nik: string
  gender: string
  birth_place: string
  birth_date: string
  nationality: string
  occupation: string
  education_level: string
  ktp_address: string
  domicile_same_as_ktp: boolean
  domicile_address?: string
  province: string
  city: string
  district: string
  village: string
  postal_code?: string
}

function sessionCookieOptions() {
  return {
    httpOnly: true,
    sameSite: 'lax' as const,
    path: '/',
    secure: process.env.NODE_ENV === 'production',
  }
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
    const { token, user } = (await res.json()) as AuthResponse
    const session: Session = { id: user.id, role: user.role, name: user.name, token }
    setCookie(COOKIE, JSON.stringify(session), sessionCookieOptions())
    return session
  })

export const registerFn = createServerFn({ method: 'POST' })
  .validator((data: RegisterPayload) => data)
  .handler(async ({ data }): Promise<Session> => {
    const res = await fetch(`${API_BASE_URL}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    })
    if (!res.ok) {
      const message = await res.text().catch(() => '')
      throw new Error(message || 'registration failed')
    }
    const { token, user } = (await res.json()) as AuthResponse
    const session: Session = { id: user.id, role: user.role, name: user.name, token }
    setCookie(COOKIE, JSON.stringify(session), sessionCookieOptions())
    return session
  })

export const logoutFn = createServerFn({ method: 'POST' }).handler(async () => {
  deleteCookie(COOKIE, { path: '/' })
})
