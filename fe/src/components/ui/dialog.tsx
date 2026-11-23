import * as React from 'react'
import { createPortal } from 'react-dom'
import { X } from 'lucide-react'
import { cn } from '@/lib/utils'

// Adaptasi komponen shadcn "dialog" ke design system proyek (palet slate),
// tanpa dependensi radix. API sama: Dialog, DialogTrigger, DialogClose,
// DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter.

type DialogContextValue = {
  open: boolean
  setOpen: (open: boolean) => void
}

const DialogContext = React.createContext<DialogContextValue | null>(null)

function useDialog() {
  const context = React.useContext(DialogContext)
  if (!context) throw new Error('Bagian Dialog harus dipakai di dalam <Dialog>')
  return context
}

export type DialogProps = {
  open?: boolean
  onOpenChange?: (open: boolean) => void
  defaultOpen?: boolean
  children: React.ReactNode
}

export function Dialog({ open, onOpenChange, defaultOpen = false, children }: DialogProps) {
  const [internalOpen, setInternalOpen] = React.useState(defaultOpen)
  const isControlled = open !== undefined
  const openValue = isControlled ? open : internalOpen

  const setOpen = React.useCallback(
    (next: boolean) => {
      if (!isControlled) setInternalOpen(next)
      onOpenChange?.(next)
    },
    [isControlled, onOpenChange],
  )

  return <DialogContext.Provider value={{ open: openValue, setOpen }}>{children}</DialogContext.Provider>
}

export function DialogTrigger({ className, onClick, ...props }: React.ComponentProps<'button'>) {
  const { setOpen } = useDialog()
  return (
    <button
      type="button"
      {...props}
      onClick={(event) => {
        onClick?.(event)
        setOpen(true)
      }}
      className={className}
    />
  )
}

export function DialogClose({ className, onClick, ...props }: React.ComponentProps<'button'>) {
  const { setOpen } = useDialog()
  return (
    <button
      type="button"
      {...props}
      onClick={(event) => {
        onClick?.(event)
        setOpen(false)
      }}
      className={className}
    />
  )
}

export function DialogContent({
  className,
  children,
  showCloseButton = true,
  ...props
}: React.ComponentProps<'div'> & { showCloseButton?: boolean }) {
  const { open, setOpen } = useDialog()
  const [mounted, setMounted] = React.useState(false)

  React.useEffect(() => setMounted(true), [])

  React.useEffect(() => {
    if (!open) return
    function onKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') setOpen(false)
    }
    document.addEventListener('keydown', onKeyDown)
    const previousOverflow = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    return () => {
      document.removeEventListener('keydown', onKeyDown)
      document.body.style.overflow = previousOverflow
    }
  }, [open, setOpen])

  if (!mounted || !open) return null

  return createPortal(
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <button
        type="button"
        aria-label="Tutup dialog"
        className="absolute inset-0 bg-slate-950/50"
        onClick={() => setOpen(false)}
      />
      <div
        role="dialog"
        aria-modal="true"
        {...props}
        className={cn(
          'relative z-10 grid max-h-[85vh] w-full max-w-lg gap-4 overflow-y-auto rounded-lg border border-slate-200 bg-white p-6 shadow-lg',
          className,
        )}
      >
        {children}
        {showCloseButton ? (
          <button
            type="button"
            aria-label="Tutup"
            onClick={() => setOpen(false)}
            className="absolute right-4 top-4 rounded-md p-1 text-slate-500 transition-colors hover:bg-slate-100 hover:text-slate-900"
          >
            <X className="h-4 w-4" />
          </button>
        ) : null}
      </div>
    </div>,
    document.body,
  )
}

export function DialogHeader({ className, ...props }: React.ComponentProps<'div'>) {
  return <div className={cn('flex flex-col gap-1.5 text-left', className)} {...props} />
}

export function DialogFooter({ className, ...props }: React.ComponentProps<'div'>) {
  return <div className={cn('flex flex-col-reverse gap-2 sm:flex-row sm:justify-end', className)} {...props} />
}

export function DialogTitle({ className, ...props }: React.ComponentProps<'h2'>) {
  return <h2 className={cn('text-lg font-semibold leading-none text-slate-900', className)} {...props} />
}

export function DialogDescription({ className, ...props }: React.ComponentProps<'p'>) {
  return <p className={cn('text-sm text-slate-500', className)} {...props} />
}
