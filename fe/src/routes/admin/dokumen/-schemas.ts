// Schema & aturan validasi form master jenis dokumen.
//
// Field mengikuti body POST/PATCH `/document-types` (modul documentrequirements).

import { z } from 'zod'

export const documentTypeFormSchema = z.object({
  code: z.string().trim().min(1, 'Kode wajib diisi').max(50, 'Maksimal 50 karakter'),
  label: z.string().trim().min(1, 'Label wajib diisi').max(255, 'Maksimal 255 karakter'),
})

export type DocumentTypeFormValues = z.infer<typeof documentTypeFormSchema>

export type DocumentTypeFormValidation =
  | { ok: true; data: DocumentTypeFormValues }
  | { ok: false; errors: Record<string, string> }

/** Validasi form jenis dokumen; mengembalikan data ternormalisasi bila valid. */
export function validateDocumentTypeForm(values: unknown): DocumentTypeFormValidation {
  const result = documentTypeFormSchema.safeParse(values)
  if (result.success) return { ok: true, data: result.data }

  const errors: Record<string, string> = {}
  for (const issue of result.error.issues) {
    const key = String(issue.path[0] ?? '')
    if (key && !(key in errors)) errors[key] = issue.message
  }
  return { ok: false, errors }
}
