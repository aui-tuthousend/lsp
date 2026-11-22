import { registerFn } from '@/lib/auth/session'
import { uploadFile, type UploadPurpose } from '@/lib/upload'
import type { CreateAsesiProfiles } from './-types'

async function uploadIfPresent(file: File | null, purpose: UploadPurpose) {
  if (!file) return undefined
  return uploadFile(file, purpose)
}

/**
 * Unggah dokumen (bila ada) lalu kirim pendaftaran asesi (akun + profil) ke
 * backend. Memakai server function agar cookie sesi ikut di-set, konsisten
 * dengan login. Dokumen dikirim sebagai object key hasil unggah.
 */
export async function registerAsesi(values: CreateAsesiProfiles) {
  const [photo, ktpFile, cvFile, npwpFile] = await Promise.all([
    uploadIfPresent(values.photo, 'photo'),
    uploadIfPresent(values.ktp_file, 'ktp_file'),
    uploadIfPresent(values.cv_file, 'cv_file'),
    uploadIfPresent(values.npwp_file, 'npwp_file'),
  ])

  return registerFn({
    data: {
      email: values.email.trim(),
      password: values.password,
      full_name: values.full_name.trim(),
      phone: values.phone.trim(),
      nik: values.nik.trim(),
      gender: values.gender,
      birth_place: values.birth_place.trim(),
      birth_date: values.birth_date,
      nationality: values.nationality.trim(),
      occupation: values.occupation.trim(),
      education_level: values.education_level,
      ktp_address: values.ktp_address.trim(),
      domicile_same_as_ktp: values.domicile_same_as_ktp,
      domicile_address: values.domicile_same_as_ktp ? undefined : values.domicile_address.trim() || undefined,
      province: values.province.trim(),
      city: values.city.trim(),
      district: values.district.trim(),
      village: values.village.trim(),
      postal_code: values.postal_code.trim() || undefined,
      photo,
      ktp_file: ktpFile,
      cv_file: cvFile,
      npwp_file: npwpFile,
    },
  })
}
