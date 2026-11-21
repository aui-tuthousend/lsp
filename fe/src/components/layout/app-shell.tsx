import { Outlet, useRouter } from '@tanstack/react-router'
import { Button } from '@/components/ui/button'
import { logoutFn } from '@/lib/auth/session'
import type { Session } from '@/lib/auth/types'
import { useNotifications } from '@/lib/sse/use-notifications'
import { Sidebar } from './sidebar'

export function AppShell({ session }: { session: Session }) {
  const router = useRouter()
  useNotifications(session)

  async function logout() {
    await logoutFn()
    await router.invalidate()
    await router.navigate({ to: '/login' })
  }

  return (
    <div className="flex min-h-screen bg-slate-50">
      <Sidebar role={session.role} />
      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex h-16 items-center justify-between border-b border-slate-200 bg-white px-4 lg:px-6">
          <div>
            <div className="text-sm font-medium">{session.name}</div>
            <div className="text-xs text-slate-500">{session.role}</div>
          </div>
          <Button variant="outline" onClick={logout}>Logout</Button>
        </header>
        <main className="flex-1 p-4 lg:p-6">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
