// Schema & aturan validasi form pendaftaran asesi (4 step).
//
// - Schema tiap step dipakai untuk memvalidasi saat user menekan "Lanjut".
// - `createAsesiProfilesSchema` adalah schema gabungan seluruh form.

import { z } from 'zod'

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const NIK_RE = /^\d{16}$/
const PHONE_RE = /^[0-9+\-\s()]{5,20}$/
const POSTAL_RE = /^\d{5}$/

const MAX_FILE_BYTES = 2 * 1024 * 1024 // 2MB

function formatAccept(types: string[]) {
  return types.map((type) => type.split('/')[1]?.toUpperCase() ?? type).join('/')
}

/** Dokumen opsional dengan batas ukuran & tipe file. */
function optionalFile(label: string, accepted: string[]) {
  return z
    .instanceof(File, { message: `${label} tidak valid` })
    .refine((file) => file.size <= MAX_FILE_BYTES, { message: `${label} maksimal 2MB` })
    .refine((file) => accepted.includes(file.type), { message: `${label} harus ${formatAccept(accepted)}` })
    .nullable()
    .optional()
}

// ── Step 1 — Data Pribadi ────────────────────────────────────────────────────
export const dataPribadiSchema = z.object({
  full_name: z.string().trim().min(1, 'Nama lengkap wajib diisi').max(100, 'Maksimal 255 karakter'),
  nik: z.string().trim().regex(NIK_RE, 'NIK harus terdiri dari 16 digit angka'),
  phone: z.string().trim().regex(PHONE_RE, 'Nomor HP tidak valid (5-20 karakter)'),
  gender: z.enum(['LAKI-LAKI', 'PEREMPUAN'], 'Pilih jenis kelamin'),
  birth_place: z.string().trim().min(1, 'Tempat lahir wajib diisi').max(255, 'Maksimal 255 karakter'),
  birth_date: z
    .string()
    .min(1, 'Tanggal lahir wajib diisi')
    .refine((value) => !Number.isNaN(Date.parse(value)), 'Format tanggal tidak valid')
    .refine((value) => new Date(value) <= new Date(), 'Tanggal lahir tidak boleh di masa depan'),
  nationality: z.string().trim().min(1, 'Kewarganegaraan wajib diisi').max(100, 'Maksimal 100 karakter'),
  occupation: z.string().trim().min(1, 'Pekerjaan wajib diisi').max(255, 'Maksimal 255 karakter'),
  education_level: z.enum(
    ['SD', 'SMP', 'SMA_SMK', 'D1', 'D2', 'D3', 'D4', 'S1', 'S2', 'S3'],
    'Pilih pendidikan terakhir',
  ),
})

// ── Step 2 — Alamat & Domisili ───────────────────────────────────────────────
const alamatDomisiliBase = z.object({
  ktp_address: z.string().trim().min(1, 'Alamat KTP wajib diisi'),
  domicile_same_as_ktp: z.boolean(),
  domicile_address: z.string().trim(),
  postal_code: z.string().trim().regex(POSTAL_RE, 'Kode pos harus 5 digit angka'),
  province: z.string().trim().min(1, 'Provinsi belum terisi'),
  city: z.string().trim().min(1, 'Kota/Kabupaten belum terisi'),
  district: z.string().trim().min(1, 'Kecamatan belum terisi'),
  village: z.string().trim().min(1, 'Desa/Kelurahan belum terisi'),
})

export const alamatDomisiliSchema = alamatDomisiliBase.superRefine((value, ctx) => {
  if (!value.domicile_same_as_ktp && value.domicile_address.trim().length === 0) {
    ctx.addIssue({
      code: 'custom',
      path: ['domicile_address'],
      message: 'Alamat domisili wajib diisi bila berbeda dari KTP',
    })
  }
})

// ── Step 3 — Unggah Dokumen Pendukung ────────────────────────────────────────
export const dokumenPendukungSchema = z.object({
  photo: optionalFile('Pas foto', ['image/jpeg', 'image/png']),
  ktp_file: optionalFile('Scan KTP', ['image/jpeg', 'image/png', 'application/pdf']),
  cv_file: optionalFile('CV', ['application/pdf']),
  npwp_file: optionalFile('NPWP', ['image/jpeg', 'image/png', 'application/pdf']),
})

// ── Step 4 — Email & Keamanan Akun ───────────────────────────────────────────
const akunKeamananBase = z.object({
  email: z.string().trim().regex(EMAIL_RE, 'Format email tidak valid'),
  password: z.string().min(8, 'Password minimal 8 karakter').max(72, 'Password maksimal 72 karakter'),
  password_confirmation: z.string().min(1, 'Konfirmasi password wajib diisi'),
})

export const akunKeamananSchema = akunKeamananBase.superRefine((value, ctx) => {
  if (value.password !== value.password_confirmation) {
    ctx.addIssue({
      code: 'custom',
      path: ['password_confirmation'],
      message: 'Konfirmasi password tidak cocok',
    })
  }
})

/** Schema gabungan seluruh form (semua step). */
export const createAsesiProfilesSchema = z
  .object({
    ...dataPribadiSchema.shape,
    ...alamatDomisiliBase.shape,
    ...dokumenPendukungSchema.shape,
    ...akunKeamananBase.shape,
  })
  .superRefine((value, ctx) => {
    if (!value.domicile_same_as_ktp && value.domicile_address.trim().length === 0) {
      ctx.addIssue({
        code: 'custom',
        path: ['domicile_address'],
        message: 'Alamat domisili wajib diisi bila berbeda dari KTP',
      })
    }
    if (value.password !== value.password_confirmation) {
      ctx.addIssue({
        code: 'custom',
        path: ['password_confirmation'],
        message: 'Konfirmasi password tidak cocok',
      })
    }
  })

/** Schema per step, urut sesuai urutan langkah pada form. */
export const stepSchemas = [
  dataPribadiSchema,
  alamatDomisiliSchema,
  dokumenPendukungSchema,
  akunKeamananSchema,
] as const

export type StepValidation = {
  ok: boolean
  errors: Record<string, string>
}

function toErrorMap(error: z.ZodError): Record<string, string> {
  const map: Record<string, string> = {}
  for (const issue of error.issues) {
    const key = String(issue.path[0] ?? '')
    if (key && !(key in map)) map[key] = issue.message
  }
  return map
}

/** Validasi satu step berdasarkan index-nya (0..3). */
export function validateStep(step: number, values: unknown): StepValidation {
  const schema = stepSchemas[step]
  if (!schema) return { ok: true, errors: {} }
  const result = schema.safeParse(values)
  return result.success ? { ok: true, errors: {} } : { ok: false, errors: toErrorMap(result.error) }
}

/** Validasi seluruh form sekaligus. */
export function validateAll(values: unknown): StepValidation {
  const result = createAsesiProfilesSchema.safeParse(values)
  return result.success ? { ok: true, errors: {} } : { ok: false, errors: toErrorMap(result.error) }
}
