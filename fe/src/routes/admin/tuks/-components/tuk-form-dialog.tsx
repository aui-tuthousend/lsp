import { useEffect, useState, type FormEvent } from 'react'
import { Loader2, Save } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import type { Session } from '@/lib/auth/types'
import { useCreateTuk, useUpdateTuk } from '../-hooks'
import { validateTUKForm, type TUKFormValues } from '../-schemas'
import type { CreateTUKBody, TUK, UpdateTUKBody } from '../-types'

const EMPTY_FORM: TUKFormValues = {
  name: '',
  postal_code: '',
  phone: '',
  email: '',
  province: '',
  city: '',
  address: '',
  maps_url: '',
  logo: '',
  letterhead: '',
}

function toForm(tuk: TUK): TUKFormValues {
  return {
    name: tuk.name,
    postal_code: tuk.postal_code,
    phone: tuk.phone,
    email: tuk.email,
    province: tuk.province,
    city: tuk.city,
    address: tuk.address,
    maps_url: tuk.maps_url ?? '',
    logo: tuk.logo ?? '',
    letterhead: tuk.letterhead ?? '',
  }
}

/** Bentuk payload yang siap dikirim ke backend. */
function toPayload(values: TUKFormValues): CreateTUKBody {
  return {
    name: values.name,
    postal_code: values.postal_code,
    phone: values.phone,
    email: values.email,
    province: values.province,
    city: values.city,
    address: values.address,
    maps_url: values.maps_url || undefined,
    logo: values.logo || undefined,
    letterhead: values.letterhead || undefined,
  }
}

/** Bandingkan dengan nilai awal; kirim hanya field yang berubah (partial update). */
function diffPatch(values: TUKFormValues, baseline: TUKFormValues): UpdateTUKBody {
  const next = toPayload(values)
  const prev = toPayload(baseline)
  const patch: Record<string, unknown> = {}
  for (const key of Object.keys(next) as (keyof UpdateTUKBody)[]) {
    if (next[key] !== prev[key]) patch[key] = next[key]
  }
  return patch as UpdateTUKBody
}

export function TUKFormDialog({
  open,
  onOpenChange,
  tuk,
  session,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  tuk?: TUK | null
  session: Session | null
}) {
  const isEdit = Boolean(tuk)
  const create = useCreateTuk(session)
  const update = useUpdateTuk(session)
  const pending = create.isPending || update.isPending

  const [values, setValues] = useState<TUKFormValues>(EMPTY_FORM)
  const [baseline, setBaseline] = useState<TUKFormValues>(EMPTY_FORM)
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    if (!open) return
    const form = tuk ? toForm(tuk) : EMPTY_FORM
    setValues(form)
    setBaseline(form)
    setErrors({})
  }, [open, tuk])

  function set<K extends keyof TUKFormValues>(key: K, value: TUKFormValues[K]) {
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
    const result = validateTUKForm(values)
    if (!result.ok) {
      setErrors(result.errors)
      toast.error('Lengkapi data yang belum valid')
      return
    }

    if (tuk) {
      const patch = diffPatch(result.data, baseline)
      if (Object.keys(patch).length === 0) {
        toast.info('Tidak ada perubahan untuk disimpan')
        return
      }
      update.mutate({ id: tuk.id, body: patch }, { onSuccess: () => onOpenChange(false) })
      return
    }

    create.mutate(toPayload(result.data), { onSuccess: () => onOpenChange(false) })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEdit ? 'Ubah TUK' : 'Tambah TUK'}</DialogTitle>
          <DialogDescription>
            {isEdit ? 'Perbarui data Tempat Uji Kompetensi.' : 'Lengkapi data Tempat Uji Kompetensi baru.'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={submit} className="space-y-4">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <Field label="Nama TUK" htmlFor="name" required error={errors.name} className="sm:col-span-2">
              <Input id="name" maxLength={100} value={values.name} onChange={(e) => set('name', e.target.value)} placeholder="mis. TUK Jakarta Pusat" />
            </Field>

            <Field label="Email" htmlFor="email" required error={errors.email}>
              <Input id="email" type="email" maxLength={100} value={values.email} onChange={(e) => set('email', e.target.value)} placeholder="tuk@contoh.id" />
            </Field>

            <Field label="Telepon" htmlFor="phone" required error={errors.phone}>
              <Input id="phone" inputMode="tel" maxLength={13} value={values.phone} onChange={(e) => set('phone', e.target.value)} placeholder="021xxxxxxx" />
            </Field>

            <Field label="Provinsi" htmlFor="province" required error={errors.province}>
              <Input id="province" maxLength={50} value={values.province} onChange={(e) => set('province', e.target.value)} />
            </Field>

            <Field label="Kota / Kabupaten" htmlFor="city" required error={errors.city}>
              <Input id="city" maxLength={50} value={values.city} onChange={(e) => set('city', e.target.value)} />
            </Field>

            <Field label="Kode Pos" htmlFor="postal_code" required error={errors.postal_code}>
              <Input id="postal_code" inputMode="numeric" maxLength={5} value={values.postal_code} onChange={(e) => set('postal_code', e.target.value)} />
            </Field>

            <Field label="Alamat" htmlFor="address" required error={errors.address} className="sm:col-span-2">
              <Input id="address" value={values.address} onChange={(e) => set('address', e.target.value)} placeholder="Jalan, nomor, kelurahan" />
            </Field>

            <Field label="Link Maps" htmlFor="maps_url" error={errors.maps_url} hint="Opsional · URL Google Maps" className="sm:col-span-2">
              <Input id="maps_url" value={values.maps_url} onChange={(e) => set('maps_url', e.target.value)} placeholder="https://maps.google.com/..." />
            </Field>

            <Field label="URL Logo" htmlFor="logo" error={errors.logo} hint="Opsional">
              <Input id="logo" value={values.logo} onChange={(e) => set('logo', e.target.value)} placeholder="https://..." />
            </Field>

            <Field label="URL Kop Surat" htmlFor="letterhead" error={errors.letterhead} hint="Opsional">
              <Input id="letterhead" value={values.letterhead} onChange={(e) => set('letterhead', e.target.value)} placeholder="https://..." />
            </Field>
          </div>

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
