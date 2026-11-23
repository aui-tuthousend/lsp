import { useEffect, useState, type FormEvent } from 'react'
import { Loader2, Save } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
import type { Session } from '@/lib/auth/types'
import { useAsesorAccount, useCreateAsesor, useTukOptions, useUpdateAsesor } from '../-hooks'
import { validateAsesorForm, type AsesorFormMode, type AsesorFormValues } from '../-schemas'
import type { AsesorAccount, AsesorProfile, CreateAsesorBody, UpdateAsesorAccountBody, UpdateAsesorProfileBody } from '../-types'

const EMPTY_FORM: AsesorFormValues = {
  name: '',
  email: '',
  password: '',
  phone: '',
  nomor_registrasi: '',
  tuk_id: '',
  province: '',
  city: '',
  daily_assignment_limit: 0,
}

/** Susun nilai form dari profil (list) + akun (detail user). */
function buildForm(profile: AsesorProfile | null, account: AsesorAccount | null): AsesorFormValues {
  return {
    name: account?.name ?? profile?.name ?? '',
    email: account?.email ?? profile?.email ?? '',
    password: '',
    phone: account?.phone ?? '',
    nomor_registrasi: account?.nomor_registrasi ?? '',
    tuk_id: profile?.tuk_id ?? '',
    province: profile?.province ?? '',
    city: profile?.city ?? '',
    daily_assignment_limit: profile?.daily_assignment_limit ?? 0,
  }
}

function createAccountBody(values: AsesorFormValues): CreateAsesorBody {
  return {
    email: values.email,
    password: values.password,
    name: values.name,
    phone: values.phone || undefined,
    nomor_registrasi: values.nomor_registrasi || undefined,
  }
}

function createProfileBody(values: AsesorFormValues): UpdateAsesorProfileBody {
  const body: UpdateAsesorProfileBody = { daily_assignment_limit: values.daily_assignment_limit }
  if (values.tuk_id) body.tuk_id = values.tuk_id
  if (values.province) body.province = values.province
  if (values.city) body.city = values.city
  return body
}

/** Bandingkan dengan nilai awal; kirim hanya field akun yang berubah (partial update). */
function accountPatch(values: AsesorFormValues, baseline: AsesorFormValues): UpdateAsesorAccountBody {
  const patch: UpdateAsesorAccountBody = {}
  if (values.name !== baseline.name) patch.name = values.name
  if (values.phone !== baseline.phone) patch.phone = values.phone
  if (values.nomor_registrasi !== baseline.nomor_registrasi) patch.nomor_registrasi = values.nomor_registrasi
  if (values.password) patch.password = values.password
  return patch
}

/** Bandingkan dengan nilai awal; kirim hanya field profil yang berubah (partial update). */
function profilePatch(values: AsesorFormValues, baseline: AsesorFormValues): UpdateAsesorProfileBody {
  const patch: UpdateAsesorProfileBody = {}
  if (values.tuk_id && values.tuk_id !== baseline.tuk_id) patch.tuk_id = values.tuk_id
  if (values.province !== baseline.province) patch.province = values.province
  if (values.city !== baseline.city) patch.city = values.city
  if (values.daily_assignment_limit !== baseline.daily_assignment_limit) patch.daily_assignment_limit = values.daily_assignment_limit
  return patch
}

