import { createFileRoute } from '@tanstack/react-router'
import { AppShell } from '@/components/layout/app-shell'
import { requireRole } from '@/lib/auth/guards'

export const Route = createFileRoute('/validator')({
  beforeLoad: ({ context }) => requireRole(context, 'validator'),
  component: ValidatorLayout,
})

function ValidatorLayout() {
  const { auth } = Route.useRouteContext()
  return <AppShell session={auth!} />
}
