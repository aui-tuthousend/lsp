# Setup Frontend (TanStack Start + shadcn) — Rencana

> Status: **rencana / belum diimplementasikan.**
> Dokumen ini jadi acuan saat mulai membangun `fe/`.
> Backend ada di `be/` (Go + Huma), FE ada di `fe/` (standalone Bun project, dependency terisolasi di `fe/node_modules`).

## 0. Prinsip

- **Selalu pakai bun** untuk install, script, dan menjalankan tooling (`bun`, `bunx`). Jangan `npm`/`yarn`/`pnpm`.
- Satu **registry navigasi** + **satu komponen `AppShell`** yang dipakai oleh tiap layout role (routing folder-based).
- **Session** disimpan di cookie httpOnly (token JWT) dan dibaca lewat **middleware** TanStack.
- **SSE** satu koneksi per tab; server mem-multiplex topik (`user:<id>` + `role:<role>`).
- Semua data server lewat **TanStack Query** (staleTime eksplisit), bukan fetch ad-hoc.

---

## 1. Tooling & versi

| Tool | Keterangan |
|---|---|
| bun | package manager + runtime script |
| TanStack Start | framework full-stack (Router + Vite) |
| TanStack Query | data fetching/caching |
| shadcn/ui | komponen (Tailwind) |
| Tailwind CSS | styling |

Script `fe/package.json` (draft):

```json
{
  "scripts": {
    "dev": "bun --bun vite dev",
    "build": "bun --bun vite build",
    "start": "bun --bun vite start",
    "typecheck": "bun --bun tsc --noEmit",
    "lint": "bun --bun eslint .",
    "format": "bun --bun prettier --write ."
  }
}
```

> Nama script `vite dev`/`vite build` mengikuti TanStack Start (Vite-based). Sesuaikan bila scaffolder menghasilkan nama berbeda.

---

## 2. Scaffold TanStack Start

```bash
cd /home/iverson/code/lsp
bunx create-start-app@latest fe      # gunakan scaffolder resmi TanStack Start
cd fe
bun install
```

> Nama paket scaffolder bisa berubah — cek dokumentasi resmi TanStack Start. Inti: hasilnya harus punya `src/routes/`, `vite.config.ts`, dan entry Start.

---

## 3. shadcn/ui

```bash
cd fe
bunx --bun shadcn@latest init
bunx --bun shadcn@latest add sidebar button dropdown-menu avatar input form table dialog sheet badge card separator skeleton sonner
```

- Komponen masuk ke `src/components/ui/`.
- Blok `Sidebar` shadcn (`SidebarProvider`, `Sidebar`, `SidebarMenu`, `SidebarMenuButton`) dipakai **sekali** di `AppShell`.

---

## 4. Struktur folder (target)

**Routing: folder-based, mengikuti pola Next.js.** Setiap role punya segmen path sendiri (`/admin`, `/asesor`, `/asesi`, `/validator`, `/manager`). Di TanStack Router, `routes/<segmen>/route.tsx` = layout untuk path tersebut, dan file di dalamnya jadi child.

