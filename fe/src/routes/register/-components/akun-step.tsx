import { useState } from 'react'
import { Eye, EyeOff } from 'lucide-react'
import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils'
import { inputErrorClass, type StepProps } from './form'

export function AkunStep({ values, errors, onChange }: StepProps) {
  const [showPassword, setShowPassword] = useState(false)

  return (
    <div className="space-y-4">
      <Field label="Email" htmlFor="email" required error={errors.email} hint="Email aktif untuk verifikasi dan login">
        <Input
          id="email"
          type="email"
          autoComplete="email"
          value={values.email}
          onChange={(event) => onChange('email', event.target.value)}
          placeholder="nama@email.com"
          className={inputErrorClass(Boolean(errors.email))}
        />
      </Field>

      <Field label="Password" htmlFor="password" required error={errors.password}>
        <div className="relative">
          <Input
            id="password"
            type={showPassword ? 'text' : 'password'}
            autoComplete="new-password"
            value={values.password}
            onChange={(event) => onChange('password', event.target.value)}
            placeholder="Minimal 8 karakter"
            className={cn('pr-10', inputErrorClass(Boolean(errors.password)))}
          />
          <button
            type="button"
            onClick={() => setShowPassword((prev) => !prev)}
            aria-label={showPassword ? 'Sembunyikan password' : 'Tampilkan password'}
            className="absolute inset-y-0 right-0 flex items-center px-3 text-slate-400 transition hover:text-slate-600"
          >
            {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
          </button>
        </div>
        <PasswordStrength value={values.password} />
      </Field>

      <Field label="Konfirmasi Password" htmlFor="password_confirmation" required error={errors.password_confirmation}>
        <Input
          id="password_confirmation"
          type={showPassword ? 'text' : 'password'}
          autoComplete="new-password"
          value={values.password_confirmation}
          onChange={(event) => onChange('password_confirmation', event.target.value)}
          placeholder="Ulangi password"
          className={inputErrorClass(Boolean(errors.password_confirmation))}
        />
      </Field>
    </div>
  )
}

function passwordScore(value: string) {
  let score = 0
  if (value.length >= 8) score += 1
  if (/[a-z]/.test(value) && /[A-Z]/.test(value)) score += 1
  if (/\d/.test(value)) score += 1
  if (/[^A-Za-z0-9]/.test(value)) score += 1
  return score
}

function PasswordStrength({ value }: { value: string }) {
  if (!value) return null
  const score = passwordScore(value)
  const labels = ['Sangat lemah', 'Lemah', 'Cukup', 'Kuat', 'Sangat kuat']
  const colors = ['bg-rose-400', 'bg-orange-400', 'bg-amber-400', 'bg-emerald-400', 'bg-emerald-500']
  return (
    <div className="space-y-1">
      <div className="flex gap-1">
        {[0, 1, 2, 3].map((index) => (
          <span key={index} className={cn('h-1 flex-1 rounded-full', index < score ? colors[score] : 'bg-slate-200')} />
        ))}
      </div>
      <p className="text-xs text-slate-400">Kekuatan password: {labels[score]}</p>
    </div>
  )
}
