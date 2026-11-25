import { Field } from '@/components/ui/field'
import { Input } from '@/components/ui/input'
import { Select } from '@/components/ui/select'
import { MODE_OPTIONS } from '../-helpers'
import type { AssessmentMode, TUK } from '../-types'

export type StepScheduleProps = {
  mode: AssessmentMode | ''
  scheduledDate: string
  minDate: string
  tukId: string
  tuks: TUK[]
  tuksLoading: boolean
  errors: Record<string, string>
  onModeChange: (mode: AssessmentMode | '') => void
  onDateChange: (date: string) => void
  onTukChange: (tukId: string) => void
}

/** Step 4 — pilih jadwal asesmen (metode, tanggal, dan lokasi TUK bila offline). */
export function StepSchedule({
  mode,
  scheduledDate,
  minDate,
  tukId,
  tuks,
  tuksLoading,
  errors,
  onModeChange,
  onDateChange,
  onTukChange,
}: StepScheduleProps) {
  return (
    <div className="space-y-4">
      <Field label="Metode Pelaksanaan" htmlFor="assessment-mode" required error={errors.mode}>
        <Select
          id="assessment-mode"
          value={mode}
          onChange={(event) => onModeChange(event.target.value as AssessmentMode | '')}
        >
          <option value="">— Pilih metode —</option>
          {MODE_OPTIONS.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </Select>
      </Field>

      {mode === 'OFFLINE' ? (
        <Field
          label="Tempat Uji Kompetensi (TUK)"
          htmlFor="tuk"
          required
          error={errors.tuk_id}
          hint="Pilih lokasi pelaksanaan asesmen."
        >
          <Select id="tuk" value={tukId} disabled={tuksLoading} onChange={(event) => onTukChange(event.target.value)}>
            <option value="">{tuksLoading ? 'Memuat TUK...' : '— Pilih TUK —'}</option>
            {tuks.map((tuk) => (
              <option key={tuk.id} value={tuk.id}>{`${tuk.name} - ${tuk.province}`}</option>
            ))}
          </Select>
        </Field>
      ) : null}

      <Field
        label="Tanggal yang Diajukan"
        htmlFor="scheduled-date"
        required
        error={errors.scheduled_date}
        hint="Pilih tanggal rencana pelaksanaan asesmen."
      >
        <Input
          id="scheduled-date"
          type="date"
          value={scheduledDate}
          min={minDate}
          onChange={(event) => onDateChange(event.target.value)}
        />
      </Field>
    </div>
  )
}