```
fe/
├── src/
│   ├── routes/
│   │   ├── __root.tsx                     # root + router context (queryClient, auth)
│   │   ├── index.tsx                      # "/" → redirect ke home sesuai role
│   │   ├── login.tsx                      # "/login"  (publik, guest-only)
│   │   ├── register.tsx                   # "/register" (publik, guest-only)
│   │   ├── verify.$token.tsx              # "/verify/:token" (publik, landing QR asesi)
│   │   ├── scan.$token.tsx                # "/scan/:token" (fullscreen, tanpa sidebar)
│   │   │
│   │   ├── admin/
│   │   │   ├── route.tsx                  # layout "/admin"  → guard: role admin
│   │   │   ├── index.tsx                  # "/admin"        (dashboard admin)
│   │   │   ├── tuks.tsx                   # "/admin/tuks"
│   │   │   ├── asesor-profiles.tsx        # "/admin/asesor-profiles"
│   │   │   └── users.tsx                  # "/admin/users"
│   │   ├── asesor/
│   │   │   ├── route.tsx                  # layout "/asesor" → guard: role asesor
│   │   │   ├── index.tsx                  # "/asesor"
│   │   │   ├── assignments.tsx            # "/asesor/assignments"
│   │   │   └── assignments.$assignmentId.tsx # "/asesor/assignments/:id" (assigned-only)
│   │   ├── asesi/
│   │   │   ├── route.tsx                  # layout "/asesi"  → guard: role asesi
│   │   │   └── index.tsx
│   │   ├── validator/
│   │   │   ├── route.tsx                  # layout "/validator" → guard: role validator
│   │   │   └── index.tsx
│   │   └── manager/
│   │       ├── route.tsx                  # layout "/manager" → guard: role manager_sertifikasi
│   │       └── index.tsx
│   ├── components/
│   │   ├── ui/                            # shadcn
│   │   └── layout/
│   │       ├── app-shell.tsx              # Sidebar + Header + <Outlet/> (dipakai tiap layout role)
│   │       ├── sidebar.tsx
│   │       └── nav.ts                     # registry navigasi + navFor(role)
│   ├── lib/
│   │   ├── api.ts                         # fetch wrapper (Bearer token)
│   │   ├── query.ts                       # QueryClient + default options
│   │   ├── query-keys.ts                  # query key factory
│   │   ├── auth/
│   │   │   ├── session.ts                 # server fn: login/logout/getSession
│   │   │   ├── guards.ts                  # helper beforeLoad: requireGuest/requireAuth/requireRole
│   │   │   ├── middleware.ts              # createMiddleware (server function)
│   │   │   └── types.ts                   # Session type + homeForRole(role)
│   │   └── sse/
│   │       ├── use-notifications.ts       # hook EventSource
│   │       └── events.ts                  # tipe event + mapping
│   └── start.ts                           # registerGlobalMiddleware (opsional)
├── .env.local
└── package.json
```

Catatan: `AppShell` adalah **satu komponen** yang dirender oleh setiap `routes/<role>/route.tsx`. Jadi tiap role punya layout route sendiri (agar guard & path jelas), tapi shell/sidebar-nya tetap satu — konsisten dengan keputusan "1 sidebar + registry nav".

---

## 5. Session & middleware (createMiddleware)

### 5.1 Bentuk session

```ts
// src/lib/auth/types.ts
export type Role = "admin" | "asesi" | "asesor" | "validator" | "manager_sertifikasi"

export type Session = {
  id: string      // user id
  role: Role
  name: string
  token: string   // JWT
}
```

Session hanya menyimpan 4 field ini (sesuai permintaan). Token JWT juga disimpan di cookie httpOnly supaya aman dari XSS; field `token` dipakai FE untuk header `Authorization` dan `?token=` SSE.

### 5.2 Server functions (cookie httpOnly)

```ts
// src/lib/auth/session.ts
import { createServerFn } from "@tanstack/react-start"
import { getCookie, setCookie, deleteCookie } from "@tanstack/react-start/server"
import type { Session } from "./types"

const COOKIE = "lsp_session"

export const getSessionFn = createServerFn({ method: "GET" }).handler(async (): Promise<Session | null> => {
  const raw = getCookie(COOKIE)
  if (!raw) return null
  try { return JSON.parse(raw) as Session } catch { return null }
})

export const loginFn = createServerFn({ method: "POST" })
  .validator((d: { email: string; password: string }) => d)
  .handler(async ({ data }): Promise<Session> => {
    const res = await fetch(`${process.env.API_BASE_URL}/auth/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(data),
    })
    if (!res.ok) throw new Error("login failed")
    const { token, user } = await res.json()
    const session: Session = { id: user.id, role: user.role, name: user.name, token }
    setCookie(COOKIE, JSON.stringify(session), {
      httpOnly: true, sameSite: "lax", path: "/", secure: process.env.NODE_ENV === "production",
    })
    return session
  })

