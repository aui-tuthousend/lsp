import { useId, useState } from 'react'
import { Loader2, Search } from 'lucide-react'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'

/** Field alamat yang dikelola blok ini. */
export type AddressFieldKey = 'postal_code' | 'province' | 'city' | 'district' | 'village'

/** Nilai alamat yang dibutuhkan; `village` opsional (tidak semua tabel memilikinya). */
export type AddressValues = {
  postal_code: string
  province: string
  city: string
  district: string
  village?: string
}

export type AddressFieldsProps = {
  values: AddressValues
  errors?: Partial<Record<AddressFieldKey, string>>
  onChange: (key: AddressFieldKey, value: string) => void
  /** Tampilkan field Desa/Kelurahan (default true). */
  withVillage?: boolean
  /**
   * Kunci provinsi/kota/kecamatan (read-only) agar konsisten dengan hasil lookup
   * kode pos; hanya desa yang bisa diubah. Tetap ada tombol "Isi manual" sebagai
   * jalan keluar bila kode pos tidak ditemukan. Default false.
   */
  lockRegion?: boolean
  className?: string
}

const POSTAL_RE = /^\d{5}$/
const KODEPOS_ENDPOINT = 'https://www.emsifa.com/api-wilayah-indonesia/v2/postal-codes'

type EmsifaVillage = {
  id: string
  name: string
  postal_code: string
  province?: { id: string; name: string }
  regency?: { id: string; name: string }
  district?: { id: string; name: string }
}

type PostalLookupResult = {
  province: string
  city: string
  district: string
  villages: string[]
}

/**
 * Cari provinsi/kota/kecamatan/desa dari kode pos (5 digit).
 * Sumber: API statis wilayah Indonesia (emsifa v2). Satu kode pos bisa dipakai
 * banyak desa/kelurahan, sehingga desa dikembalikan sebagai daftar.
 */
async function lookupPostalCode(postalCode: string): Promise<PostalLookupResult | null> {
  const code = postalCode.trim()
  if (!POSTAL_RE.test(code)) return null

  const res = await fetch(`${KODEPOS_ENDPOINT}/${code}.json`)
  if (res.status === 404) return null
  if (!res.ok) throw new Error(`Pencarian kode pos gagal (${res.status})`)

  const json = (await res.json()) as { data?: EmsifaVillage[] }
  const items = json.data ?? []
  if (items.length === 0) return null

  const first = items[0]
  const villages = Array.from(new Set(items.map((item) => item.name).filter(Boolean)))

  return {
    province: first.province?.name ?? '',
    city: first.regency?.name ?? '',
    district: first.district?.name ?? '',
    villages,
  }
}

function errorClass(hasError?: boolean) {
  return hasError ? 'border-rose-300 focus:border-rose-400 focus:ring-rose-100' : undefined
}

/**
 * Blok alamat: isi kode pos lalu tekan "Cari wilayah" — provinsi/kota/kecamatan/
 * desa terisi otomatis. Bila `lockRegion` aktif, provinsi/kota/kecamatan
 * dikunci (read-only) dan hanya desa yang bisa diubah, dengan tombol
 * "Isi manual" sebagai fallback. Dipakai di form pendaftaran, profil asesi,
 * permohonan APL-01, dan form TUK.
 */
