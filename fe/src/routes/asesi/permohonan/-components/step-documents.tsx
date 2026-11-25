import { Loader2 } from 'lucide-react'
import { FileField } from './file-field'
import type { SchemaDocumentRequirement } from '../-types'

export type StepDocumentsProps = {
  requirements: SchemaDocumentRequirement[]
  loading: boolean
  /** Berkas baru yang dipilih, per `document_type_id`. */
  files: Record<string, File | null>
  /** Object key tersimpan (dari detail permohonan), per `document_type_id`. */
  existing: Record<string, string | null>
  errors: Record<string, string>
  onPick: (documentTypeId: string, file: File | null) => void
}

/** Step 3 — unggah bukti kelengkapan sesuai skema terpilih. */
export function StepDocuments({ requirements, loading, files, existing, errors, onPick }: StepDocumentsProps) {
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
      {requirements.map((requirement) => (
        <FileField
          key={requirement.id}
          id={`doc-${requirement.document_type_id}`}
          label={requirement.document_type.label}
          accept="image/jpeg,image/png,application/pdf"
          required={requirement.is_required}
          file={files[requirement.document_type_id] ?? null}
          hasExisting={Boolean(existing[requirement.document_type_id])}
          error={errors[requirement.document_type_id]}
          onPick={(file) => onPick(requirement.document_type_id, file)}
        />
      ))}
    </div>
  )
}
