import * as React from 'react'
import { useRef } from 'react'
import { Camera, CreditCard, FileText, Info, Receipt, X } from 'lucide-react'
import {
  Attachment,
  AttachmentAction,
  AttachmentActions,
  AttachmentContent,
  AttachmentDescription,
  AttachmentMedia,
  AttachmentTitle,
  AttachmentTrigger,
} from '@/components/ui/attachment'
import type { StepProps } from './form'

type DocField = 'photo' | 'ktp_file' | 'cv_file' | 'npwp_file'

type DocConfig = {
  field: DocField
  label: string
  icon: React.ComponentType<{ className?: string }>
  accept: string
  hint: string
}

const DOCS: DocConfig[] = [
  { field: 'photo', label: 'Pas Foto', icon: Camera, accept: 'image/png,image/jpeg', hint: 'JPG/PNG, maks 2MB' },
  { field: 'ktp_file', label: 'Scan KTP', icon: CreditCard, accept: 'image/png,image/jpeg,application/pdf', hint: 'JPG/PNG/PDF, maks 2MB' },
  { field: 'cv_file', label: 'Curriculum Vitae', icon: FileText, accept: 'application/pdf', hint: 'PDF, maks 2MB' },
  { field: 'npwp_file', label: 'NPWP', icon: Receipt, accept: 'image/png,image/jpeg,application/pdf', hint: 'JPG/PNG/PDF, maks 2MB' },
]

function formatSize(bytes: number) {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${Math.round(bytes / 1024)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`
}

export function DokumenStep({ values, errors, onChange }: StepProps) {
  return (
    <div className="space-y-3">
      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
        {DOCS.map((doc) => (
          <DocAttachment
            key={doc.field}
            doc={doc}
            file={values[doc.field]}
            error={errors[doc.field]}
            onChange={(file) => onChange(doc.field, file)}
          />
        ))}
      </div>
      <p className="flex items-start gap-1.5 text-xs text-slate-400">
        <Info className="mt-0.5 h-3.5 w-3.5 shrink-0" />
        Semua dokumen bersifat opsional untuk saat ini dan dapat dilengkapi setelah pendaftaran.
      </p>
    </div>
  )
}

function DocAttachment({
  doc,
  file,
  error,
  onChange,
}: {
  doc: DocConfig
  file: File | null
  error?: string
  onChange: (file: File | null) => void
}) {
  const inputRef = useRef<HTMLInputElement>(null)
  const Icon = doc.icon
  const state = error ? 'error' : file ? 'done' : 'idle'

  return (
    <div className="min-w-0">
      <Attachment state={state} className="w-full">
        <AttachmentMedia>
          <Icon />
        </AttachmentMedia>
        <AttachmentContent>
          <AttachmentTitle>{doc.label}</AttachmentTitle>
          <AttachmentDescription>{error ?? (file ? `${file.name} · ${formatSize(file.size)}` : doc.hint)}</AttachmentDescription>
        </AttachmentContent>
        {file ? (
          <AttachmentActions>
            <AttachmentAction aria-label={`Hapus ${doc.label}`} onClick={() => onChange(null)}>
              <X />
            </AttachmentAction>
          </AttachmentActions>
        ) : null}
        <AttachmentTrigger aria-label={`Pilih ${doc.label}`} onClick={() => inputRef.current?.click()} />
      </Attachment>
      <input
        ref={inputRef}
        type="file"
        accept={doc.accept}
        className="hidden"
        onChange={(event) => onChange(event.target.files?.[0] ?? null)}
      />
    </div>
  )
}
