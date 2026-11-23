import * as React from 'react'
import { cn } from '@/lib/utils'

/** Select native yang diselaraskan dengan gaya `Input`. */
export function Select({ className, children, ...props }: React.SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <select
      className={cn(
        'flex h-10 w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-slate-300 disabled:cursor-not-allowed disabled:bg-slate-50',
        className,
      )}
      {...props}
    >
      {children}
    </select>
  )
}
