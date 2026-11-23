import {
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
  type ColumnDef,
  type PaginationState,
} from '@tanstack/react-table'
import { ChevronLeft, ChevronRight, Eye, Pencil, Trash2 } from 'lucide-react'
import { useMemo } from 'react'

export type PagedResponse<T> = {
  data: T[]
  meta: {
    page: number
    limit: number
    total: number
  }
}

type Primitive = string | number | boolean | bigint | symbol | null | undefined | Date

export type DataPath<T> = T extends Primitive
  ? never
  : {
      [Key in keyof T & string]: T[Key] extends Primitive | readonly unknown[]
        ? Key
        : Key | `${Key}.${DataPath<T[Key]>}`
    }[keyof T & string]

type PathValue<T, Path extends string> = Path extends `${infer Head}.${infer Tail}`
  ? Head extends keyof T
    ? PathValue<T[Head], Tail>
    : never
  : Path extends keyof T
    ? T[Path]
    : never

export type DataTableAttribute<T> = {
  [Path in DataPath<T>]: {
    accessor: Path
    header: string
    format?: (value: PathValue<T, Path>, row: T) => string
  }
}[DataPath<T>]

type DataTableProps<T> = {
  data: PagedResponse<T>
  attribute: readonly DataTableAttribute<T>[]
  getRowId: (row: T) => string
  onPageChange: (page: number, limit: number) => void
  onRead: (row: T) => void
  onEdit: (row: T) => void
  onDelete: (row: T) => void
  label?: string
  pageSizeOptions?: number[]
}

function getPathValue<T>(row: T, path: string): unknown {
  return path.split('.').reduce<unknown>((value, key) => {
    if (value === null || typeof value !== 'object') return undefined
    return (value as Record<string, unknown>)[key]
  }, row)
}

function formatValue(value: unknown): string {
  if (value === null || value === undefined || value === '') return '—'
  if (value instanceof Date) return new Intl.DateTimeFormat('id-ID').format(value)
  if (typeof value === 'number') {
    return new Intl.NumberFormat('id-ID', { maximumFractionDigits: 2 }).format(value)
  }
  if (typeof value === 'boolean') return value ? 'Ya' : 'Tidak'
  return String(value)
}

