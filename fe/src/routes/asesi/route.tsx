import { createFileRoute } from '@tanstack/react-router'
import { AppShell } from '@/components/layout/app-shell'
import { requireRole } from '@/lib/auth/guards'

export const Route = createFileRoute('/asesi')({
  beforeLoad: ({ context }) => requireRole(context, 'asesi'),
  component: AsesiLayout,
})

function AsesiLayout() {
  const { auth } = Route.useRouteContext()
  return <AppShell session={auth!} />
}
