import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useEffect, useMemo, useRef, useState } from 'react'
import { toast } from 'sonner'
import { CalendarDays, Check, ChevronLeft, ChevronRight, FileText, Layers, Loader2, Lock, User } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Stepper } from '@/components/ui/stepper'
import type { Session } from '@/lib/auth/types'
import { useAppEvent } from '@/lib/sse/use-app-event'
import { uploadFile } from '@/lib/upload'
import { cn } from '@/lib/utils'
import {
  EMPTY_FORM,
  EMPTY_OCCUPATION_FORM,
  diffPatch,
  toForm,
  toOccupationBody,
  toOccupationForm,
  validate,
  validateOccupation,
  type OccupationFormValues,
  type ProfileFormValues,
} from '../profil/-store'
import { useAsesiProfileByUser, useOccupationByUser, useUpdateAsesiProfile, useUpsertOccupation } from '../profil/-hooks'
import { useApplication, usePutData, usePutDocuments, usePutSchedule, useSubmit, useUpdateStep, useSchemaOptions, useSchemaRequirements, useTukOptions } from './-hooks'
import { readVerifyToken, writeVerifyToken } from './-helpers'
import { QrPanel } from './-components/qr-panel'
import { StepCertification } from './-components/step-certification'
import { StepDocuments } from './-components/step-documents'
import { StepIdentity } from './-components/step-identity'
import { StepSchedule } from './-components/step-schedule'
import type { ApplicationDetail, AssessmentMode, AssessmentPurpose } from './-types'

export const Route = createFileRoute('/asesi/permohonan/$applicationId')({
  component: PermohonanPage,
})

const STEPS = [
  { title: 'Data Pribadi', description: 'Data pribadi, pekerjaan, dan berkas identitas.', icon: User },
  { title: 'Data Sertifikasi', description: 'Pilih skema dan tujuan asesmen.', icon: Layers },
  { title: 'Bukti Kelengkapan', description: 'Unggah dokumen sesuai skema yang dipilih.', icon: FileText },
  { title: 'Jadwal Asesmen', description: 'Ajukan metode dan tanggal pelaksanaan.', icon: CalendarDays },
] as const

function PermohonanPage() {
  const { applicationId } = Route.useParams()
  const { auth } = Route.useRouteContext()
  const detailQuery = useApplication(applicationId, auth)

  if (detailQuery.isLoading) {
    return (
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Loader2 className="h-4 w-4 animate-spin" />
        Memuat permohonan...
      </div>
    )
  }

  if (detailQuery.isError || !detailQuery.data) {
    return (
      <Card>
        <CardContent className="p-6 text-sm text-rose-600">Permohonan tidak ditemukan atau gagal dimuat.</CardContent>
      </Card>
    )
  }

  return <PermohonanWizard detail={detailQuery.data} session={auth!} />
}

