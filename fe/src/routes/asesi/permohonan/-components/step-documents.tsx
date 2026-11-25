import { Loader2 } from 'lucide-react'
import { toast } from 'sonner'
import { FileField } from './file-field'
import type { SchemaDocumentRequirement } from '../-types'

const FALLBACK_ACCEPT = 'image/jpeg,image/png,application/pdf'

function formatMime(mime: string): string {
  if (mime === 'application/pdf') return 'PDF'
  if (mime === 'image/jpeg') return 'JPG'
  if (mime === 'image/png') return 'PNG'
  return mime
}

/** Nilai `accept` untuk input file berdasarkan tipe berkas yang diizinkan jenis dokumen. */
function acceptFor(requirement: SchemaDocumentRequirement): string {
  const allowed = requirement.document_type.allowed_mime_types ?? []
  return allowed.length > 0 ? allowed.join(',') : FALLBACK_ACCEPT
}

export type StepDocumentsProps = {
  requirements: SchemaDocumentRequirement[]
  loading: boolean
  /** Berkas baru yang dipilih, per `document_type_id`. */
  files: Record<string, File | null>
  /** Object key tersimpan (dari detail permohonan), per `document_type_id`. */
  existing: Record<string, string | null>
  errors: Record<string, string>
  /** Nonaktifkan kontrol saat proses unggah/simpan berjalan. */
  disabled?: boolean
  onPick: (documentTypeId: string, file: File | null) => void
}

/** Step 3 — unggah bukti kelengkapan sesuai skema terpilih. */
export function StepDocuments({ requirements, loading, files, existing, errors, disabled, onPick }: StepDocumentsProps) {
  if (loading) {
    return (
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <Loader2 className="h-4 w-4 animate-spin" />
        Memuat persyaratan dokumen...
      </div>
    )
  }

  if (requirements.length === 0) {
    return <p className="text-sm text-slate-500">Skema ini tidak memerlukan dokumen tambahan.</p>
  }

  return (
    <div className="space-y-5">
      {requirements.map((requirement) => {
        const allowed = requirement.document_type.allowed_mime_types ?? []
        return (
          <FileField
            key={requirement.id}
            id={`doc-${requirement.document_type_id}`}
            label={requirement.document_type.label}
            accept={acceptFor(requirement)}
            required={requirement.is_required}
            file={files[requirement.document_type_id] ?? null}
            hasExisting={Boolean(existing[requirement.document_type_id])}
            error={errors[requirement.document_type_id]}
            disabled={disabled}
            onPick={(file) => {
              if (file && allowed.length > 0 && !allowed.includes(file.type)) {
                toast.error(`${requirement.document_type.label} harus ${allowed.map(formatMime).join(', ')}`)
                return
              }
              onPick(requirement.document_type_id, file)
            }}
          />
        )
      })}
    </div>
  )
}
