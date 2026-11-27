import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useMemo, useState } from 'react'
import { useQueryClient } from '@tanstack/react-query'
import { toast } from 'sonner'
import { CalendarClock, Check, Eye, FileText, Loader2, Lock, Video } from 'lucide-react'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import type { Session } from '@/lib/auth/types'
import { useAppEvent } from '@/lib/sse/use-app-event'
import { uploadFile } from '@/lib/upload'
import { cn } from '@/lib/utils'
import { QrPanel } from '../permohonan/-components/qr-panel'
import { useApplication, useSchemaRequirements } from '../permohonan/-hooks'
import { UnitAnswer, EVIDENCE_ACCEPT, type UnitAnswerValue } from './-components/unit-answer'
import { DECISION_LABELS, formatDate, formatDateTime, purposeLabel, readVerifyToken, writeVerifyToken } from './-helpers'
import { apl02Keys } from './-queryKeys'
import { useAssessment, usePutUnits, useSchemaTree, useSubmitAssessment } from './-hooks'
import type { AssessmentDetail, ElementNode } from './-types'

export const Route = createFileRoute('/asesi/asesmen-mandiri/$assessmentId')({
  component: AssessmentPage,
})

const MAX_FILE_BYTES = 2 * 1024 * 1024 // 2MB
const EVIDENCE_MIME = EVIDENCE_ACCEPT.split(',')

function AssessmentPage() {
  const { assessmentId } = Route.useParams()
  const { auth } = Route.useRouteContext()
  const query = useAssessment(assessmentId, auth)

  if (query.isLoading) {
    return (
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Loader2 className="h-4 w-4 animate-spin" />
        Memuat asesmen mandiri...
      </div>
    )
  }

  if (query.isError || !query.data) {
    return (
      <Card>
        <CardContent className="p-6 text-sm text-rose-600">Asesmen mandiri tidak ditemukan atau gagal dimuat.</CardContent>
      </Card>
    )
  }

  return <AssessmentForm detail={query.data} session={auth!} />
}

