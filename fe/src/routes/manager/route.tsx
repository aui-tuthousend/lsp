import { createFileRoute } from '@tanstack/react-router'
import { AppShell } from '@/components/layout/app-shell'
import { requireRole } from '@/lib/auth/guards'

export const Route = createFileRoute('/manager')({
  beforeLoad: ({ context }) => requireRole(context, 'manager_sertifikasi'),
  component: ManagerLayout,
})

function ManagerLayout() {
  const { auth } = Route.useRouteContext()
  return <AppShell session={auth!} />
}
