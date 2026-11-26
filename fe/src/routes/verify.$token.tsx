import { createFileRoute, Link } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { CheckCircle2, Loader2, ShieldCheck, XCircle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { ApiError } from '@/lib/api'
import { apl01Service } from './asesi/permohonan/-service'
import { apl02Service } from './asesi/asesmen-mandiri/-service'
import type { SignatureStatus } from './asesi/permohonan/-types'

export const Route = createFileRoute('/verify/$token')({
  component: VerifyPage,
})

type VerifyState = 'checking' | 'idle' | 'submitting' | 'done' | 'already' | 'invalid' | 'error'
type VerifyKind = 'apl01' | 'apl02'

/** Keterangan QR ini memverifikasi formulir/step apa. */
const VERIFY_LABELS: Record<VerifyKind, { title: string; target: string; form: string }> = {
  apl01: {
    title: 'Verifikasi Permohonan',
    target: 'Formulir APL-01 — Langkah 4: Jadwal Asesmen',
    form: 'Formulir APL-01',
  },
  apl02: {
    title: 'Verifikasi Asesmen Mandiri',
    target: 'Formulir APL-02 — Asesmen Mandiri',
    form: 'Formulir APL-02',
  },
}

/**
 * Laman verifikasi QR. Dibuka dari hasil scan QR (`/verify/{token}`). Token bisa
 * milik APL-01 maupun APL-02; jenisnya dideteksi dari endpoint status. Bila token
 * sudah dipakai, halaman hanya menampilkan keterangan tanpa tombol verifikasi.
 */
function VerifyPage() {
  const { token } = Route.useParams()
  const { auth } = Route.useRouteContext()
  const [state, setState] = useState<VerifyState>('checking')
  const [kind, setKind] = useState<VerifyKind>('apl01')
  const [message, setMessage] = useState('')
  const [doneCode, setDoneCode] = useState('')
  const [confirmOpen, setConfirmOpen] = useState(false)

  useEffect(() => {
    if (!auth) return
    let active = true
    setState('checking')

    function applyStatus(status: SignatureStatus, code: string | undefined) {
      if (status === 'dipakai') {
        setMessage(code ? `Sudah diverifikasi (${code}).` : '')
        setState('already')
      } else if (status === 'menunggu') {
        setState('idle')
      } else {
        setState('invalid')
      }
    }

    void (async () => {
      try {
        const status = await apl01Service.verifyStatus(token, auth)
        if (!active) return
        setKind('apl01')
        applyStatus(status.status, status.application_code)
      } catch (error) {
        if (!active) return
        if (!(error instanceof ApiError) || error.status !== 404) {
          setState('error')
          return
        }
        try {
          const status = await apl02Service.verifyStatus(token, auth)
          if (!active) return
          setKind('apl02')
          applyStatus(status.status, status.code)
        } catch (error2) {
          if (!active) return
          setState(error2 instanceof ApiError && error2.status === 404 ? 'invalid' : 'error')
        }
      }
    })()

    return () => {
      active = false
    }
  }, [auth, token])

  async function confirm() {
    setState('submitting')
    try {
      if (kind === 'apl02') {
        const detail = await apl02Service.verify(token, auth)
        setDoneCode(detail.assessment.code)
      } else {
        const detail = await apl01Service.verify(token, auth)
        setDoneCode(detail.application.code)
      }
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
          : 'Gagal memverifikasi. Silakan coba lagi.',
      )
      setState('error')
    }
  }

  const label = VERIFY_LABELS[kind]

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-50 p-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <ShieldCheck className="h-5 w-5" />
            {label.title}
          </CardTitle>
          <p className="text-sm text-slate-500">{label.target}</p>
        </CardHeader>
        <CardContent className="space-y-4">
          {!auth ? (
            <>
              <p className="text-sm text-slate-600">
                Anda perlu masuk sebagai asesi untuk memverifikasi ini.
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
              {state === 'done' && doneCode ? (
                <p className="text-sm text-slate-500">{doneCode} telah diverifikasi.</p>
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
              <p className="text-sm text-rose-600">{message || 'Gagal memverifikasi. Silakan coba lagi.'}</p>
            </div>
          ) : (
            <>
              <p className="text-sm text-slate-600">
                QR ini memverifikasi seluruh data pengajuan{' '}
                <span className="font-medium text-slate-700">{label.form}</span> Anda. Tekan tombol di bawah untuk
                menyetujui dan memverifikasi.
              </p>
              <Button className="w-full gap-2" onClick={() => setConfirmOpen(true)}>
                <CheckCircle2 className="h-4 w-4" />
                Setuju &amp; Verifikasi
              </Button>
              <Dialog open={confirmOpen} onOpenChange={setConfirmOpen}>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Konfirmasi Verifikasi</DialogTitle>
                    <DialogDescription>
                      Apakah Anda sudah memastikan data yang Anda masukkan sudah benar?
                    </DialogDescription>
                  </DialogHeader>
                  <DialogFooter>
                    <Button variant="outline" onClick={() => setConfirmOpen(false)}>
                      Batal
                    </Button>
                    <Button
                      className="gap-2"
                      onClick={() => {
                        setConfirmOpen(false)
                        void confirm()
                      }}
                    >
                      <CheckCircle2 className="h-4 w-4" />
                      Ya, Verifikasi
                    </Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
            </>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
