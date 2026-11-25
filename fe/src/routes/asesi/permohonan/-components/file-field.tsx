import { useRef } from 'react'
import { Paperclip, X } from 'lucide-react'
import { Button } from '@/components/ui/button'

export type FileFieldProps = {
  id: string
  label: string
  accept?: string
  file: File | null
  /** True bila sudah ada berkas tersimpan di server (object key). */
  hasExisting?: boolean
  required?: boolean
  error?: string
  /** Nonaktifkan kontrol saat proses unggah/simpan berjalan. */
  disabled?: boolean
  onPick: (file: File | null) => void
}

/** Kontrol unggah berkas sederhana (native input, dipicu oleh tombol). */
export function FileField({ id, label, accept, file, hasExisting, required, error, disabled, onPick }: FileFieldProps) {
  const inputRef = useRef<HTMLInputElement>(null)

  function clear() {
    onPick(null)
    if (inputRef.current) inputRef.current.value = ''
  }

  return (
    <div className="space-y-1.5">
      <label htmlFor={id} className="flex items-center gap-1 text-sm font-medium text-slate-700">
        {label}
        {required ? <span className="text-rose-500">*</span> : null}
      </label>
      <div className="flex flex-wrap items-center gap-2">
        <input
          ref={inputRef}
          id={id}
          type="file"
          accept={accept}
          className="hidden"
          disabled={disabled}
          onChange={(event) => onPick(event.target.files?.[0] ?? null)}
        />
        <Button
          type="button"
          variant="outline"
          size="sm"
          className="gap-1.5"
          disabled={disabled}
          onClick={() => inputRef.current?.click()}
        >
          <Paperclip className="h-4 w-4" />
          {file ? 'Ganti berkas' : 'Pilih berkas'}
        </Button>
        {file ? (
          <span className="flex items-center gap-1 text-xs text-slate-600">
            {file.name}
            <button
              type="button"
              aria-label="Batalkan berkas"
              disabled={disabled}
              className="text-slate-400 transition-colors hover:text-rose-600 disabled:opacity-50 disabled:hover:text-slate-400"
              onClick={clear}
            >
              <X className="h-3.5 w-3.5" />
            </button>
          </span>
        ) : hasExisting ? (
          <span className="text-xs text-emerald-600">Sudah diunggah</span>
        ) : (
          <span className="text-xs text-slate-400">Belum ada berkas</span>
        )}
      </div>
      {error ? <p className="text-xs font-medium text-rose-600">{error}</p> : null}
    </div>
  )
}
