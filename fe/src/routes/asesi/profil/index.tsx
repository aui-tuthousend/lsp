import { useEffect, useState, type FormEvent } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { Loader2, Save } from 'lucide-react'
import { toast } from 'sonner'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'
import { useAsesiProfileByUser, useMe, useUpdateAsesiProfile } from './-hooks'
import {
  EDUCATION_OPTIONS,
  EMPTY_FORM,
  GENDER_OPTIONS,
  diffPatch,
  toForm,
  validate,
  type ProfileFormValues,
} from './-store'
import type { EducationLevel } from './-types'

export const Route = createFileRoute('/asesi/profil/')({
  component: ProfilPage,
})

const controlClass =
  'flex h-10 w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-slate-400 focus:ring-2 focus:ring-slate-200'

function ProfilPage() {
  const { auth } = Route.useRouteContext()
  const me = useMe(auth)
  const profileQuery = useAsesiProfileByUser(auth?.id ?? '', auth)
  const update = useUpdateAsesiProfile(auth)

  const [values, setValues] = useState<ProfileFormValues>(EMPTY_FORM)
  const [baseline, setBaseline] = useState<ProfileFormValues>(EMPTY_FORM)
  const today = new Date().toISOString().slice(0, 10)

  useEffect(() => {
    if (!profileQuery.data) return
    const form = toForm(profileQuery.data)
    setValues(form)
    setBaseline(form)
  }, [profileQuery.data])

  const patch = diffPatch(values, baseline)
  const isDirty = Object.keys(patch).length > 0

  function set<K extends keyof ProfileFormValues>(key: K, value: ProfileFormValues[K]) {
    setValues((prev) => ({ ...prev, [key]: value }))
  }

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const message = validate(values)
    if (message) {
      toast.error(message)
      return
    }
    const id = profileQuery.data?.id
    if (!id) {
      toast.error('Profil belum tersedia')
      return
    }
    if (!isDirty) {
      toast.info('Tidak ada perubahan untuk disimpan')
      return
    }
    update.mutate({ id, body: patch })
  }

  const account = me.data

  return (
    <form onSubmit={submit} className="space-y-4">
      {/* Tabel users — data akun (accordion, hanya baca). */}
      <Card className="px-6">
        <Accordion type="single" collapsible defaultValue="account">
          <AccordionItem value="account">
            <AccordionTrigger className="text-base">Data Akun</AccordionTrigger>
            <AccordionContent>
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
                <InfoItem label="Email" value={account?.email} />
                <InfoItem label="Nama Akun" value={account?.name} />
                <InfoItem label="Role" value={account?.role} />
                <InfoItem label="Nomor Registrasi" value={account?.nomor_registrasi} />
                <InfoItem label="Telepon" value={account?.phone} />
                <InfoItem label="Verifikasi Email" value={account?.email_verified_at ? 'Terverifikasi' : 'Belum terverifikasi'} />
              </div>
            </AccordionContent>
          </AccordionItem>
        </Accordion>
      </Card>

      {/* Tabel asesi_profiles — data profil, bisa diubah lewat PATCH. */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Data Profil Asesi</CardTitle>
          <p className="text-xs text-slate-500">Berasal dari tabel asesi_profiles.</p>
        </CardHeader>
        <CardContent>
          {profileQuery.isLoading ? (
            <div className="flex items-center gap-2 py-6 text-sm text-slate-500">
              <Loader2 className="h-4 w-4 animate-spin" />
              Memuat profil...
            </div>
          ) : !profileQuery.data ? (
            <p className="py-6 text-sm text-slate-500">Profil asesi belum tersedia.</p>
          ) : (
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <Field label="Nama Lengkap" htmlFor="full_name" required className="sm:col-span-2">
                <Input id="full_name" maxLength={255} value={values.full_name} onChange={(e) => set('full_name', e.target.value)} />
              </Field>

              <Field label="NIK" htmlFor="nik" required hint="16 digit sesuai KTP">
                <Input
                  id="nik"
                  inputMode="numeric"
                  maxLength={16}
                  value={values.nik}
                  onChange={(e) => set('nik', e.target.value.replace(/\D/g, ''))}
                />
              </Field>

              <Field label="Nomor HP" htmlFor="phone" required>
                <Input id="phone" inputMode="tel" maxLength={20} value={values.phone} onChange={(e) => set('phone', e.target.value)} />
              </Field>

              <Field label="Jenis Kelamin" required className="sm:col-span-2">
                <div className="grid grid-cols-2 gap-2">
                  {GENDER_OPTIONS.map((option) => (
                    <button
                      key={option.value}
                      type="button"
                      onClick={() => set('gender', option.value)}
                      className={cn(
                        'flex items-center justify-center rounded-lg border px-3 py-2.5 text-sm font-medium transition',
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

              <Field label="Tempat Lahir" htmlFor="birth_place" required>
                <Input id="birth_place" value={values.birth_place} onChange={(e) => set('birth_place', e.target.value)} />
              </Field>

              <Field label="Tanggal Lahir" htmlFor="birth_date" required>
                <Input id="birth_date" type="date" max={today} value={values.birth_date} onChange={(e) => set('birth_date', e.target.value)} />
              </Field>

              <Field label="Kewarganegaraan" htmlFor="nationality" required>
                <Input id="nationality" value={values.nationality} onChange={(e) => set('nationality', e.target.value)} />
              </Field>

              <Field label="Pekerjaan" htmlFor="occupation" required>
                <Input id="occupation" value={values.occupation} onChange={(e) => set('occupation', e.target.value)} />
              </Field>

              <Field label="Pendidikan Terakhir" htmlFor="education_level" required className="sm:col-span-2">
                <select
                  id="education_level"
                  value={values.education_level}
                  onChange={(e) => set('education_level', e.target.value as EducationLevel | '')}
                  className={cn(controlClass, 'appearance-none pr-9')}
                >
                  <option value="">Pilih pendidikan</option>
                  {EDUCATION_OPTIONS.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>
              </Field>

              <Field label="Alamat sesuai KTP" htmlFor="ktp_address" required className="sm:col-span-2">
                <textarea
                  id="ktp_address"
                  rows={3}
                  value={values.ktp_address}
                  onChange={(e) => set('ktp_address', e.target.value)}
                  className={cn(controlClass, 'h-auto resize-none py-2')}
                />
              </Field>

              <div className="flex items-center justify-between gap-3 rounded-lg border border-slate-200 bg-slate-50 px-3 py-3 sm:col-span-2">
                <div className="min-w-0">
                  <p className="text-sm font-medium text-slate-700">Domisili sama dengan KTP</p>
                  <p className="text-xs text-slate-500">Matikan bila alamat tinggal berbeda dari KTP.</p>
                </div>
                <button
                  type="button"
                  role="switch"
                  aria-checked={values.domicile_same_as_ktp}
                  onClick={() => set('domicile_same_as_ktp', !values.domicile_same_as_ktp)}
                  className={cn(
                    'relative inline-flex h-6 w-11 shrink-0 items-center rounded-full transition-colors',
                    values.domicile_same_as_ktp ? 'bg-emerald-500' : 'bg-slate-300',
                  )}
                >
                  <span
                    className={cn(
                      'inline-block h-5 w-5 transform rounded-full bg-white shadow transition-transform',
                      values.domicile_same_as_ktp ? 'translate-x-5' : 'translate-x-0.5',
                    )}
                  />
                </button>
              </div>

              {!values.domicile_same_as_ktp ? (
                <Field label="Alamat Domisili" htmlFor="domicile_address" required className="sm:col-span-2">
                  <textarea
                    id="domicile_address"
                    rows={3}
                    value={values.domicile_address}
                    onChange={(e) => set('domicile_address', e.target.value)}
                    className={cn(controlClass, 'h-auto resize-none py-2')}
                  />
                </Field>
              ) : null}

              <Field label="Kode Pos" htmlFor="postal_code">
                <Input
                  id="postal_code"
                  inputMode="numeric"
                  maxLength={5}
                  value={values.postal_code}
                  onChange={(e) => set('postal_code', e.target.value.replace(/\D/g, ''))}
                />
              </Field>

              <Field label="Provinsi" htmlFor="province" required>
                <Input id="province" value={values.province} onChange={(e) => set('province', e.target.value)} />
              </Field>

              <Field label="Kota / Kabupaten" htmlFor="city" required>
                <Input id="city" value={values.city} onChange={(e) => set('city', e.target.value)} />
              </Field>

              <Field label="Kecamatan" htmlFor="district" required>
                <Input id="district" value={values.district} onChange={(e) => set('district', e.target.value)} />
              </Field>

              <Field label="Desa / Kelurahan" htmlFor="village" required>
                <Input id="village" value={values.village} onChange={(e) => set('village', e.target.value)} />
              </Field>
            </div>
          )}
        </CardContent>
      </Card>

      <div className="flex justify-end">
        <Button type="submit" disabled={update.isPending || !profileQuery.data || !isDirty} className="gap-2">
          {update.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
          {update.isPending ? 'Menyimpan...' : 'Simpan Perubahan'}
        </Button>
      </div>
    </form>
  )
}

function InfoItem({ label, value }: { label: string; value?: string | null }) {
  return (
    <div className="min-w-0">
      <div className="text-xs text-slate-500">{label}</div>
      <div className="truncate text-sm font-medium text-slate-900">{value ? value : '-'}</div>
    </div>
  )
}
