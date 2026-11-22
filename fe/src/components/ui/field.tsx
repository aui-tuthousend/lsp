import * as React from 'react'
import { cn } from '@/lib/utils'

export type FieldProps = {
  label: string
  htmlFor?: string
  error?: string
  hint?: string
  required?: boolean
  className?: string
  children: React.ReactNode
}

/** Pembungkus label + kontrol form + pesan error/hint. */
export function Field({ label, htmlFor, error, hint, required, className, children }: FieldProps) {
  return (
    <div className={cn('space-y-1.5', className)}>
      <label htmlFor={htmlFor} className="flex items-center gap-1 text-sm font-medium text-slate-700">
        {label}
        {required ? <span className="text-rose-500">*</span> : null}
      </label>
      {children}
      {error ? (
        <p className="text-xs font-medium text-rose-600">{error}</p>
      ) : hint ? (
        <p className="text-xs text-slate-400">{hint}</p>
      ) : null}
    </div>
  )
}
