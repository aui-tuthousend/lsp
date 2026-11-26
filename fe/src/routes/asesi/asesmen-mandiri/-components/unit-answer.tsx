import { useRef } from 'react'
import { Check, Eye, Paperclip, X } from 'lucide-react'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { cn } from '@/lib/utils'
import type { AssessmentUnit, ElementNode } from '../-types'

export const EVIDENCE_ACCEPT = 'image/jpeg,image/png,application/pdf'

/** Satu berkas bukti pada satu unit (tersimpan di server atau masih lokal). */
export type UnitEvidenceItem = {
  /** ID unik untuk key React & penghapusan. */
  id: string
  /** Object key bila berkas sudah tersimpan di server. */
  key?: string
  /** URL presigned untuk pratinjau (bila tersimpan). */
  viewUrl?: string
  /** Berkas lokal yang belum diunggah (diunggah saat Simpan). */
  file?: File
}

/** Nilai jawaban satu unit (dipakai state form). */
export type UnitAnswerValue = {
  isCompetent: boolean | null
  files: UnitEvidenceItem[]
}

export type UnitAnswerProps = {
  unit: AssessmentUnit
  /** Elemen & KUK unit (dari struktur skema). */
  elements: ElementNode[]
  value: UnitAnswerValue
  disabled?: boolean
  error?: string
  onChange: (value: UnitAnswerValue) => void
  onAddFile: (file: File) => void
  onRemoveFile: (id: string) => void
}

function fileNameFromKey(key: string): string {
  const parts = key.split('/')
  return parts[parts.length - 1] || key
}

/** Satu unit kompetensi: kode/title, elemen & KUK, pernyataan kompetensi, bukti. */
export function UnitAnswer({ unit, elements, value, disabled, error, onChange, onAddFile, onRemoveFile }: UnitAnswerProps) {
  const inputRef = useRef<HTMLInputElement>(null)
  const competent = value.isCompetent === true

  function setCompetent(next: boolean) {
    // Saat dinyatakan tidak kompeten, berkas bukti dikosongkan.
    onChange(next ? { isCompetent: true, files: value.files } : { isCompetent: false, files: [] })
  }

  const statementBadge =
    value.isCompetent === true
      ? { label: 'Kompeten', className: 'bg-emerald-100 text-emerald-700' }
      : value.isCompetent === false
        ? { label: 'Tidak Kompeten', className: 'bg-rose-100 text-rose-700' }
        : { label: 'Belum dijawab', className: 'bg-slate-100 text-slate-600' }

  return (
    <Card>
      <CardContent className="space-y-3 p-4">
        <div className="space-y-0.5">
          <p className="text-sm font-semibold text-slate-900">{unit.unit_code}</p>
          <p className="text-sm text-slate-600">{unit.unit_title}</p>
        </div>

        <Accordion type="single" collapsible className="border-t border-slate-100 pt-1">
          <AccordionItem value="elements">
            <AccordionTrigger>Lihat Elemen dan Kriteria Unjuk Kerja</AccordionTrigger>
            <AccordionContent>
              {elements.length === 0 ? (
                <p className="text-sm text-slate-500">Belum ada elemen pada unit ini.</p>
              ) : (
                <div className="space-y-3">
                  {elements.map((element, index) => (
                    <div key={element.id} className="rounded-md border border-slate-200">
                      <div className="border-b border-slate-100 bg-slate-50 px-3 py-2 text-sm font-medium text-slate-900">
                        {element.title || `Elemen ${index + 1}`}
                      </div>
                      {element.criteria.length === 0 ? (
                        <p className="px-3 py-2 text-xs text-slate-400">Belum ada KUK.</p>
                      ) : (
                        <ul className="divide-y divide-slate-100">
                          {element.criteria.map((criteria) => (
                            <li key={criteria.id} className="flex items-start gap-2 px-3 py-2">
                              <span className="text-xs text-slate-400">{criteria.sort_order}.</span>
                              <span className="text-sm text-slate-700">{criteria.content}</span>
                            </li>
                          ))}
                        </ul>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </AccordionContent>
          </AccordionItem>
        </Accordion>

        <div className="space-y-2 border-t border-slate-100 pt-3">
          <div className="flex flex-wrap items-center gap-2">
            <span className="text-sm font-medium text-slate-700">Pernyataan Kompetensi</span>
            <span
              className={cn(
                'inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium',
                statementBadge.className,
              )}
            >
              {statementBadge.label}
            </span>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <Button
              type="button"
              size="sm"
              variant={competent ? 'default' : 'outline'}
              className="gap-1.5"
              disabled={disabled}
              onClick={() => setCompetent(true)}
            >
              <Check className="h-4 w-4" />
              Kompeten
            </Button>
            <Button
              type="button"
              size="sm"
              variant={value.isCompetent === false ? 'default' : 'outline'}
              className="gap-1.5"
              disabled={disabled}
              onClick={() => setCompetent(false)}
            >
              <X className="h-4 w-4" />
              Tidak Kompeten
            </Button>
          </div>
        </div>

        {competent ? (
          <div className="space-y-2 border-t border-slate-100 pt-3">
            <div className="flex flex-wrap items-center gap-2">
              <input
                ref={inputRef}
                type="file"
                accept={EVIDENCE_ACCEPT}
                className="hidden"
                disabled={disabled}
                onChange={(event) => {
                  const file = event.target.files?.[0]
                  if (file) onAddFile(file)
                  event.target.value = ''
                }}
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
                Tambah bukti
              </Button>
              <span className="text-xs text-slate-400">JPG/PNG/PDF, maks 2MB</span>
            </div>

            {value.files.length > 0 ? (
              <ul className="space-y-1.5">
                {value.files.map((item) => (
                  <li
                    key={item.id}
                    className="flex items-center justify-between gap-2 rounded-md border border-slate-200 bg-white px-3 py-2"
                  >
                    <span className="truncate text-xs text-slate-600">
                      {item.file ? item.file.name : fileNameFromKey(item.key ?? '')}
                    </span>
                    <span className="flex shrink-0 items-center gap-1.5">
                      {item.viewUrl ? (
                        <a
                          href={item.viewUrl}
                          target="_blank"
                          rel="noreferrer"
                          aria-disabled={disabled}
                          className={cn(
                            'inline-flex h-7 items-center gap-1 rounded-md border border-slate-200 px-2 text-xs font-medium transition-colors hover:bg-slate-50',
                            disabled && 'pointer-events-none opacity-50',
                          )}
                        >
                          <Eye className="h-3.5 w-3.5" />
                          Lihat
                        </a>
                      ) : item.file ? (
                        <span className="text-xs text-slate-400">Belum diunggah</span>
                      ) : null}
                      <button
                        type="button"
                        aria-label="Hapus berkas"
                        disabled={disabled}
                        className="text-slate-400 transition-colors hover:text-rose-600 disabled:opacity-50 disabled:hover:text-slate-400"
                        onClick={() => onRemoveFile(item.id)}
                      >
                        <X className="h-3.5 w-3.5" />
                      </button>
                    </span>
                  </li>
                ))}
              </ul>
            ) : (
              <p className="text-xs text-slate-400">Belum ada bukti diunggah.</p>
            )}
          </div>
        ) : null}

        {error ? <p className="text-xs font-medium text-rose-600">{error}</p> : null}
      </CardContent>
    </Card>
  )
}