export const logoutFn = createServerFn({ method: "POST" }).handler(async () => {
  deleteCookie(COOKIE, { path: "/" })
})
```

### 5.3 Middleware

```ts
// src/lib/auth/middleware.ts
import { createMiddleware } from "@tanstack/react-router"
import { redirect } from "@tanstack/react-router"
import { getSessionFn } from "./session"
import { homeForRole, type Role } from "./types"

// Memuat session ke context; dipakai server function maupun route guard.
export const authMiddleware = createMiddleware({ type: "function" }).server(async ({ next }) => {
  const session = await getSessionFn()
  if (!session) throw redirect({ to: "/login" })
  return next({ context: { session } })
})

// Guard role: batasi akses berdasarkan role.
export const requireRole = (...roles: Role[]) =>
  createMiddleware({ type: "function" }).server(async ({ next }) => {
    const session = await getSessionFn()
    if (!session) throw redirect({ to: "/login" })
    if (!roles.includes(session.role)) throw redirect({ to: homeForRole(session.role) })
    return next({ context: { session } })
  })
```

Contoh pemakaian di server function terproteksi:

```ts
export const listTuksFn = createServerFn({ method: "GET" })
  .middleware([requireRole("admin")])
  .handler(async ({ context }) => {
    const res = await fetch(`${process.env.API_BASE_URL}/tuks`, {
      headers: { Authorization: `Bearer ${context.session.token}` },
    })
    return res.json()
  })
```

### 5.4 Integrasi ke Router

```ts
// src/routes/__root.tsx (ringkas)
import { createRootRouteWithContext, Outlet } from "@tanstack/react-router"
import type { QueryClient } from "@tanstack/react-query"
import type { Session } from "@/lib/auth/types"

export type RouterContext = { queryClient: QueryClient; auth: Session | null }

export const Route = createRootRouteWithContext<RouterContext>()({
  beforeLoad: async () => ({ auth: await getSessionFn() }),
  component: () => <Outlet />,
})
```

### 5.5 Guard `beforeLoad` (3 kasus)

Semua guard dijalankan di **`beforeLoad`** pada layout/route yang relevan. Karena `context.auth` sudah diisi di `__root.tsx`, guard cukup membacanya.

Helper terpusat:

```ts
// src/lib/auth/types.ts (tambahan)
export function homeForRole(role: Role): string {
  switch (role) {
    case "admin": return "/admin"
    case "asesor": return "/asesor"
    case "asesi": return "/asesi"
    case "validator": return "/validator"
    case "manager_sertifikasi": return "/manager"
  }
}
```

```ts
// src/lib/auth/guards.ts
import { redirect } from "@tanstack/react-router"
import type { Role } from "./types"
import { homeForRole } from "./types"
import type { RouterContext } from "@/routes/__root"

// Kasus 1 — guest-only: sudah login tidak boleh ke /login, /register.
export function requireGuest(ctx: RouterContext) {
  if (ctx.auth) throw redirect({ to: homeForRole(ctx.auth.role) })
}

// Kasus 2 — auth-only: belum login tidak boleh ke route terproteksi.
export function requireAuth(ctx: RouterContext, redirectTo?: string) {
  if (!ctx.auth) {
    throw redirect({ to: "/login", search: redirectTo ? { redirect: redirectTo } : undefined })
  }
}

