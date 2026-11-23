import { useEffect, useState, type FormEvent } from 'react'
import { Loader2, Save } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import type { Session } from '@/lib/auth/types'
import { useCreateUnit, useUpdateUnit } from '../-hooks'
import { parseSortOrder, validateUnitForm, type UnitFormValues } from '../-schemas'
import type { CompetencyUnit, CreateUnitBody, UpdateUnitBody } from '../-types'

const EMPTY_FORM: UnitFormValues = { code: '', title: '', sort_order: '0' }

function toForm(unit: CompetencyUnit): UnitFormValues {
  return { code: unit.code, title: unit.title, sort_order: String(unit.sort_order) }
}

function toPayload(values: UnitFormValues): CreateUnitBody {
  return { code: values.code, title: values.title, sort_order: parseSortOrder(values.sort_order) }
}

/** Bandingkan dengan nilai awal; kirim hanya field yang berubah (partial update). */
function diffPatch(values: UnitFormValues, baseline: UnitFormValues): UpdateUnitBody {
  const next = toPayload(values)
  const prev = toPayload(baseline)
  const patch: Record<string, unknown> = {}
  for (const key of Object.keys(next) as (keyof UpdateUnitBody)[]) {
    if (next[key] !== prev[key]) patch[key] = next[key]
  }
  return patch as UpdateUnitBody
}

export function UnitFormDialog({
  open,
  onOpenChange,
  schemaId,
  unit,
  session,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  schemaId: string
  unit?: CompetencyUnit | null
  session: Session | null
}) {
  const isEdit = Boolean(unit)
  const create = useCreateUnit(schemaId, session)
  const update = useUpdateUnit(schemaId, session)
  const pending = create.isPending || update.isPending

  const [values, setValues] = useState<UnitFormValues>(EMPTY_FORM)
  const [baseline, setBaseline] = useState<UnitFormValues>(EMPTY_FORM)
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    if (!open) return
    const form = unit ? toForm(unit) : EMPTY_FORM
    setValues(form)
    setBaseline(form)
    setErrors({})
  }, [open, unit])

  function set<K extends keyof UnitFormValues>(key: K, value: UnitFormValues[K]) {
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
    const result = validateUnitForm(values)
    if (!result.ok) {
      setErrors(result.errors)
      toast.error('Lengkapi data yang belum valid')
      return
    }

    if (unit) {
      const patch = diffPatch(result.data, baseline)
      if (Object.keys(patch).length === 0) {
        toast.info('Tidak ada perubahan untuk disimpan')
        return
      }
      update.mutate({ id: unit.id, body: patch }, { onSuccess: () => onOpenChange(false) })
      return
    }

    create.mutate(toPayload(result.data), { onSuccess: () => onOpenChange(false) })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEdit ? 'Ubah Unit Kompetensi' : 'Tambah Unit Kompetensi'}</DialogTitle>
          <DialogDescription>
            {isEdit ? 'Perbarui data unit kompetensi.' : 'Tambahkan unit kompetensi pada skema ini.'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={submit} className="space-y-4">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <Field label="Kode Unit" htmlFor="code" required error={errors.code} hint="mis. N.80PAM00.001.2">
              <Input id="code" maxLength={50} value={values.code} onChange={(e) => set('code', e.target.value)} />
            </Field>

            <Field label="Urutan" htmlFor="sort_order" error={errors.sort_order} hint="Opsional">
              <Input id="sort_order" inputMode="numeric" value={values.sort_order} onChange={(e) => set('sort_order', e.target.value)} />
            </Field>

            <Field label="Judul Unit" htmlFor="title" required error={errors.title} className="sm:col-span-2">
              <Input id="title" maxLength={255} value={values.title} onChange={(e) => set('title', e.target.value)} />
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
