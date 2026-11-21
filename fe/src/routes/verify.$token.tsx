import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/verify/$token')({
  component: VerifyPage,
})

function VerifyPage() {
  const { token } = Route.useParams()
  return <div className="p-6">Halaman verifikasi QR: {token}</div>
}
