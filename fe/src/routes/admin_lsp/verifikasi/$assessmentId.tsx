import { createFileRoute, Link } from '@tanstack/react-router'
import { useMemo, useState } from 'react'
import { toast } from 'sonner'
import { ArrowLeft, Check, Loader2, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import type { Session } from '@/lib/auth/types'
import { ApplicationStatusBadge } from '../../asesi/permohonan/-components/application-status-badge'
import { formatDate, formatDateTime, textareaClass } from '../../asesi/permohonan/-helpers'
import { useApplication, useSchemaRequirements, useTukOptions } from '../../asesi/permohonan/-hooks'
import type { Application, ApplicationDetail } from '../../asesi/permohonan/-types'
import { AssessmentStatusBadge } from '../../asesi/asesmen-mandiri/-components/assessment-status-badge'
import { useAssessment, useSchemaTree } from '../../asesi/asesmen-mandiri/-hooks'
import { purposeLabel } from '../../asesi/asesmen-mandiri/-helpers'
import type { AssessmentDetail, ElementNode } from '../../asesi/asesmen-mandiri/-types'
import { useAsesiProfileByUser } from '../../asesi/profil/-hooks'
import { EDUCATION_OPTIONS, GENDER_OPTIONS } from '../../asesi/profil/-store'
import type { AsesiProfile } from '../../asesi/profil/-types'
import { DocumentsAccordion } from '../-components/documents-accordion'
import { UnitSummary } from '../-components/unit-summary'
import { useDecide } from '../-hooks'

export const Route = createFileRoute('/admin_lsp/verifikasi/$assessmentId')({ component: VerifyPage })

const GENDER_LABELS = new Map(GENDER_OPTIONS.map((option) => [option.value, option.label]))
const EDUCATION_LABELS = new Map(EDUCATION_OPTIONS.map((option) => [option.value, option.label]))
const MODE_LABELS: Record<string, string> = { ONLINE: 'Online', OFFLINE: 'Offline' }

function VerifyPage() {
  const { assessmentId } = Route.useParams()
  const { auth } = Route.useRouteContext()
  const assessmentQuery = useAssessment(assessmentId, auth)

  if (assessmentQuery.isLoading) {
    return <Loading label="Memuat asesmen mandiri..." />
  }
  if (assessmentQuery.isError || !assessmentQuery.data) {
    return <ErrorCard text="Asesmen mandiri tidak ditemukan atau gagal dimuat." />
  }

  return <VerificationView detail={assessmentQuery.data} session={auth!} />
}

function VerificationView({ detail, session }: { detail: AssessmentDetail; session: Session }) {
  const applicationId = detail.assessment.application_id

  // Semua hook dipanggil tanpa syarat (aturan hooks); data turunan menyusul.
  const applicationQuery = useApplication(applicationId, session)
  const userId = applicationQuery.data?.application.user_id ?? ''
  const profileQuery = useAsesiProfileByUser(userId, session)
  const schemaTree = useSchemaTree(detail.application.certification_schema_id, session)
  const requirements = useSchemaRequirements(detail.application.certification_schema_id, session)
  const tukOptions = useTukOptions(session)

  const elementsByUnit = useMemo(() => {
    const map: Record<string, ElementNode[]> = {}
    for (const unit of schemaTree.data?.units ?? []) map[unit.id] = unit.elements
    return map
  }, [schemaTree.data])

  if (applicationQuery.isLoading) {
    return <Loading label="Memuat permohonan..." />
  }
  if (applicationQuery.isError || !applicationQuery.data) {
    return <ErrorCard text="Permohonan tidak ditemukan atau gagal dimuat." />
  }

  const applicationDetail: ApplicationDetail = applicationQuery.data
  const application = applicationDetail.application
  const schedule = applicationDetail.schedule
  const tukName = schedule?.tuk_id
    ? (tukOptions.data?.data ?? []).find((tuk) => tuk.id === schedule.tuk_id)?.name ?? schedule.tuk_id
    : '—'

  return (
    <div className="space-y-5">
      <div className="space-y-2">
        <Link
          to="/admin_lsp"
          className="inline-flex items-center gap-1.5 text-sm text-slate-500 transition-colors hover:text-slate-900"
        >
          <ArrowLeft className="h-4 w-4" />
          Kembali ke daftar
        </Link>
        <div className="flex flex-wrap items-center gap-2">
          <h1 className="text-lg font-semibold text-slate-900">Verifikasi Pengajuan</h1>
          <ApplicationStatusBadge status={application.status} />
          <AssessmentStatusBadge status={detail.assessment.status} />
        </div>
        <p className="text-sm text-slate-500">
          {profileQuery.data?.full_name ? `${profileQuery.data.full_name} · ` : ''}
          {application.code} · {detail.assessment.code}
        </p>
      </div>

      <Tabs defaultValue="apl01">
        <TabsList>
          <TabsTrigger value="apl01">APL-01</TabsTrigger>
          <TabsTrigger value="apl02">APL-02</TabsTrigger>
        </TabsList>

        <TabsContent value="apl01" className="space-y-4">
          <IdentityCard profile={profileQuery.data} loading={profileQuery.isLoading} />
          <CertificationCard
            schemaName={detail.application.schema_name}
            schemaCode={detail.application.schema_code}
            purpose={detail.application.purpose}
          />
          <ScheduleCard
            mode={schedule?.mode}
            date={schedule?.scheduled_date}
            tukName={tukName}
          />
          <DocumentsAccordion
            documents={applicationDetail.documents}
            requirements={requirements.data?.data ?? []}
            loading={requirements.isLoading}
          />
        </TabsContent>

        <TabsContent value="apl02" className="space-y-4">
          <div className="flex flex-wrap items-baseline gap-2">
            <h2 className="text-base font-semibold text-slate-900">Pernyataan Kompetensi per Unit</h2>
            <span className="text-sm text-slate-500">{detail.units.length} unit kompetensi</span>
          </div>
          {detail.units.length === 0 ? (
            <Card>
              <CardContent className="p-4 text-sm text-slate-500">Belum ada unit kompetensi.</CardContent>
            </Card>
          ) : (
            <div className="space-y-4">
              {detail.units.map((unit) => (
                <UnitSummary key={unit.id} unit={unit} elements={elementsByUnit[unit.competency_unit_id] ?? []} />
              ))}
            </div>
          )}
        </TabsContent>
      </Tabs>

      <DecisionPanel application={application} session={session} />
    </div>
  )
}

function IdentityCard({ profile, loading }: { profile?: AsesiProfile; loading: boolean }) {
  return (
    <Card>
      <CardContent className="p-4 sm:p-6">
        <h2 className="text-base font-semibold text-slate-900">Identitas Asesi</h2>
        {loading ? (
          <Loading label="Memuat profil..." />
        ) : !profile ? (
          <p className="mt-3 text-sm text-slate-500">Profil asesi belum diisi.</p>
        ) : (
          <dl className="mt-4 grid gap-4 sm:grid-cols-2">
            <InfoItem label="Nama Lengkap" value={profile.full_name} />
            <InfoItem label="NIK" value={profile.nik} />
            <InfoItem label="Jenis Kelamin" value={GENDER_LABELS.get(profile.gender) ?? profile.gender} />
            <InfoItem label="Tempat / Tanggal Lahir" value={`${profile.birth_place || '—'}, ${formatDate(profile.birth_date)}`} />
            <InfoItem label="Pendidikan Terakhir" value={EDUCATION_LABELS.get(profile.education_level) ?? profile.education_level} />
            <InfoItem label="Pekerjaan" value={profile.occupation} />
            <InfoItem label="Telepon" value={profile.phone} />
            <InfoItem label="Kewarganegaraan" value={profile.nationality} />
            <InfoItem
              label="Alamat Domisili"
              value={profile.domicile_same_as_ktp ? profile.ktp_address : (profile.domicile_address ?? profile.ktp_address)}
              className="sm:col-span-2"
            />
          </dl>
        )}
      </CardContent>
    </Card>
  )
}

function CertificationCard({
  schemaName,
  schemaCode,
  purpose,
}: {
  schemaName?: string
  schemaCode?: string
  purpose?: string
}) {
  return (
    <Card>
      <CardContent className="p-4 sm:p-6">
        <h2 className="text-base font-semibold text-slate-900">Data Sertifikasi</h2>
        <dl className="mt-4 grid gap-4 sm:grid-cols-2">
          <InfoItem label="Nama Skema" value={schemaName ?? '—'} />
          <InfoItem label="Nomor Skema" value={schemaCode ?? '—'} />
          <InfoItem label="Tujuan Asesmen" value={purposeLabel(purpose)} className="sm:col-span-2" />
        </dl>
      </CardContent>
    </Card>
  )
}

function ScheduleCard({ mode, date, tukName }: { mode?: string; date?: string; tukName: string }) {
  return (
    <Card>
      <CardContent className="p-4 sm:p-6">
        <h2 className="text-base font-semibold text-slate-900">Jadwal Asesmen</h2>
        <dl className="mt-4 grid gap-4 sm:grid-cols-3">
          <InfoItem label="Metode" value={mode ? (MODE_LABELS[mode] ?? mode) : '—'} />
          <InfoItem label="Tanggal" value={formatDate(date)} />
          <InfoItem label="TUK" value={tukName} />
        </dl>
      </CardContent>
    </Card>
  )
}

function DecisionPanel({ application, session }: { application: Application; session: Session }) {
  const decide = useDecide(application.id, session)
  const [approveOpen, setApproveOpen] = useState(false)
  const [rejectOpen, setRejectOpen] = useState(false)
  const [note, setNote] = useState('')

  if (application.status === 'APPROVED' || application.status === 'REJECTED') {
    return (
      <Card>
        <CardHeader className="flex flex-row items-center justify-between space-y-0">
          <CardTitle>Keputusan Verifikasi</CardTitle>
          <ApplicationStatusBadge status={application.status} />
        </CardHeader>
        <CardContent className="space-y-2">
          {application.verification_note ? (
            <p className="whitespace-pre-line text-sm text-slate-600">{application.verification_note}</p>
          ) : (
            <p className="text-sm text-slate-400">Tanpa catatan.</p>
          )}
          {application.verified_at ? (
            <p className="text-xs text-slate-400">Diputuskan {formatDateTime(application.verified_at)}</p>
          ) : null}
        </CardContent>
      </Card>
    )
  }

  if (application.status !== 'SUBMITTED') {
    return (
      <Card>
        <CardContent className="p-4 text-sm text-slate-500">
          Permohonan belum diajukan sehingga belum dapat diputuskan.
        </CardContent>
      </Card>
    )
  }

  function submit(decision: 'APPROVED' | 'REJECTED') {
    const trimmed = note.trim()
    if (decision === 'REJECTED' && trimmed === '') {
      toast.error('Catatan wajib diisi untuk penolakan')
      return
    }
    decide.mutate(
      { decision, note: trimmed || undefined },
      {
        onSuccess: () => {
          setApproveOpen(false)
          setRejectOpen(false)
          setNote('')
        },
      },
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Keputusan Verifikasi</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        <p className="text-sm text-slate-500">
          Pastikan seluruh data APL-01 &amp; APL-02 sudah sesuai sebelum mengambil keputusan.
        </p>
        <div className="flex flex-wrap gap-2">
          <Button
            className="gap-1.5"
            disabled={decide.isPending}
            onClick={() => {
              setNote('')
              setApproveOpen(true)
            }}
          >
            <Check className="h-4 w-4" />
            Setujui
          </Button>
          <Button
            variant="outline"
            className="gap-1.5 text-rose-600 hover:bg-rose-50"
            disabled={decide.isPending}
            onClick={() => {
              setNote('')
              setRejectOpen(true)
            }}
          >
            <X className="h-4 w-4" />
            Tolak
          </Button>
        </div>
      </CardContent>

      <Dialog open={approveOpen} onOpenChange={setApproveOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Setujui Pengajuan</DialogTitle>
            <DialogDescription>
              Permohonan akan disetujui. Anda dapat menambahkan catatan (opsional).
            </DialogDescription>
          </DialogHeader>
          <textarea
            className={textareaClass}
            value={note}
            onChange={(event) => setNote(event.target.value)}
            placeholder="Catatan (opsional)"
          />
          <DialogFooter>
            <Button variant="outline" onClick={() => setApproveOpen(false)}>
              Batal
            </Button>
            <Button className="gap-1.5" disabled={decide.isPending} onClick={() => submit('APPROVED')}>
              <Check className="h-4 w-4" />
              {decide.isPending ? 'Memproses...' : 'Setujui'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={rejectOpen} onOpenChange={setRejectOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Tolak Pengajuan</DialogTitle>
            <DialogDescription>
              Permohonan akan ditolak. Catatan wajib diisi sebagai alasan penolakan.
            </DialogDescription>
          </DialogHeader>
          <textarea
            className={textareaClass}
            value={note}
            onChange={(event) => setNote(event.target.value)}
            placeholder="Alasan penolakan (wajib)"
          />
          <DialogFooter>
            <Button variant="outline" onClick={() => setRejectOpen(false)}>
              Batal
            </Button>
            <Button
              className="gap-1.5 bg-rose-600 hover:bg-rose-700"
              disabled={decide.isPending}
              onClick={() => submit('REJECTED')}
            >
              <X className="h-4 w-4" />
              {decide.isPending ? 'Memproses...' : 'Tolak'}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Card>
  )
}

function InfoItem({ label, value, className }: { label: string; value: string; className?: string }) {
  return (
    <div className={className}>
      <dt className="text-xs font-medium text-slate-500">{label}</dt>
      <dd className="mt-0.5 text-sm text-slate-900">{value || '—'}</dd>
    </div>
  )
}

function Loading({ label }: { label: string }) {
  return (
    <div className="flex items-center gap-2 text-sm text-slate-500">
      <Loader2 className="h-4 w-4 animate-spin" />
      {label}
    </div>
  )
}

function ErrorCard({ text }: { text: string }) {
  return (
    <Card>
      <CardContent className="p-6 text-sm text-rose-600">{text}</CardContent>
    </Card>
  )
}
