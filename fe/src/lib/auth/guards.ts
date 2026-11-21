import { redirect } from '@tanstack/react-router'
import type { RouterContext } from '@/routes/__root'
import type { Role } from './types'
import { homeForRole } from './types'

export function requireGuest(ctx: RouterContext) {
  if (ctx.auth) throw redirect({ to: homeForRole(ctx.auth.role) })
}

export function requireAuth(ctx: RouterContext, redirectTo?: string) {
  if (!ctx.auth) {
    throw redirect({ to: '/login', search: redirectTo ? { redirect: redirectTo } : undefined })
  }
}

export function requireRole(ctx: RouterContext, ...roles: Role[]) {
  requireAuth(ctx)
  if (ctx.auth && !roles.includes(ctx.auth.role)) {
    throw redirect({ to: homeForRole(ctx.auth.role) })
  }
}
