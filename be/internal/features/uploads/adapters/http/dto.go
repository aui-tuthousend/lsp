package http

type PresignInput struct {
	Body PresignBody `json:"body"`
}

// PresignBody meminta URL unggah untuk satu dokumen.
type PresignBody struct {
	Purpose     string `json:"purpose" enum:"photo,ktp_file,cv_file,npwp_file"`
	NIK         string `json:"nik" pattern:"^[0-9]{16}$"`
	ContentType string `json:"content_type" minLength:"1" maxLength:"100"`
}

type PresignOutput struct {
	Body PresignResponse `json:"body"`
}

type PresignResponse struct {
	ObjectKey string `json:"object_key"`
	UploadURL string `json:"upload_url"`
	ExpiresIn int    `json:"expires_in"`
}
