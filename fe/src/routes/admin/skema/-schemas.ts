// Schema & aturan validasi form master skema sertifikasi.
//
// Field mengikuti body POST/PATCH modul certificationschema. `sort_order` di
// form disimpan sebagai string lalu dikonversi ke number saat membentuk payload.

import { z } from 'zod'

const SORT_ORDER_RE = /^\d+$/

export const schemaFormSchema = z.object({
  code: z.string().trim().min(1, 'Kode wajib diisi').max(50, 'Maksimal 50 karakter'),
  name: z.string().trim().min(1, 'Nama wajib diisi').max(255, 'Maksimal 255 karakter'),
  schema_code: z.string().trim().max(100, 'Maksimal 100 karakter'),
  sort_order: z.string().trim().refine((value) => value === '' || SORT_ORDER_RE.test(value), 'Harus bilangan bulat ≥ 0'),
  is_active: z.boolean(),
})

export type SchemaFormValues = z.infer<typeof schemaFormSchema>

export const unitFormSchema = z.object({
  code: z.string().trim().min(1, 'Kode unit wajib diisi').max(50, 'Maksimal 50 karakter'),
  title: z.string().trim().min(1, 'Judul unit wajib diisi').max(255, 'Maksimal 255 karakter'),
  sort_order: z.string().trim().refine((value) => value === '' || SORT_ORDER_RE.test(value), 'Harus bilangan bulat ≥ 0'),
})

export type UnitFormValues = z.infer<typeof unitFormSchema>

export const elementFormSchema = z.object({
  title: z.string().trim().max(255, 'Maksimal 255 karakter'),
  sort_order: z.string().trim().refine((value) => value === '' || SORT_ORDER_RE.test(value), 'Harus bilangan bulat ≥ 0'),
})

export type ElementFormValues = z.infer<typeof elementFormSchema>

export const criteriaFormSchema = z.object({
  content: z.string().trim().min(1, 'Isi KUK wajib diisi'),
  sort_order: z.string().trim().refine((value) => value === '' || SORT_ORDER_RE.test(value), 'Harus bilangan bulat ≥ 0'),
})

export type CriteriaFormValues = z.infer<typeof criteriaFormSchema>

export type FormValidation<T> = { ok: true; data: T } | { ok: false; errors: Record<string, string> }

/** Ubah error zod menjadi map field -> pesan (pesan pertama per field). */
function toErrorMap(issues: { path: PropertyKey[]; message: string }[]): Record<string, string> {
  const errors: Record<string, string> = {}
  for (const issue of issues) {
    const key = String(issue.path[0] ?? '')
    if (key && !(key in errors)) errors[key] = issue.message
  }
  return errors
}

export function validateSchemaForm(values: unknown): FormValidation<SchemaFormValues> {
  const result = schemaFormSchema.safeParse(values)
  return result.success ? { ok: true, data: result.data } : { ok: false, errors: toErrorMap(result.error.issues) }
}

export function validateUnitForm(values: unknown): FormValidation<UnitFormValues> {
  const result = unitFormSchema.safeParse(values)
  return result.success ? { ok: true, data: result.data } : { ok: false, errors: toErrorMap(result.error.issues) }
}

export function validateElementForm(values: unknown): FormValidation<ElementFormValues> {
  const result = elementFormSchema.safeParse(values)
  return result.success ? { ok: true, data: result.data } : { ok: false, errors: toErrorMap(result.error.issues) }
}

export function validateCriteriaForm(values: unknown): FormValidation<CriteriaFormValues> {
  const result = criteriaFormSchema.safeParse(values)
  return result.success ? { ok: true, data: result.data } : { ok: false, errors: toErrorMap(result.error.issues) }
}

/** Konversi sort_order string ke number (default 0). */
export function parseSortOrder(value: string): number {
  const trimmed = value.trim()
  return trimmed === '' ? 0 : Number.parseInt(trimmed, 10)
}
