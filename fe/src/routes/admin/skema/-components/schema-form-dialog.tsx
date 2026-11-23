import { useEffect, useState, type FormEvent } from 'react'
import { Loader2, Save } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
import type { Session } from '@/lib/auth/types'
import { useCreateSchema, useUpdateSchema } from '../-hooks'
import { parseSortOrder, validateSchemaForm, type SchemaFormValues } from '../-schemas'
import type { CertificationSchema, CreateSchemaBody, UpdateSchemaBody } from '../-types'

const EMPTY_FORM: SchemaFormValues = {
  code: '',
  name: '',
  schema_code: '',
  sort_order: '0',
  is_active: true,
}

function toForm(schema: CertificationSchema): SchemaFormValues {
  return {
    code: schema.code,
    name: schema.name,
    schema_code: schema.schema_code ?? '',
    sort_order: String(schema.sort_order),
    is_active: schema.is_active,
  }
}

/** Bentuk payload yang siap dikirim ke backend. */
function toPayload(values: SchemaFormValues): CreateSchemaBody {
  return {
    code: values.code,
    name: values.name,
    schema_code: values.schema_code || undefined,
    sort_order: parseSortOrder(values.sort_order),
    is_active: values.is_active,
  }
}

/** Bandingkan dengan nilai awal; kirim hanya field yang berubah (partial update). */
function diffPatch(values: SchemaFormValues, baseline: SchemaFormValues): UpdateSchemaBody {
  const next = toPayload(values)
  const prev = toPayload(baseline)
  const patch: Record<string, unknown> = {}
  for (const key of Object.keys(next) as (keyof UpdateSchemaBody)[]) {
    if (next[key] !== prev[key]) patch[key] = next[key]
  }
  return patch as UpdateSchemaBody
}

export function SchemaFormDialog({
  open,
  onOpenChange,
  schema,
  session,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  schema?: CertificationSchema | null
  session: Session | null
}) {
  const isEdit = Boolean(schema)
  const create = useCreateSchema(session)
  const update = useUpdateSchema(session)
  const pending = create.isPending || update.isPending

  const [values, setValues] = useState<SchemaFormValues>(EMPTY_FORM)
  const [baseline, setBaseline] = useState<SchemaFormValues>(EMPTY_FORM)
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    if (!open) return
    const form = schema ? toForm(schema) : EMPTY_FORM
    setValues(form)
    setBaseline(form)
    setErrors({})
  }, [open, schema])

  function set<K extends keyof SchemaFormValues>(key: K, value: SchemaFormValues[K]) {
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
    const result = validateSchemaForm(values)
    if (!result.ok) {
      setErrors(result.errors)
      toast.error('Lengkapi data yang belum valid')
      return
    }

    if (schema) {
      const patch = diffPatch(result.data, baseline)
      if (Object.keys(patch).length === 0) {
        toast.info('Tidak ada perubahan untuk disimpan')
        return
      }
      update.mutate({ id: schema.id, body: patch }, { onSuccess: () => onOpenChange(false) })
      return
    }

    create.mutate(toPayload(result.data), { onSuccess: () => onOpenChange(false) })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEdit ? 'Ubah Skema Sertifikasi' : 'Tambah Skema Sertifikasi'}</DialogTitle>
          <DialogDescription>
            {isEdit ? 'Perbarui data master skema sertifikasi.' : 'Lengkapi data skema sertifikasi baru.'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={submit} className="space-y-4">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <Field label="Kode" htmlFor="code" required error={errors.code} hint="Slug internal, mis. gada_pratama">
              <Input id="code" maxLength={50} value={values.code} onChange={(e) => set('code', e.target.value)} placeholder="gada_pratama" />
            </Field>

            <Field label="Nama" htmlFor="name" required error={errors.name}>
              <Input id="name" maxLength={255} value={values.name} onChange={(e) => set('name', e.target.value)} placeholder="mis. Gada Pratama" />
            </Field>

            <Field label="Nomor Skema (BNSP)" htmlFor="schema_code" error={errors.schema_code} hint="Opsional">
              <Input
                id="schema_code"
                maxLength={100}
                value={values.schema_code}
                onChange={(e) => set('schema_code', e.target.value)}
                placeholder="mis. SS/OKU/04/GU/X/2019"
              />
            </Field>

            <Field label="Urutan" htmlFor="sort_order" error={errors.sort_order} hint="Opsional · kecil tampil lebih dulu">
              <Input id="sort_order" inputMode="numeric" value={values.sort_order} onChange={(e) => set('sort_order', e.target.value)} />
            </Field>

            <Field label="Status" htmlFor="is_active" error={errors.is_active}>
              <Select id="is_active" value={values.is_active ? 'true' : 'false'} onChange={(e) => set('is_active', e.target.value === 'true')}>
                <option value="true">Aktif</option>
                <option value="false">Tidak aktif</option>
              </Select>
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
