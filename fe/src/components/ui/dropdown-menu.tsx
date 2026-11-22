import * as React from 'react'
import { cn } from '@/lib/utils'

// Dropdown minimal tanpa dependensi eksternal: menutup saat klik di luar,
// tekan Escape, atau memilih item.

type DropdownContextValue = {
  open: boolean
  setOpen: (open: boolean) => void
  contentId: string
}

const DropdownContext = React.createContext<DropdownContextValue | null>(null)

function useDropdown() {
  const context = React.useContext(DropdownContext)
  if (!context) throw new Error('Bagian DropdownMenu harus dipakai di dalam <DropdownMenu>')
  return context
}

export function DropdownMenu({ children, className }: { children: React.ReactNode; className?: string }) {
  const [open, setOpen] = React.useState(false)
  const contentId = React.useId()
  const rootRef = React.useRef<HTMLDivElement>(null)

  React.useEffect(() => {
    if (!open) return
    function onPointerDown(event: PointerEvent) {
      if (rootRef.current && !rootRef.current.contains(event.target as Node)) setOpen(false)
    }
    function onKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') setOpen(false)
    }
    document.addEventListener('pointerdown', onPointerDown)
    document.addEventListener('keydown', onKeyDown)
    return () => {
      document.removeEventListener('pointerdown', onPointerDown)
      document.removeEventListener('keydown', onKeyDown)
    }
  }, [open])

  return (
    <DropdownContext.Provider value={{ open, setOpen, contentId }}>
      <div ref={rootRef} className={cn('relative', className)}>
        {children}
      </div>
    </DropdownContext.Provider>
  )
}

export function DropdownMenuTrigger({ className, onClick, ...props }: React.ComponentProps<'button'>) {
  const { open, setOpen, contentId } = useDropdown()
  return (
    <button
      type="button"
      aria-haspopup="menu"
      aria-expanded={open}
      aria-controls={open ? contentId : undefined}
      {...props}
      onClick={(event) => {
        onClick?.(event)
        setOpen(!open)
      }}
      className={className}
    />
  )
}

export function DropdownMenuContent({
  className,
  align = 'end',
  ...props
}: React.ComponentProps<'div'> & { align?: 'start' | 'end' }) {
  const { open, contentId } = useDropdown()
  if (!open) return null
  return (
    <div
      id={contentId}
      role="menu"
      {...props}
      className={cn(
        'absolute top-full z-50 mt-2 min-w-48 rounded-lg border border-slate-200 bg-white p-1 shadow-lg',
        align === 'end' ? 'right-0' : 'left-0',
        className,
      )}
    />
  )
}

export function DropdownMenuItem({
  className,
  onSelect,
  onClick,
  ...props
}: React.ComponentProps<'button'> & { onSelect?: () => void }) {
  const { setOpen } = useDropdown()
  return (
    <button
      type="button"
      role="menuitem"
      {...props}
      onClick={(event) => {
        onClick?.(event)
        onSelect?.()
        setOpen(false)
      }}
      className={cn(
        'flex w-full items-center gap-2 rounded-md px-2.5 py-2 text-left text-sm font-medium text-slate-700 transition-colors hover:bg-slate-100',
        className,
      )}
    />
  )
}

export function DropdownMenuSeparator({ className }: { className?: string }) {
  return <div role="separator" className={cn('-mx-1 my-1 h-px bg-slate-100', className)} />
}
