import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ApplicationList } from './permohonan/-components/application-list'
import { useCreateDraft, useMyApplications } from './permohonan/-hooks'

export const Route = createFileRoute('/asesi/')({
  component: AsesiHome,
})

const PAGE_SIZE = 20

function AsesiHome() {
  const { auth } = Route.useRouteContext()
  const navigate = useNavigate()
  const { data, isLoading, isError } = useMyApplications({ page: 1, limit: PAGE_SIZE }, auth)
  const createDraft = useCreateDraft(auth)

  function startNew() {
    createDraft.mutate(undefined, {
      onSuccess: (detail) => {
        void navigate({ to: '/asesi/permohonan/$applicationId', params: { applicationId: detail.application.id } })
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
          <CardTitle>Permohonan Saya</CardTitle>
        </CardHeader>
        <CardContent className="p-4">
          <ApplicationList items={data?.data ?? []} isLoading={isLoading} isError={isError} />
        </CardContent>
      </Card>
    </div>
  )
}
