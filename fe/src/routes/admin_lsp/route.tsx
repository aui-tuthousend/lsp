import { createFileRoute } from '@tanstack/react-router'
import { AppShell } from '@/components/layout/app-shell'
import { requireRole } from '@/lib/auth/guards'

export const Route = createFileRoute('/admin_lsp')({
  beforeLoad: ({ context }) => requireRole(context, 'admin_lsp'),
  component: AdminLspLayout,
})

function AdminLspLayout() {
  const { auth } = Route.useRouteContext()
  return <AppShell session={auth!} />
}
