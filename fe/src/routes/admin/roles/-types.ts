// Tipe untuk halaman master data role (/admin/roles).
//
// Field mengikuti DTO backend (modul users) dan body POST/PATCH /roles.

/** Role seperti dikembalikan backend. */
export type Role = {
  id: string
  code: string
  name: string
  description?: string | null
}

/** Body pembuatan role (POST /roles). */
export type CreateRoleBody = {
  code: string
  name: string
  description?: string
}

/** Body update parsial (PATCH /roles/{id}) — kirim field yang berubah saja. */
export type UpdateRoleBody = Partial<CreateRoleBody>

/** Respons daftar role: { data: Role[] }. */
export type RoleListResponse = {
  data: Role[]
}
