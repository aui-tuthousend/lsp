import { createFileRoute } from '@tanstack/react-router'
import { AppShell } from '@/components/layout/app-shell'
import { requireRole } from '@/lib/auth/guards'

export const Route = createFileRoute('/asesor')({
  beforeLoad: ({ context }) => requireRole(context, 'asesor'),
  component: AsesorLayout,
})

function AsesorLayout() {
  const { auth } = Route.useRouteContext()
  return <AppShell session={auth!} />
}
