import * as React from 'react'
import { HeadContent, Outlet, Scripts, createRootRouteWithContext } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'
import { QueryClientProvider } from '@tanstack/react-query'
import type { QueryClient } from '@tanstack/react-query'
import type { Session } from '@/lib/auth/types'
import { getSessionFn } from '@/lib/auth/session'
import '../styles.css'

export type RouterContext = {
  queryClient: QueryClient
  auth: Session | null
}

export const Route = createRootRouteWithContext<RouterContext>()({
  beforeLoad: async () => ({ auth: await getSessionFn() }),
  component: RootComponent,
})

function RootComponent() {
  const { queryClient } = Route.useRouteContext()
  return (
    <RootDocument>
      <QueryClientProvider client={queryClient}>
        <Outlet />
        <TanStackRouterDevtools position="bottom-right" />
      </QueryClientProvider>
    </RootDocument>
  )
}

function RootDocument({ children }: { children: React.ReactNode }) {
  return (
    <html lang="id">
      <head>
        <HeadContent />
      </head>
      <body>
        {children}
        <Scripts />
      </body>
    </html>
  )
}
