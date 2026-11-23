// Tipe untuk master data TUK (Tempat Uji Kompetensi).
// Field mengikuti DTO backend `TUKDTO` (modul tuk) dan body POST/PATCH `/tuks`.

/** TUK seperti dikembalikan backend. */
export type TUK = {
  id: string
  name: string
  logo?: string | null
  postal_code: string
  phone: string
  email: string
  province: string
  city: string
  address: string
  maps_url?: string | null
  letterhead?: string | null
  created_at: string
  updated_at: string
  deleted_at?: string | null
}

/** Body pembuatan TUK (POST /tuks). */
export type CreateTUKBody = {
  name: string
  logo?: string
  postal_code: string
  phone: string
  email: string
  province: string
  city: string
  address: string
  maps_url?: string
  letterhead?: string
}

/** Body update parsial (PATCH /tuks/{id}) — kirim field yang berubah saja. */
export type UpdateTUKBody = Partial<CreateTUKBody>

export type ListTUKParams = {
  page?: number
  limit?: number
}
