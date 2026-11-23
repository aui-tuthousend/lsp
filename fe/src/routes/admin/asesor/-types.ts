// Tipe untuk halaman master asesor (/admin/asesor).
//
// Menggabungkan dua sumber data backend:
// - Akun user (tabel `users`): POST /users, PATCH /users/{id}, DELETE /users/{id}
// - Profil asesor (tabel `asesor_profiles`): GET/PATCH /asesor-profiles/{user_id}
//
// Sertifikasi asesor tidak dikelola di sini (diisi sendiri oleh asesor).

/** Profil asesor seperti dikembalikan backend (list memuat name/email/tuk_name). */
export type AsesorProfile = {
  user_id: string
  name: string
  email: string
  tuk_id?: string | null
  tuk_name?: string | null
  province?: string | null
  city?: string | null
  daily_assignment_limit: number
  created_at?: string
  updated_at?: string
}

/** Akun user seperti dikembalikan POST /users dan PATCH /users/{id}. */
export type AsesorAccount = {
  id: string
  email: string
  name: string
  role: string
  isvalid: boolean
  phone?: string | null
  nomor_registrasi?: string | null
  created_at: string
  updated_at: string
}

/** Body pembuatan akun asesor (POST /users dengan role asesor). */
export type CreateAsesorBody = {
  email: string
  password: string
  name: string
  phone?: string
  nomor_registrasi?: string
}

/** Body update parsial akun (PATCH /users/{id}) — kirim field yang berubah saja. */
export type UpdateAsesorAccountBody = {
  name?: string
  phone?: string
  nomor_registrasi?: string
  password?: string
  isvalid?: boolean
}

/** Body update parsial profil asesor (PATCH /asesor-profiles/{user_id}). */
export type UpdateAsesorProfileBody = {
  tuk_id?: string
  province?: string
  city?: string
  daily_assignment_limit?: number
}

export type ListAsesorParams = {
  page?: number
  limit?: number
}

/** Opsi TUK untuk dropdown afiliasi. */
export type TukOption = {
  id: string
  name: string
}
