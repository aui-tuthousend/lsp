import * as React from 'react'
import { Camera, CreditCard, FileText, Info, Receipt, Trash2 } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { StepProps } from './form'

export function DokumenStep({ values, errors, onChange }: StepProps) {
  return (
    <div className="space-y-3">
      <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
        <FileField
          id="photo"
          label="Pas Foto"
          icon={Camera}
          accept="image/png,image/jpeg"
          hint="JPG/PNG, maks 2MB"
          value={values.photo}
          onChange={(file) => onChange('photo', file)}
          error={errors.photo}
        />
        <FileField
          id="ktp_file"
          label="Scan KTP"
          icon={CreditCard}
          accept="image/png,image/jpeg,application/pdf"
          hint="JPG/PNG/PDF, maks 2MB"
          value={values.ktp_file}
          onChange={(file) => onChange('ktp_file', file)}
          error={errors.ktp_file}
        />
        <FileField
          id="cv_file"
          label="Curriculum Vitae"
          icon={FileText}
          accept="application/pdf"
          hint="PDF, maks 2MB"
          value={values.cv_file}
          onChange={(file) => onChange('cv_file', file)}
          error={errors.cv_file}
        />
        <FileField
          id="npwp_file"
          label="NPWP"
          icon={Receipt}
          accept="image/png,image/jpeg,application/pdf"
          hint="JPG/PNG/PDF, maks 2MB"
          value={values.npwp_file}
          onChange={(file) => onChange('npwp_file', file)}
          error={errors.npwp_file}
        />
      </div>
      <p className="flex items-start gap-1.5 text-xs text-slate-400">
        <Info className="mt-0.5 h-3.5 w-3.5 shrink-0" />
        Semua dokumen bersifat opsional untuk saat ini dan dapat dilengkapi setelah pendaftaran.
      </p>
    </div>
  )
}

function FileField({
  id,
  label,
  icon: Icon,
  accept,
  hint,
  value,
  onChange,
  error,
}: {
  id: string
  label: string
  icon: React.ComponentType<{ className?: string }>
  accept: string
  hint: string
  value: File | null
  onChange: (file: File | null) => void
  error?: string
}) {
  return (
    <div>
      <label
        htmlFor={id}
        className={cn(
          'group flex cursor-pointer flex-col items-center justify-center gap-1.5 rounded-xl border-2 border-dashed p-4 text-center transition',
          value ? 'border-emerald-300 bg-emerald-50/60' : 'border-slate-200 bg-slate-50 hover:border-slate-300 hover:bg-slate-100',
          error && !value && 'border-rose-300 bg-rose-50',
        )}
      >
        <Icon className={cn('h-5 w-5', value ? 'text-emerald-500' : 'text-slate-400 group-hover:text-slate-600')} />
        <span className="text-sm font-medium text-slate-700">{label}</span>
        <span className="max-w-full truncate text-xs text-slate-500">{value ? value.name : hint}</span>
        <input id={id} type="file" accept={accept} className="hidden" onChange={(event) => onChange(event.target.files?.[0] ?? null)} />
      </label>
      <div className="mt-1 flex items-center justify-between gap-2 px-0.5">
        {error ? <p className="text-xs font-medium text-rose-600">{error}</p> : <span />}
        {value ? (
          <button type="button" onClick={() => onChange(null)} className="inline-flex items-center gap-1 text-xs text-rose-600 hover:underline">
            <Trash2 className="h-3.5 w-3.5" />
            Hapus
          </button>
        ) : null}
      </div>
    </div>
  )
}
