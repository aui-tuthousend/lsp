import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'
import type { EducationLevel, Gender } from '../-types'
import { controlClass, inputErrorClass, type StepProps } from './form'

const GENDER_OPTIONS: { value: Gender; label: string }[] = [
  { value: 'LAKI-LAKI', label: 'Laki-laki' },
  { value: 'PEREMPUAN', label: 'Perempuan' },
]

const EDUCATION_OPTIONS: { value: EducationLevel; label: string }[] = [
  { value: 'SD', label: 'SD' },
  { value: 'SMP', label: 'SMP' },
  { value: 'SMA_SMK', label: 'SMA/SMK' },
  { value: 'D1', label: 'D1' },
  { value: 'D2', label: 'D2' },
  { value: 'D3', label: 'D3' },
  { value: 'D4', label: 'D4' },
  { value: 'S1', label: 'S1' },
  { value: 'S2', label: 'S2' },
  { value: 'S3', label: 'S3' },
]

export function DataPribadiStep({ values, errors, onChange, today }: StepProps & { today: string }) {
  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
      <Field label="Nama Lengkap" htmlFor="full_name" required error={errors.full_name} className="sm:col-span-2">
        <Input
          id="full_name"
          maxLength={100}
          value={values.full_name}
          onChange={(event) => onChange('full_name', event.target.value)}
          placeholder="Sesuai KTP"
          className={inputErrorClass(Boolean(errors.full_name))}
        />
      </Field>

      <Field label="NIK" htmlFor="nik" required error={errors.nik} hint="16 digit sesuai KTP">
        <Input
          id="nik"
          inputMode="numeric"
          maxLength={16}
          value={values.nik}
          onChange={(event) => onChange('nik', event.target.value.replace(/\D/g, ''))}
          placeholder="3271xxxxxxxxxxxx"
          className={inputErrorClass(Boolean(errors.nik))}
        />
      </Field>

      <Field label="Nomor HP" htmlFor="phone" required error={errors.phone}>
        <Input
          id="phone"
          inputMode="tel"
          maxLength={12}
          value={values.phone}
          onChange={(event) => onChange('phone', event.target.value)}
          placeholder="08xxxxxxxxxx"
          className={inputErrorClass(Boolean(errors.phone))}
        />
      </Field>

      <Field label="Jenis Kelamin" required error={errors.gender} className="sm:col-span-2">
        <div className="grid grid-cols-2 gap-2">
          {GENDER_OPTIONS.map((option) => (
            <button
              key={option.value}
              type="button"
              onClick={() => onChange('gender', option.value)}
              className={cn(
                'flex items-center justify-center gap-2 rounded-lg border px-3 py-2.5 text-sm font-medium transition',
                values.gender === option.value
                  ? 'border-slate-900 bg-slate-900 text-white'
                  : 'border-slate-200 bg-white text-slate-600 hover:bg-slate-50',
              )}
            >
              {option.label}
            </button>
          ))}
        </div>
      </Field>

      <Field label="Tempat Lahir" htmlFor="birth_place" required error={errors.birth_place}>
        <Input
          id="birth_place"
          value={values.birth_place}
          onChange={(event) => onChange('birth_place', event.target.value)}
          placeholder="Kota kelahiran"
          className={inputErrorClass(Boolean(errors.birth_place))}
        />
      </Field>

      <Field label="Tanggal Lahir" htmlFor="birth_date" required error={errors.birth_date}>
        <Input
          id="birth_date"
          type="date"
          max={today}
          value={values.birth_date}
          onChange={(event) => onChange('birth_date', event.target.value)}
          className={inputErrorClass(Boolean(errors.birth_date))}
        />
      </Field>

      <Field label="Kewarganegaraan" htmlFor="nationality" required error={errors.nationality}>
        <Input
          id="nationality"
          value={values.nationality}
          onChange={(event) => onChange('nationality', event.target.value)}
          className={inputErrorClass(Boolean(errors.nationality))}
        />
      </Field>

      <Field label="Pekerjaan" htmlFor="occupation" required error={errors.occupation}>
        <Input
          id="occupation"
          value={values.occupation}
          onChange={(event) => onChange('occupation', event.target.value)}
          placeholder="mis. Karyawan Swasta"
          className={inputErrorClass(Boolean(errors.occupation))}
        />
      </Field>

      <Field label="Pendidikan Terakhir" htmlFor="education_level" required error={errors.education_level} className="sm:col-span-2">
        <select
          id="education_level"
          value={values.education_level}
          onChange={(event) => onChange('education_level', event.target.value as EducationLevel | '')}
          className={cn(controlClass, 'appearance-none pr-9', inputErrorClass(Boolean(errors.education_level)))}
        >
          <option value="">Pilih pendidikan</option>
          {EDUCATION_OPTIONS.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
      </Field>
    </div>
  )
}
