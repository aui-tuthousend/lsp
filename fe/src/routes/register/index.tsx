import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { requireGuest } from '@/lib/auth/guards'

export const Route = createFileRoute('/register/')({
  beforeLoad: ({ context }) => requireGuest(context),
  component: RegisterPage,
})

function RegisterPage() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-50 p-4">
      <Card className="w-full max-w-md">
        <CardHeader><CardTitle>Register Asesi</CardTitle></CardHeader>
        <CardContent className="text-sm text-slate-600">Form register akan diisi pada tahap halaman asesi.</CardContent>
      </Card>
    </div>
  )
}
