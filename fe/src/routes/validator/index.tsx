import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export const Route = createFileRoute('/validator/')({ component: ValidatorHome })

function ValidatorHome() {
  return <Card><CardHeader><CardTitle>Dashboard Validator</CardTitle></CardHeader><CardContent>Daftar validasi dokumen.</CardContent></Card>
}
