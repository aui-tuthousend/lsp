import { useEffect, useState, type FormEvent } from 'react'
import { Loader2, Save } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import type { Session } from '@/lib/auth/types'
import { useCreateDocumentType, useUpdateDocumentType } from '../-hooks'
import { validateDocumentTypeForm, type DocumentTypeFormValues } from '../-schemas'
import type { CreateDocumentTypeBody, DocumentType, UpdateDocumentTypeBody } from '../-types'

const EMPTY_FORM: DocumentTypeFormValues = {
  code: '',
  label: '',
}

function toForm(docType: DocumentType): DocumentTypeFormValues {
  return {
    code: docType.code,
    label: docType.label,
  }
}

/** Bentuk payload yang siap dikirim ke backend. */
function toPayload(values: DocumentTypeFormValues): CreateDocumentTypeBody {
  return {
    code: values.code,
    label: values.label,
  }
}

/** Bandingkan dengan nilai awal; kirim hanya field yang berubah (partial update). */
function diffPatch(values: DocumentTypeFormValues, baseline: DocumentTypeFormValues): UpdateDocumentTypeBody {
  const next = toPayload(values)
  const prev = toPayload(baseline)
  const patch: Record<string, unknown> = {}
  for (const key of Object.keys(next) as (keyof UpdateDocumentTypeBody)[]) {
    if (next[key] !== prev[key]) patch[key] = next[key]
  }
  return patch as UpdateDocumentTypeBody
}

export function DocumentTypeFormDialog({
  open,
  onOpenChange,
  docType,
  session,
}: {
  open: boolean
  onOpenChange: (open: boolean) => void
  docType?: DocumentType | null
  session: Session | null
}) {
  const isEdit = Boolean(docType)
  const create = useCreateDocumentType(session)
  const update = useUpdateDocumentType(session)
  const pending = create.isPending || update.isPending

  const [values, setValues] = useState<DocumentTypeFormValues>(EMPTY_FORM)
  const [baseline, setBaseline] = useState<DocumentTypeFormValues>(EMPTY_FORM)
  const [errors, setErrors] = useState<Record<string, string>>({})

  useEffect(() => {
    if (!open) return
    const form = docType ? toForm(docType) : EMPTY_FORM
    setValues(form)
    setBaseline(form)
    setErrors({})
  }, [open, docType])

  function set<K extends keyof DocumentTypeFormValues>(key: K, value: DocumentTypeFormValues[K]) {
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
    const result = validateDocumentTypeForm(values)
    if (!result.ok) {
      setErrors(result.errors)
      toast.error('Lengkapi data yang belum valid')
      return
    }

    if (docType) {
      const patch = diffPatch(result.data, baseline)
      if (Object.keys(patch).length === 0) {
        toast.info('Tidak ada perubahan untuk disimpan')
        return
      }
      update.mutate({ id: docType.id, body: patch }, { onSuccess: () => onOpenChange(false) })
      return
    }

    create.mutate(toPayload(result.data), { onSuccess: () => onOpenChange(false) })
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>{isEdit ? 'Ubah Jenis Dokumen' : 'Tambah Jenis Dokumen'}</DialogTitle>
          <DialogDescription>
            {isEdit ? 'Perbarui data master jenis dokumen.' : 'Lengkapi data jenis dokumen baru.'}
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={submit} className="space-y-4">
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <Field label="Kode" htmlFor="code" required error={errors.code} hint="mis. ktp, ijazah">
              <Input id="code" maxLength={50} value={values.code} onChange={(e) => set('code', e.target.value)} placeholder="ktp" />
            </Field>

            <Field label="Label" htmlFor="label" required error={errors.label}>
              <Input id="label" maxLength={255} value={values.label} onChange={(e) => set('label', e.target.value)} placeholder="mis. KTP" />
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
