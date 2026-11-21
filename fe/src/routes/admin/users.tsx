import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export const Route = createFileRoute('/admin/users')({ component: UsersPage })

function UsersPage() {
  return <Card><CardHeader><CardTitle>Users</CardTitle></CardHeader><CardContent>Halaman user admin.</CardContent></Card>
}
