// Schema & aturan validasi form master user.
//
// Role yang tersedia dibatasi ke validator, manager_sertifikasi, & admin_lsp.
// Email & password wajib saat create; saat edit password opsional (reset).

import { z } from 'zod'

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const PHONE_RE = /^[0-9+\-\s()]{5,20}$/

/** Role yang boleh dikelola di halaman ini. */
export const MANAGED_ROLES = ['validator', 'manager_sertifikasi', 'admin_lsp'] as const
export type ManagedRole = (typeof MANAGED_ROLES)[number]

export const ROLE_OPTIONS: { value: ManagedRole; label: string }[] = [
  { value: 'validator', label: 'Validator' },
  { value: 'manager_sertifikasi', label: 'Manager Sertifikasi' },
  { value: 'admin_lsp', label: 'Admin LSP' },
]

export function roleLabel(role: string): string {
  return ROLE_OPTIONS.find((option) => option.value === role)?.label ?? role
}

const sharedFields = {
  name: z.string().trim().min(2, 'Nama minimal 2 karakter').max(255, 'Maksimal 255 karakter'),
  role: z
    .string()
    .refine((value): value is ManagedRole => (MANAGED_ROLES as readonly string[]).includes(value), 'Pilih role yang valid'),
  phone: z.string().trim().refine((value) => value === '' || PHONE_RE.test(value), 'Nomor telepon tidak valid (5-20 karakter)'),
  nomor_registrasi: z.string().trim().max(50, 'Maksimal 50 karakter'),
}

export const userCreateSchema = z.object({
  ...sharedFields,
  email: z.string().trim().max(255, 'Maksimal 255 karakter').regex(EMAIL_RE, 'Format email tidak valid'),
  password: z.string().min(8, 'Minimal 8 karakter').max(72, 'Maksimal 72 karakter'),
})

export const userEditSchema = z.object({
  ...sharedFields,
  email: z.string().trim(),
  password: z.string().refine((value) => value === '' || (value.length >= 8 && value.length <= 72), 'Password 8-72 karakter'),
})

export type UserFormValues = {
  name: string
  email: string
  password: string
  role: string
  phone: string
  nomor_registrasi: string
}

export type UserFormMode = 'create' | 'edit'

export type UserFormValidation =
  | { ok: true; data: UserFormValues }
  | { ok: false; errors: Record<string, string> }

/** Validasi form user; mengembalikan data yang sudah dinormalisasi bila valid. */
export function validateUserForm(values: UserFormValues, mode: UserFormMode): UserFormValidation {
  const schema = mode === 'create' ? userCreateSchema : userEditSchema
  const result = schema.safeParse(values)
  if (result.success) return { ok: true, data: result.data as UserFormValues }

  const errors: Record<string, string> = {}
  for (const issue of result.error.issues) {
    const key = String(issue.path[0] ?? '')
    if (key && !(key in errors)) errors[key] = issue.message
  }
  return { ok: false, errors }
}
