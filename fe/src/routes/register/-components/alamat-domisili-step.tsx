import { Loader2, Search } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'
import { controlClass, inputErrorClass, type StepProps } from './form'

export function AlamatDomisiliStep({
  values,
  errors,
  onChange,
  postalLoading,
  onLookupPostal,
}: StepProps & { postalLoading: boolean; onLookupPostal: () => void }) {
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

      <div className="rounded-xl border border-dashed border-slate-300 bg-slate-50/70 p-3 sm:p-4">
        <div className="mb-3 flex items-center gap-2 text-sm font-medium text-slate-700">
          <Search className="h-4 w-4 text-slate-400" />
          Isi kode pos, wilayah terisi otomatis
        </div>
        <div className="flex flex-col gap-2 sm:flex-row">
          <Input
            id="postal_code"
            inputMode="numeric"
            maxLength={5}
            value={values.postal_code}
            onChange={(event) => onChange('postal_code', event.target.value.replace(/\D/g, ''))}
            placeholder="mis. 40115"
            className={cn('sm:max-w-40', inputErrorClass(Boolean(errors.postal_code)))}
          />
          <Button type="button" onClick={onLookupPostal} disabled={postalLoading} className="gap-2 sm:w-auto">
            {postalLoading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Search className="h-4 w-4" />}
            {postalLoading ? 'Mencari...' : 'Cari wilayah'}
          </Button>
        </div>
        {errors.postal_code ? <p className="mt-1.5 text-xs font-medium text-rose-600">{errors.postal_code}</p> : null}

        <div className="mt-3 grid grid-cols-1 gap-3 sm:grid-cols-2">
          <Field label="Provinsi" error={errors.province} hint="Terisi otomatis, bisa diubah">
            <Input
              value={values.province}
              onChange={(event) => onChange('province', event.target.value)}
              placeholder="Terisi otomatis"
              className={inputErrorClass(Boolean(errors.province))}
            />
          </Field>
          <Field label="Kota / Kabupaten" error={errors.city} hint="Terisi otomatis, bisa diubah">
            <Input
              value={values.city}
              onChange={(event) => onChange('city', event.target.value)}
              placeholder="Terisi otomatis"
              className={inputErrorClass(Boolean(errors.city))}
            />
          </Field>
          <Field label="Kecamatan" error={errors.district} hint="Terisi otomatis, bisa diubah">
            <Input
              value={values.district}
              onChange={(event) => onChange('district', event.target.value)}
              placeholder="Terisi otomatis"
              className={inputErrorClass(Boolean(errors.district))}
            />
          </Field>
          <Field label="Desa / Kelurahan" error={errors.village} hint="Terisi otomatis, bisa diubah">
            <Input
              value={values.village}
              onChange={(event) => onChange('village', event.target.value)}
              placeholder="Terisi otomatis"
              className={inputErrorClass(Boolean(errors.village))}
            />
          </Field>
        </div>
      </div>
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
