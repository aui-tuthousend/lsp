import { useEffect, useRef, useState } from 'react'
import { Eye, FileText, Loader2, Upload } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import type { Session } from '@/lib/auth/types'
import { cn } from '@/lib/utils'
import { useReplaceProfileFile, type ProfileFileField } from '../-hooks'
import type { AsesiProfile } from '../-types'

const MAX_FILE_BYTES = 2 * 1024 * 1024 // 2MB

type FileMeta = {
  field: ProfileFileField
  label: string
  accept: string
  mime: string[]
  hint: string
}

// Format berkas: foto JPG/PNG, dokumen lain wajib PDF.
const FILES: FileMeta[] = [
  { field: 'photo', label: 'Pas Foto', accept: 'image/png,image/jpeg', mime: ['image/jpeg', 'image/png'], hint: 'JPG/PNG, maks 2MB' },
  { field: 'ktp_file', label: 'Scan KTP', accept: 'application/pdf', mime: ['application/pdf'], hint: 'PDF, maks 2MB' },
  { field: 'cv_file', label: 'Curriculum Vitae', accept: 'application/pdf', mime: ['application/pdf'], hint: 'PDF, maks 2MB' },
  { field: 'npwp_file', label: 'NPWP', accept: 'application/pdf', mime: ['application/pdf'], hint: 'PDF, maks 2MB' },
]

const LINK_CLASS =
  'inline-flex h-9 items-center justify-center gap-1.5 rounded-md border border-slate-200 bg-white px-3 text-sm font-medium transition-colors hover:bg-slate-50'

type SupportingFilesProps = {
  session: Session | null
  profile: AsesiProfile
  /**
   * Dipanggil saat status unggah berubah. Dipakai pemanggil (mis. wizard) untuk
   * menonaktifkan tombol "Lanjut" selama unggah berjalan agar data konsisten.
   */
  onBusyChange?: (busy: boolean) => void
}

/** Berkas pendukung profil (pas foto, KTP, CV, NPWP) — lihat & ganti berkas. */
export function SupportingFiles({ session, profile, onBusyChange }: SupportingFilesProps) {
  const replace = useReplaceProfileFile(session)
  const [activeField, setActiveField] = useState<ProfileFileField | null>(null)
  const busy = replace.isPending

  useEffect(() => {
    onBusyChange?.(busy)
  }, [busy, onBusyChange])

  function handlePick(meta: FileMeta, file: File | null) {
    if (!file || busy) return
    if (file.size > MAX_FILE_BYTES) {
      toast.error(`${meta.label} maksimal 2MB`)
      return
    }
    if (!meta.mime.includes(file.type)) {
      toast.error(`${meta.label} harus ${meta.accept === 'application/pdf' ? 'PDF' : 'JPG/PNG'}`)
      return
    }
    setActiveField(meta.field)
    replace.mutate({ id: profile.id, field: meta.field, file }, { onSettled: () => setActiveField(null) })
  }

  return (
    <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
      {FILES.map((meta) => (
        <FileRow
          key={meta.field}
          meta={meta}
          url={profile[meta.field] ?? null}
          uploading={activeField === meta.field}
          disabled={busy}
          onPick={(file) => handlePick(meta, file)}
        />
      ))}
    </div>
  )
}

function FileRow({
  meta,
  url,
  uploading,
  disabled,
  onPick,
}: {
  meta: FileMeta
  url: string | null
  uploading: boolean
  disabled: boolean
  onPick: (file: File | null) => void
}) {
  const inputRef = useRef<HTMLInputElement>(null)

  return (
    <div className="flex items-center justify-between gap-3 rounded-lg border border-slate-200 bg-white px-3 py-3">
      <div className="min-w-0">
        <p className="flex items-center gap-1.5 text-sm font-medium text-slate-700">
          <FileText className="h-4 w-4 text-slate-400" />
          {meta.label}
        </p>
        <p className="mt-0.5 truncate text-xs text-slate-500">{url ? 'Sudah diunggah' : meta.hint}</p>
      </div>
      <div className="flex shrink-0 items-center gap-1.5">
        {url ? (
          <a
            href={url}
            target="_blank"
            rel="noreferrer"
            aria-disabled={disabled}
            className={cn(LINK_CLASS, disabled && 'pointer-events-none opacity-50')}
          >
            <Eye className="h-4 w-4" />
            Lihat
          </a>
        ) : null}
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="gap-1.5"
          disabled={disabled}
          onClick={() => inputRef.current?.click()}
        >
          {uploading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Upload className="h-4 w-4" />}
          {uploading ? 'Mengunggah...' : url ? 'Ganti' : 'Unggah'}
        </Button>
      </div>
      <input
        ref={inputRef}
        type="file"
        accept={meta.accept}
        className="hidden"
        onChange={(event) => {
          onPick(event.target.files?.[0] ?? null)
          event.target.value = ''
        }}
      />
    </div>
  )
}
