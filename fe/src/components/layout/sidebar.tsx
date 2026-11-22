import { Link, useRouterState } from '@tanstack/react-router'
import type { Role } from '@/lib/auth/types'
import { cn } from '@/lib/utils'
import { navFor } from './nav'

type SidebarProps = {
  role: Role
  collapsed?: boolean
  onNavigate?: () => void
  className?: string
}

export function Sidebar({ role, collapsed = false, onNavigate, className }: SidebarProps) {
  const pathname = useRouterState({ select: (state) => state.location.pathname })
  const sections = navFor(role)

  return (
    <aside
      className={cn(
        'flex h-full shrink-0 flex-col border-r border-slate-200 bg-white transition-[width] duration-200',
        collapsed ? 'w-20' : 'w-72',
        className,
      )}
    >
      <div className={cn('flex h-16 items-center border-b border-slate-200 px-4', collapsed && 'justify-center px-0')}>
        {collapsed ? (
          <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-slate-900 text-sm font-semibold text-white">L</div>
        ) : (
          <div>
            <div className="font-semibold">LSP</div>
            <div className="text-xs text-slate-500">Management System</div>
          </div>
        )}
      </div>
      <nav className={cn('flex-1 space-y-6 overflow-y-auto p-4', collapsed && 'px-3')}>
        {sections.map((section) => (
          <div key={section.title ?? 'section'}>
            {section.title && !collapsed ? (
              <div className="mb-2 px-3 text-xs font-medium uppercase tracking-wide text-slate-500">{section.title}</div>
            ) : null}
            <div className="space-y-1">
              {section.items.map((item) => {
                const Icon = item.icon
                const active = item.to === '/' ? pathname === item.to : pathname.startsWith(item.to)
                return (
                  <Link
                    key={item.to}
                    to={item.to}
                    title={collapsed ? item.label : undefined}
                    onClick={onNavigate}
                    className={cn(
                      'flex min-h-10 items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100',
                      collapsed && 'justify-center px-0',
                      active && 'bg-slate-100 text-slate-950',
                    )}
                  >
                    <Icon className="h-4 w-4 shrink-0" />
                    {!collapsed ? <span className="truncate">{item.label}</span> : null}
                  </Link>
                )
              })}
            </div>
          </div>
        ))}
      </nav>
    </aside>
  )
}
