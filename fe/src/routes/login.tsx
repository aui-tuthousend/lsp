import { createFileRoute, useRouter } from '@tanstack/react-router'
import { useState } from 'react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { requireGuest } from '@/lib/auth/guards'
import { loginFn } from '@/lib/auth/session'
import { homeForRole } from '@/lib/auth/types'

export const Route = createFileRoute('/login')({
  beforeLoad: ({ context }) => requireGuest(context),
  component: LoginPage,
})

function LoginPage() {
  const router = useRouter()
  const [email, setEmail] = useState('admin@lsp.local')
  const [password, setPassword] = useState('password')
  const [loading, setLoading] = useState(false)

  async function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setLoading(true)
    try {
      const session = await loginFn({ data: { email, password } })
      await router.invalidate()
      toast.success('Login berhasil', { description: `Selamat datang, ${session.name}` })
      await router.navigate({ to: homeForRole(session.role), replace: true })
    } catch {
      toast.error('Login gagal', { description: 'Email atau password tidak valid' })
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-50 p-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Login LSP</CardTitle>
        </CardHeader>
        <CardContent>
          <form className="space-y-4" onSubmit={submit}>
            <Input value={email} onChange={(event) => setEmail(event.target.value)} placeholder="Email" type="email" />
            <Input value={password} onChange={(event) => setPassword(event.target.value)} placeholder="Password" type="password" />
            <Button className="w-full" type="submit" disabled={loading}>
              {loading ? 'Memproses...' : 'Masuk'}
            </Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
