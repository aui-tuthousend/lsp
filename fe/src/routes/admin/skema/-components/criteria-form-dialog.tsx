import { useEffect, useState, type FormEvent } from 'react'
import { Loader2, Save } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import type { Session } from '@/lib/auth/types'
import { useCreateCriteria, useUpdateCriteria } from '../-hooks'
import { parseSortOrder, validateCriteriaForm, type CriteriaFormValues } from '../-schemas'
import type { CompetencyUnitCriteria, CreateCriteriaBody, UpdateCriteriaBody } from '../-types'

const EMPTY_FORM: CriteriaFormValues = { content: '', sort_order: '0' }

function toForm(criteria: CompetencyUnitCriteria): CriteriaFormValues {
  return { content: criteria.content, sort_order: String(criteria.sort_order) }
}

function toPayload(values: CriteriaFormValues): CreateCriteriaBody {
  return { content: values.content, sort_order: parseSortOrder(values.sort_order) }
}

/** Bandingkan dengan nilai awal; kirim hanya field yang berubah (partial update). */
function diffPatch(values: CriteriaFormValues, baseline: CriteriaFormValues): UpdateCriteriaBody {
  const next = toPayload(values)
  const prev = toPayload(baseline)
  const patch: Record<string, unknown> = {}
  for (const key of Object.keys(next) as (keyof UpdateCriteriaBody)[]) {
    if (next[key] !== prev[key]) patch[key] = next[key]
  }
  return patch as UpdateCriteriaBody
}

export function CriteriaFormDialog({
  open,
  onOpenChange,
  schemaId,
  elementId,
  criteria,
  session,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  schemaId: string
  elementId: string
  criteria?: CompetencyUnitCriteria | null
  session: Session | null
}) {
  const isEdit = Boolean(criteria)
  const create = useCreateCriteria(schemaId, session)
  const update = useUpdateCriteria(schemaId, session)
  const pending = create.isPending || update.isPending

  const [values, setValues] = useState<CriteriaFormValues>(EMPTY_FORM)
  const [baseline, setBaseline] = useState<CriteriaFormValues>(EMPTY_FORM)
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    if (!open) return
    const form = criteria ? toForm(criteria) : EMPTY_FORM
    setValues(form)
    setBaseline(form)
    setErrors({})
  }, [open, criteria])

  function set<K extends keyof CriteriaFormValues>(key: K, value: CriteriaFormValues[K]) {
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
    const result = validateCriteriaForm(values)
    if (!result.ok) {
      setErrors(result.errors)
      toast.error('Lengkapi data yang belum valid')
      return
    }

    if (criteria) {
      const patch = diffPatch(result.data, baseline)
      if (Object.keys(patch).length === 0) {
        toast.info('Tidak ada perubahan untuk disimpan')
        return
      }
      update.mutate({ id: criteria.id, body: patch }, { onSuccess: () => onOpenChange(false) })
      return
    }

    create.mutate({ elementId, body: toPayload(result.data) }, { onSuccess: () => onOpenChange(false) })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEdit ? 'Ubah KUK' : 'Tambah KUK'}</DialogTitle>
          <DialogDescription>
            {isEdit ? 'Perbarui kriteria unjuk kerja (KUK).' : 'Tambahkan kriteria unjuk kerja (KUK) pada elemen ini.'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={submit} className="space-y-4">
          <Field label="Isi KUK" htmlFor="content" required error={errors.content}>
            <textarea
              id="content"
              rows={3}
              value={values.content}
              onChange={(e) => set('content', e.target.value)}
              className="flex w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-slate-300"
            />
          </Field>

          <Field label="Urutan" htmlFor="sort_order" error={errors.sort_order} hint="Opsional">
            <Input id="sort_order" inputMode="numeric" value={values.sort_order} onChange={(e) => set('sort_order', e.target.value)} />
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
