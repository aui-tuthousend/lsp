import { Eye, FileText, Loader2 } from 'lucide-react'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Card, CardContent } from '@/components/ui/card'
import type { SchemaDocumentRequirement } from '../../admin/skema/-types'
import type { ApplicationDocument } from '../../asesi/permohonan/-types'

type DocumentsAccordionProps = {
  documents: ApplicationDocument[]
  requirements: SchemaDocumentRequirement[]
  loading?: boolean
}

/**
 * Pratinjau (read-only) dokumen bukti kelengkapan APL-01 untuk admin LSP. Tiap
 * dokumen ditampilkan sebagai baris dengan tautan pratinjau (presigned URL).
 */
export function DocumentsAccordion({ documents, requirements, loading }: DocumentsAccordionProps) {
  const { hasFile, viewUrl } = buildFileIndex(documents)
  const uploaded = requirements.filter((requirement) => hasFile.has(requirement.document_type_id))

  return (
    <Card>
      <CardContent className="p-0">
        <Accordion type="single" collapsible defaultValue="documents" className="px-4">
          <AccordionItem value="documents">
            <AccordionTrigger>
              <span className="flex flex-wrap items-baseline gap-2">
                <span className="flex items-center gap-1.5">
                  <FileText className="h-4 w-4 text-slate-500" />
                  Dokumen APL-01
                </span>
                <span className="text-xs font-normal text-slate-400">{uploaded.length} berkas terunggah</span>
              </span>
            </AccordionTrigger>
            <AccordionContent>
              {loading ? (
                <div className="flex items-center gap-2 text-sm text-slate-500">
                  <Loader2 className="h-4 w-4 animate-spin" />
                  Memuat dokumen...
                </div>
              ) : uploaded.length === 0 ? (
                <p className="text-sm text-slate-500">Belum ada dokumen yang diunggah pada APL-01.</p>
              ) : (
                <div className="grid grid-cols-1 gap-3 sm:grid-cols-2">
                  {uploaded.map((requirement) => (
                    <DocumentRow
                      key={requirement.document_type_id}
                      label={requirement.document_type.label}
                      url={viewUrl[requirement.document_type_id]}
                    />
                  ))}
                </div>
              )}
            </AccordionContent>
          </AccordionItem>
        </Accordion>
      </CardContent>
    </Card>
  )
}

function buildFileIndex(documents: ApplicationDocument[]) {
  const hasFile = new Set<string>()
  const viewUrl: Record<string, string> = {}
  for (const document of documents) {
    if (!document.file_url) continue
    hasFile.add(document.document_type_id)
    if (document.file_view_url) viewUrl[document.document_type_id] = document.file_view_url
  }
  return { hasFile, viewUrl }
}

/** Baris dokumen read-only: label + tombol pratinjau. */
function DocumentRow({ label, url }: { label: string; url?: string }) {
  return (
    <div className="flex items-center justify-between gap-3 rounded-lg border border-slate-200 bg-white px-3 py-3">
      <div className="min-w-0">
        <p className="flex items-center gap-1.5 text-sm font-medium text-slate-700">
          <FileText className="h-4 w-4 text-slate-400" />
          {label}
        </p>
        <p className="mt-0.5 truncate text-xs text-slate-500">Sudah diunggah</p>
      </div>
      {url ? (
        <a
          href={url}
          target="_blank"
          rel="noreferrer"
          className="inline-flex h-9 shrink-0 items-center justify-center gap-1.5 rounded-md border border-slate-200 bg-white px-3 text-sm font-medium transition-colors hover:bg-slate-50"
        >
          <Eye className="h-4 w-4" />
          Lihat
        </a>
      ) : null}
    </div>
  )
}
