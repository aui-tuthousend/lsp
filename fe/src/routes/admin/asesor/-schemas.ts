// Schema & aturan validasi form master asesor.
//
// Form memuat dua bagian: data akun (name/email/password/phone/nomor_registrasi)
// dan data profil asesor (tuk_id/province/city/daily_assignment_limit).
// Email & password wajib hanya saat create; saat edit password opsional (reset).

import { z } from 'zod'

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const PHONE_RE = /^[0-9+\-\s()]{5,20}$/

const sharedFields = {
  name: z.string().trim().min(2, 'Nama minimal 2 karakter').max(255, 'Maksimal 255 karakter'),
  phone: z.string().trim().refine((value) => value === '' || PHONE_RE.test(value), 'Nomor telepon tidak valid (5-20 karakter)'),
  nomor_registrasi: z.string().trim().max(50, 'Maksimal 50 karakter'),
  tuk_id: z.string().trim(),
  province: z.string().trim().max(120, 'Maksimal 120 karakter'),
  city: z.string().trim().max(120, 'Maksimal 120 karakter'),
  daily_assignment_limit: z.number().int('Harus bilangan bulat').min(0, 'Minimal 0'),
}

export const asesorCreateSchema = z.object({
  ...sharedFields,
  email: z.string().trim().max(255, 'Maksimal 255 karakter').regex(EMAIL_RE, 'Format email tidak valid'),
  password: z.string().min(8, 'Minimal 8 karakter').max(72, 'Maksimal 72 karakter'),
})

export const asesorEditSchema = z.object({
  ...sharedFields,
  email: z.string().trim(),
  password: z.string().refine((value) => value === '' || (value.length >= 8 && value.length <= 72), 'Password 8-72 karakter'),
})

export type AsesorFormValues = {
  name: string
  email: string
  password: string
  phone: string
  nomor_registrasi: string
  tuk_id: string
  province: string
  city: string
  daily_assignment_limit: number
}

export type AsesorFormMode = 'create' | 'edit'

export type AsesorFormValidation =
  | { ok: true; data: AsesorFormValues }
  | { ok: false; errors: Record<string, string> }

/** Validasi form asesor; mengembalikan data yang sudah dinormalisasi bila valid. */
export function validateAsesorForm(values: AsesorFormValues, mode: AsesorFormMode): AsesorFormValidation {
  const schema = mode === 'create' ? asesorCreateSchema : asesorEditSchema
  const result = schema.safeParse(values)
  if (result.success) return { ok: true, data: result.data as AsesorFormValues }

  const errors: Record<string, string> = {}
  for (const issue of result.error.issues) {
    const key = String(issue.path[0] ?? '')
    if (key && !(key in errors)) errors[key] = issue.message
  }
  return { ok: false, errors }
}
