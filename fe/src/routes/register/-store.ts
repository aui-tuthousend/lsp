// Draft pendaftaran asesi: state form + langkah aktif yang dipersist ke
// localStorage agar tidak hilang saat refresh / tab ditutup.
//
// Catatan keamanan & keterbatasan:
// - `password`/`password_confirmation` dan field file TIDAK pernah disimpan.
// - Karena pakai SSR, hidrasi ditunda (`skipHydration`) lalu dijalankan di
//   client agar HTML server dan render pertama client tetap sama.

import { create } from 'zustand'
import { createJSONStorage, persist } from 'zustand/middleware'
import type { CreateAsesiProfiles } from './-types'

const INITIAL_VALUES: CreateAsesiProfiles = {
  full_name: '',
  nik: '',
  phone: '',
  gender: '',
  birth_place: '',
  birth_date: '',
  nationality: 'Indonesia',
  occupation: '',
  education_level: '',
  ktp_address: '',
  domicile_same_as_ktp: true,
  domicile_address: '',
  postal_code: '',
  province: '',
  city: '',
  district: '',
  village: '',
  photo: null,
  ktp_file: null,
  cv_file: null,
  npwp_file: null,
  email: '',
  password: '',
  password_confirmation: '',
}

type RegisterDraftState = {
  step: number
  values: CreateAsesiProfiles
  setStep: (step: number) => void
  update: <K extends keyof CreateAsesiProfiles>(key: K, value: CreateAsesiProfiles[K]) => void
  patch: (partial: Partial<CreateAsesiProfiles>) => void
  reset: () => void
}

export const useRegisterDraft = create<RegisterDraftState>()(
  persist(
    (set) => ({
      step: 0,
      values: INITIAL_VALUES,
      setStep: (step) => set({ step }),
      update: (key, value) => set((state) => ({ values: { ...state.values, [key]: value } })),
      patch: (partial) => set((state) => ({ values: { ...state.values, ...partial } })),
      reset: () => set({ step: 0, values: INITIAL_VALUES }),
    }),
    {
      name: 'lsp-register-draft',
      version: 1,
      skipHydration: true,
      storage: createJSONStorage(() => localStorage),
      // Hanya field teks yang disimpan; password & file dikosongkan.
      partialize: (state) => ({
        step: state.step,
        values: {
          ...state.values,
          password: '',
          password_confirmation: '',
          photo: null,
          ktp_file: null,
          cv_file: null,
          npwp_file: null,
        },
      }),
    },
  ),
)