export function AddressFields({
  values,
  errors,
  onChange,
  withVillage = true,
  lockRegion = false,
  className,
}: AddressFieldsProps) {
  const [loading, setLoading] = useState(false)
  const [manual, setManual] = useState(false)
  const [villageOptions, setVillageOptions] = useState<string[]>([])
  const villageListId = useId()

  const regionLocked = lockRegion && !manual

  async function lookup() {
    const code = values.postal_code.trim()
    if (!POSTAL_RE.test(code)) {
      toast.error('Kode pos harus 5 digit angka')
      return
    }
    setLoading(true)
    try {
      const result = await lookupPostalCode(code)
      if (!result) {
        toast.error('Kode pos tidak ditemukan', {
          description: 'Pastikan 5 digit kode pos sudah benar, atau isi wilayah secara manual.',
        })
        return
      }
      onChange('province', result.province)
      onChange('city', result.city)
      onChange('district', result.district)
      if (withVillage) {
        setVillageOptions(result.villages)
        onChange('village', result.villages.length === 1 ? result.villages[0] : '')
      }
      const region = [result.district, result.city, result.province].filter(Boolean).join(', ')
      toast.success('Wilayah terisi otomatis', {
        description:
          withVillage && result.villages.length > 1
            ? `${region} — pilih desa/kelurahan`
            : [withVillage ? result.villages[0] : '', region].filter(Boolean).join(', '),
      })
    } catch {
      toast.error('Gagal mencari kode pos', { description: 'Periksa koneksi lalu coba lagi.' })
    } finally {
      setLoading(false)
    }
  }

  const regionHint = regionLocked ? 'Terisi otomatis dari kode pos' : 'Terisi otomatis, bisa diubah'
  const lockedClass = regionLocked ? 'cursor-not-allowed bg-slate-100 text-slate-500' : undefined

  return (
    <div className={cn('rounded-xl border border-dashed border-slate-300 bg-slate-50/70 p-3 sm:p-4', className)}>
      <div className="mb-3 flex items-center justify-between gap-2">
        <div className="flex items-center gap-2 text-sm font-medium text-slate-700">
          <Search className="h-4 w-4 text-slate-400" />
          Isi kode pos, wilayah terisi otomatis
        </div>
        {lockRegion ? (
          <button
            type="button"
            onClick={() => setManual((value) => !value)}
            className="text-xs font-medium text-slate-500 underline underline-offset-2 hover:text-slate-700"
          >
            {manual ? 'Kunci otomatis' : 'Isi manual'}
          </button>
        ) : null}
      </div>
      <div className="flex flex-col gap-2 sm:flex-row">
        <Input
          inputMode="numeric"
          maxLength={5}
          value={values.postal_code}
          onChange={(event) => onChange('postal_code', event.target.value.replace(/\D/g, ''))}
          placeholder="mis. 40115"
          className={cn('sm:max-w-40', errorClass(Boolean(errors?.postal_code)))}
        />
        <Button type="button" onClick={() => void lookup()} disabled={loading} className="gap-2 sm:w-auto">
          {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <Search className="h-4 w-4" />}
          {loading ? 'Mencari...' : 'Cari wilayah'}
        </Button>
      </div>
      {errors?.postal_code ? <p className="mt-1.5 text-xs font-medium text-rose-600">{errors.postal_code}</p> : null}

      <div className="mt-3 grid grid-cols-1 gap-3 sm:grid-cols-2">
        <Field label="Provinsi" error={errors?.province} hint={regionHint}>
          <Input
            value={values.province}
            onChange={(event) => onChange('province', event.target.value)}
            readOnly={regionLocked}
            placeholder="Terisi otomatis"
            className={cn(errorClass(Boolean(errors?.province)), lockedClass)}
          />
        </Field>
        <Field label="Kota / Kabupaten" error={errors?.city} hint={regionHint}>
          <Input
            value={values.city}
            onChange={(event) => onChange('city', event.target.value)}
            readOnly={regionLocked}
            placeholder="Terisi otomatis"
            className={cn(errorClass(Boolean(errors?.city)), lockedClass)}
          />
        </Field>
        <Field label="Kecamatan" error={errors?.district} hint={regionHint}>
          <Input
            value={values.district}
            onChange={(event) => onChange('district', event.target.value)}
            readOnly={regionLocked}
            placeholder="Terisi otomatis"
            className={cn(errorClass(Boolean(errors?.district)), lockedClass)}
          />
        </Field>
        {withVillage ? (
          <Field
            label="Desa / Kelurahan"
            error={errors?.village}
            hint={villageOptions.length > 1 ? 'Pilih dari hasil pencarian atau ketik manual' : 'Terisi otomatis, bisa diubah'}
          >
            <Input
              value={values.village ?? ''}
              onChange={(event) => onChange('village', event.target.value)}
              list={villageOptions.length > 0 ? villageListId : undefined}
              placeholder="Terisi otomatis"
              className={errorClass(Boolean(errors?.village))}
            />
            {villageOptions.length > 0 ? (
              <datalist id={villageListId}>
                {villageOptions.map((village) => (
                  <option key={village} value={village} />
                ))}
              </datalist>
            ) : null}
          </Field>
        ) : null}
      </div>
    </div>
  )
}
