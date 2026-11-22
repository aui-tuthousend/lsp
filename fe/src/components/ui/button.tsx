import * as React from 'react'
import { cn } from '@/lib/utils'

export type ButtonProps = React.ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: 'default' | 'outline' | 'ghost'
  size?: 'default' | 'sm' | 'icon'
}

export function Button({ className, variant = 'default', size = 'default', ...props }: ButtonProps) {
  return (
    <button
      className={cn(
        'inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors disabled:pointer-events-none disabled:opacity-50',
        size === 'default' && 'h-10 px-4 py-2',
        size === 'sm' && 'h-9 px-3',
        size === 'icon' && 'h-8 w-8 p-0',
        variant === 'default' && 'bg-slate-900 text-white hover:bg-slate-800',
        variant === 'outline' && 'border border-slate-200 bg-white hover:bg-slate-50',
        variant === 'ghost' && 'hover:bg-slate-100',
        className,
      )}
      {...props}
    />
  )
}