// Kasus 3 — role-only: role tidak sesuai diarahkan ke home role-nya.
export function requireRole(ctx: RouterContext, ...roles: Role[]) {
  requireAuth(ctx)
  if (ctx.auth && !roles.includes(ctx.auth.role)) {
    throw redirect({ to: homeForRole(ctx.auth.role) })
  }
}
```

**Kasus 1 — sudah login tidak boleh ke `/login`:**

```ts
// src/routes/login.tsx
export const Route = createFileRoute("/login")({
  beforeLoad: ({ context }) => requireGuest(context),
  component: LoginPage,
})
```

**Kasus 2 — belum login tidak boleh ke route terproteksi** (dipasang di tiap layout role):

```ts
// src/routes/asesor/route.tsx
export const Route = createFileRoute("/asesor")({
  beforeLoad: ({ context, location }) => requireRole(context, "asesor", /* + role lain bila perlu */),
  component: AppShell,
})
```

**Kasus 3 — role salah** sudah ditangani `requireRole` di atas: mis. asesi membuka `/admin` → diarahkan ke `/asesi`.

`/` (index) mengarahkan ke home sesuai role:

```ts
// src/routes/index.tsx
export const Route = createFileRoute("/")({
  beforeLoad: ({ context }) => {
    if (!context.auth) throw redirect({ to: "/login" })
    throw redirect({ to: homeForRole(context.auth.role) })
  },
})
```

### 5.6 Otorisasi tingkat resource (asesor yang ditugaskan)

Role guard hanya memastikan "dia seorang asesor". Untuk **halaman detail asesmen yang hanya boleh dibuka asesor yang ditugaskan**, perlu cek kepemilikan (ownership) terhadap resource, bukan sekadar role.

```ts
// src/routes/asesor/assignments.$assignmentId.tsx
import { createFileRoute, notFound, redirect } from "@tanstack/react-router"

export const Route = createFileRoute("/asesor/assignments/$assignmentId")({
  beforeLoad: async ({ context, params }) => {
    requireRole(context, "asesor")

    const assignment = await getAssignmentFn({ data: { id: params.assignmentId } })

    // Tidak ada → 404 (jangan bocorkan keberadaan resource).
    if (!assignment) throw notFound()

    // Ada tapi bukan milik asesor ini → tendang ke daftar penugasan.
    if (assignment.asesorId !== context.auth!.id) {
      throw redirect({ to: "/asesor/assignments" })
    }
  },
  loader: ({ params }) => getAssignmentDetailFn({ data: { id: params.assignmentId } }),
  component: AssignmentDetailPage,
})
```

Aturan penting:

- **Role ≠ ownership.** `requireRole` cek peran; cek `assignment.asesorId === session.id` cek kepemilikan. Keduanya diperlukan.
- **Guard FE hanya untuk UX** (redirect cepat, sembunyikan UI). **Backend tetap wajib otorisasi**: endpoint detail asesmen harus memverifikasi bahwa pemanggil adalah asesor yang ditugaskan, jika tidak → `403` (atau `404` agar tidak membocorkan keberadaan data). Jangan pernah mengandalkan guard klien sebagai pengamanan.
- Pilihan `404` vs `403`: pakai `404` bila keberadaan resource sensitif; `403` bila tidak masalah diketahui.

> Catatan umum: nama export `createMiddleware`/`createServerFn`/`getCookie` bisa berasal dari `@tanstack/react-router` atau `@tanstack/react-start` tergantung versi. Sesuaikan dengan versi yang ter-scaffold.

---

## 6. TanStack Query (staleTime & default)

```ts
// src/lib/query.ts
import { QueryClient } from "@tanstack/react-query"

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 60_000,          // 1 menit: data "segar" tanpa refetch
      gcTime: 5 * 60_000,         // 5 menit cache
      retry: 1,
      refetchOnWindowFocus: false,
      refetchOnReconnect: true,
    },
    mutations: { retry: 0 },
  },
})
```

Panduan `staleTime` per jenis data:

| Data | staleTime | Alasan |
|---|---|---|
| `auth/me`, daftar role | `Infinity` | hampir statis per sesi |
| master TUK, profil asesor | `5 * 60_000` | jarang berubah |
| daftar penugasan/notifikasi | `15_000` | cukup dinamis |
| detail realtime (status scan) | `0` | selalu ambil terbaru; SSE yang push update |

Query key factory:

```ts
// src/lib/query-keys.ts
export const qk = {
  me: ["me"] as const,
  tuks: (p?: object) => ["tuks", p] as const,
  asesorProfiles: (p?: object) => ["asesor-profiles", p] as const,
  assignments: (p?: object) => ["assignments", p] as const,
  notifications: ["notifications"] as const,
}
```

Pemakaian:

```ts
const { data } = useQuery({
  queryKey: qk.tuks({ page: 1 }),
  queryFn: () => listTuks({ page: 1 }),
  staleTime: 5 * 60_000,
})
```

Wiring `queryClient` ke router context (di root) supaya bisa `ensureQueryData` di `loader`.

---

## 7. SSE — notifikasi realtime per role

### 7.1 Konsep

- Satu `EventSource` per tab ke `GET /events`.
- Server menentukan topik dari claims: `user:<id>` **dan** `role:<role>`.
- FE tidak perlu buka koneksi per role — cukup satu, lalu dispatch berdasarkan `data.type` / `data.topic`.
- Nama event di stream selalu `notification`; jenis sebenarnya ada di `data.type`.

### 7.2 Token

`EventSource` tidak bisa mengatur header, jadi token dikirim via query: `?token=<jwt>`. (Backend sudah mendukung `?token=`.)

```ts
const url = `${import.meta.env.VITE_API_BASE_URL}/events?token=${encodeURIComponent(session.token)}`
```

### 7.3 Tipe event

```ts
// src/lib/sse/events.ts
export type AppEvent = {
  type: string                 // mis. "verification.scanned", "apl01.approved"
  topic: string                // mis. "user:<id>" | "role:<role>"
  data?: Record<string, unknown>
}

