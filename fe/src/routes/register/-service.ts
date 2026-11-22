import { registerFn } from '@/lib/auth/session'
import type { CreateAsesiProfiles } from './-types'

/**
 * Kirim pendaftaran asesi (akun + profil) ke backend.
 * Memakai server function agar cookie sesi ikut di-set, konsisten dengan login.
 */
export function registerAsesi(values: CreateAsesiProfiles) {
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
    },
  })
}
