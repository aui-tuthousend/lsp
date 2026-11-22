// Tipe data pendaftaran asesi untuk form multi-step di halaman /register.
//
// Field mengikuti skema backend `asesi_profiles` (lihat migrasi 000005) dengan
// tambahan data akun (email & password) yang dipakai untuk membuat user.
// Catatan: field enum bertipe `'' ` saat belum dipilih agar form bisa mulai
// kosong; nilai valid akan dipastikan oleh schema di `-schemas.ts`.

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

/** Step 1 — Data Pribadi */
export type DataPribadi = {
  full_name: string
  nik: string
  phone: string
  gender: Gender | ''
  birth_place: string
  birth_date: string
  nationality: string
  occupation: string
  education_level: EducationLevel | ''
}

/** Step 2 — Alamat & Domisili */
export type AlamatDomisili = {
  ktp_address: string
  domicile_same_as_ktp: boolean
  domicile_address: string
  postal_code: string
  province: string
  city: string
  district: string
  village: string
}

/** Step 3 — Unggah Dokumen Pendukung */
export type DokumenPendukung = {
  photo: File | null
  ktp_file: File | null
  cv_file: File | null
  npwp_file: File | null
}

/** Step 4 — Email & Keamanan Akun */
export type AkunKeamanan = {
  email: string
  password: string
  password_confirmation: string
}

/** Payload lengkap pendaftaran asesi (gabungan 4 step). */
export type CreateAsesiProfiles = DataPribadi & AlamatDomisili & DokumenPendukung & AkunKeamanan

/** Hasil pencarian kode pos → wilayah (provinsi/kota/kecamatan/desa). */
export type PostalCodeResult = {
  province: string
  city: string
  district: string
  village: string
}
