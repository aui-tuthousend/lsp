import { Link, useRouterState } from '@tanstack/react-router'
import type { Role } from '@/lib/auth/types'
import { cn } from '@/lib/utils'
import { navFor } from './nav'

export function Sidebar({ role }: { role: Role }) {
  const pathname = useRouterState({ select: (state) => state.location.pathname })
  const sections = navFor(role)

  return (
    <aside className="hidden w-72 shrink-0 border-r border-slate-200 bg-white lg:block">
      <div className="flex h-16 items-center border-b border-slate-200 px-6">
        <div>
          <div className="font-semibold">LSP</div>
          <div className="text-xs text-slate-500">Management System</div>
        </div>
      </div>
      <nav className="space-y-6 p-4">
        {sections.map((section) => (
          <div key={section.title ?? 'section'}>
            {section.title ? <div className="mb-2 px-3 text-xs font-medium uppercase tracking-wide text-slate-500">{section.title}</div> : null}
            <div className="space-y-1">
              {section.items.map((item) => {
                const Icon = item.icon
                const active = item.to === '/' ? pathname === item.to : pathname.startsWith(item.to)
                return (
                  <Link
                    key={item.to}
                    to={item.to}
                    className={cn(
                      'flex items-center gap-3 rounded-md px-3 py-2 text-sm font-medium text-slate-700 hover:bg-slate-100',
                      active && 'bg-slate-100 text-slate-950',
                    )}
                  >
                    <Icon className="h-4 w-4" />
                    {item.label}
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
