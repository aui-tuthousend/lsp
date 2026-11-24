// Tipe untuk master jenis dokumen (/admin/dokumen).
//
// Field mengikuti DTO backend modul documentrequirements (`DocumentTypeDTO`)
// dan body POST/PATCH `/document-types`.

/** Jenis dokumen seperti dikembalikan backend. */
export type DocumentType = {
  id: string
  code: string
  label: string
  created_at: string
  updated_at: string
}

/** Body pembuatan jenis dokumen (POST /document-types). */
export type CreateDocumentTypeBody = {
  code: string
  label: string
}

/** Body update parsial (PATCH /document-types/{id}). */
export type UpdateDocumentTypeBody = Partial<CreateDocumentTypeBody>

export type ListDocumentTypeParams = {
  page?: number
  limit?: number
}
