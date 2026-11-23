// Schema & aturan validasi form master data role.

import { z } from 'zod'

const CODE_RE = /^[a-z0-9_]+$/

export const roleFormSchema = z.object({
  code: z
    .string()
    .trim()
    .min(1, 'Kode wajib diisi')
    .max(50, 'Maksimal 50 karakter')
    .regex(CODE_RE, 'Kode hanya huruf kecil, angka, dan underscore'),
  name: z.string().trim().min(1, 'Nama wajib diisi').max(100, 'Maksimal 100 karakter'),
  description: z.string().trim().max(255, 'Maksimal 255 karakter'),
})

export type RoleFormValues = z.infer<typeof roleFormSchema>

export type RoleFormValidation =
  | { ok: true; data: RoleFormValues }
  | { ok: false; errors: Record<string, string> }

/** Validasi form role; mengembalikan data yang sudah dinormalisasi bila valid. */
export function validateRoleForm(values: unknown): RoleFormValidation {
  const result = roleFormSchema.safeParse(values)
  if (result.success) return { ok: true, data: result.data }

  const errors: Record<string, string> = {}
  for (const issue of result.error.issues) {
    const key = String(issue.path[0] ?? '')
    if (key && !(key in errors)) errors[key] = issue.message
  }
  return { ok: false, errors }
}
