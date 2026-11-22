import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useMemo, useState } from 'react'
import { toast } from 'sonner'
import { Check, ChevronLeft, ChevronRight, MapPin, ShieldCheck, Upload, User } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Stepper } from '@/components/ui/stepper'
import { requireGuest } from '@/lib/auth/guards'
import { cn } from '@/lib/utils'
import { AkunStep } from './-components/akun-step'
import { AlamatDomisiliStep } from './-components/alamat-domisili-step'
import { DataPribadiStep } from './-components/data-pribadi-step'
import { DokumenStep } from './-components/dokumen-step'
import { useRegisterAsesi } from './-hooks'
import { lookupPostalCode, validateAll, validateStep } from './-schemas'
import { useRegisterDraft } from './-store'
import type { CreateAsesiProfiles } from './-types'

export const Route = createFileRoute('/register/')({
  beforeLoad: ({ context }) => requireGuest(context),
  component: RegisterPage,
})

const STEPS = [
  { title: 'Data Pribadi', description: 'Identitas diri sesuai dokumen resmi.', icon: User, accent: 'from-sky-500 to-blue-600' },
  { title: 'Alamat & Domisili', description: 'Alamat KTP dan tempat tinggal saat ini.', icon: MapPin, accent: 'from-emerald-500 to-teal-600' },
  { title: 'Dokumen Pendukung', description: 'Unggah berkas pendukung pendaftaran.', icon: Upload, accent: 'from-amber-500 to-orange-600' },
  { title: 'Email & Keamanan Akun', description: 'Kredensial untuk masuk ke akun Anda.', icon: ShieldCheck, accent: 'from-violet-500 to-purple-600' },
] as const

function RegisterPage() {
  const step = useRegisterDraft((state) => state.step)
  const setStep = useRegisterDraft((state) => state.setStep)
  const values = useRegisterDraft((state) => state.values)
  const storeUpdate = useRegisterDraft((state) => state.update)
  const patch = useRegisterDraft((state) => state.patch)
  const [errors, setErrors] = useState<Record<string, string>>({})
  const [postalLoading, setPostalLoading] = useState(false)
  const registerMutation = useRegisterAsesi()

  const today = useMemo(() => new Date().toISOString().slice(0, 10), [])
  const current = STEPS[step]

  // Hidrasi draft dari localStorage setelah render pertama di client.
  useEffect(() => {
    void useRegisterDraft.persist.rehydrate()
  }, [])

  useEffect(() => {
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }, [step])

  function update<K extends keyof CreateAsesiProfiles>(key: K, value: CreateAsesiProfiles[K]) {
    storeUpdate(key, value)
    setErrors((prev) => {
      if (!(key in prev)) return prev
      const next = { ...prev }
      delete next[key as string]
      return next
    })
  }

  function goNext() {
    const { ok, errors: stepErrors } = validateStep(step, values)
    if (!ok) {
      setErrors(stepErrors)
      toast.error('Lengkapi data yang belum valid')
      return
    }
    setErrors({})
    setStep(Math.min(step + 1, STEPS.length - 1))
  }

  function goBack() {
    setErrors({})
    setStep(Math.max(step - 1, 0))
  }

  async function handleLookupPostal() {
    const code = values.postal_code.trim()
    if (!/^\d{5}$/.test(code)) {
      setErrors((prev) => ({ ...prev, postal_code: 'Kode pos harus 5 digit angka' }))
      return
    }
    setPostalLoading(true)
    try {
      const result = await lookupPostalCode(code)
      if (!result) {
        toast.error('Kode pos tidak ditemukan', { description: 'Pastikan 5 digit kode pos sudah benar.' })
        return
      }
      patch(result)
      setErrors((prev) => {
        const next = { ...prev }
        delete next.postal_code
        delete next.province
        delete next.city
        delete next.district
        delete next.village
        return next
      })
      toast.success('Wilayah terisi otomatis', {
        description: `${result.village}, ${result.district}, ${result.city}, ${result.province}`,
      })
    } catch {
      toast.error('Gagal mencari kode pos', { description: 'Periksa koneksi lalu coba lagi.' })
    } finally {
      setPostalLoading(false)
    }
  }

  function handleFinish() {
    if (!validateAll(values).ok) {
      for (let i = 0; i < STEPS.length; i++) {
        const result = validateStep(i, values)
        if (!result.ok) {
          setStep(i)
          setErrors(result.errors)
          break
        }
      }
      toast.error('Masih ada data yang belum valid')
      return
    }
    registerMutation.mutate(values)
  }

  function renderStep() {
    switch (step) {
      case 0:
        return <DataPribadiStep values={values} errors={errors} onChange={update} today={today} />
      case 1:
        return <AlamatDomisiliStep values={values} errors={errors} onChange={update} postalLoading={postalLoading} onLookupPostal={handleLookupPostal} />
      case 2:
        return <DokumenStep values={values} errors={errors} onChange={update} />
      default:
        return <AkunStep values={values} errors={errors} onChange={update} />
    }
  }

  const CurrentIcon = current.icon

  return (
    <div className="min-h-screen bg-gradient-to-b from-slate-50 via-white to-slate-50">
      <header className="border-b border-slate-200/70 bg-white/80 backdrop-blur">
        <div className="mx-auto flex max-w-3xl items-center gap-3 px-4 py-4 sm:px-6">
          <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-slate-900 text-white">
            <ShieldCheck className="h-5 w-5" />
          </div>
          <div className="min-w-0">
            <p className="truncate text-sm font-semibold text-slate-900">Pendaftaran Asesi</p>
            <p className="truncate text-xs text-slate-500">LSP Management System</p>
          </div>
          <span className="ml-auto hidden shrink-0 rounded-full bg-slate-100 px-3 py-1 text-xs font-medium text-slate-600 sm:inline">
            Langkah {step + 1} / {STEPS.length}
          </span>
        </div>
      </header>

      <main className="mx-auto w-full max-w-3xl px-3 py-6 sm:px-6 sm:py-10">
        <Stepper steps={STEPS} current={step} />

        <Card className="mt-6 overflow-hidden">
          <div className={cn('h-1.5 w-full bg-gradient-to-r', current.accent)} />
          <CardContent className="p-4 sm:p-6">
            <div className="mb-5 flex items-start gap-3">
              <span className={cn('flex h-11 w-11 shrink-0 items-center justify-center rounded-xl bg-gradient-to-br text-white shadow-sm', current.accent)}>
                <CurrentIcon className="h-5 w-5" />
              </span>
              <div className="min-w-0">
                <h1 className="text-lg font-semibold text-slate-900 sm:text-xl">{current.title}</h1>
                <p className="text-sm text-slate-500">{current.description}</p>
              </div>
            </div>

            {renderStep()}
          </CardContent>
        </Card>

        <div className="mt-4 flex items-center justify-between gap-3">
          <Button type="button" variant="outline" onClick={goBack} disabled={step === 0} className="gap-1.5">
            <ChevronLeft className="h-4 w-4" />
            Sebelumnya
          </Button>
          {step < STEPS.length - 1 ? (
            <Button type="button" onClick={goNext} className="gap-1.5">
              Lanjut
              <ChevronRight className="h-4 w-4" />
            </Button>
          ) : (
            <Button type="button" onClick={handleFinish} disabled={registerMutation.isPending} className="gap-1.5">
              <Check className="h-4 w-4" />
              {registerMutation.isPending ? 'Mendaftar...' : 'Selesai'}
            </Button>
          )}
        </div>
      </main>
    </div>
  )
}
