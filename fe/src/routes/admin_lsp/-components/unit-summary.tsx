import { Check, Eye, X } from 'lucide-react'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Card, CardContent } from '@/components/ui/card'
import { cn } from '@/lib/utils'
import type { AssessmentUnit, ElementNode } from '../../asesi/asesmen-mandiri/-types'

function fileNameFromKey(key: string): string {
  const parts = key.split('/')
  return parts[parts.length - 1] || key
}

/**
 * Ringkasan read-only satu unit kompetensi APL-02 untuk admin LSP: kode/title,
 * elemen & KUK, pernyataan kompetensi asesi, dan berkas bukti.
 */
export function UnitSummary({ unit, elements }: { unit: AssessmentUnit; elements: ElementNode[] }) {
  const badge =
    unit.is_competent === true
      ? { label: 'Kompeten', className: 'bg-emerald-100 text-emerald-700', icon: Check }
      : unit.is_competent === false
        ? { label: 'Tidak Kompeten', className: 'bg-rose-100 text-rose-700', icon: X }
        : { label: 'Belum dijawab', className: 'bg-slate-100 text-slate-600', icon: null }
  const BadgeIcon = badge.icon

  return (
    <Card>
      <CardContent className="space-y-3 p-4">
        <div className="flex flex-wrap items-start justify-between gap-2">
          <div className="space-y-0.5">
            <p className="text-sm font-semibold text-slate-900">{unit.unit_code}</p>
            <p className="text-sm text-slate-600">{unit.unit_title}</p>
          </div>
          <span
            className={cn(
              'inline-flex shrink-0 items-center gap-1 rounded-full px-2.5 py-0.5 text-xs font-medium',
              badge.className,
            )}
          >
            {BadgeIcon ? <BadgeIcon className="h-3.5 w-3.5" /> : null}
            {badge.label}
          </span>
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
          <p className="text-xs font-medium text-slate-500">Bukti yang diunggah</p>
          {unit.files.length === 0 ? (
            <p className="text-xs text-slate-400">Belum ada bukti diunggah.</p>
          ) : (
            <ul className="space-y-1.5">
              {unit.files.map((file) => (
                <li
                  key={file.id}
                  className="flex items-center justify-between gap-2 rounded-md border border-slate-200 bg-white px-3 py-2"
                >
                  <span className="truncate text-xs text-slate-600">{fileNameFromKey(file.file_url)}</span>
                  {file.file_view_url ? (
                    <a
                      href={file.file_view_url}
                      target="_blank"
                      rel="noreferrer"
                      className="inline-flex h-7 shrink-0 items-center gap-1 rounded-md border border-slate-200 px-2 text-xs font-medium transition-colors hover:bg-slate-50"
                    >
                      <Eye className="h-3.5 w-3.5" />
                      Lihat
                    </a>
                  ) : null}
                </li>
              ))}
            </ul>
          )}
        </div>
      </CardContent>
    </Card>
  )
}
