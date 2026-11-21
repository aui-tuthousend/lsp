import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export const Route = createFileRoute('/asesi/')({ component: AsesiHome })

function AsesiHome() {
  return <Card><CardHeader><CardTitle>Dashboard Asesi</CardTitle></CardHeader><CardContent>Status pendaftaran dan asesmen asesi.</CardContent></Card>
}
