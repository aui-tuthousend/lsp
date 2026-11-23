import { useEffect, useState, type FormEvent } from 'react'
import { Loader2, Save } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import type { Session } from '@/lib/auth/types'
import { useCreateRole, useUpdateRole } from '../-hooks'
import { validateRoleForm, type RoleFormValues } from '../-schemas'
import type { CreateRoleBody, Role, UpdateRoleBody } from '../-types'

const EMPTY_FORM: RoleFormValues = {
  code: '',
  name: '',
  description: '',
}

function toForm(role: Role): RoleFormValues {
  return {
    code: role.code,
    name: role.name,
    description: role.description ?? '',
  }
}

/** Bentuk payload yang siap dikirim ke backend. */
function toPayload(values: RoleFormValues): CreateRoleBody {
  return {
    code: values.code,
    name: values.name,
    description: values.description || undefined,
  }
}

/** Bandingkan dengan nilai awal; kirim hanya field yang berubah (partial update). */
function diffPatch(values: RoleFormValues, baseline: RoleFormValues): UpdateRoleBody {
  const next = toPayload(values)
  const prev = toPayload(baseline)
  const patch: Record<string, unknown> = {}
  for (const key of Object.keys(next) as (keyof UpdateRoleBody)[]) {
    if (next[key] !== prev[key]) patch[key] = next[key]
  }
  return patch as UpdateRoleBody
}

export function RoleFormDialog({
  open,
  onOpenChange,
  role,
  session,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  role?: Role | null
  session: Session | null
}) {
  const isEdit = Boolean(role)
  const create = useCreateRole(session)
  const update = useUpdateRole(session)
  const pending = create.isPending || update.isPending

  const [values, setValues] = useState<RoleFormValues>(EMPTY_FORM)
  const [baseline, setBaseline] = useState<RoleFormValues>(EMPTY_FORM)
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    if (!open) return
    const form = role ? toForm(role) : EMPTY_FORM
    setValues(form)
    setBaseline(form)
    setErrors({})
  }, [open, role])

  function set<K extends keyof RoleFormValues>(key: K, value: RoleFormValues[K]) {
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
    const result = validateRoleForm(values)
    if (!result.ok) {
      setErrors(result.errors)
      toast.error('Lengkapi data yang belum valid')
      return
    }

    if (role) {
      const patch = diffPatch(result.data, baseline)
      if (Object.keys(patch).length === 0) {
        toast.info('Tidak ada perubahan untuk disimpan')
        return
      }
      update.mutate({ id: role.id, body: patch }, { onSuccess: () => onOpenChange(false) })
      return
    }

    create.mutate(toPayload(result.data), { onSuccess: () => onOpenChange(false) })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEdit ? 'Ubah Role' : 'Tambah Role'}</DialogTitle>
          <DialogDescription>
            {isEdit ? 'Perbarui data master role.' : 'Buat master data role baru.'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={submit} className="space-y-4">
          <Field
            label="Kode"
            htmlFor="code"
            required
            error={errors.code}
            hint="Huruf kecil, angka, underscore — mis. admin, manager_sertifikasi"
          >
            <Input id="code" maxLength={50} value={values.code} onChange={(e) => set('code', e.target.value)} placeholder="mis. validator" />
          </Field>

          <Field label="Nama" htmlFor="name" required error={errors.name}>
            <Input id="name" maxLength={100} value={values.name} onChange={(e) => set('name', e.target.value)} placeholder="mis. Validator" />
          </Field>

          <Field label="Deskripsi" htmlFor="description" error={errors.description} hint="Opsional">
            <Input id="description" maxLength={255} value={values.description} onChange={(e) => set('description', e.target.value)} />
          </Field>

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
