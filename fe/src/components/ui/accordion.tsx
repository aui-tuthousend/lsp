import * as React from 'react'
import { ChevronDown } from 'lucide-react'
import { cn } from '@/lib/utils'

// Adaptasi komponen shadcn "accordion" ke design system proyek (palet slate),
// tanpa dependensi radix. API sama: Accordion, AccordionItem, AccordionTrigger,
// AccordionContent.

type AccordionContextValue = {
  openValues: string[]
  toggle: (value: string) => void
}

const AccordionContext = React.createContext<AccordionContextValue | null>(null)

function useAccordion() {
  const context = React.useContext(AccordionContext)
  if (!context) throw new Error('Bagian Accordion harus dipakai di dalam <Accordion>')
  return context
}

export type AccordionProps = React.ComponentProps<'div'> & {
  type?: 'single' | 'multiple'
  collapsible?: boolean
  defaultValue?: string | string[]
}

export function Accordion({ type = 'single', collapsible = false, defaultValue, className, children, ...props }: AccordionProps) {
  const [openValues, setOpenValues] = React.useState<string[]>(() => {
    if (defaultValue === undefined) return []
    return Array.isArray(defaultValue) ? defaultValue : [defaultValue]
  })

  function toggle(value: string) {
    setOpenValues((prev) => {
      const isOpen = prev.includes(value)
      if (type === 'multiple') {
        return isOpen ? prev.filter((item) => item !== value) : [...prev, value]
      }
      if (isOpen) return collapsible ? [] : prev
      return [value]
    })
  }

  return (
    <AccordionContext.Provider value={{ openValues, toggle }}>
      <div data-slot="accordion" className={className} {...props}>
        {children}
      </div>
    </AccordionContext.Provider>
  )
}

type AccordionItemContextValue = { value: string; open: boolean }

const AccordionItemContext = React.createContext<AccordionItemContextValue | null>(null)

function useAccordionItem() {
  const context = React.useContext(AccordionItemContext)
  if (!context) throw new Error('Bagian AccordionItem harus dipakai di dalam <AccordionItem>')
  return context
}

export function AccordionItem({ value, className, children, ...props }: React.ComponentProps<'div'> & { value: string }) {
  const { openValues } = useAccordion()
  const open = openValues.includes(value)

  return (
    <AccordionItemContext.Provider value={{ value, open }}>
      <div data-slot="accordion-item" data-state={open ? 'open' : 'closed'} className={cn('border-b last:border-b-0', className)} {...props}>
        {children}
      </div>
    </AccordionItemContext.Provider>
  )
}

export function AccordionTrigger({ className, children, ...props }: React.ComponentProps<'button'>) {
  const { toggle } = useAccordion()
  const { value, open } = useAccordionItem()

  return (
    <h3 className="flex">
      <button
        type="button"
        data-slot="accordion-trigger"
        aria-expanded={open}
        {...props}
        onClick={(event) => {
          props.onClick?.(event)
          toggle(value)
        }}
        className={cn(
          'flex flex-1 items-start justify-between gap-4 rounded-md py-4 text-left text-sm font-medium text-slate-900 outline-none transition-all hover:underline focus-visible:ring-2 focus-visible:ring-slate-300',
          className,
        )}
      >
        {children}
        <ChevronDown
          className={cn('pointer-events-none mt-0.5 h-4 w-4 shrink-0 text-slate-500 transition-transform duration-200', open && 'rotate-180')}
        />
      </button>
    </h3>
  )
}

export function AccordionContent({ className, children, ...props }: React.ComponentProps<'div'>) {
  const { open } = useAccordionItem()

  return (
    <div
      data-slot="accordion-content"
      data-state={open ? 'open' : 'closed'}
      className={cn('grid transition-[grid-template-rows] duration-200', open ? 'grid-rows-[1fr]' : 'grid-rows-[0fr]')}
      {...props}
    >
      <div className="overflow-hidden">
        <div className={cn('pb-4 pt-0', className)}>{children}</div>
      </div>
    </div>
  )
}