export function AsesorFormDialog({
  open,
  onOpenChange,
  asesor,
  session,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  asesor?: AsesorProfile | null
  session: Session | null
}) {
  const isEdit = Boolean(asesor)
  const mode: AsesorFormMode = isEdit ? 'edit' : 'create'
  const create = useCreateAsesor(session)
  const update = useUpdateAsesor(session)
  const pending = create.isPending || update.isPending

  const { data: tukPage } = useTukOptions(session)
  const tukOptions = tukPage?.data ?? []
  const { data: account } = useAsesorAccount(asesor?.user_id ?? '', session)

  const [values, setValues] = useState<AsesorFormValues>(EMPTY_FORM)
  const [baseline, setBaseline] = useState<AsesorFormValues>(EMPTY_FORM)
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    if (!open) return
    const form = buildForm(asesor ?? null, account ?? null)
    setValues(form)
    setBaseline(form)
    setErrors({})
  }, [open, asesor, account])

  function set<K extends keyof AsesorFormValues>(key: K, value: AsesorFormValues[K]) {
    setValues((prev) => ({ ...prev, [key]: value }))
    setErrors((prev) => {
      if (!(key in prev)) return prev
      const next = { ...prev }
      delete next[key as string]
      return next
    })
  }

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const result = validateAsesorForm(values, mode)
    if (!result.ok) {
      setErrors(result.errors)
      toast.error('Lengkapi data yang belum valid')
      return
    }
    const data = result.data

    if (asesor) {
      const account = accountPatch(data, baseline)
      const profile = profilePatch(data, baseline)
      if (Object.keys(account).length === 0 && Object.keys(profile).length === 0) {
        toast.info('Tidak ada perubahan untuk disimpan')
        return
      }
      update.mutate({ userId: asesor.user_id, account, profile }, { onSuccess: () => onOpenChange(false) })
      return
    }

    create.mutate(
      { account: createAccountBody(data), profile: createProfileBody(data) },
      { onSuccess: () => onOpenChange(false) },
    )
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEdit ? 'Ubah Asesor' : 'Tambah Asesor'}</DialogTitle>
          <DialogDescription>
            {isEdit ? 'Perbarui akun dan profil asesor.' : 'Buat akun asesor beserta profilnya.'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={submit} className="space-y-5">
          <section className="space-y-4">
            <h3 className="text-sm font-semibold text-slate-900">Data Akun</h3>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <Field label="Nama" htmlFor="name" required error={errors.name} className="sm:col-span-2">
                <Input id="name" maxLength={255} value={values.name} onChange={(e) => set('name', e.target.value)} placeholder="mis. Budi Santoso" />
              </Field>

              <Field label="Email" htmlFor="email" required error={errors.email}>
                <Input
                  id="email"
                  type="email"
                  maxLength={255}
                  value={values.email}
                  onChange={(e) => set('email', e.target.value)}
                  placeholder="asesor@contoh.id"
                  disabled={isEdit}
                />
              </Field>

              <Field
                label="Password"
                htmlFor="password"
                required={!isEdit}
                error={errors.password}
                hint={isEdit ? 'Kosongkan bila tidak diubah' : undefined}
              >
                <Input
                  id="password"
                  type="password"
                  maxLength={72}
                  value={values.password}
                  onChange={(e) => set('password', e.target.value)}
                  placeholder="Minimal 8 karakter"
                  autoComplete="new-password"
                />
              </Field>

              <Field label="Telepon" htmlFor="phone" error={errors.phone} hint="Opsional">
                <Input id="phone" inputMode="tel" maxLength={20} value={values.phone} onChange={(e) => set('phone', e.target.value)} placeholder="08xxxxxxxxxx" />
              </Field>

              <Field label="Nomor Registrasi" htmlFor="nomor_registrasi" error={errors.nomor_registrasi} hint="Opsional">
                <Input
                  id="nomor_registrasi"
                  maxLength={50}
                  value={values.nomor_registrasi}
                  onChange={(e) => set('nomor_registrasi', e.target.value)}
                />
              </Field>
            </div>
          </section>

          <section className="space-y-4 border-t border-slate-100 pt-4">
            <h3 className="text-sm font-semibold text-slate-900">Profil Asesor</h3>
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <Field label="TUK" htmlFor="tuk_id" error={errors.tuk_id} hint="Opsional" className="sm:col-span-2">
                <Select id="tuk_id" value={values.tuk_id} onChange={(e) => set('tuk_id', e.target.value)}>
                  <option value="">— Tanpa TUK —</option>
                  {tukOptions.map((tuk) => (
                    <option key={tuk.id} value={tuk.id}>
                      {tuk.name}
                    </option>
                  ))}
                </Select>
              </Field>

              <Field label="Provinsi" htmlFor="province" error={errors.province} hint="Opsional">
                <Input id="province" maxLength={120} value={values.province} onChange={(e) => set('province', e.target.value)} />
              </Field>

              <Field label="Kota / Kabupaten" htmlFor="city" error={errors.city} hint="Opsional">
                <Input id="city" maxLength={120} value={values.city} onChange={(e) => set('city', e.target.value)} />
              </Field>

              <Field
                label="Batas Penugasan Harian"
                htmlFor="daily_assignment_limit"
                error={errors.daily_assignment_limit}
                hint="0 = tanpa batas"
              >
                <Input
                  id="daily_assignment_limit"
                  type="number"
                  min={0}
                  value={values.daily_assignment_limit}
                  onChange={(e) => set('daily_assignment_limit', Number(e.target.value))}
                />
              </Field>
            </div>
          </section>

          <DialogFooter>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={pending}>
              Batal
            </Button>
            <Button type="submit" disabled={pending} className="gap-2">
              {pending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
              {pending ? 'Menyimpan...' : 'Simpan'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
