import { apiFetch } from './api'

/** Tujuan unggah, sesuai nama field dokumen di profil asesi. */
export type UploadPurpose = 'photo' | 'ktp_file' | 'cv_file' | 'npwp_file'

type PresignResponse = {
  object_key: string
  upload_url: string
  expires_in: number
}

/**
 * Unggah file langsung ke object storage lewat presigned URL, lalu kembalikan
 * object key-nya. Backend hanya menyimpan key ini (bucket privat).
 */
export async function uploadFile(file: File, purpose: UploadPurpose): Promise<string> {
  const presign = await apiFetch<PresignResponse>('/uploads/presign', null, {
    method: 'POST',
    body: JSON.stringify({ purpose, content_type: file.type }),
  })

  const res = await fetch(presign.upload_url, {
    method: 'PUT',
    headers: { 'Content-Type': file.type },
    body: file,
  })
  if (!res.ok) throw new Error(`Upload gagal (${res.status})`)

  return presign.object_key
}
