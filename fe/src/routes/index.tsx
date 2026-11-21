import { createFileRoute, redirect } from '@tanstack/react-router'
import { homeForRole } from '@/lib/auth/types'

export const Route = createFileRoute('/')({
  beforeLoad: ({ context }) => {
    if (!context.auth) throw redirect({ to: '/login' })
    throw redirect({ to: homeForRole(context.auth.role) })
  },
})
