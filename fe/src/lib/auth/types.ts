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
