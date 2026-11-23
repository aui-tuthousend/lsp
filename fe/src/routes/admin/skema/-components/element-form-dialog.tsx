import { useEffect, useState, type FormEvent } from 'react'
import { Loader2, Save } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import type { Session } from '@/lib/auth/types'
import { useCreateElement, useUpdateElement } from '../-hooks'
import { parseSortOrder, validateElementForm, type ElementFormValues } from '../-schemas'
import type { CompetencyUnitElement, CreateElementBody, UpdateElementBody } from '../-types'

const EMPTY_FORM: ElementFormValues = { title: '', sort_order: '0' }

function toForm(element: CompetencyUnitElement): ElementFormValues {
  return { title: element.title ?? '', sort_order: String(element.sort_order) }
}

function toPayload(values: ElementFormValues): CreateElementBody {
  return { title: values.title || undefined, sort_order: parseSortOrder(values.sort_order) }
}

/** Bandingkan dengan nilai awal; kirim hanya field yang berubah (partial update). */
function diffPatch(values: ElementFormValues, baseline: ElementFormValues): UpdateElementBody {
  const next = toPayload(values)
  const prev = toPayload(baseline)
  const patch: Record<string, unknown> = {}
  for (const key of Object.keys(next) as (keyof UpdateElementBody)[]) {
    if (next[key] !== prev[key] && next[key] !== undefined) patch[key] = next[key]
  }
  return patch as UpdateElementBody
}

export function ElementFormDialog({
  open,
  onOpenChange,
  schemaId,
  unitId,
  element,
  session,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  schemaId: string
  unitId: string
  element?: CompetencyUnitElement | null
  session: Session | null
}) {
  const isEdit = Boolean(element)
  const create = useCreateElement(schemaId, session)
  const update = useUpdateElement(schemaId, session)
  const pending = create.isPending || update.isPending

  const [values, setValues] = useState<ElementFormValues>(EMPTY_FORM)
  const [baseline, setBaseline] = useState<ElementFormValues>(EMPTY_FORM)
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    if (!open) return
    const form = element ? toForm(element) : EMPTY_FORM
    setValues(form)
    setBaseline(form)
    setErrors({})
  }, [open, element])

  function set<K extends keyof ElementFormValues>(key: K, value: ElementFormValues[K]) {
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
    const result = validateElementForm(values)
    if (!result.ok) {
      setErrors(result.errors)
      toast.error('Lengkapi data yang belum valid')
      return
    }

    if (element) {
      const patch = diffPatch(result.data, baseline)
      if (Object.keys(patch).length === 0) {
        toast.info('Tidak ada perubahan untuk disimpan')
        return
      }
      update.mutate({ id: element.id, body: patch }, { onSuccess: () => onOpenChange(false) })
      return
    }

    create.mutate({ unitId, body: toPayload(result.data) }, { onSuccess: () => onOpenChange(false) })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEdit ? 'Ubah Elemen' : 'Tambah Elemen'}</DialogTitle>
          <DialogDescription>
            {isEdit ? 'Perbarui elemen unit kompetensi.' : 'Tambahkan elemen pada unit kompetensi ini.'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={submit} className="space-y-4">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <Field label="Judul Elemen" htmlFor="title" error={errors.title} hint="Opsional" className="sm:col-span-2">
              <Input id="title" maxLength={255} value={values.title} onChange={(e) => set('title', e.target.value)} />
            </Field>

            <Field label="Urutan" htmlFor="sort_order" error={errors.sort_order} hint="Opsional">
              <Input id="sort_order" inputMode="numeric" value={values.sort_order} onChange={(e) => set('sort_order', e.target.value)} />
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