function AssessmentForm({ detail, session }: { detail: AssessmentDetail; session: Session }) {
  const assessmentId = detail.assessment.id
  const status = detail.assessment.status
  // Terkunci begitu diajukan atau QR terbit (meski belum discan) → read-only.
  const isReadOnly = detail.locked
  // Terkunci karena QR sudah terbit tapi belum discan (status masih DRAFT).
  const pendingQR = isReadOnly && status === 'DRAFT'

  const queryClient = useQueryClient()
  const putUnits = usePutUnits(assessmentId, session)
  const submit = useSubmitAssessment(assessmentId, session)
  const schemaTree = useSchemaTree(detail.application.certification_schema_id, session)

  const [answers, setAnswers] = useState<Record<string, UnitAnswerValue>>({})
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [actionPending, setActionPending] = useState(false)
  const [verifyToken, setVerifyToken] = useState<string | null>(null)
  const [confirmOpen, setConfirmOpen] = useState(false)
  const [hydratedId, setHydratedId] = useState<string | null>(null)

  // Elemen & KUK per unit kompetensi (dari struktur skema).
  const elementsByUnit = useMemo(() => {
    const map: Record<string, ElementNode[]> = {}
    for (const unit of schemaTree.data?.units ?? []) map[unit.id] = unit.elements
    return map
  }, [schemaTree.data])

  // Hidrasi jawaban dari server sekali per asesmen (agar berkas baru tidak ter-reset).
  useEffect(() => {
    if (hydratedId === detail.assessment.id) return
    const map: Record<string, UnitAnswerValue> = {}
    for (const unit of detail.units) {
      map[unit.id] = {
        isCompetent: unit.is_competent,
        files: unit.files.map((file) => ({ id: file.id, key: file.file_url, viewUrl: file.file_view_url })),
      }
    }
    setAnswers(map)
    setHydratedId(detail.assessment.id)
  }, [detail, hydratedId])

  // Pulihkan token QR tersimpan (server hanya menyimpan hash) agar QR tetap sama.
  useEffect(() => {
    setVerifyToken(readVerifyToken(assessmentId))
  }, [assessmentId])

  // Setelah asesi memverifikasi via QR di perangkat lain, segarkan status.
  useAppEvent('apl02.submitted', (event) => {
    if (event.data?.assessment_id !== assessmentId) return
    toast.success('Asesmen mandiri telah diverifikasi')
    void queryClient.invalidateQueries({ queryKey: apl02Keys.detail(assessmentId) })
  })

  const saving = actionPending || putUnits.isPending || submit.isPending

  // Berkas baru hanya ditambahkan ke state; diunggah sekaligus saat Simpan/Ajukan.
  function addFile(unitId: string, file: File) {
    if (file.size > MAX_FILE_BYTES) {
      toast.error('Berkas maksimal 2MB')
      return
    }
    if (!EVIDENCE_MIME.includes(file.type)) {
      toast.error('Berkas harus JPG/PNG/PDF')
      return
    }
    setAnswers((prev) => ({
      ...prev,
      [unitId]: {
        ...prev[unitId],
        files: [...(prev[unitId]?.files ?? []), { id: crypto.randomUUID(), file }],
      },
    }))
  }

  function removeFile(unitId: string, id: string) {
    setAnswers((prev) => ({
      ...prev,
      [unitId]: { ...prev[unitId], files: (prev[unitId]?.files ?? []).filter((item) => item.id !== id) },
    }))
  }

  function validate(): boolean {
    const next: Record<string, string> = {}
    for (const unit of detail.units) {
      const answer = answers[unit.id]
      if (!answer || answer.isCompetent === null) {
        next[unit.id] = 'Belum dijawab'
      } else if (answer.isCompetent && answer.files.length === 0) {
        next[unit.id] = 'Wajib unggah minimal 1 bukti'
      }
    }
    setErrors(next)
    return Object.keys(next).length === 0
  }

  async function save(): Promise<boolean> {
    if (!validate()) {
      toast.error('Lengkapi seluruh unit terlebih dahulu')
      return false
    }
    // Unggah semua berkas baru sekaligus, lalu kirim object key-nya.
    const items = await Promise.all(
      detail.units.map(async (unit) => {
        const value = answers[unit.id]
        const files: string[] = []
        for (const item of value?.files ?? []) {
          if (item.key) files.push(item.key)
          else if (item.file) files.push(await uploadFile(item.file, 'evidence'))
        }
        return {
          competency_unit_id: unit.competency_unit_id,
          is_competent: value?.isCompetent ?? null,
          files,
        }
      }),
    )
    await putUnits.mutateAsync({ items })
    return true
  }

  async function handleSave() {
    if (isReadOnly || actionPending) return
    setActionPending(true)
    try {
      await save()
    } catch {
      // error sudah ditangani toast pada mutation
    } finally {
      setActionPending(false)
    }
  }

  async function handleSubmit() {
    if (isReadOnly || actionPending) return
    setActionPending(true)
    try {
      if (!(await save())) return
      const result = await submit.mutateAsync(readVerifyToken(assessmentId))
      setVerifyToken(result.verify_token)
      writeVerifyToken(assessmentId, result.verify_token)
      window.scrollTo({ top: 0, behavior: 'smooth' })
    } catch {
      // error sudah ditangani toast
    } finally {
      setActionPending(false)
    }
  }

  // Terbitkan ulang / tampilkan kembali QR verifikasi tanpa mengubah jawaban
  // (dipakai saat token tersimpan hilang/kedaluwarsa).
  async function handleReissueQR() {
    if (actionPending) return
    setActionPending(true)
    try {
      const result = await submit.mutateAsync(readVerifyToken(assessmentId))
      setVerifyToken(result.verify_token)
      writeVerifyToken(assessmentId, result.verify_token)
      window.scrollTo({ top: 0, behavior: 'smooth' })
    } catch {
      // error sudah ditangani toast
    } finally {
      setActionPending(false)
    }
  }

  const verifyBadge =
    status === 'DRAFT'
      ? { label: 'Menunggu dipindai', className: 'bg-amber-100 text-amber-700' }
      : { label: 'Sudah diverifikasi', className: 'bg-emerald-100 text-emerald-700' }

  return (
    <div className="space-y-5">
      <div>
        <h1 className="text-lg font-semibold text-slate-900">Asesmen Mandiri</h1>
        <p className="text-sm text-slate-500">
          {detail.assessment.code} · {detail.units.length} unit kompetensi
        </p>
        {isReadOnly ? (
          <p className="mt-2 inline-flex items-center gap-1.5 rounded-full bg-slate-100 px-2.5 py-1 text-xs font-medium text-slate-600">
            <Lock className="h-3.5 w-3.5" />
            {pendingQR ? 'Formulir terkunci — menunggu verifikasi QR' : 'Asesmen sudah diajukan — hanya dapat dilihat'}
          </p>
        ) : null}
      </div>

      <Card>
        <CardContent className="p-4 sm:p-6">
          <h2 className="text-base font-semibold text-slate-900">Informasi Asesmen</h2>
          <dl className="mt-4 grid gap-4 sm:grid-cols-2">
            <InfoItem label="Nama Skema" value={detail.application.schema_name} />
            <InfoItem label="Nomor Skema" value={detail.application.schema_code} />
            <InfoItem label="Tujuan Asesmen" value={purposeLabel(detail.application.purpose)} />
            <InfoItem label="Tanggal Diajukan" value={formatDate(detail.application.data_created_at)} />
          </dl>
        </CardContent>
      </Card>

      <ApplicationDocuments applicationId={detail.assessment.application_id} session={session} />

      <Card>
        <CardContent className="space-y-4 p-4 sm:p-6">
          <div>
            <h2 className="text-base font-semibold text-slate-900">Pernyataan Kompetensi per Unit</h2>
            <p className="text-sm text-slate-500">{detail.units.length} unit kompetensi</p>
          </div>
          <div className="space-y-4">
            {detail.units.map((unit) => (
              <UnitAnswer
                key={unit.id}
                unit={unit}
                elements={elementsByUnit[unit.competency_unit_id] ?? []}
                value={answers[unit.id] ?? { isCompetent: null, files: [] }}
                disabled={isReadOnly || saving}
                error={errors[unit.id]}
                onChange={(value) => setAnswers((prev) => ({ ...prev, [unit.id]: value }))}
                onAddFile={(file) => addFile(unit.id, file)}
                onRemoveFile={(id) => removeFile(unit.id, id)}
              />
            ))}
          </div>
        </CardContent>
      </Card>

      {verifyToken && status !== 'REVIEWED' ? (
        <Card>
          <CardContent className="space-y-4 p-4 sm:p-6">
            <div className="flex flex-col items-center gap-2 text-center">
              <h2 className="text-base font-semibold text-slate-900">Verifikasi Asesmen Mandiri</h2>
              <p className="text-xs text-slate-500">Formulir APL-02 — Asesmen Mandiri</p>
              <span className={cn('inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium', verifyBadge.className)}>
                {verifyBadge.label}
              </span>
            </div>
            <QrPanel token={verifyToken} />
          </CardContent>
        </Card>
      ) : null}

      {!isReadOnly ? (
        <>
          <div className="flex items-center justify-end gap-3">
            <Button type="button" variant="outline" onClick={() => void handleSave()} disabled={saving}>
              {putUnits.isPending ? 'Menyimpan...' : 'Simpan'}
            </Button>
            <Button type="button" onClick={() => setConfirmOpen(true)} disabled={saving} className="gap-1.5">
              <Check className="h-4 w-4" />
              {actionPending || submit.isPending ? 'Memproses...' : 'Ajukan Asesmen Mandiri'}
            </Button>
          </div>

          <Dialog open={confirmOpen} onOpenChange={setConfirmOpen}>
            <DialogContent>
              <DialogHeader>
                <DialogTitle>Konfirmasi Pengajuan Asesmen Mandiri</DialogTitle>
                <DialogDescription>
                  Apakah Anda sudah memastikan seluruh pernyataan kompetensi dan berkas bukti yang Anda masukkan sudah
                  benar? Setelah diajukan, asesmen mandiri tidak dapat diubah lagi.
                </DialogDescription>
              </DialogHeader>
              <DialogFooter>
                <Button variant="outline" onClick={() => setConfirmOpen(false)}>
                  Batal
                </Button>
                <Button
                  className="gap-1.5"
                  onClick={() => {
                    setConfirmOpen(false)
                    void handleSubmit()
                  }}
                >
                  <Check className="h-4 w-4" />
                  Ya, Ajukan
                </Button>
              </DialogFooter>
            </DialogContent>
          </Dialog>
        </>
      ) : pendingQR && !verifyToken ? (
        <div className="flex items-center justify-end">
          <Button type="button" variant="outline" onClick={() => void handleReissueQR()} disabled={saving}>
            Tampilkan QR Verifikasi
          </Button>
        </div>
      ) : null}

      {status === 'REVIEWED' ? (
        <Card>
          <CardContent className="space-y-3 p-4 sm:p-6">
            <h2 className="text-base font-semibold text-slate-900">Hasil Tinjauan Asesor</h2>
            {detail.assessment.asesor_decision ? (
              <p className="text-sm text-slate-600">
                Rekomendasi:{' '}
                <span className="font-medium text-slate-900">{DECISION_LABELS[detail.assessment.asesor_decision]}</span>
              </p>
            ) : null}
            {detail.assessment.asesor_note ? (
              <p className="whitespace-pre-line text-sm text-slate-600">{detail.assessment.asesor_note}</p>
            ) : null}
            {detail.assessment.zoom_link ? (
              <div className="flex flex-wrap items-center gap-2 text-sm text-slate-600">
                <Video className="h-4 w-4 text-slate-400" />
                <a
                  href={detail.assessment.zoom_link}
                  target="_blank"
                  rel="noreferrer"
                  className="font-medium text-slate-900 underline underline-offset-2"
                >
                  Tautan pertemuan
                </a>
              </div>
            ) : null}
            {detail.assessment.zoom_date || detail.assessment.zoom_time ? (
              <p className="flex items-center gap-2 text-sm text-slate-600">
                <CalendarClock className="h-4 w-4 text-slate-400" />
                {formatDate(detail.assessment.zoom_date)}
                {detail.assessment.zoom_time ? ` · ${detail.assessment.zoom_time}` : ''}
              </p>
            ) : null}
            {detail.assessment.reviewed_at ? (
              <p className="text-xs text-slate-400">Ditinjau {formatDateTime(detail.assessment.reviewed_at)}</p>
            ) : null}
          </CardContent>
        </Card>
      ) : null}
    </div>
  )
}

