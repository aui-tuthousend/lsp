export type Role = 'admin' | 'asesi' | 'asesor' | 'validator' | 'manager_sertifikasi'

export type Session = {
  id: string
  role: Role
  name: string
  token: string
}

export function homeForRole(role: Role): string {
  switch (role) {
    case 'admin':
      return '/admin'
    case 'asesor':
      return '/asesor'
    case 'asesi':
      return '/asesi'
    case 'validator':
      return '/validator'
    case 'manager_sertifikasi':
      return '/manager'
  }
}

/** Path halaman profil per role. `null` bila belum tersedia. */
export function profileForRole(role: Role): string | null {
  switch (role) {
    case 'asesi':
      return '/asesi/profil'
    default:
      return null
  }
}
