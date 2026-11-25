import type { ReactNode } from 'react'
import { AddressFields } from '@/components/ui/address-fields'
import { Checkbox } from '@/components/ui/checkbox'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
import type { Session } from '@/lib/auth/types'
import { SupportingFiles } from '../../profil/-components/supporting-files'
import {
  EDUCATION_OPTIONS,
  GENDER_OPTIONS,
  type OccupationFormValues,
  type ProfileFormValues,
} from '../../profil/-store'
import type { AsesiProfile, EducationLevel, Gender } from '../../profil/-types'
import { textareaClass } from '../-helpers'

export type StepIdentityProps = {
  profileValues: ProfileFormValues
  occupationValues: OccupationFormValues
  profile: AsesiProfile
  session: Session | null
  onProfileChange: <K extends keyof ProfileFormValues>(key: K, value: ProfileFormValues[K]) => void
  onOccupationChange: <K extends keyof OccupationFormValues>(key: K, value: OccupationFormValues[K]) => void
  /** Diteruskan ke blok berkas: true saat ada unggahan berjalan. */
  onFilesBusyChange?: (busy: boolean) => void
}

function Section({ title, description, children }: { title: string; description?: string; children: ReactNode }) {
  return (
    <section className="space-y-3">
      <div>
        <h3 className="text-sm font-semibold text-slate-900">{title}</h3>
        {description ? <p className="text-xs text-slate-500">{description}</p> : null}
      </div>
      {children}
    </section>
  )
}

/** Step 1 — data pribadi (profil asesi), pekerjaan, dan berkas identitas. */
export function StepIdentity({
  profileValues,
  occupationValues,
  profile,
  session,
  onProfileChange,
  onOccupationChange,
  onFilesBusyChange,
}: StepIdentityProps) {
  return (
    <div className="space-y-8">
      <Section title="Data Pribadi" description="Identitas sesuai dokumen resmi. Anda dapat mengubahnya di sini.">
        <div className="grid gap-4 sm:grid-cols-2">
          <Field label="Nama Lengkap" htmlFor="full_name" required>
            <Input id="full_name" value={profileValues.full_name} onChange={(e) => onProfileChange('full_name', e.target.value)} />
          </Field>
          <Field label="NIK" htmlFor="nik" required>
            <Input id="nik" value={profileValues.nik} onChange={(e) => onProfileChange('nik', e.target.value)} />
          </Field>
          <Field label="Nomor HP" htmlFor="phone" required>
            <Input id="phone" value={profileValues.phone} onChange={(e) => onProfileChange('phone', e.target.value)} />
          </Field>
          <Field label="Jenis Kelamin" htmlFor="gender" required>
            <Select id="gender" value={profileValues.gender} onChange={(e) => onProfileChange('gender', e.target.value as Gender | '')}>
              <option value="">— Pilih —</option>
              {GENDER_OPTIONS.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </Select>
          </Field>
          <Field label="Tempat Lahir" htmlFor="birth_place" required>
            <Input id="birth_place" value={profileValues.birth_place} onChange={(e) => onProfileChange('birth_place', e.target.value)} />
          </Field>
          <Field label="Tanggal Lahir" htmlFor="birth_date" required>
            <Input id="birth_date" type="date" value={profileValues.birth_date} onChange={(e) => onProfileChange('birth_date', e.target.value)} />
          </Field>
          <Field label="Kewarganegaraan" htmlFor="nationality" required>
            <Input id="nationality" value={profileValues.nationality} onChange={(e) => onProfileChange('nationality', e.target.value)} />
          </Field>
          <Field label="Pendidikan Terakhir" htmlFor="education_level" required>
            <Select
              id="education_level"
              value={profileValues.education_level}
              onChange={(e) => onProfileChange('education_level', e.target.value as EducationLevel | '')}
            >
              <option value="">— Pilih —</option>
              {EDUCATION_OPTIONS.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </Select>
          </Field>
          <Field label="Pekerjaan" htmlFor="occupation" required>
            <Input id="occupation" value={profileValues.occupation} onChange={(e) => onProfileChange('occupation', e.target.value)} />
          </Field>
        </div>
      </Section>

      <Section title="Alamat & Domisili">
        <div className="grid gap-4">
          <Field label="Alamat KTP" htmlFor="ktp_address" required>
            <textarea id="ktp_address" className={textareaClass} value={profileValues.ktp_address} onChange={(e) => onProfileChange('ktp_address', e.target.value)} />
          </Field>
          <Checkbox
            id="domicile_same_as_ktp"
            label="Alamat domisili sama dengan alamat KTP"
            checked={profileValues.domicile_same_as_ktp}
            onChange={(e) => onProfileChange('domicile_same_as_ktp', e.target.checked)}
          />
          {!profileValues.domicile_same_as_ktp ? (
            <Field label="Alamat Domisili" htmlFor="domicile_address" required>
              <textarea id="domicile_address" className={textareaClass} value={profileValues.domicile_address} onChange={(e) => onProfileChange('domicile_address', e.target.value)} />
            </Field>
          ) : null}
          <AddressFields values={profileValues} onChange={(key, value) => onProfileChange(key, value)} lockRegion />
        </div>
      </Section>

      <Section title="Pekerjaan Sekarang" description="Data pekerjaan saat ini (dapat dikosongkan bila belum bekerja).">
        <div className="grid gap-4 sm:grid-cols-2">
          <Field label="Nama Perusahaan" htmlFor="company_name">
            <Input id="company_name" value={occupationValues.company_name} onChange={(e) => onOccupationChange('company_name', e.target.value)} />
          </Field>
          <Field label="Jabatan" htmlFor="position">
            <Input id="position" value={occupationValues.position} onChange={(e) => onOccupationChange('position', e.target.value)} />
          </Field>
          <Field label="Lama Bekerja (bulan)" htmlFor="tenure_months">
            <Input id="tenure_months" type="number" min={0} value={occupationValues.tenure_months} onChange={(e) => onOccupationChange('tenure_months', e.target.value)} />
          </Field>
          <Field label="Jumlah Bawahan" htmlFor="subordinates_count">
            <Input id="subordinates_count" type="number" min={0} value={occupationValues.subordinates_count} onChange={(e) => onOccupationChange('subordinates_count', e.target.value)} />
          </Field>
          <Field label="Email Kantor" htmlFor="office_email">
            <Input id="office_email" type="email" value={occupationValues.office_email} onChange={(e) => onOccupationChange('office_email', e.target.value)} />
          </Field>
          <Field label="Telepon Kantor" htmlFor="office_phone">
            <Input id="office_phone" value={occupationValues.office_phone} onChange={(e) => onOccupationChange('office_phone', e.target.value)} />
          </Field>
          <AddressFields
            values={occupationValues}
            withVillage={false}
            lockRegion
            onChange={(key, value) => {
              if (key !== 'village') onOccupationChange(key, value)
            }}
            className="sm:col-span-2"
          />
          <Field label="Alamat Kantor" htmlFor="occ_address" className="sm:col-span-2">
            <textarea id="occ_address" className={textareaClass} value={occupationValues.address} onChange={(e) => onOccupationChange('address', e.target.value)} />
          </Field>
        </div>
      </Section>

      <Section
        title="Dokumen Pendukung"
        description="Pas foto (JPG/PNG) dan dokumen lain (PDF). Mengganti berkas akan menghapus berkas lama."
      >
        <SupportingFiles session={session} profile={profile} onBusyChange={onFilesBusyChange} />
      </Section>
    </div>
  )
}