// Pemetaan event → query key yang perlu di-invalidate
export const eventInvalidation: Record<string, (e: AppEvent) => unknown[]> = {
  "verification.scanned": (e) => [["assignments"], ["verification", e.data?.session_id]],
  "apl01.approved": () => [["apl01"]],
  "apl02.submitted": () => [["apl02"]],
}
```

### 7.4 Hook

```ts
// src/lib/sse/use-notifications.ts
import { useEffect, useRef } from "react"
import { useQueryClient } from "@tanstack/react-query"
import type { Session } from "@/lib/auth/types"
import { eventInvalidation, type AppEvent } from "./events"

export function useNotifications(session: Session | null) {
  const qc = useQueryClient()
  const esRef = useRef<EventSource | null>(null)

  useEffect(() => {
    if (!session) return
    const url = `${import.meta.env.VITE_API_BASE_URL}/events?token=${encodeURIComponent(session.token)}`
    const es = new EventSource(url)
    esRef.current = es

    es.addEventListener("notification", (raw) => {
      const e = JSON.parse((raw as MessageEvent).data) as AppEvent
      if (e.type === "connected") return

      // invalidasi cache terkait
      const keys = eventInvalidation[e.type]?.(e) ?? []
      keys.forEach((k) => qc.invalidateQueries({ queryKey: k as unknown[] }))

      // toast opsional
      // toast.info(e.data?.message as string)
    })

    es.onerror = () => { /* EventSource auto-reconnect; biarkan browser */ }

    return () => { es.close(); esRef.current = null }
  }, [session, qc])
}
```

Dipasang sekali di `AppShell`:

```ts
function AppShell() {
  const { auth } = Route.useRouteContext()
  useNotifications(auth)
  return (/* Sidebar + Header + Outlet */)
}
```

### 7.5 Alur

```
Asesi scan QR  ─► POST /verifications/{token}/scan
                        │ (1 transaksi)
                        ├─ UPDATE verification_sessions
                        └─ INSERT outbox_events (payload: user_id = asesor)
                                │
        worker  ── FetchPendingForUpdate ─► MarkProcessed ─► COMMIT ─► publish
                                │
                    Redis Pub/Sub  channel: user:<asesor_id>
                                │
        API SSE handler (subscriber) ─► EventSource (browser asesor)
                                │
        useNotifications ─► invalidateQueries(["assignments"]) ─► UI refresh