/** Satu baris info read-only pada kartu "Informasi Asesmen". */
function InfoItem({ label, value }: { label: string; value: string }) {
  return (
    <div className="space-y-0.5">
      <dt className="text-xs font-medium text-slate-500">{label}</dt>
      <dd className="text-sm text-slate-900">{value || '—'}</dd>
    </div>
  )
}

/**
 * Accordion pratinjau (read-only) dokumen bukti kelengkapan yang diunggah pada
 * APL-01. Tiap dokumen ditampilkan sebagai baris bergaya "berkas pendukung profil".
 */
function ApplicationDocuments({ applicationId, session }: { applicationId: string; session: Session }) {
  const { data: application, isLoading: applicationLoading } = useApplication(applicationId, session)
  const schemaId = application?.data?.certification_schema_id ?? ''
  const { data: requirements, isLoading: requirementsLoading } = useSchemaRequirements(schemaId, session)

  const { hasFile, viewUrl } = useMemo(() => {
    const hasFile = new Set<string>()
    const viewUrl: Record<string, string> = {}
    for (const document of application?.documents ?? []) {
      if (!document.file_url) continue
      hasFile.add(document.document_type_id)
      if (document.file_view_url) viewUrl[document.document_type_id] = document.file_view_url
    }
    return { hasFile, viewUrl }
  }, [application])

  const uploaded = (requirements?.data ?? []).filter((requirement) => hasFile.has(requirement.document_type_id))
  const loading = applicationLoading || (Boolean(schemaId) && requirementsLoading)

  return (
    <Card>
      <CardContent className="p-0">
        <Accordion type="single" collapsible defaultValue="documents" className="px-4">
          <AccordionItem value="documents">
            <AccordionTrigger>
              <span className="flex flex-wrap items-baseline gap-2">
                <span className="flex items-center gap-1.5">
                  <FileText className="h-4 w-4 text-slate-500" />
                  Dokumen Relvevan APL-01
                </span>
                <span className="text-xs font-normal text-slate-400">{uploaded.length} berkas terunggah</span>
              </span>
            </AccordionTrigger>
            <AccordionContent>
              {loading ? (
                <div className="flex items-center gap-2 text-sm text-slate-500">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  Memuat dokumen...
                </div>
              ) : uploaded.length === 0 ? (
                <p className="text-sm text-slate-500">Belum ada dokumen yang diunggah pada APL-01.</p>
              ) : (
                <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                  {uploaded.map((requirement) => (
                    <DocumentRow
                      key={requirement.document_type_id}
                      label={requirement.document_type.label}
                      url={viewUrl[requirement.document_type_id]}
                    />
                  ))}
                </div>
              )}
            </AccordionContent>
          </AccordionItem>
        </Accordion>
      </CardContent>
    </Card>
  )
}

/** Baris dokumen read-only bergaya "berkas pendukung profil" (label + Lihat). */
function DocumentRow({ label, url }: { label: string; url?: string }) {
  return (
    <div className="flex items-center justify-between gap-3 rounded-lg border border-slate-200 bg-white px-3 py-3">
      <div className="min-w-0">
        <p className="flex items-center gap-1.5 text-sm font-medium text-slate-700">
          <FileText className="h-4 w-4 text-slate-400" />
          {label}
        </p>
        <p className="mt-0.5 truncate text-xs text-slate-500">Sudah diunggah</p>
      </div>
      {url ? (
        <a
          href={url}
          target="_blank"
          rel="noreferrer"
          className="inline-flex h-9 shrink-0 items-center justify-center gap-1.5 rounded-md border border-slate-200 bg-white px-3 text-sm font-medium transition-colors hover:bg-slate-50"
        >
          <Eye className="h-4 w-4" />
          Lihat
        </a>
      ) : null}
    </div>
  )
}
