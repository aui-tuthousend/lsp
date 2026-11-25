// Schema & aturan validasi form master TUK.
//
// Field mengikuti body `POST /tuks` dan `PATCH /tuks/{id}` (modul tuk).

import { z } from 'zod'

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const PHONE_RE = /^[0-9+\-\s()]{5,20}$/
const HTTP_URL_RE = /^https?:\/\/\S+$/

export const tukFormSchema = z.object({
  name: z.string().trim().min(2, 'Nama minimal 2 karakter').max(255, 'Maksimal 255 karakter'),
  postal_code: z.string().trim().min(3, 'Kode pos minimal 3 karakter').max(10, 'Maksimal 10 karakter'),
  phone: z.string().trim().regex(PHONE_RE, 'Nomor telepon tidak valid (5-20 karakter)'),
  email: z.string().trim().max(255, 'Maksimal 255 karakter').regex(EMAIL_RE, 'Format email tidak valid'),
  province: z.string().trim().min(2, 'Provinsi minimal 2 karakter').max(120, 'Maksimal 120 karakter'),
  city: z.string().trim().min(2, 'Kota/Kabupaten minimal 2 karakter').max(120, 'Maksimal 120 karakter'),
  district: z.string().trim().min(2, 'Kecamatan minimal 2 karakter').max(120, 'Maksimal 120 karakter'),
  village: z.string().trim().min(2, 'Desa/Kelurahan minimal 2 karakter').max(120, 'Maksimal 120 karakter'),
  address: z.string().trim().min(5, 'Alamat minimal 5 karakter'),
  maps_url: z.string().trim().refine((value) => value === '' || HTTP_URL_RE.test(value), 'Link maps tidak valid'),
  logo: z.string().trim().refine((value) => value === '' || HTTP_URL_RE.test(value), 'URL logo tidak valid'),
  letterhead: z.string().trim().refine((value) => value === '' || HTTP_URL_RE.test(value), 'URL kop surat tidak valid'),
})

export type TUKFormValues = z.infer<typeof tukFormSchema>

export type TUKFormValidation =
  | { ok: true; data: TUKFormValues }
  | { ok: false; errors: Record<string, string> }

/** Validasi form TUK; mengembalikan data yang sudah dinormalisasi bila valid. */
export function validateTUKForm(values: unknown): TUKFormValidation {
  const result = tukFormSchema.safeParse(values)
  if (result.success) return { ok: true, data: result.data }

  const errors: Record<string, string> = {}
  for (const issue of result.error.issues) {
    const key = String(issue.path[0] ?? '')
    if (key && !(key in errors)) errors[key] = issue.message
  }
  return { ok: false, errors }
}
