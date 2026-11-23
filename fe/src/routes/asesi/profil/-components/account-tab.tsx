import { useEffect, useState, type FormEvent } from 'react'
import { Loader2, Save } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import type { Session } from '@/lib/auth/types'
import { useMe, useUpdateMe } from '../-hooks'
import { EMPTY_ACCOUNT_FORM, diffAccount, toAccountForm, validateAccount, type AccountFormValues } from '../-store'

/** Tab "Data Akun" — atribut dari tabel `users` (nama & telepon bisa diubah). */
export function AccountTab({ session }: { session: Session | null }) {
  const me = useMe(session)
  const update = useUpdateMe(session)

  const [values, setValues] = useState<AccountFormValues>(EMPTY_ACCOUNT_FORM)
  const [baseline, setBaseline] = useState<AccountFormValues>(EMPTY_ACCOUNT_FORM)

  useEffect(() => {
    if (!me.data) return
    const form = toAccountForm(me.data)
    setValues(form)
    setBaseline(form)
  }, [me.data])

  const patch = diffAccount(values, baseline)
  const isDirty = Object.keys(patch).length > 0
  const account = me.data

  function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const message = validateAccount(values)
    if (message) {
      toast.error(message)
      return
    }
    if (!isDirty) {
      toast.info('Tidak ada perubahan untuk disimpan')
      return
    }
    update.mutate(patch)
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Data Akun</CardTitle>
        <p className="text-xs text-slate-500">Berasal dari tabel users.</p>
      </CardHeader>
      <CardContent>
        {me.isLoading ? (
          <div className="flex items-center gap-2 py-6 text-sm text-slate-500">
            <Loader2 className="h-4 w-4 animate-spin" />
            Memuat data akun...
          </div>
        ) : (
          <form onSubmit={submit} className="space-y-4">
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
              <Field label="Nama Akun" htmlFor="account_name" required>
                <Input id="account_name" maxLength={255} value={values.name} onChange={(e) => setValues((prev) => ({ ...prev, name: e.target.value }))} />
              </Field>

              <Field label="Telepon" htmlFor="account_phone">
                <Input
                  id="account_phone"
                  inputMode="tel"
                  maxLength={20}
                  value={values.phone}
                  onChange={(e) => setValues((prev) => ({ ...prev, phone: e.target.value }))}
                />
              </Field>

              <InfoItem label="Email" value={account?.email} />
              <InfoItem label="Role" value={account?.role} />
              <InfoItem label="Nomor Registrasi" value={account?.nomor_registrasi} />
              <InfoItem label="Verifikasi Email" value={account?.email_verified_at ? 'Terverifikasi' : 'Belum terverifikasi'} />
            </div>

            <div className="flex justify-end">
              <Button type="submit" disabled={update.isPending || !isDirty} className="gap-2">
                {update.isPending ? <Loader2 className="h-4 w-4 animate-spin" /> : <Save className="h-4 w-4" />}
                {update.isPending ? 'Menyimpan...' : 'Simpan'}
              </Button>
            </div>
          </form>
        )}
      </CardContent>
    </Card>
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
