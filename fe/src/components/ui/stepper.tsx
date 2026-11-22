import * as React from 'react'
import { Check } from 'lucide-react'
import { cn } from '@/lib/utils'

export type StepperStep = {
  title: string
  icon: React.ComponentType<{ className?: string }>
}

export function Stepper({ steps, current }: { steps: readonly StepperStep[]; current: number }) {
  return (
    <ol className="flex items-center gap-2">
      {steps.map((step, index) => {
        const state = index < current ? 'done' : index === current ? 'active' : 'todo'
        const StepIcon = step.icon
        return (
          <li key={step.title} className="flex min-w-0 flex-1 items-center gap-2">
            <div className="flex min-w-0 items-center gap-2">
              <span
                className={cn(
                  'flex h-9 w-9 shrink-0 items-center justify-center rounded-full border text-sm font-semibold transition',
                  state === 'done' && 'border-emerald-500 bg-emerald-500 text-white',
                  state === 'active' && 'border-slate-900 bg-slate-900 text-white',
                  state === 'todo' && 'border-slate-200 bg-white text-slate-400',
                )}
              >
                {state === 'done' ? <Check className="h-4 w-4" /> : <StepIcon className="h-4 w-4" />}
              </span>
              <span className="hidden min-w-0 flex-col sm:flex">
                <span className={cn('truncate text-xs font-semibold', state === 'todo' ? 'text-slate-400' : 'text-slate-900')}>
                  {step.title}
                </span>
                <span className="text-[11px] text-slate-400">Langkah {index + 1}</span>
              </span>
            </div>
            {index < steps.length - 1 ? (
              <span className={cn('h-px flex-1 rounded', index < current ? 'bg-emerald-400' : 'bg-slate-200')} />
            ) : null}
          </li>
        )
      })}
    </ol>
  )
}
