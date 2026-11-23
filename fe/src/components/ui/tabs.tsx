import * as React from 'react'
import { cn } from '@/lib/utils'

// Adaptasi komponen shadcn "tabs" ke design system proyek (palet slate),
// tanpa dependensi radix. API sama: Tabs, TabsList, TabsTrigger, TabsContent.
// Konten tab tetap ter-mount (disembunyikan) agar state form tidak hilang
// saat berpindah tab.

type TabsContextValue = { value: string; setValue: (value: string) => void }

const TabsContext = React.createContext<TabsContextValue | null>(null)

function useTabs() {
  const context = React.useContext(TabsContext)
  if (!context) throw new Error('Bagian Tabs harus dipakai di dalam <Tabs>')
  return context
}

export type TabsProps = React.ComponentProps<'div'> & {
  value?: string
  defaultValue?: string
  onValueChange?: (value: string) => void
}

export function Tabs({ value, defaultValue, onValueChange, className, children, ...props }: TabsProps) {
  const [internal, setInternal] = React.useState(defaultValue ?? '')
  const current = value ?? internal

  function setValue(next: string) {
    if (value === undefined) setInternal(next)
    onValueChange?.(next)
  }

  return (
    <TabsContext.Provider value={{ value: current, setValue }}>
      <div data-slot="tabs" className={cn('flex flex-col gap-4', className)} {...props}>
        {children}
      </div>
    </TabsContext.Provider>
  )
}

export function TabsList({ className, ...props }: React.ComponentProps<'div'>) {
  return (
    <div
      data-slot="tabs-list"
      role="tablist"
      className={cn('inline-flex w-full items-center gap-1 overflow-x-auto rounded-lg bg-slate-100 p-1 sm:w-auto', className)}
      {...props}
    />
  )
}

export function TabsTrigger({ value, className, children, onClick, ...props }: React.ComponentProps<'button'> & { value: string }) {
  const { value: active, setValue } = useTabs()
  const selected = active === value

  return (
    <button
      type="button"
      role="tab"
      aria-selected={selected}
      data-state={selected ? 'active' : 'inactive'}
      onClick={(event) => {
        onClick?.(event)
        setValue(value)
      }}
      className={cn(
        'inline-flex flex-1 items-center justify-center whitespace-nowrap rounded-md px-3 py-1.5 text-sm font-medium text-slate-500 transition outline-none hover:text-slate-700 focus-visible:ring-2 focus-visible:ring-slate-300 sm:flex-none',
        selected && 'bg-white text-slate-900 shadow-sm',
        className,
      )}
      {...props}
    >
      {children}
    </button>
  )
}

export function TabsContent({ value, className, children, ...props }: React.ComponentProps<'div'> & { value: string }) {
  const { value: active } = useTabs()
  const selected = active === value

  return (
    <div data-slot="tabs-content" role="tabpanel" data-state={selected ? 'active' : 'inactive'} hidden={!selected} className={cn('outline-none', className)} {...props}>
      {children}
    </div>
  )
}
