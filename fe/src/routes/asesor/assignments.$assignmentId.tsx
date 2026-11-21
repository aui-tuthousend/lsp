import { createFileRoute, redirect } from '@tanstack/react-router'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { requireRole } from '@/lib/auth/guards'

export const Route = createFileRoute('/asesor/assignments/$assignmentId')({
  beforeLoad: async ({ context, params }) => {
    requireRole(context, 'asesor')

    // TODO: ganti dengan server function getAssignmentFn saat endpoint assignment tersedia.
    const assignment = { id: params.assignmentId, asesorId: context.auth!.id }
    if (assignment.asesorId !== context.auth!.id) {
      throw redirect({ to: '/asesor/assignments' })
    }
  },
  component: AssignmentDetailPage,
})

function AssignmentDetailPage() {
  const { assignmentId } = Route.useParams()
  return <Card><CardHeader><CardTitle>Detail Asesmen</CardTitle></CardHeader><CardContent>ID: {assignmentId}</CardContent></Card>
}