```

> Redis Pub/Sub bersifat fire-and-forget; status tetap disimpan di DB. Saat reload, query biasa akan menarik status terbaru. SSE hanya untuk update live.

---

## 8. Navigasi (registry + guard)

```ts
// src/components/layout/nav.ts
import type { Role } from "@/lib/auth/types"

type NavItem = { label: string; to: string; icon: LucideIcon; roles: Role[] }
type NavSection = { title?: string; items: NavItem[] }

const ALL: Role[] = ["admin", "asesi", "asesor", "validator", "manager_sertifikasi"]

export const navSections: NavSection[] = [
  { title: "Umum", items: [
    { label: "Dashboard", to: "/", icon: Home, roles: ALL },   // "/" redirect ke home role
  ]},
  { title: "Penugasan", items: [
    { label: "Penugasan Saya", to: "/asesor/assignments", icon: ClipboardList, roles: ["asesor"] },
  ]},
  { title: "Master Data", items: [
    { label: "TUK", to: "/admin/tuks", icon: Building2, roles: ["admin"] },
    { label: "Profil Asesor", to: "/admin/asesor-profiles", icon: Users, roles: ["admin"] },
    { label: "Users", to: "/admin/users", icon: UserCog, roles: ["admin"] },
  ]},
]

export function navFor(role: Role) {
  return navSections
    .map((s) => ({ ...s, items: s.items.filter((i) => i.roles.includes(role)) }))
    .filter((s) => s.items.length > 0)
}
```

`Sidebar` merender `navFor(session.role)`. **Menu tersembunyi bukan pengamanan** — guard route tetap wajib.

---

## 9. Environment

`.env.local` (FE):

```
VITE_API_BASE_URL=http://localhost:8080
API_BASE_URL=http://localhost:8080        # dipakai server function (SSR)
```

`.env.example` FE menyusul dengan kunci yang sama.

---

## 10. Checklist implementasi

- [ ] Scaffold TanStack Start di `fe/` (bun)
- [ ] Init shadcn + komponen dasar
- [ ] `queryClient` + `query-keys` + wiring ke router context
- [ ] `Session` type + `homeForRole` + server functions (login/logout/getSession) + cookie httpOnly
- [ ] `authMiddleware` + `requireRole` (createMiddleware, server function)
- [ ] Root context (`__root.tsx`) memuat session
- [ ] Routing folder-based: `/admin`, `/asesor`, `/asesi`, `/validator`, `/manager`
- [ ] Layout publik: `/login`, `/register`, `/verify/:token`, `/scan/:token`
- [ ] `AppShell` + `Sidebar` + `navFor(role)` (dipakai tiap layout role)
- [ ] Guard `beforeLoad` — 3 kasus: `requireGuest`, `requireAuth`, `requireRole` + redirect `/` per role
- [ ] Otorisasi resource: `/asesor/assignments/:id` hanya untuk asesor yang ditugaskan (+ backend 403/404)
- [ ] `useNotifications` (SSE) + `eventInvalidation`
- [ ] Halaman awal: login, home per role, `/admin/tuks`, `/admin/asesor-profiles`
- [ ] `.env.example` FE

## 11. Yang perlu diputuskan menyusul

- Apakah FE memakai cookie session untuk **semua** request (same-origin proxy) atau header `Authorization` + `?token=` untuk SSE.
- Perlu `manager_sertifikasi` punya halaman approval khusus? (memengaruhi nav section).
- Strategi SSR vs CSR untuk halaman terproteksi (Start bisa SSR; guard tetap dijalankan di server via middleware).
