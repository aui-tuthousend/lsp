import type { CreateAsesiProfiles } from '../-types'

export const controlClass =
  'flex h-10 w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-slate-400 focus:ring-2 focus:ring-slate-200'

export function inputErrorClass(hasError?: boolean) {
  return hasError ? 'border-rose-300 focus:border-rose-400 focus:ring-rose-100' : undefined
}

export type StepFieldChange = <K extends keyof CreateAsesiProfiles>(key: K, value: CreateAsesiProfiles[K]) => void

export type StepProps = {
  values: CreateAsesiProfiles
  errors: Record<string, string>
  onChange: StepFieldChange
}
