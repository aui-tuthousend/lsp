import { AddressFields } from '@/components/ui/address-fields'
import { Field } from '@/components/ui/field'
import { cn } from '@/lib/utils'
import { controlClass, inputErrorClass, type StepProps } from './form'

export function AlamatDomisiliStep({ values, errors, onChange }: StepProps) {
  return (
    <div className="space-y-4">
      <Field label="Alamat sesuai KTP" htmlFor="ktp_address" required error={errors.ktp_address}>
        <textarea
          id="ktp_address"
          rows={3}
          value={values.ktp_address}
          onChange={(event) => onChange('ktp_address', event.target.value)}
          placeholder="Jalan, RT/RW, nomor rumah"
          className={cn(controlClass, 'h-auto resize-none py-2', inputErrorClass(Boolean(errors.ktp_address)))}
        />
      </Field>

      <div className="flex items-center justify-between gap-3 rounded-lg border border-slate-200 bg-slate-50 px-3 py-3">
        <div className="min-w-0">
          <p className="text-sm font-medium text-slate-700">Domisili sama dengan KTP</p>
          <p className="text-xs text-slate-500">Matikan bila alamat tinggal berbeda dari KTP.</p>
        </div>
        <Toggle checked={values.domicile_same_as_ktp} onChange={(value) => onChange('domicile_same_as_ktp', value)} />
      </div>

      {!values.domicile_same_as_ktp ? (
        <Field label="Alamat Domisili" htmlFor="domicile_address" required error={errors.domicile_address}>
          <textarea
            id="domicile_address"
            rows={3}
            value={values.domicile_address}
            onChange={(event) => onChange('domicile_address', event.target.value)}
            placeholder="Alamat tempat tinggal saat ini"
            className={cn(controlClass, 'h-auto resize-none py-2', inputErrorClass(Boolean(errors.domicile_address)))}
          />
        </Field>
      ) : null}

      <AddressFields values={values} errors={errors} onChange={(key, value) => onChange(key, value)} lockRegion />
    </div>
  )
}

function Toggle({ checked, onChange }: { checked: boolean; onChange: (value: boolean) => void }) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      onClick={() => onChange(!checked)}
      className={cn('relative inline-flex h-6 w-11 shrink-0 items-center rounded-full transition-colors', checked ? 'bg-emerald-500' : 'bg-slate-300')}
    >
      <span className={cn('inline-block h-5 w-5 transform rounded-full bg-white shadow transition-transform', checked ? 'translate-x-5' : 'translate-x-0.5')} />
    </button>
  )
}
