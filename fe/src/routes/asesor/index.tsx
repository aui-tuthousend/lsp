import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export const Route = createFileRoute('/asesor/')({ component: AsesorHome })

function AsesorHome() {
  return <Card><CardHeader><CardTitle>Dashboard Asesor</CardTitle></CardHeader><CardContent>Ringkasan penugasan asesor.</CardContent></Card>
}
