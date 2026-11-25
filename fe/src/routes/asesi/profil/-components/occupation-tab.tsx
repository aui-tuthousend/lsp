import { useEffect, useState, type FormEvent } from 'react'
import { Loader2, Save } from 'lucide-react'
import { toast } from 'sonner'
import { AddressFields } from '@/components/ui/address-fields'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'
import type { Session } from '@/lib/auth/types'
import { useOccupationByUser, useUpsertOccupation } from '../-hooks'
import {
  EMPTY_OCCUPATION_FORM,
  controlClass,
  occupationDirty,
  toOccupationBody,
  toOccupationForm,
  validateOccupation,
  type OccupationFormValues,
} from '../-store'

/** Tab "Pekerjaan" — atribut dari tabel `asesi_occupations` (upsert). */
export function OccupationTab({ session }: { session: Session | null }) {
  const userId = session?.id ?? ''
  const occupationQuery = useOccupationByUser(userId, session)
  const upsert = useUpsertOccupation(session)

  const [values, setValues] = useState<OccupationFormValues>(EMPTY_OCCUPATION_FORM)
  const [baseline, setBaseline] = useState<OccupationFormValues>(EMPTY_OCCUPATION_FORM)

  useEffect(() => {
    const form = occupationQuery.data ? toOccupationForm(occupationQuery.data) : EMPTY_OCCUPATION_FORM
    setValues(form)
    setBaseline(form)
  }, [occupationQuery.data])

  const isDirty = occupationDirty(values, baseline)

  function set<K extends keyof OccupationFormValues>(key: K, value: OccupationFormValues[K]) {
    setValues((prev) => ({ ...prev, [key]: value }))
  }

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const message = validateOccupation(values)
    if (message) {
      toast.error(message)
      return
    }
    if (!isDirty) {
      toast.info('Tidak ada perubahan untuk disimpan')
      return
    }
    if (!userId) {
      toast.error('Sesi tidak valid')
      return
    }
    upsert.mutate({ userId, body: toOccupationBody(values) })
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Data Pekerjaan</CardTitle>
        <p className="text-xs text-slate-500">Berasal dari tabel asesi_occupations.</p>
      </CardHeader>
      <CardContent>
        {occupationQuery.isLoading ? (
          <div className="flex items-center gap-2 py-6 text-sm text-slate-500">
            <Loader2 className="h-4 w-4 animate-spin" />
            Memuat data pekerjaan...
          </div>
        ) : (
          <form onSubmit={submit} className="space-y-4">
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <Field label="Nama Perusahaan" htmlFor="company_name" required>
                <Input id="company_name" maxLength={100} value={values.company_name} onChange={(e) => set('company_name', e.target.value)} />
              </Field>

              <Field label="Jabatan" htmlFor="position" required>
                <Input id="position" maxLength={100} value={values.position} onChange={(e) => set('position', e.target.value)} />
              </Field>

              <Field label="Lama Bekerja (bulan)" htmlFor="tenure_months">
                <Input
                  id="tenure_months"
                  inputMode="numeric"
                  value={values.tenure_months}
                  onChange={(e) => set('tenure_months', e.target.value.replace(/\D/g, ''))}
                />
              </Field>

              <Field label="Jumlah Bawahan" htmlFor="subordinates_count">
                <Input
                  id="subordinates_count"
                  inputMode="numeric"
                  value={values.subordinates_count}
                  onChange={(e) => set('subordinates_count', e.target.value.replace(/\D/g, ''))}
                />
              </Field>

              <Field label="Alamat" htmlFor="address" required className="sm:col-span-2">
                <textarea
                  id="address"
                  rows={3}
                  value={values.address}
                  onChange={(e) => set('address', e.target.value)}
                  className={cn(controlClass, 'h-auto resize-none py-2')}
                />
              </Field>

              <AddressFields
                values={values}
                withVillage={false}
                lockRegion
                onChange={(key, value) => {
                  if (key !== 'village') set(key, value)
                }}
                className="sm:col-span-2"
              />

              <Field label="Email Kantor" htmlFor="office_email">
                <Input id="office_email" type="email" maxLength={255} value={values.office_email} onChange={(e) => set('office_email', e.target.value)} />
              </Field>

              <Field label="Telepon Kantor" htmlFor="office_phone">
                <Input id="office_phone" inputMode="tel" maxLength={20} value={values.office_phone} onChange={(e) => set('office_phone', e.target.value)} />
              </Field>
            </div>

            <div className="flex justify-end">
              <Button type="submit" disabled={upsert.isPending || !isDirty} className="gap-2">
                {upsert.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
                {upsert.isPending ? 'Menyimpan...' : 'Simpan'}
              </Button>
            </div>
          </form>
        )}
      </CardContent>
    </Card>
  )
}
