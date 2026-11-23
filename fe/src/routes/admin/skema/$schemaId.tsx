import { useState } from 'react'
import { createFileRoute, Link } from '@tanstack/react-router'
import { ArrowLeft, Loader2, Pencil, Plus, Trash2 } from 'lucide-react'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { ConfirmDialog } from '@/components/ui/confirm-dialog'
import { CriteriaFormDialog } from './-components/criteria-form-dialog'
import { ElementFormDialog } from './-components/element-form-dialog'
import { SchemaFormDialog } from './-components/schema-form-dialog'
import { UnitFormDialog } from './-components/unit-form-dialog'
import { useDeleteCriteria, useDeleteElement, useDeleteUnit, useSchemaTree } from './-hooks'
import type { CompetencyUnit, CompetencyUnitCriteria, CompetencyUnitElement } from './-types'

export const Route = createFileRoute('/admin/skema/$schemaId')({
  component: SchemaDetailPage,
})

function SchemaDetailPage() {
  const { schemaId } = Route.useParams()
  const { auth } = Route.useRouteContext()
  const { data: tree, isLoading, isError } = useSchemaTree(schemaId, auth)

  // Skema
  const [schemaFormOpen, setSchemaFormOpen] = useState(false)

  // Unit
  const [unitForm, setUnitForm] = useState<{ open: boolean; unit: CompetencyUnit | null }>({ open: false, unit: null })
  const [deletingUnit, setDeletingUnit] = useState<CompetencyUnit | null>(null)
  const removeUnit = useDeleteUnit(schemaId, auth)

  // Elemen
  const [elementForm, setElementForm] = useState<{ open: boolean; unitId: string; element: CompetencyUnitElement | null }>({
    open: false,
    unitId: '',
    element: null,
  })
  const [deletingElement, setDeletingElement] = useState<CompetencyUnitElement | null>(null)
  const removeElement = useDeleteElement(schemaId, auth)

  // KUK
  const [criteriaForm, setCriteriaForm] = useState<{ open: boolean; elementId: string; criteria: CompetencyUnitCriteria | null }>({
    open: false,
    elementId: '',
    criteria: null,
  })
  const [deletingCriteria, setDeletingCriteria] = useState<CompetencyUnitCriteria | null>(null)
  const removeCriteria = useDeleteCriteria(schemaId, auth)

  if (isLoading) {
    return (
      <div className="flex items-center gap-2 p-6 text-sm text-slate-500">
        <Loader2 className="h-4 w-4 animate-spin" />
        Memuat...
      </div>
    )
  }

  if (isError || !tree) {
    return (
      <div className="space-y-4">
        <Link to="/admin/skema" className="inline-flex items-center gap-1.5 text-sm text-slate-600 hover:text-slate-900">
          <ArrowLeft className="h-4 w-4" />
          Kembali ke daftar skema
        </Link>
        <p className="text-sm text-rose-600">Gagal memuat skema sertifikasi.</p>
      </div>
    )
  }

  const units = tree.units ?? []

  return (
    <div className="space-y-4">
      <Link to="/admin/skema" className="inline-flex items-center gap-1.5 text-sm text-slate-600 hover:text-slate-900">
        <ArrowLeft className="h-4 w-4" />
        Kembali ke daftar skema
      </Link>

      <div className="flex flex-wrap items-start justify-between gap-3">
        <div className="space-y-1">
          <div className="flex items-center gap-2">
            <h1 className="text-lg font-semibold text-slate-900">{tree.name}</h1>
            <span
              className={
                tree.is_active
                  ? 'rounded-full bg-emerald-50 px-2 py-0.5 text-xs font-medium text-emerald-700'
                  : 'rounded-full bg-slate-100 px-2 py-0.5 text-xs font-medium text-slate-500'
              }
            >
              {tree.is_active ? 'Aktif' : 'Tidak aktif'}
            </span>
          </div>
          <p className="text-sm text-slate-500">
            <code className="rounded bg-slate-100 px-1.5 py-0.5 text-xs font-medium text-slate-700">{tree.code}</code>
            {tree.schema_code ? <span className="ml-2">Nomor: {tree.schema_code}</span> : null}
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" className="gap-2" onClick={() => setSchemaFormOpen(true)}>
            <Pencil className="h-4 w-4" />
            Ubah Skema
          </Button>
          <Button className="gap-2" onClick={() => setUnitForm({ open: true, unit: null })}>
            <Plus className="h-4 w-4" />
            Tambah Unit
          </Button>
        </div>
      </div>

      <Card>
        <CardContent className="p-0">
          {units.length === 0 ? (
            <p className="p-6 text-sm text-slate-500">Belum ada unit kompetensi pada skema ini.</p>
          ) : (
            <Accordion type="multiple" className="px-4">
              {units.map((unit) => (
                <AccordionItem key={unit.id} value={unit.id}>
                  <div className="flex items-center gap-2">
                    <div className="flex-1">
                      <AccordionTrigger>
                        <span className="flex flex-wrap items-baseline gap-2">
                          <code className="rounded bg-slate-100 px-1.5 py-0.5 text-xs font-medium text-slate-700">{unit.code}</code>
                          <span className="text-slate-900">{unit.title}</span>
                          <span className="text-xs font-normal text-slate-400">urutan {unit.sort_order}</span>
                        </span>
                      </AccordionTrigger>
                    </div>
                    <div className="flex shrink-0 items-center gap-1">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="gap-1.5"
                        onClick={() => setElementForm({ open: true, unitId: unit.id, element: null })}
                      >
                        <Plus className="h-4 w-4" />
                        Elemen
                      </Button>
                      <Button variant="ghost" size="icon" aria-label={`Ubah unit ${unit.title}`} onClick={() => setUnitForm({ open: true, unit })}>
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        aria-label={`Hapus unit ${unit.title}`}
                        className="text-rose-600 hover:bg-rose-50"
                        onClick={() => setDeletingUnit(unit)}
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>

                  <AccordionContent className="space-y-3">
                    {unit.elements.length === 0 ? (
                      <p className="text-sm text-slate-500">Belum ada elemen pada unit ini.</p>
                    ) : (
                      unit.elements.map((element) => (
                        <div key={element.id} className="rounded-md border border-slate-200">
                          <div className="flex items-center justify-between gap-2 border-b border-slate-100 bg-slate-50 px-3 py-2">
                            <div className="text-sm">
                              <span className="font-medium text-slate-900">{element.title || 'Elemen'}</span>
                              <span className="ml-2 text-xs text-slate-400">urutan {element.sort_order}</span>
                            </div>
                            <div className="flex shrink-0 items-center gap-1">
                              <Button
                                variant="ghost"
                                size="sm"
                                className="gap-1.5"
                                onClick={() => setCriteriaForm({ open: true, elementId: element.id, criteria: null })}
                              >
                                <Plus className="h-4 w-4" />
                                KUK
                              </Button>
                              <Button
                                variant="ghost"
                                size="icon"
                                aria-label="Ubah elemen"
                                onClick={() => setElementForm({ open: true, unitId: unit.id, element })}
                              >
                                <Pencil className="h-4 w-4" />
                              </Button>
                              <Button
                                variant="ghost"
                                size="icon"
                                aria-label="Hapus elemen"
                                className="text-rose-600 hover:bg-rose-50"
                                onClick={() => setDeletingElement(element)}
                              >
                                <Trash2 className="h-4 w-4" />
                              </Button>
                            </div>
                          </div>

                          <ul className="divide-y divide-slate-100">
                            {element.criteria.length === 0 ? (
                              <li className="px-3 py-2 text-xs text-slate-400">Belum ada KUK.</li>
                            ) : (
                              element.criteria.map((criteria) => (
                                <li key={criteria.id} className="flex items-start justify-between gap-3 px-3 py-2">
                                  <p className="text-sm text-slate-700">
                                    <span className="mr-2 text-xs text-slate-400">{criteria.sort_order}.</span>
                                    {criteria.content}
                                  </p>
                                  <div className="flex shrink-0 items-center gap-1">
                                    <Button
                                      variant="ghost"
                                      size="icon"
                                      aria-label="Ubah KUK"
                                      onClick={() => setCriteriaForm({ open: true, elementId: element.id, criteria })}
                                    >
                                      <Pencil className="h-4 w-4" />
                                    </Button>
                                    <Button
                                      variant="ghost"
                                      size="icon"
                                      aria-label="Hapus KUK"
                                      className="text-rose-600 hover:bg-rose-50"
                                      onClick={() => setDeletingCriteria(criteria)}
                                    >
                                      <Trash2 className="h-4 w-4" />
                                    </Button>
                                  </div>
                                </li>
                              ))
                            )}
                          </ul>
                        </div>
                      ))
                    )}
                  </AccordionContent>
                </AccordionItem>
              ))}
            </Accordion>
          )}
        </CardContent>
      </Card>

      <SchemaFormDialog open={schemaFormOpen} onOpenChange={setSchemaFormOpen} schema={tree} session={auth} />

      <UnitFormDialog
        open={unitForm.open}
        onOpenChange={(open) => setUnitForm((prev) => ({ ...prev, open }))}
        schemaId={schemaId}
        unit={unitForm.unit}
        session={auth}
      />

      <ElementFormDialog
        open={elementForm.open}
        onOpenChange={(open) => setElementForm((prev) => ({ ...prev, open }))}
        schemaId={schemaId}
        unitId={elementForm.unitId}
        element={elementForm.element}
        session={auth}
      />

      <CriteriaFormDialog
        open={criteriaForm.open}
        onOpenChange={(open) => setCriteriaForm((prev) => ({ ...prev, open }))}
        schemaId={schemaId}
        elementId={criteriaForm.elementId}
        criteria={criteriaForm.criteria}
        session={auth}
      />

      <ConfirmDialog
        open={Boolean(deletingUnit)}
        onOpenChange={(open) => (!open ? setDeletingUnit(null) : undefined)}
        title="Hapus Unit Kompetensi"
        description={`Yakin ingin menghapus unit "${deletingUnit?.title}"? Seluruh elemen dan KUK di dalamnya akan ikut terhapus.`}
        confirmLabel="Hapus"
        pendingLabel="Menghapus..."
        destructive
        pending={removeUnit.isPending}
        onConfirm={() => {
          if (!deletingUnit) return
          removeUnit.mutate(deletingUnit.id, { onSuccess: () => setDeletingUnit(null) })
        }}
      />

      <ConfirmDialog
        open={Boolean(deletingElement)}
        onOpenChange={(open) => (!open ? setDeletingElement(null) : undefined)}
        title="Hapus Elemen"
        description={`Yakin ingin menghapus elemen "${deletingElement?.title || 'ini'}"? Seluruh KUK di dalamnya akan ikut terhapus.`}
        confirmLabel="Hapus"
        pendingLabel="Menghapus..."
        destructive
        pending={removeElement.isPending}
        onConfirm={() => {
          if (!deletingElement) return
          removeElement.mutate(deletingElement.id, { onSuccess: () => setDeletingElement(null) })
        }}
      />

      <ConfirmDialog
        open={Boolean(deletingCriteria)}
        onOpenChange={(open) => (!open ? setDeletingCriteria(null) : undefined)}
        title="Hapus KUK"
        description="Yakin ingin menghapus kriteria unjuk kerja ini?"
        confirmLabel="Hapus"
        pendingLabel="Menghapus..."
        destructive
        pending={removeCriteria.isPending}
        onConfirm={() => {
          if (!deletingCriteria) return
          removeCriteria.mutate(deletingCriteria.id, { onSuccess: () => setDeletingCriteria(null) })
        }}
      />
    </div>
  )
}
