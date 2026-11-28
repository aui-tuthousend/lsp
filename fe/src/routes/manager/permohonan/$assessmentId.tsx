import { createFileRoute } from '@tanstack/react-router'
import { ApplicationReview } from '../../admin_lsp/-components/application-review'
import { AsesorAssignmentPanel } from '../-components/asesor-assignment-panel'

export const Route = createFileRoute('/manager/permohonan/$assessmentId')({ component: ManagerDetailPage })

function ManagerDetailPage() {
  const { assessmentId } = Route.useParams()
  const { auth } = Route.useRouteContext()
  const session = auth!

  return (
    <ApplicationReview
      assessmentId={assessmentId}
      session={session}
      backTo="/manager"
      backLabel="Kembali ke daftar"
      title="Plotting Asesor"
      footer={({ applicationDetail }) => (
        <AsesorAssignmentPanel
          application={applicationDetail.application}
          assignment={applicationDetail.assignment ?? null}
          session={session}
        />
      )}
    />
  )
}
