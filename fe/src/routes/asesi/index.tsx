import { useState } from 'react'
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { ApplicationList } from './permohonan/-components/application-list'
import { useCreateDraft, useMyApplications } from './permohonan/-hooks'
import { STATUS_LABELS } from './permohonan/-helpers'
import type { Application, ApplicationStatus } from './permohonan/-types'
import { useCreateOrGetAssessment } from './asesmen-mandiri/-hooks'

export const Route = createFileRoute('/asesi/')({
  component: AsesiHome,
})

const PAGE_SIZE = 20

const STATUS_FILTERS: { value: ApplicationStatus | 'ALL'; label: string }[] = [
  { value: 'ALL', label: 'Semua' },
  { value: 'DRAFT', label: STATUS_LABELS.DRAFT },
  { value: 'SUBMITTED', label: STATUS_LABELS.SUBMITTED },
  { value: 'APPROVED', label: STATUS_LABELS.APPROVED },
  { value: 'REJECTED', label: STATUS_LABELS.REJECTED },
]

function AsesiHome() {
  const { auth } = Route.useRouteContext()
  const navigate = useNavigate()
  const [status, setStatus] = useState<ApplicationStatus | 'ALL'>('ALL')
  const { data, isLoading, isError } = useMyApplications(
    { page: 1, limit: PAGE_SIZE, status: status === 'ALL' ? undefined : status },
    auth,
  )
  const createDraft = useCreateDraft(auth)
  const openAssessment = useCreateOrGetAssessment(auth)

  function startNew() {
    createDraft.mutate(undefined, {
      onSuccess: (detail) => {
        void navigate({ to: '/asesi/permohonan/$applicationId', params: { applicationId: detail.application.id } })
      },
    })
  }

  function startAssessment(application: Application) {
    openAssessment.mutate(application.id, {
      onSuccess: (detail) => {
        void navigate({ to: '/asesi/asesmen-mandiri/$assessmentId', params: { assessmentId: detail.assessment.id } })
      },
    })
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="text-lg font-semibold text-slate-900">Dashboard Asesi</h1>
          <p className="text-sm text-slate-500">Permohonan sertifikasi kompetensi Anda.</p>
        </div>
        <Button onClick={startNew} disabled={createDraft.isPending} className="gap-2">
          <Plus className="h-4 w-4" />
          {createDraft.isPending ? 'Membuat...' : 'Buat Permohonan Baru'}
        </Button>
      </div>

      <Card>
        <CardHeader>
          <div className="flex flex-wrap items-center justify-between gap-3">
            <CardTitle>Permohonan Saya</CardTitle>
            <Tabs
              value={status}
              onValueChange={(value) => setStatus(value as ApplicationStatus | 'ALL')}
              className="w-full min-w-0 sm:w-auto"
            >
              <TabsList>
                {STATUS_FILTERS.map((filter) => (
                  <TabsTrigger key={filter.value} value={filter.value}>
                    {filter.label}
                  </TabsTrigger>
                ))}
              </TabsList>
            </Tabs>
          </div>
        </CardHeader>
        <CardContent className="p-4">
          <ApplicationList
            items={data?.data ?? []}
            isLoading={isLoading}
            isError={isError}
            onOpenAssessment={startAssessment}
            emptyMessage={
              status === 'ALL'
                ? undefined
                : `Tidak ada permohonan berstatus ${STATUS_LABELS[status].toLowerCase()}.`
            }
          />
        </CardContent>
      </Card>
    </div>
  )
}
