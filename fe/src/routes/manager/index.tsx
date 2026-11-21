import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export const Route = createFileRoute('/manager/')({ component: ManagerHome })

function ManagerHome() {
  return <Card><CardHeader><CardTitle>Dashboard Manager Sertifikasi</CardTitle></CardHeader><CardContent>Ringkasan sertifikasi.</CardContent></Card>
}
