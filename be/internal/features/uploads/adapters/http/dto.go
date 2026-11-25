package http

type PresignInput struct {
	Body PresignBody `json:"body"`
}

// PresignBody meminta URL unggah untuk satu dokumen.
type PresignBody struct {
	Purpose        string  `json:"purpose" enum:"photo,ktp_file,cv_file,npwp_file,document"`
	ContentType    string  `json:"content_type" minLength:"1" maxLength:"100"`
	DocumentTypeID *string `json:"document_type_id,omitempty" doc:"Wajib saat purpose=document agar tipe berkas dibatasi sesuai jenis dokumen"`
}

type PresignOutput struct {
	Body PresignResponse `json:"body"`
}

type PresignResponse struct {
	ObjectKey string `json:"object_key"`
	UploadURL string `json:"upload_url"`
	ExpiresIn int    `json:"expires_in"`
}
