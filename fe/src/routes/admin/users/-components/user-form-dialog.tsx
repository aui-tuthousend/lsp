import { useEffect, useState, type FormEvent } from 'react'
import { Loader2, Save } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
import type { Session } from '@/lib/auth/types'
import { useCreateUser, useUpdateUser } from '../-hooks'
import { ROLE_OPTIONS, validateUserForm, type UserFormMode, type UserFormValues } from '../-schemas'
import type { CreateUserBody, UpdateUserBody, User } from '../-types'

const EMPTY_FORM: UserFormValues = {
  name: '',
  email: '',
  password: '',
  role: ROLE_OPTIONS[0].value,
  phone: '',
  nomor_registrasi: '',
}

function toForm(user: User): UserFormValues {
  return {
    name: user.name,
    email: user.email,
    password: '',
    role: user.role,
    phone: user.phone ?? '',
    nomor_registrasi: user.nomor_registrasi ?? '',
  }
}

/** Bentuk payload create (POST /users). */
function toPayload(values: UserFormValues): CreateUserBody {
  return {
    name: values.name,
    email: values.email,
    password: values.password,
    role: values.role,
    phone: values.phone || undefined,
    nomor_registrasi: values.nomor_registrasi || undefined,
  }
}

/** Bandingkan dengan nilai awal; kirim hanya field yang berubah (partial update). Role tidak diubah. */
function diffPatch(values: UserFormValues, baseline: UserFormValues): UpdateUserBody {
  const patch: UpdateUserBody = {}
  if (values.name !== baseline.name) patch.name = values.name
  if (values.phone !== baseline.phone) patch.phone = values.phone
  if (values.nomor_registrasi !== baseline.nomor_registrasi) patch.nomor_registrasi = values.nomor_registrasi
  if (values.password) patch.password = values.password
  return patch
}

export function UserFormDialog({
  open,
  onOpenChange,
  user,
  session,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  user?: User | null
  session: Session | null
}) {
  const isEdit = Boolean(user)
  const mode: UserFormMode = isEdit ? 'edit' : 'create'
  const create = useCreateUser(session)
  const update = useUpdateUser(session)
  const pending = create.isPending || update.isPending

  const [values, setValues] = useState<UserFormValues>(EMPTY_FORM)
  const [baseline, setBaseline] = useState<UserFormValues>(EMPTY_FORM)
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    if (!open) return
    const form = user ? toForm(user) : EMPTY_FORM
    setValues(form)
    setBaseline(form)
    setErrors({})
  }, [open, user])

  function set<K extends keyof UserFormValues>(key: K, value: UserFormValues[K]) {
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
    const result = validateUserForm(values, mode)
    if (!result.ok) {
      setErrors(result.errors)
      toast.error('Lengkapi data yang belum valid')
      return
    }

    if (user) {
      const patch = diffPatch(result.data, baseline)
      if (Object.keys(patch).length === 0) {
        toast.info('Tidak ada perubahan untuk disimpan')
        return
      }
      update.mutate({ id: user.id, body: patch }, { onSuccess: () => onOpenChange(false) })
      return
    }

    create.mutate(toPayload(result.data), { onSuccess: () => onOpenChange(false) })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEdit ? 'Ubah User' : 'Tambah User'}</DialogTitle>
          <DialogDescription>
            {isEdit ? 'Perbarui data akun user.' : 'Buat akun user dengan role validator atau manager sertifikasi.'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={submit} className="space-y-4">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <Field label="Nama" htmlFor="name" required error={errors.name} className="sm:col-span-2">
              <Input id="name" maxLength={255} value={values.name} onChange={(e) => set('name', e.target.value)} placeholder="mis. Andi Pratama" />
            </Field>

            <Field label="Email" htmlFor="email" required error={errors.email}>
              <Input
                id="email"
                type="email"
                maxLength={255}
                value={values.email}
                onChange={(e) => set('email', e.target.value)}
                placeholder="user@contoh.id"
                disabled={isEdit}
              />
            </Field>

            <Field label="Role" htmlFor="role" required error={errors.role} hint={isEdit ? 'Role tidak dapat diubah' : undefined}>
              <Select id="role" value={values.role} onChange={(e) => set('role', e.target.value)} disabled={isEdit}>
                {ROLE_OPTIONS.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </Select>
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
