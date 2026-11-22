// Base URL API untuk request yang dijalankan di browser (fetch data & SSE).
//
// Prioritas:
// 1. VITE_API_BASE_URL bila di-set.
// 2. Turunkan dari host yang sedang dibuka (hostname:8080) supaya saat dev server
//    diakses dari device lain (mis. via IP LAN) request tetap mengarah ke mesin dev.
// 3. Fallback localhost untuk SSR/lingkungan tanpa window.
function resolveApiBaseUrl(): string {
  const fromEnv = import.meta.env.VITE_API_BASE_URL
  if (fromEnv) return fromEnv
  if (typeof window !== 'undefined') {
    return `${window.location.protocol}//${window.location.hostname}:8080`
  }
  return 'http://localhost:8080'
}

export const API_BASE_URL = resolveApiBaseUrl()
