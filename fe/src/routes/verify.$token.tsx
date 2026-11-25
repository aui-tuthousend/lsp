import { createFileRoute, Link } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { CheckCircle2, Loader2, ShieldCheck, XCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { ApiError } from '@/lib/api'
import { apl01Service } from './asesi/permohonan/-service'
import type { ApplicationDetail } from './asesi/permohonan/-types'

export const Route = createFileRoute('/verify/$token')({
  component: VerifyPage,
})

type VerifyState = 'checking' | 'idle' | 'submitting' | 'done' | 'already' | 'invalid' | 'error'

/**
 * Laman verifikasi QR APL-01. Dibuka dari hasil scan QR (`/verify/{token}`).
 * Saat dibuka, status token dicek lebih dulu: bila sudah dipakai, halaman hanya
 * menampilkan keterangan "sudah terverifikasi" tanpa tombol verifikasi.
 */
function VerifyPage() {
  const { token } = Route.useParams()
  const { auth } = Route.useRouteContext()
  const [state, setState] = useState<VerifyState>('checking')
  const [message, setMessage] = useState('')
  const [result, setResult] = useState<ApplicationDetail | null>(null)

  useEffect(() => {
    if (!auth) return
    let active = true
    setState('checking')
    apl01Service
      .verifyStatus(token, auth)
      .then((status) => {
        if (!active) return
        if (status.status === 'dipakai') {
          setMessage(status.application_code ? `Permohonan ${status.application_code} sudah diverifikasi.` : '')
          setState('already')
        } else if (status.status === 'menunggu') {
          setState('idle')
        } else {
          setState('invalid')
        }
      })
      .catch((error) => {
        if (!active) return
        setState(error instanceof ApiError && error.status === 404 ? 'invalid' : 'error')
      })
    return () => {
      active = false
    }
  }, [auth, token])

  async function confirm() {
    setState('submitting')
    try {
      const detail = await apl01Service.verify(token, auth)
      setResult(detail)
      setState('done')
    } catch (error) {
      if (error instanceof ApiError && error.status === 409) {
        setMessage('')
        setState('already')
        return
      }
      setMessage(
        error instanceof ApiError && error.status === 400
          ? 'Token tidak valid atau sudah kedaluwarsa.'
          : 'Gagal memverifikasi permohonan. Silakan coba lagi.',
      )
      setState('error')
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-50 p-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <ShieldCheck className="h-5 w-5" />
            Verifikasi Permohonan
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {!auth ? (
            <>
              <p className="text-sm text-slate-600">
                Anda perlu masuk sebagai asesi untuk memverifikasi permohonan ini.
              </p>
              <Link
                to="/login"
                className="inline-flex h-10 w-full items-center justify-center rounded-md bg-slate-900 text-sm font-medium text-white transition-colors hover:bg-slate-800"
              >
                Masuk
              </Link>
            </>
          ) : state === 'checking' || state === 'submitting' ? (
            <div className="flex items-center justify-center gap-2 py-6 text-sm text-slate-500">
              <Loader2 className="h-4 w-4 animate-spin" />
              {state === 'submitting' ? 'Memverifikasi...' : 'Memeriksa status token...'}
            </div>
          ) : state === 'done' || state === 'already' ? (
            <div className="space-y-2 text-center">
              <CheckCircle2 className="mx-auto h-10 w-10 text-emerald-500" />
              <p className="text-sm font-medium text-slate-900">
                {state === 'done' ? 'Anda telah menyetujui' : 'Token sudah terverifikasi'}
              </p>
              {state === 'done' && result ? (
                <p className="text-sm text-slate-500">Permohonan {result.application.code} telah diverifikasi.</p>
              ) : null}
              {state === 'already' && message ? <p className="text-sm text-slate-500">{message}</p> : null}
              <p className="text-xs text-slate-400">Anda bisa menutup halaman ini.</p>
            </div>
          ) : state === 'invalid' ? (
            <div className="space-y-2 text-center">
              <XCircle className="mx-auto h-10 w-10 text-rose-500" />
              <p className="text-sm text-rose-600">Token tidak valid atau sudah kedaluwarsa.</p>
            </div>
          ) : state === 'error' ? (
            <div className="space-y-2 text-center">
              <XCircle className="mx-auto h-10 w-10 text-rose-500" />
              <p className="text-sm text-rose-600">{message}</p>
            </div>
          ) : (
            <>
              <p className="text-sm text-slate-600">
                Tekan tombol di bawah untuk menyetujui dan memverifikasi permohonan APL-01 Anda.
              </p>
              <Button className="w-full gap-2" onClick={() => void confirm()}>
                <CheckCircle2 className="h-4 w-4" />
                Setuju &amp; Verifikasi
              </Button>
            </>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
