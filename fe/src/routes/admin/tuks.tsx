import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { apiFetch, type PagedResponse } from '@/lib/api'
import { qk } from '@/lib/query-keys'

type TUK = { id: string; name: string; city: string; province: string; email: string }

export const Route = createFileRoute('/admin/tuks')({ component: TuksPage })

function TuksPage() {
  const { auth } = Route.useRouteContext()
  const { data, isLoading } = useQuery({
    queryKey: qk.tuks({ page: 1 }),
    queryFn: () => apiFetch<PagedResponse<TUK>>('/tuks', auth),
    staleTime: 5 * 60_000,
  })

  return (
    <Card>
      <CardHeader><CardTitle>Master TUK</CardTitle></CardHeader>
      <CardContent>
        {isLoading ? 'Memuat...' : null}
        <div className="space-y-2">
          {data?.data.map((tuk) => (
            <div key={tuk.id} className="rounded-md border border-slate-200 p-3">
              <div className="font-medium">{tuk.name}</div>
              <div className="text-sm text-slate-500">{tuk.city}, {tuk.province} · {tuk.email}</div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}
