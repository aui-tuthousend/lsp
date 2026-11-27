import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { ArrowRight, Inbox, Layers, Loader2, User } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { ApplicationStatusBadge } from '../asesi/permohonan/-components/application-status-badge'
import { formatDateTime } from '../asesi/permohonan/-helpers'
import { useQueue } from './-hooks'
import type { ApplicationStatus, QueueItem } from './-types'

export const Route = createFileRoute('/admin_lsp/')({ component: AdminLspHome })

const PAGE_SIZE = 20

// Tab antrian berbasis status permohonan APL-01 (keputusan admin), bukan status
// asesmen mandiri: item yang sudah diputuskan berpindah ke tabnya.
const TABS: { value: ApplicationStatus; label: string; empty: string }[] = [
  { value: 'SUBMITTED', label: 'Menunggu Verifikasi', empty: 'Tidak ada pengajuan yang menunggu verifikasi.' },
  { value: 'APPROVED', label: 'Disetujui', empty: 'Belum ada pengajuan yang disetujui.' },
  { value: 'REJECTED', label: 'Ditolak', empty: 'Belum ada pengajuan yang ditolak.' },
]

function AdminLspHome() {
  const { auth } = Route.useRouteContext()
  const [status, setStatus] = useState<ApplicationStatus>('SUBMITTED')
  const { data, isLoading, isError } = useQueue({ page: 1, limit: PAGE_SIZE, application_status: status }, auth)
  const items = data?.data ?? []
  const active = TABS.find((tab) => tab.value === status) ?? TABS[0]

  return (
    <div className="space-y-4">
      <div>
        <h1 className="text-lg font-semibold text-slate-900">Verifikasi Pengajuan</h1>
        <p className="text-sm text-slate-500">
          Daftar pengajuan asesi yang menunggu keputusan. Tinjau data APL-01 &amp; APL-02 lalu setujui atau tolak.
        </p>
      </div>

      <Tabs value={status} onValueChange={(value) => setStatus(value as ApplicationStatus)}>
        <TabsList>
          {TABS.map((tab) => (
            <TabsTrigger key={tab.value} value={tab.value}>
              {tab.label}
            </TabsTrigger>
          ))}
        </TabsList>
      </Tabs>

      <Card>
        <CardHeader>
          <CardTitle>{active.label}</CardTitle>
        </CardHeader>
        <CardContent className="p-4">
          {isLoading ? (
            <div className="flex items-center gap-2 py-6 text-sm text-slate-500">
              <Loader2 className="h-4 w-4 animate-spin" />
              Memuat pengajuan...
            </div>
          ) : isError ? (
            <p className="py-6 text-sm text-rose-600">Gagal memuat daftar pengajuan.</p>
          ) : items.length === 0 ? (
            <div className="flex flex-col items-center gap-2 py-8 text-center">
              <Inbox className="h-8 w-8 text-slate-300" />
              <p className="text-sm text-slate-500">{active.empty}</p>
            </div>
          ) : (
            <div className="grid gap-3">
              {items.map((item) => (
                <QueueRow key={item.id} item={item} />
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

function QueueRow({ item }: { item: QueueItem }) {
  const navigate = useNavigate()

  return (
    <div className="flex flex-wrap items-center justify-between gap-3 rounded-lg border border-slate-200 bg-white p-4">
      <div className="min-w-0 space-y-1">
        <div className="flex flex-wrap items-center gap-2">
          <span className="flex items-center gap-1.5 font-medium text-slate-900">
            <User className="h-4 w-4 text-slate-400" />
            {item.asesi_name || 'Asesi'}
          </span>
          {item.application_status ? <ApplicationStatusBadge status={item.application_status} /> : null}
        </div>
        <p className="text-xs text-slate-500">
          {item.application_code || item.code}
          <span className="text-slate-400"> · {item.code}</span>
        </p>
        {item.schema_name ? (
          <p className="flex flex-wrap items-center gap-1.5 text-xs text-slate-600">
            <Layers className="h-3.5 w-3.5 text-slate-400" />
            <span className="font-medium">{item.schema_name}</span>
            {item.schema_code ? <span className="text-slate-400">· {item.schema_code}</span> : null}
          </p>
        ) : null}
        <p className="text-xs text-slate-400">Diajukan {formatDateTime(item.submitted_at ?? item.created_at)}</p>
      </div>
      <Button
        variant="outline"
        size="sm"
        className="gap-1.5"
        onClick={() => navigate({ to: '/admin_lsp/verifikasi/$assessmentId', params: { assessmentId: item.id } })}
      >
        Tinjau
        <ArrowRight className="h-4 w-4" />
      </Button>
    </div>
  )
}
