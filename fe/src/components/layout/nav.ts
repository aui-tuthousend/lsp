import { Building2, ClipboardCheck, ClipboardList, FileCheck, FileText, Home, Layers, LayoutDashboard, ShieldCheck, UserCog, Users } from 'lucide-react'
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

const ALL: Role[] = ['admin', 'admin_lsp', 'asesi', 'asesor', 'validator', 'manager_sertifikasi']

// Role yang beranda-nya memakai item "Dashboard" umum (via '/'). manager_sertifikasi
// punya menu tersendiri sehingga tidak ikut item umum ini.
const GENERAL_DASHBOARD: Role[] = ['admin', 'admin_lsp', 'asesi', 'asesor', 'validator']

export const navSections: NavSection[] = [
  {
    title: 'Umum',
    items: [
      { label: 'Dashboard', to: '/', icon: Home, roles: GENERAL_DASHBOARD },
      { label: 'Materi Uji Kompetensi', to: '/admin_lsp/muk', icon: FileCheck, roles: ['admin_lsp'] },
    ],
  },
  {
    title: 'Manajemen Sertifikasi',
    items: [{ label: 'Dashboard', to: '/manager', icon: LayoutDashboard, roles: ['manager_sertifikasi'] }],
  },
  {
    title: 'Sertifikasi',
    items: [
      { label: 'Asesmen Mandiri', to: '/asesi/asesmen-mandiri', icon: ClipboardCheck, roles: ['asesi'] },
    ],
  },
  {
    title: 'Penugasan',
    items: [{ label: 'Penugasan Saya', to: '/asesor/assignments', icon: ClipboardList, roles: ['asesor'] }],
  },
  {
    title: 'Master Data',
    items: [
      { label: 'TUK', to: '/admin/tuks', icon: Building2, roles: ['admin'] },
      { label: 'Asesor', to: '/admin/asesor', icon: Users, roles: ['admin'] },
      { label: 'Skema Sertifikasi', to: '/admin/skema', icon: Layers, roles: ['admin'] },
      { label: 'Jenis Dokumen', to: '/admin/dokumen', icon: FileText, roles: ['admin'] },
      { label: 'Users', to: '/admin/users', icon: UserCog, roles: ['admin'] },
      { label: 'Roles', to: '/admin/roles', icon: ShieldCheck, roles: ['admin'] },
    ],
  },
]

export function navFor(role: Role) {
  return navSections
    .map((section) => ({ ...section, items: section.items.filter((item) => item.roles.includes(role)) }))
    .filter((section) => section.items.length > 0)
}
