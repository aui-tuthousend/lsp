import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export const Route = createFileRoute('/admin/')({ component: AdminHome })

function AdminHome() {
  return <DashboardCard title="Dashboard Admin" description="Ringkasan aktivitas sistem LSP." />
}

function DashboardCard({ title, description }: { title: string; description: string }) {
  return <Card><CardHeader><CardTitle>{title}</CardTitle></CardHeader><CardContent>{description}</CardContent></Card>
}
