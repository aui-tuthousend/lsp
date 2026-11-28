import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useState } from 'react'
import { ArrowRight, Inbox, Layers, Loader2, User, UserCheck } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { ApplicationStatusBadge } from '../asesi/permohonan/-components/application-status-badge'
import { formatDateTime } from '../asesi/permohonan/-helpers'
import { AssignmentStatusBadge } from './-components/assignment-status-badge'
import { useManagerQueue } from './-hooks'
import type { QueueItem } from './-types'

export const Route = createFileRoute('/manager/')({ component: ManagerHome })

const PAGE_SIZE = 20

// Tab antrian berbasis penugasan asesor: memisahkan pengajuan yang belum
// di-plotting dari yang sudah ditugaskan.
type AssignFilter = 'UNASSIGNED' | 'ASSIGNED'

const TABS: { value: AssignFilter; label: string; empty: string }[] = [
  { value: 'UNASSIGNED', label: 'Belum Ditugaskan', empty: 'Semua pengajuan sudah ditugaskan ke asesor.' },
  { value: 'ASSIGNED', label: 'Sudah Ditugaskan', empty: 'Belum ada pengajuan yang ditugaskan ke asesor.' },
]

function ManagerHome() {
  const { auth } = Route.useRouteContext()
  const [filter, setFilter] = useState<AssignFilter>('UNASSIGNED')
  const { data, isLoading, isError } = useManagerQueue(
    { page: 1, limit: PAGE_SIZE, application_status: 'APPROVED', assigned: filter === 'ASSIGNED' },
    auth,
  )
  const items = data?.data ?? []
  const active = TABS.find((tab) => tab.value === filter) ?? TABS[0]

  return (
    <div className="space-y-4">
      <div>
        <h1 className="text-lg font-semibold text-slate-900">Pengajuan Disetujui</h1>
        <p className="text-sm text-slate-500">
          Daftar permohonan yang sudah disetujui admin LSP. Buka detail untuk menugaskan asesor.
        </p>
      </div>

      <Tabs value={filter} onValueChange={(value) => setFilter(value as AssignFilter)}>
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
        {item.asesor_assignment_status ? (
          <div className="flex flex-wrap items-center gap-1.5">
            <UserCheck className="h-3.5 w-3.5 text-slate-400" />
            <span className="text-xs text-slate-600">{item.asesor_name || 'Asesor'}</span>
            <AssignmentStatusBadge status={item.asesor_assignment_status} />
          </div>
        ) : (
          <p className="text-xs font-medium text-amber-600">Belum ditugaskan ke asesor</p>
        )}
      </div>
      <Button
        variant="outline"
        size="sm"
        className="gap-1.5"
        onClick={() => navigate({ to: '/manager/permohonan/$assessmentId', params: { assessmentId: item.id } })}
      >
        Tinjau
        <ArrowRight className="h-4 w-4" />
      </Button>
    </div>
  )
}
