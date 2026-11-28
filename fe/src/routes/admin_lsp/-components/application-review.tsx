import { Link } from '@tanstack/react-router'
import { useMemo, type ReactNode } from 'react'
import { ArrowLeft, Loader2 } from 'lucide-react'
import { Card, CardContent } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import type { Session } from '@/lib/auth/types'
import { ApplicationStatusBadge } from '../../asesi/permohonan/-components/application-status-badge'
import { formatDate } from '../../asesi/permohonan/-helpers'
import { useApplication, useSchemaRequirements, useTukOptions } from '../../asesi/permohonan/-hooks'
import type { ApplicationDetail } from '../../asesi/permohonan/-types'
import { AssessmentStatusBadge } from '../../asesi/asesmen-mandiri/-components/assessment-status-badge'
import { useAssessment, useSchemaTree } from '../../asesi/asesmen-mandiri/-hooks'
import { purposeLabel } from '../../asesi/asesmen-mandiri/-helpers'
import type { AssessmentDetail, ElementNode } from '../../asesi/asesmen-mandiri/-types'
import { useAsesiProfileByUser } from '../../asesi/profil/-hooks'
import { EDUCATION_OPTIONS, GENDER_OPTIONS } from '../../asesi/profil/-store'
import type { AsesiProfile } from '../../asesi/profil/-types'
import { DocumentsAccordion } from './documents-accordion'
import { UnitSummary } from './unit-summary'

const GENDER_LABELS = new Map(GENDER_OPTIONS.map((option) => [option.value, option.label]))
const EDUCATION_LABELS = new Map(EDUCATION_OPTIONS.map((option) => [option.value, option.label]))
const MODE_LABELS: Record<string, string> = { ONLINE: 'Online', OFFLINE: 'Offline' }

export type ReviewContext = {
  applicationDetail: ApplicationDetail
  profile?: AsesiProfile
}

type ApplicationReviewProps = {
  /** ID asesmen mandiri (APL-02); permohonan APL-01 diturunkan darinya. */
  assessmentId: string
  session: Session
  /** Tautan "kembali" di atas halaman. */
  backTo: string
  backLabel: string
  /** Judul halaman (mis. "Verifikasi Pengajuan"). */
  title: string
  /** Panel aksi di bawah (keputusan admin / plotting asesor). */
  footer?: (ctx: ReviewContext) => ReactNode
}

type ApplicationReviewInnerProps = Omit<ApplicationReviewProps, 'assessmentId'> & {
  detail: AssessmentDetail
}

/**
 * Tampilan read-only pengajuan (APL-01 & APL-02) yang dipakai bersama oleh
 * admin_lsp dan manager_sertifikasi. Bagian bawah diserahkan lewat `footer`.
 */
export function ApplicationReview({ assessmentId, session, ...rest }: ApplicationReviewProps) {
  const assessmentQuery = useAssessment(assessmentId, session)

  if (assessmentQuery.isLoading) {
    return <Loading label="Memuat asesmen mandiri..." />
  }
  if (assessmentQuery.isError || !assessmentQuery.data) {
    return <ErrorCard text="Asesmen mandiri tidak ditemukan atau gagal dimuat." />
  }

  return <ApplicationReviewInner detail={assessmentQuery.data} session={session} {...rest} />
}

function ApplicationReviewInner({ detail, session, backTo, backLabel, title, footer }: ApplicationReviewInnerProps) {
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
          to={backTo}
          className="inline-flex items-center gap-1.5 text-sm text-slate-500 transition-colors hover:text-slate-900"
        >
          <ArrowLeft className="h-4 w-4" />
          {backLabel}
        </Link>
        <div className="flex flex-wrap items-center gap-2">
          <h1 className="text-lg font-semibold text-slate-900">{title}</h1>
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
          <ScheduleCard mode={schedule?.mode} date={schedule?.scheduled_date} tukName={tukName} />
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

      {footer?.({ applicationDetail, profile: profileQuery.data })}
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
            <InfoItem
              label="Tempat / Tanggal Lahir"
              value={`${profile.birth_place || '—'}, ${formatDate(profile.birth_date)}`}
            />
            <InfoItem
              label="Pendidikan Terakhir"
              value={EDUCATION_LABELS.get(profile.education_level) ?? profile.education_level}
            />
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
