import { createRouter as createTanStackRouter } from '@tanstack/react-router'
import { routeTree } from './routeTree.gen'
import { queryClient } from './lib/query'

export function getRouter() {
  return createTanStackRouter({
    routeTree,
    context: {
      queryClient,
      auth: null,
    },
    defaultPreload: 'intent',
    scrollRestoration: true,
  })
}

declare module '@tanstack/react-router' {
  interface Register {
    router: ReturnType<typeof getRouter>
  }
}
