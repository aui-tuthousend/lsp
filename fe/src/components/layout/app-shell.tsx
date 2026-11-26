import { useQueryClient } from '@tanstack/react-query'
import { Outlet, useRouter } from '@tanstack/react-router'
import { LogOut, Menu, PanelLeftClose, PanelLeftOpen, User, X } from 'lucide-react'
import { useState } from 'react'
import { toast } from 'sonner'
import { Avatar } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { logoutFn } from '@/lib/auth/session'
import type { Session } from '@/lib/auth/types'
import { profileForRole } from '@/lib/auth/types'
import { useNotifications } from '@/lib/sse/use-notifications'
import { cn } from '@/lib/utils'
import { Sidebar } from './sidebar'

export function AppShell({ session }: { session: Session }) {
  const router = useRouter()
  const queryClient = useQueryClient()
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false)
  const [mobileSidebarOpen, setMobileSidebarOpen] = useState(false)
  useNotifications(session)

  async function logout() {
    await logoutFn()
    // Bersihkan cache agar data milik akun sebelumnya tidak tersisa saat berganti akun.
    queryClient.clear()
    await router.invalidate()
    toast.success('Berhasil logout')
    await router.navigate({ to: '/login', replace: true })
  }

  function openProfile() {
    const path = profileForRole(session.role)
    if (!path) {
      toast.info('Halaman profil belum tersedia')
      return
    }
    void router.navigate({ to: path })
  }

  return (
    <div className="min-h-screen bg-slate-50 text-slate-950">
      {mobileSidebarOpen ? (
        <div className="fixed inset-0 z-40 lg:hidden" role="dialog" aria-modal="true">
          <button
            type="button"
            aria-label="Tutup menu"
            className="absolute inset-0 bg-slate-950/40"
            onClick={() => setMobileSidebarOpen(false)}
          />
          <div className="relative z-10 h-full w-[min(18rem,85vw)] bg-white shadow-xl">
            <button
              type="button"
              aria-label="Tutup menu"
              className="absolute right-3 top-3 inline-flex h-10 w-10 items-center justify-center rounded-md text-slate-500 hover:bg-slate-100 hover:text-slate-900"
              onClick={() => setMobileSidebarOpen(false)}
            >
              <X className="h-5 w-5" />
            </button>
            <Sidebar role={session.role} onNavigate={() => setMobileSidebarOpen(false)} />
          </div>
        </div>
      ) : null}

      <div className="flex min-h-screen min-w-0">
        <div className="hidden lg:block">
          <Sidebar role={session.role} collapsed={sidebarCollapsed} />
        </div>

        <div className="flex min-w-0 flex-1 flex-col">
          <header className="sticky top-0 z-30 flex h-16 items-center justify-between gap-3 border-b border-slate-200 bg-white/95 px-3 backdrop-blur sm:px-4 lg:px-6">
            <div className="flex min-w-0 items-center gap-2 sm:gap-3">
              <Button
                type="button"
                variant="ghost"
                className="h-10 w-10 shrink-0 px-0 lg:hidden"
                aria-label="Buka menu"
                onClick={() => setMobileSidebarOpen(true)}
              >
                <Menu className="h-5 w-5" />
              </Button>
              <Button
                type="button"
                variant="ghost"
                className="hidden h-10 w-10 shrink-0 px-0 lg:inline-flex"
                aria-label={sidebarCollapsed ? 'Buka sidebar' : 'Tutup sidebar'}
                onClick={() => setSidebarCollapsed((value) => !value)}
              >
                {sidebarCollapsed ? <PanelLeftOpen className="h-5 w-5" /> : <PanelLeftClose className="h-5 w-5" />}
              </Button>
              <div className="min-w-0">
                <div className="truncate text-sm font-medium">{session.name}</div>
                <div className="truncate text-xs text-slate-500">{session.role}</div>
              </div>
            </div>
            <DropdownMenu className="shrink-0">
              <DropdownMenuTrigger
                aria-label="Menu akun"
                className="rounded-full outline-none transition focus-visible:ring-2 focus-visible:ring-slate-300"
              >
                <Avatar name={session.name} />
              </DropdownMenuTrigger>
              <DropdownMenuContent className="w-56">
                <div className="px-2.5 py-2">
                  <div className="truncate text-sm font-medium text-slate-900">{session.name}</div>
                  <div className="truncate text-xs text-slate-500">{session.role}</div>
                </div>
                <DropdownMenuSeparator />
                <DropdownMenuItem onSelect={openProfile}>
                  <User className="h-4 w-4" />
                  Profil
                </DropdownMenuItem>
                <DropdownMenuItem
                  onSelect={() => void logout()}
                  className="text-rose-600 hover:bg-rose-50"
                >
                  <LogOut className="h-4 w-4" />
                  Logout
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </header>
          <main className={cn('min-w-0 flex-1 overflow-x-hidden p-3 sm:p-4 lg:p-6')}>
            <div className="mx-auto w-full max-w-7xl">
              <Outlet />
            </div>
          </main>
        </div>
      </div>
    </div>
  )
}