function PermohonanWizard({ detail, session }: { detail: ApplicationDetail; session: Session }) {
  const navigate = useNavigate()
  const applicationId = detail.application.id

  const profileQuery = useAsesiProfileByUser(session.id, session)
  const occupationQuery = useOccupationByUser(session.id, session)
  const schemaOptions = useSchemaOptions(session)

  const updateStep = useUpdateStep(applicationId, session)
  const putData = usePutData(applicationId, session)
  const putDocuments = usePutDocuments(applicationId, session)
  const putSchedule = usePutSchedule(applicationId, session)
  const submit = useSubmit(applicationId, session)
  const updateProfile = useUpdateAsesiProfile(session)
  const upsertOccupation = useUpsertOccupation(session)

  // State lokal wizard (sumber kebenaran tetap server).
  const [step, setStep] = useState(0)
  const [schemaId, setSchemaId] = useState('')
  const [purpose, setPurpose] = useState<AssessmentPurpose | ''>('')
  const [mode, setMode] = useState<AssessmentMode | ''>('')
  const [scheduledDate, setScheduledDate] = useState('')
  const [tukId, setTukId] = useState('')
  const [docFiles, setDocFiles] = useState<Record<string, File | null>>({})
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [verifyToken, setVerifyToken] = useState<string | null>(null)

  const [profileValues, setProfileValues] = useState<ProfileFormValues>(EMPTY_FORM)
  const [occupationValues, setOccupationValues] = useState<OccupationFormValues>(EMPTY_OCCUPATION_FORM)
  const [profileBaseline, setProfileBaseline] = useState<ProfileFormValues>(EMPTY_FORM)
  const [filesUploading, setFilesUploading] = useState(false)
  // Menandai aksi Lanjut/Ajukan berjalan, termasuk fase unggah berkas yang
  // tidak tercatat sebagai mutation (agar tombol terkunci selama proses).
  const [actionPending, setActionPending] = useState(false)

  const [hydratedId, setHydratedId] = useState<string | null>(null)
  const hydratedProfileId = useRef<string | null>(null)
  const requirements = useSchemaRequirements(schemaId, session)
  const tukOptions = useTukOptions(session)
  const minDate = useMemo(() => new Date().toISOString().slice(0, 10), [])

  // Permohonan yang sudah diajukan (bukan draft) ditampilkan read-only.
  const isReadOnly = detail.application.status !== 'DRAFT'

  // Hidrasi state dari server sekali per permohonan.
  useEffect(() => {
    if (hydratedId === detail.application.id) return
    const readOnly = detail.application.status !== 'DRAFT'
    setStep(readOnly ? 0 : Math.min(Math.max(detail.application.current_step - 1, 0), STEPS.length - 1))
    setSchemaId(detail.data?.certification_schema_id ?? '')
    setPurpose(detail.data?.purpose ?? '')
    setMode(detail.schedule?.mode ?? '')
    setScheduledDate(detail.schedule?.scheduled_date ?? '')
    setTukId(detail.schedule?.tuk_id ?? '')
    setHydratedId(detail.application.id)
  }, [detail, hydratedId])

  // Hidrasi form profil sekali per profil (id) agar penggantian berkas tidak
  // mengembalikan perubahan field yang belum disimpan.
  useEffect(() => {
    const profile = profileQuery.data
    if (!profile || hydratedProfileId.current === profile.id) return
    hydratedProfileId.current = profile.id
    const form = toForm(profile)
    setProfileValues(form)
    setProfileBaseline(form)
  }, [profileQuery.data])

  useEffect(() => {
    if (occupationQuery.data === undefined) return
    const form = occupationQuery.data ? toOccupationForm(occupationQuery.data) : EMPTY_OCCUPATION_FORM
    setOccupationValues(form)
  }, [occupationQuery.data])

  // Pulihkan token QR tersimpan (server hanya menyimpan hash) agar QR tetap sama.
  useEffect(() => {
    setVerifyToken(readVerifyToken(applicationId))
  }, [applicationId])

  // Event SSE: setelah asesi memverifikasi via QR, arahkan ke beranda.
  useAppEvent('apl01.verified', (event) => {
    if (event.data?.application_id !== applicationId) return
    toast.success('Pendaftaran APL-01 telah selesai')
    void navigate({ to: '/asesi', replace: true })
  })

  const activeSchemas = useMemo(() => (schemaOptions.data?.data ?? []).filter((schema) => schema.is_active), [schemaOptions.data])
  const requirementItems = requirements.data?.data ?? []
  const tukItems = tukOptions.data?.data ?? []
  const existingDocuments = useMemo(() => {
    const map: Record<string, string | null> = {}
    for (const doc of detail.documents) map[doc.document_type_id] = doc.file_url ?? null
    return map
  }, [detail.documents])

  const verifyBadge = useMemo(() => {
    const signatures = detail.signatures.filter((sig) => sig.entity === 'permohonan' && sig.role === 'asesi')
    return signatures.some((sig) => sig.status === 'dipakai')
      ? { label: 'Sudah diverifikasi', className: 'bg-emerald-100 text-emerald-700' }
      : { label: 'Menunggu dipindai', className: 'bg-amber-100 text-amber-700' }
  }, [detail.signatures])

  const current = STEPS[step]
  const saving = actionPending || updateStep.isPending || putData.isPending || putDocuments.isPending || putSchedule.isPending || submit.isPending || updateProfile.isPending || upsertOccupation.isPending || filesUploading

  async function saveIdentity(): Promise<boolean> {
    const profileError = validate(profileValues)
    if (profileError) {
      toast.error(profileError)
      return false
    }
    const hasOccupation = occupationValues.company_name.trim().length > 0
    if (hasOccupation) {
      const occupationError = validateOccupation(occupationValues)
      if (occupationError) {
        toast.error(occupationError)
        return false
      }
    }
    if (!profileQuery.data) {
      toast.error('Profil belum termuat')
      return false
    }

    const patch = diffPatch(profileValues, profileBaseline)
    if (Object.keys(patch).length > 0) {
      await updateProfile.mutateAsync({ id: profileQuery.data.id, body: patch })
      setProfileBaseline(profileValues)
    }
    if (hasOccupation) {
      await upsertOccupation.mutateAsync({ userId: session.id, body: toOccupationBody(occupationValues) })
    }
    return true
  }

  async function saveCertification(): Promise<boolean> {
    if (!schemaId) {
      toast.error('Pilih skema sertifikasi')
      return false
    }
    if (!purpose) {
      toast.error('Pilih tujuan asesmen')
      return false
    }
    await putData.mutateAsync({ certification_schema_id: schemaId, purpose })
    return true
  }

  async function saveDocuments(): Promise<boolean> {
    const missing: Record<string, string> = {}
    for (const requirement of requirementItems) {
      if (!requirement.is_required) continue
      const hasNew = Boolean(docFiles[requirement.document_type_id])
      const hasExisting = Boolean(existingDocuments[requirement.document_type_id])
      if (!hasNew && !hasExisting) missing[requirement.document_type_id] = 'Dokumen wajib diunggah'
    }
    if (Object.keys(missing).length > 0) {
      setErrors(missing)
      toast.error('Lengkapi dokumen wajib terlebih dahulu')
      return false
    }
    setErrors({})

    const items = await Promise.all(
      requirementItems.map(async (requirement) => {
        const file = docFiles[requirement.document_type_id]
        const uploaded = file ? await uploadFile(file, 'document', requirement.document_type_id) : undefined
        return {
          document_type_id: requirement.document_type_id,
          file_url: uploaded ?? existingDocuments[requirement.document_type_id] ?? null,
        }
      }),
    )
    await putDocuments.mutateAsync({ items })
    setDocFiles({})
    return true
  }

  async function saveSchedule(): Promise<boolean> {
    const nextErrors: Record<string, string> = {}
    if (!mode) nextErrors.mode = 'Pilih metode pelaksanaan'
    if (!scheduledDate) nextErrors.scheduled_date = 'Pilih tanggal asesmen'
    if (mode === 'OFFLINE' && !tukId) nextErrors.tuk_id = 'Pilih TUK'
    if (Object.keys(nextErrors).length > 0) {
      setErrors(nextErrors)
      toast.error('Lengkapi jadwal asesmen terlebih dahulu')
      return false
    }
    setErrors({})
    await putSchedule.mutateAsync({
      mode: mode as AssessmentMode,
      scheduled_date: scheduledDate,
      tuk_id: mode === 'OFFLINE' ? tukId : null,
    })
    return true
  }

  async function goNext() {
    if (isReadOnly) {
      setStep((prev) => Math.min(prev + 1, STEPS.length - 1))
      window.scrollTo({ top: 0, behavior: 'smooth' })
      return
    }
    if (actionPending) return
    setActionPending(true)
    try {
      const ok =
        step === 0 ? await saveIdentity() : step === 1 ? await saveCertification() : step === 2 ? await saveDocuments() : false
      if (!ok) return
      const next = Math.min(step + 1, STEPS.length - 1)
      setStep(next)
      updateStep.mutate(next + 1)
      window.scrollTo({ top: 0, behavior: 'smooth' })
    } catch {
      // error sudah ditangani toast pada masing-masing mutation
    } finally {
      setActionPending(false)
    }
  }

  function goBack() {
    setErrors({})
    setStep((prev) => Math.max(prev - 1, 0))
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  async function handleSubmit() {
    if (isReadOnly || actionPending) return
    setActionPending(true)
    try {
      if (!(await saveSchedule())) return
      const result = await submit.mutateAsync(readVerifyToken(applicationId))
      setVerifyToken(result.verify_token)
      writeVerifyToken(applicationId, result.verify_token)
      window.scrollTo({ top: 0, behavior: 'smooth' })
    } catch {
      // error sudah ditangani toast
    } finally {
      setActionPending(false)
    }
  }

  function renderStep() {
    switch (step) {
      case 0:
        return profileQuery.isLoading || !profileQuery.data ? (
          <div className="flex items-center gap-2 text-sm text-slate-500">
            <Loader2 className="h-4 w-4 animate-spin" />
            Memuat profil...
          </div>
        ) : (
          <StepIdentity
            profileValues={profileValues}
            occupationValues={occupationValues}
            profile={profileQuery.data}
            session={session}
            onProfileChange={(key, value) => setProfileValues((prev) => ({ ...prev, [key]: value }))}
            onOccupationChange={(key, value) => setOccupationValues((prev) => ({ ...prev, [key]: value }))}
            onFilesBusyChange={setFilesUploading}
          />
        )
      case 1:
        return (
          <StepCertification
            schemas={activeSchemas}
            loading={schemaOptions.isLoading}
            schemaId={schemaId}
            purpose={purpose}
            errors={errors}
            onSchemaChange={(id) => {
              setSchemaId(id)
              setDocFiles({})
            }}
            onPurposeChange={setPurpose}
          />
        )
      case 2:
        return (
          <StepDocuments
            requirements={requirementItems}
            loading={requirements.isLoading}
            files={docFiles}
            existing={existingDocuments}
            errors={errors}
            disabled={saving}
            onPick={(documentTypeId, file) => setDocFiles((prev) => ({ ...prev, [documentTypeId]: file }))}
          />
        )
      default:
        return (
          <StepSchedule
            mode={mode}
            scheduledDate={scheduledDate}
            minDate={minDate}
            tukId={tukId}
            tuks={tukItems}
            tuksLoading={tukOptions.isLoading}
            errors={errors}
            onModeChange={(nextMode) => {
              setMode(nextMode)
              if (nextMode !== 'OFFLINE') setTukId('')
            }}
            onDateChange={setScheduledDate}
            onTukChange={setTukId}
          />
        )
    }
  }

  return (
    <div className="space-y-5">
      <div>
        <h1 className="text-lg font-semibold text-slate-900">Permohonan Sertifikasi</h1>
        <p className="text-sm text-slate-500">
          {detail.application.code} · Langkah {step + 1} / {STEPS.length}
        </p>
        {isReadOnly ? (
          <p className="mt-2 inline-flex items-center gap-1.5 rounded-full bg-slate-100 px-2.5 py-1 text-xs font-medium text-slate-600">
            <Lock className="h-3.5 w-3.5" />
            Permohonan sudah diajukan — hanya dapat dilihat
          </p>
        ) : null}
      </div>

      <Stepper steps={STEPS} current={step} />

      <Card>
        <CardContent className="space-y-5 p-4 sm:p-6">
          <div>
            <h2 className="text-base font-semibold text-slate-900">{current.title}</h2>
            <p className="text-sm text-slate-500">{current.description}</p>
          </div>

          <fieldset disabled={isReadOnly} className="m-0 min-w-0 border-0 p-0">
            {renderStep()}
          </fieldset>
        </CardContent>
      </Card>

      {step === STEPS.length - 1 && verifyToken ? (
        <Card>
          <CardContent className="space-y-4 p-4 sm:p-6">
            <div className="flex flex-col items-center gap-2 text-center">
              <h2 className="text-base font-semibold text-slate-900">Verifikasi Permohonan</h2>
              <span
                className={cn(
                  'inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium',
                  verifyBadge.className,
                )}
              >
                {verifyBadge.label}
              </span>
            </div>
            <QrPanel token={verifyToken} />
          </CardContent>
        </Card>
      ) : null}

      <div className="flex items-center justify-between gap-3">
        <Button type="button" variant="outline" onClick={goBack} disabled={step === 0 || saving} className="gap-1.5">
          <ChevronLeft className="h-4 w-4" />
          Sebelumnya
        </Button>
        {step < STEPS.length - 1 ? (
          <Button type="button" onClick={() => void goNext()} disabled={saving} className={cn('gap-1.5')}>
            {saving ? 'Menyimpan...' : 'Lanjut'}
            <ChevronRight className="h-4 w-4" />
          </Button>
        ) : (
          <Button type="button" onClick={() => void handleSubmit()} disabled={isReadOnly || saving} className="gap-1.5">
            <Check className="h-4 w-4" />
            {saving ? 'Memproses...' : 'Ajukan Permohonan'}
          </Button>
        )}
      </div>
    </div>
  )
}
