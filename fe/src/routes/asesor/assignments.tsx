import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export const Route = createFileRoute('/asesor/assignments')({ component: AssignmentsPage })

function AssignmentsPage() {
  return <Card><CardHeader><CardTitle>Penugasan Saya</CardTitle></CardHeader><CardContent>Daftar penugasan asesor.</CardContent></Card>
}
