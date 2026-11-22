import type { AsesiProfile, EducationLevel, Gender, UpdateAsesiProfileBody } from './-types'

// Konstanta, tipe, dan helper murni untuk form profil asesi. State form tetap
// page-local di komponen; di sini hanya bagian yang bebas dari React.

export const GENDER_OPTIONS: { value: Gender; label: string }[] = [
  { value: 'LAKI-LAKI', label: 'Laki-laki' },
  { value: 'PEREMPUAN', label: 'Perempuan' },
]

export const EDUCATION_OPTIONS: { value: EducationLevel; label: string }[] = [
  { value: 'SD', label: 'SD' },
  { value: 'SMP', label: 'SMP' },
  { value: 'SMA_SMK', label: 'SMA/SMK' },
  { value: 'D1', label: 'D1' },
  { value: 'D2', label: 'D2' },
  { value: 'D3', label: 'D3' },
  { value: 'D4', label: 'D4' },
  { value: 'S1', label: 'S1' },
  { value: 'S2', label: 'S2' },
  { value: 'S3', label: 'S3' },
]

/** Field yang bisa diubah, semuanya berasal dari tabel `asesi_profiles`. */
export type ProfileFormValues = {
  full_name: string
  nik: string
  phone: string
  gender: Gender | ''
  birth_place: string
  birth_date: string
  nationality: string
  occupation: string
  education_level: EducationLevel | ''
  ktp_address: string
  domicile_same_as_ktp: boolean
  domicile_address: string
  postal_code: string
  province: string
  city: string
  district: string
  village: string
}

export const EMPTY_FORM: ProfileFormValues = {
  full_name: '',
  nik: '',
  phone: '',
  gender: '',
  birth_place: '',
  birth_date: '',
  nationality: 'Indonesia',
  occupation: '',
  education_level: '',
  ktp_address: '',
  domicile_same_as_ktp: true,
  domicile_address: '',
  postal_code: '',
  province: '',
  city: '',
  district: '',
  village: '',
}

export function toForm(profile: AsesiProfile): ProfileFormValues {
  return {
    full_name: profile.full_name,
    nik: profile.nik,
    phone: profile.phone,
    gender: profile.gender,
    birth_place: profile.birth_place,
    birth_date: profile.birth_date,
    nationality: profile.nationality,
    occupation: profile.occupation,
    education_level: profile.education_level,
    ktp_address: profile.ktp_address,
    domicile_same_as_ktp: profile.domicile_same_as_ktp,
    domicile_address: profile.domicile_address ?? '',
    postal_code: profile.postal_code ?? '',
    province: profile.province,
    city: profile.city,
    district: profile.district,
    village: profile.village,
  }
}

/** Bentuk payload yang sudah dinormalisasi (trim, cast enum, dll). */
function toPatchShape(values: ProfileFormValues): UpdateAsesiProfileBody {
  return {
    full_name: values.full_name.trim(),
    nik: values.nik.trim(),
    phone: values.phone.trim(),
    gender: values.gender as Gender,
    birth_place: values.birth_place.trim(),
    birth_date: values.birth_date,
    nationality: values.nationality.trim(),
    occupation: values.occupation.trim(),
    education_level: values.education_level as EducationLevel,
    ktp_address: values.ktp_address.trim(),
    domicile_same_as_ktp: values.domicile_same_as_ktp,
    domicile_address: values.domicile_same_as_ktp ? undefined : values.domicile_address.trim(),
    postal_code: values.postal_code.trim() || undefined,
    province: values.province.trim(),
    city: values.city.trim(),
    district: values.district.trim(),
    village: values.village.trim(),
  }
}

/** Bandingkan dengan nilai awal; kirim hanya field yang berubah (partial update). */
export function diffPatch(values: ProfileFormValues, baseline: ProfileFormValues): UpdateAsesiProfileBody {
  const next = toPatchShape(values)
  const prev = toPatchShape(baseline)
  const patch: Record<string, unknown> = {}
  for (const key of Object.keys(next) as (keyof UpdateAsesiProfileBody)[]) {
    if (next[key] !== prev[key]) patch[key] = next[key]
  }
  return patch as UpdateAsesiProfileBody
}

export function validate(values: ProfileFormValues): string | null {
  if (!values.full_name.trim()) return 'Nama lengkap wajib diisi'
  if (values.nik.trim().length !== 16) return 'NIK harus 16 digit'
  if (!values.phone.trim()) return 'Nomor HP wajib diisi'
  if (!values.gender) return 'Jenis kelamin wajib dipilih'
  if (!values.birth_place.trim()) return 'Tempat lahir wajib diisi'
  if (!values.birth_date) return 'Tanggal lahir wajib diisi'
  if (!values.education_level) return 'Pendidikan terakhir wajib dipilih'
  if (!values.occupation.trim()) return 'Pekerjaan wajib diisi'
  if (!values.ktp_address.trim()) return 'Alamat KTP wajib diisi'
  if (!values.domicile_same_as_ktp && !values.domicile_address.trim()) return 'Alamat domisili wajib diisi'
  if (!values.province.trim() || !values.city.trim() || !values.district.trim() || !values.village.trim()) {
    return 'Wilayah (provinsi, kota, kecamatan, desa) wajib diisi'
  }
  return null
}
