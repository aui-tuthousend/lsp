import { redirect } from '@tanstack/react-router'
import { createMiddleware } from '@tanstack/react-start'
import { getSessionFn } from './session'
import { homeForRole, type Role } from './types'

export const authMiddleware = createMiddleware({ type: 'function' }).server(async ({ next }) => {
  const session = await getSessionFn()
  if (!session) throw redirect({ to: '/login' })
  return next({ context: { session } })
})

export const requireRoleMiddleware = (...roles: Role[]) =>
  createMiddleware({ type: 'function' }).server(async ({ next }) => {
    const session = await getSessionFn()
    if (!session) throw redirect({ to: '/login' })
    if (!roles.includes(session.role)) throw redirect({ to: homeForRole(session.role) })
    return next({ context: { session } })
  })
