import { useEffect, useState } from 'react'
import { QRCodeSVG } from 'qrcode.react'
import { Copy, ExternalLink, QrCode } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'

/**
 * Panel QR verifikasi APL-01. QR berisi tautan laman verifikasi (`/verify/{token}`)
 * yang dibuka di perangkat/laman lain; token tidak dikirim ke layanan pihak ketiga.
 */
export function QrPanel({ token }: { token: string }) {
  const [verifyUrl, setVerifyUrl] = useState('')

  useEffect(() => {
    setVerifyUrl(`${window.location.origin}/verify/${token}`)
  }, [token])

  async function copyLink() {
    try {
      await navigator.clipboard.writeText(verifyUrl)
      toast.success('Tautan verifikasi disalin')
    } catch {
      toast.error('Gagal menyalin tautan')
    }
  }

  return (
    <div className="flex flex-col items-center gap-4 text-center">
      <div className="rounded-2xl border border-slate-200 bg-white p-4">
        {verifyUrl ? (
          <QRCodeSVG value={verifyUrl} size={200} level="M" />
        ) : (
          <div className="flex h-[200px] w-[200px] items-center justify-center text-slate-300">
            <QrCode className="h-10 w-10" />
          </div>
        )}
      </div>
      <p className="max-w-sm text-sm text-slate-500">
        Scan QR ini dengan kamera perangkat lain, lalu tekan <span className="font-medium text-slate-700">Setuju</span> untuk
        memverifikasi permohonan. Halaman ini akan otomatis diperbarui.
      </p>
      <div className="flex flex-wrap items-center justify-center gap-2">
        <Button
          variant="outline"
          className="gap-2"
          disabled={!verifyUrl}
          onClick={() => window.open(verifyUrl, '_blank', 'noopener')}
        >
          <ExternalLink className="h-4 w-4" />
          Buka di laman baru
        </Button>
        <Button variant="ghost" className="gap-2" disabled={!verifyUrl} onClick={() => void copyLink()}>
          <Copy className="h-4 w-4" />
          Salin tautan
        </Button>
      </div>
    </div>
  )
}
