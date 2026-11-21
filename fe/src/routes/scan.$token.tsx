import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/scan/$token')({
  component: ScanPage,
})

function ScanPage() {
  const { token } = Route.useParams()
  return <div className="flex min-h-screen items-center justify-center p-6">Scan token: {token}</div>
}
