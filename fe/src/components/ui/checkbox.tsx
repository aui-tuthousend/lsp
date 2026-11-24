import * as React from 'react'
import { Check } from 'lucide-react'
import { cn } from '@/lib/utils'

export type CheckboxProps = Omit<React.InputHTMLAttributes<HTMLInputElement>, 'type'> & {
  label?: React.ReactNode
}

/** Checkbox sederhana (adaptasi shadcn tanpa radix) dengan label opsional. */
export function Checkbox({ className, label, id, ...props }: CheckboxProps) {
  const inputId = id ?? props.name
  return (
    <label
      htmlFor={inputId}
      className={cn(
        'inline-flex items-center gap-2 text-sm text-slate-700',
        props.disabled ? 'cursor-not-allowed opacity-60' : 'cursor-pointer',
        className,
      )}
    >
      <span className="relative inline-flex h-4 w-4 shrink-0 items-center justify-center">
        <input
          id={inputId}
          type="checkbox"
          className="peer h-4 w-4 cursor-pointer appearance-none rounded border border-slate-300 bg-white outline-none transition-colors checked:border-slate-900 checked:bg-slate-900 focus-visible:ring-2 focus-visible:ring-slate-300 disabled:cursor-not-allowed"
          {...props}
        />
        <Check className="pointer-events-none absolute h-3 w-3 text-white opacity-0 peer-checked:opacity-100" />
      </span>
      {label ? <span>{label}</span> : null}
    </label>
  )
}
