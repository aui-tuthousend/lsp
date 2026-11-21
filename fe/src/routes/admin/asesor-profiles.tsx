import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { apiFetch, type PagedResponse } from '@/lib/api'
import { qk } from '@/lib/query-keys'

type AsesorProfile = { user_id: string; name: string; email: string; tuk_name?: string; province?: string; city?: string; daily_assignment_limit: number }

export const Route = createFileRoute('/admin/asesor-profiles')({ component: AsesorProfilesPage })

function AsesorProfilesPage() {
  const { auth } = Route.useRouteContext()
  const { data, isLoading } = useQuery({
    queryKey: qk.asesorProfiles({ page: 1 }),
    queryFn: () => apiFetch<PagedResponse<AsesorProfile>>('/asesor-profiles', auth),
    staleTime: 5 * 60_000,
  })

  return (
    <Card>
      <CardHeader><CardTitle>Profil Asesor</CardTitle></CardHeader>
      <CardContent>
        {isLoading ? 'Memuat...' : null}
        <div className="space-y-2">
          {data?.data.map((item) => (
            <div key={item.user_id} className="rounded-md border border-slate-200 p-3">
              <div className="font-medium">{item.name}</div>
              <div className="text-sm text-slate-500">{item.email} · {item.tuk_name ?? 'Belum ada TUK'} · limit {item.daily_assignment_limit}</div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  )
}
