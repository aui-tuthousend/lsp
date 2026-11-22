// Tipe untuk halaman profil asesi (/asesi/profil).
//
// Field mengikuti DTO backend `AsesiProfileDTO` (modul asesi) beserta body
// `POST /asesi-profiles` dan `PATCH /asesi-profiles/{id}`.

export type Gender = 'LAKI-LAKI' | 'PEREMPUAN'

export type EducationLevel =
  | 'SD'
  | 'SMP'
  | 'SMA_SMK'
  | 'D1'
  | 'D2'
  | 'D3'
  | 'D4'
  | 'S1'
  | 'S2'
  | 'S3'

/** Profil asesi seperti dikembalikan backend (GET). */
export type AsesiProfile = {
  id: string
  user_id: string
  full_name: string
  nik: string
  phone: string
  gender: Gender
  birth_place: string
  birth_date: string
  nationality: string
  occupation: string
  education_level: EducationLevel
  ktp_address: string
  domicile_same_as_ktp: boolean
  domicile_address?: string | null
  province: string
  city: string
  district: string
  village: string
  postal_code?: string | null
  photo?: string | null
  ktp_file?: string | null
  cv_file?: string | null
  npwp_file?: string | null
  created_at: string
  updated_at: string
}

/** Body pembuatan profil (POST /asesi-profiles). */
export type CreateAsesiProfileBody = {
  user_id: string
  full_name: string
  nik: string
  phone: string
  gender: Gender
  birth_place: string
  birth_date: string
  nationality?: string
  occupation: string
  education_level: EducationLevel
  ktp_address: string
  domicile_same_as_ktp: boolean
  domicile_address?: string
  province: string
  city: string
  district: string
  village: string
  postal_code?: string
  photo?: string
  ktp_file?: string
  cv_file?: string
  npwp_file?: string
}

/** Body update parsial (PATCH /asesi-profiles/{id}) — kirim field yang berubah saja. */
export type UpdateAsesiProfileBody = Partial<Omit<CreateAsesiProfileBody, 'user_id'>>

export type ListAsesiProfilesParams = {
  page?: number
  limit?: number
}

/** Akun user seperti dikembalikan GET /auth/me (tabel `users`). */
export type UserAccount = {
  id: string
  email: string
  name: string
  role: string
  isvalid: boolean
  phone?: string | null
  nomor_registrasi?: string | null
  email_verified_at?: string | null
  created_at: string
  updated_at: string
}
