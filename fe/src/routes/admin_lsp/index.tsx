import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export const Route = createFileRoute('/admin_lsp/')({ component: AdminLspHome })

function AdminLspHome() {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Dashboard Admin LSP</CardTitle>
      </CardHeader>
      <CardContent>Ringkasan verifikasi dan pengelolaan permohonan sertifikasi.</CardContent>
    </Card>
  )
}
