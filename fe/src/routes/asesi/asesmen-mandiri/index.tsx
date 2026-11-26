import { createFileRoute } from '@tanstack/react-router'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { AssessmentList } from './-components/assessment-list'
import { useMyAssessments } from './-hooks'

export const Route = createFileRoute('/asesi/asesmen-mandiri/')({
  component: AsesmenMandiriHome,
})

const PAGE_SIZE = 20

function AsesmenMandiriHome() {
  const { auth } = Route.useRouteContext()
  const { data, isLoading, isError } = useMyAssessments({ page: 1, limit: PAGE_SIZE }, auth)

  return (
    <div className="space-y-4">
      <div>
        <h1 className="text-lg font-semibold text-slate-900">Asesmen Mandiri</h1>
        <p className="text-sm text-slate-500">Pernyataan kompetensi per unit kompetensi (FR.APL.02).</p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Daftar Asesmen Mandiri</CardTitle>
        </CardHeader>
        <CardContent className="p-4">
          <AssessmentList items={data?.data ?? []} isLoading={isLoading} isError={isError} />
        </CardContent>
      </Card>
    </div>
  )
}
