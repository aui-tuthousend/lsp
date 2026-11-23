// Tipe untuk halaman master user (/admin/users).
//
// Mengelola akun pada tabel `users` via endpoint admin:
// POST /users, GET /users, PATCH /users/{id}, DELETE /users/{id}.
// Role yang dikelola di sini hanya validator & manager_sertifikasi
// (asesor punya halaman sendiri; asesi mendaftar mandiri; admin tidak dibuat di sini).

/** User seperti dikembalikan backend (GET /users). */
export type User = {
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

/** Body pembuatan user (POST /users). */
export type CreateUserBody = {
  email: string
  password: string
  name: string
  role: string
  phone?: string
  nomor_registrasi?: string
}

/** Body update parsial (PATCH /users/{id}) — kirim field yang berubah saja. */
export type UpdateUserBody = {
  name?: string
  phone?: string
  nomor_registrasi?: string
  password?: string
  isvalid?: boolean
}

export type ListUsersParams = {
  page?: number
  limit?: number
  /** Filter kode role; dikirim sebagai query `roles` berulang. */
  roles?: string[]
}
