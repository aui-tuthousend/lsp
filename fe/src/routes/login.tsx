import { createFileRoute, useRouter } from '@tanstack/react-router'
import { useState } from 'react'
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
  const [error, setError] = useState<string | null>(null)

  async function submit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError(null)
    try {
      const session = await loginFn({ data: { email, password } })
      await router.invalidate()
      await router.navigate({ to: homeForRole(session.role) })
    } catch {
      setError('Email atau password tidak valid')
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
            {error ? <p className="text-sm text-red-600">{error}</p> : null}
            <Button className="w-full" type="submit">Masuk</Button>
          </form>
        </CardContent>
      </Card>
    </div>
  )
}
