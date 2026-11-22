import * as React from 'react'
import { cn } from '@/lib/utils'

export type AvatarProps = {
  name: string
  src?: string | null
  alt?: string
  className?: string
}

function initials(name: string): string {
  const parts = name.trim().split(/\s+/).filter(Boolean)
  if (parts.length === 0) return '?'
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase()
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
}

/** Avatar user: menampilkan foto bila `src` ada, jika tidak pakai inisial nama. */
export function Avatar({ name, src, alt, className }: AvatarProps) {
  const [broken, setBroken] = React.useState(false)
  const showImage = Boolean(src) && !broken

  return (
    <span
      className={cn(
        'inline-flex h-9 w-9 shrink-0 select-none items-center justify-center overflow-hidden rounded-full bg-slate-200 text-sm font-semibold text-slate-700',
        className,
      )}
    >
      {showImage ? (
        <img src={src as string} alt={alt ?? name} className="h-full w-full object-cover" onError={() => setBroken(true)} />
      ) : (
        initials(name)
      )}
    </span>
  )
}