export function DataTable<T>({
  data,
  attribute,
  getRowId,
  onPageChange,
  onRead,
  onEdit,
  onDelete,
  label = 'Daftar data',
  pageSizeOptions = [8, 12, 20],
}: DataTableProps<T>) {
  const columns = useMemo<ColumnDef<T, unknown>[]>(() => {
    const helper = createColumnHelper<T>()
    const dataColumns: ColumnDef<T, unknown>[] = attribute.map((column) =>
      helper.accessor((row) => getPathValue(row, column.accessor), {
        id: column.accessor,
        header: column.header,
        cell: (info) => {
          const value = info.getValue()
          const formatted = column.format
            ? column.format(value as never, info.row.original)
            : formatValue(value)
          return <span className="whitespace-nowrap">{formatted}</span>
        },
      }),
    )

    dataColumns.push(
      helper.display({
        id: 'actions',
        header: () => <span className="sr-only">Aksi</span>,
        cell: ({ row }) => {
          const item = row.original
          const rowLabel = getRowId(item)
          return (
            <div className="flex items-center justify-end gap-1">
              <button
                aria-label={`Baca data ${rowLabel}`}
                className="inline-flex size-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                onClick={() => onRead(item)}
                title="Baca"
                type="button"
              >
                <Eye aria-hidden="true" className="size-4" />
              </button>
              <button
                aria-label={`Edit data ${rowLabel}`}
                className="inline-flex size-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-muted hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                onClick={() => onEdit(item)}
                title="Edit"
                type="button"
              >
                <Pencil aria-hidden="true" className="size-4" />
              </button>
              <button
                aria-label={`Hapus data ${rowLabel}`}
                className="inline-flex size-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-destructive/10 hover:text-destructive focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
                onClick={() => onDelete(item)}
                title="Hapus"
                type="button"
              >
                <Trash2 aria-hidden="true" className="size-4" />
              </button>
            </div>
          )
        },
      }),
    )

    return dataColumns
  }, [attribute, getRowId, onDelete, onEdit, onRead])

  const pagination: PaginationState = {
    pageIndex: Math.max(data.meta.page - 1, 0),
    pageSize: data.meta.limit,
  }
  const pageCount = Math.max(Math.ceil(data.meta.total / data.meta.limit), 1)

  const table = useReactTable({
    data: data.data,
    columns,
    getRowId,
    getCoreRowModel: getCoreRowModel(),
    manualPagination: true,
    pageCount,
    state: { pagination },
    onPaginationChange: (updater) => {
      const next = typeof updater === 'function' ? updater(pagination) : updater
      if (next.pageSize !== data.meta.limit) {
        onPageChange(1, next.pageSize)
      } else {
        onPageChange(next.pageIndex + 1, next.pageSize)
      }
    },
  })

  const firstItem = data.meta.total === 0 ? 0 : (data.meta.page - 1) * data.meta.limit + 1
  const lastItem = Math.min(data.meta.page * data.meta.limit, data.meta.total)

  return (
    <section aria-label={label} className="overflow-hidden rounded-xl border border-border bg-card">
      <div className="overflow-x-auto">
        <table className="w-full min-w-[720px] border-collapse text-left text-sm">
          <caption className="sr-only">{label}</caption>
          <thead className="bg-muted/50 text-xs font-medium uppercase tracking-[0.08em] text-muted-foreground">
            {table.getHeaderGroups().map((headerGroup) => (
              <tr key={headerGroup.id}>
                {headerGroup.headers.map((header, index) => (
                  <th
                    className={`h-12 px-5 ${index === headerGroup.headers.length - 1 ? 'text-right' : ''}`}
                    key={header.id}
                    scope="col"
                  >
                    {header.isPlaceholder
                      ? null
                      : flexRender(header.column.columnDef.header, header.getContext())}
                  </th>
                ))}
              </tr>
            ))}
          </thead>
          <tbody className="divide-y divide-border">
            {table.getRowModel().rows.length > 0 ? (
              table.getRowModel().rows.map((row) => (
                <tr className="transition-colors hover:bg-muted/30" key={row.id}>
                  {row.getVisibleCells().map((cell, index) => (
                    <td
                      className={`h-[58px] px-5 text-foreground ${index === row.getVisibleCells().length - 1 ? 'text-right' : ''}`}
                      key={cell.id}
                    >
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </td>
                  ))}
                </tr>
              ))
            ) : (
              <tr>
                <td
                  className="h-28 px-5 text-center text-sm text-muted-foreground"
                  colSpan={table.getAllLeafColumns().length}
                >
                  Belum ada data untuk ditampilkan.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <div className="flex flex-col gap-4 border-t border-border px-5 py-4 sm:flex-row sm:items-center sm:justify-between">
        <p aria-live="polite" className="text-sm text-muted-foreground">
          Menampilkan <span className="font-medium text-foreground">{firstItem}–{lastItem}</span> dari{' '}
          <span className="font-medium text-foreground">{data.meta.total}</span> data
        </p>
        <div className="flex flex-wrap items-center gap-3">
          <label className="flex items-center gap-2 text-sm text-muted-foreground">
            Baris
            <select
              aria-label="Jumlah baris per halaman"
              className="h-9 rounded-md border border-input bg-background px-2.5 text-sm text-foreground outline-none focus-visible:ring-2 focus-visible:ring-ring"
              onChange={(event) => table.setPageSize(Number(event.target.value))}
              value={data.meta.limit}
            >
              {pageSizeOptions.map((size) => (
                <option key={size} value={size}>
                  {size}
                </option>
              ))}
            </select>
          </label>
          <span className="min-w-20 text-center text-sm text-muted-foreground">
            Halaman {data.meta.page} dari {pageCount}
          </span>
          <div className="flex items-center gap-1">
            <button
              aria-label="Halaman sebelumnya"
              className="inline-flex size-9 items-center justify-center rounded-md border border-input text-foreground transition-colors hover:bg-muted disabled:pointer-events-none disabled:opacity-40"
              disabled={!table.getCanPreviousPage()}
              onClick={() => table.previousPage()}
              type="button"
            >
              <ChevronLeft aria-hidden="true" className="size-4" />
            </button>
            <button
              aria-label="Halaman berikutnya"
              className="inline-flex size-9 items-center justify-center rounded-md border border-input text-foreground transition-colors hover:bg-muted disabled:pointer-events-none disabled:opacity-40"
              disabled={!table.getCanNextPage()}
              onClick={() => table.nextPage()}
              type="button"
            >
              <ChevronRight aria-hidden="true" className="size-4" />
            </button>
          </div>
        </div>
      </div>
    </section>
  )
}
