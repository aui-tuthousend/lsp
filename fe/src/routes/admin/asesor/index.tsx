import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/admin/asesor/')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/admin/asesor/"!</div>
}
