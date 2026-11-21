import { createFileRoute } from '@tanstack/react-router'
import { AppShell } from '@/components/layout/app-shell'
import { requireRole } from '@/lib/auth/guards'

export const Route = createFileRoute('/admin')({
  beforeLoad: ({ context }) => requireRole(context, 'admin'),
  component: AdminLayout,
})

function AdminLayout() {
  const { auth } = Route.useRouteContext()
  return <AppShell session={auth!} />
}
