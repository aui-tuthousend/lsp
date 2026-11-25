import { Field } from '@/components/ui/field'
import { Select } from '@/components/ui/select'
import { PURPOSE_OPTIONS } from '../-helpers'
import type { AssessmentPurpose, CertificationSchema } from '../-types'

export type StepCertificationProps = {
  schemas: CertificationSchema[]
  loading: boolean
  schemaId: string
  purpose: AssessmentPurpose | ''
  errors: Record<string, string>
  onSchemaChange: (id: string) => void
  onPurposeChange: (purpose: AssessmentPurpose | '') => void
}

/** Step 2 — pilih skema sertifikasi + tujuan asesmen. */
export function StepCertification({
  schemas,
  loading,
  schemaId,
  purpose,
  errors,
  onSchemaChange,
  onPurposeChange,
}: StepCertificationProps) {
  return (
    <div className="space-y-4">
      <Field label="Skema Sertifikasi" htmlFor="certification-schema" required error={errors.schema_id}>
        <Select
          id="certification-schema"
          value={schemaId}
          disabled={loading}
          onChange={(event) => onSchemaChange(event.target.value)}
        >
          <option value="">{loading ? 'Memuat skema...' : '— Pilih skema —'}</option>
          {schemas.map((schema) => (
            <option key={schema.id} value={schema.id}>
              {schema.name}
            </option>
          ))}
        </Select>
      </Field>

      <Field label="Tujuan Asesmen" htmlFor="assessment-purpose" required error={errors.purpose}>
        <Select
          id="assessment-purpose"
          value={purpose}
          onChange={(event) => onPurposeChange(event.target.value as AssessmentPurpose | '')}
        >
          <option value="">— Pilih tujuan —</option>
          {PURPOSE_OPTIONS.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </Select>
      </Field>
    </div>
  )
}
