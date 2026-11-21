import { Building2, ClipboardList, Home, UserCog, Users } from 'lucide-react'
import type { Role } from '@/lib/auth/types'

export type NavItem = {
  label: string
  to: string
  icon: typeof Home
  roles: Role[]
}

export type NavSection = {
  title?: string
  items: NavItem[]
}

const ALL: Role[] = ['admin', 'asesi', 'asesor', 'validator', 'manager_sertifikasi']

export const navSections: NavSection[] = [
  {
    title: 'Umum',
    items: [{ label: 'Dashboard', to: '/', icon: Home, roles: ALL }],
  },
  {
    title: 'Penugasan',
    items: [{ label: 'Penugasan Saya', to: '/asesor/assignments', icon: ClipboardList, roles: ['asesor'] }],
  },
  {
    title: 'Master Data',
    items: [
      { label: 'TUK', to: '/admin/tuks', icon: Building2, roles: ['admin'] },
      { label: 'Profil Asesor', to: '/admin/asesor-profiles', icon: Users, roles: ['admin'] },
      { label: 'Users', to: '/admin/users', icon: UserCog, roles: ['admin'] },
    ],
  },
]

export function navFor(role: Role) {
  return navSections
    .map((section) => ({ ...section, items: section.items.filter((item) => item.roles.includes(role)) }))
    .filter((section) => section.items.length > 0)
}
