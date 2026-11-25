package domain

import "time"

// FormEntity menandai formulir mana yang ditandatangani/diverifikasi.
type FormEntity string

const (
	EntityPermohonan         FormEntity = "permohonan"          // FR.APL.01
	EntityAsesmenMandiri     FormEntity = "asesmen_mandiri"     // FR.APL.02
	EntityPersetujuanAsesmen FormEntity = "persetujuan_asesmen" // FR.AK.01
	EntityPraAsesmen         FormEntity = "pra_asesmen"         // FR.AK.01
	EntityProsesAsesmen      FormEntity = "proses_asesmen"      // ujian
	EntityRekamanAsesmen     FormEntity = "rekaman_asesmen"     // FR.AK.02
	EntityUmpanBalik         FormEntity = "umpan_balik"
)

// SignerRole adalah pihak yang melakukan verifikasi QR.
type SignerRole string

const (
	SignerAsesi       SignerRole = "asesi"
	SignerAsesor      SignerRole = "asesor"
	SignerAsesorUlang SignerRole = "asesor_ulang"
)

// SignatureStatus adalah status sebuah token verifikasi.
type SignatureStatus string

const (
	SignaturePending SignatureStatus = "menunggu"
	SignatureUsed    SignatureStatus = "dipakai"
	SignatureExpired SignatureStatus = "kadaluarsa"
	SignatureRevoked SignatureStatus = "dibatalkan"
)

// FormSignature adalah token verifikasi QR. Token disimpan sebagai hash
// (sha256), bukan plaintext.
type FormSignature struct {
	ID            string
	ApplicationID string
	Entity        FormEntity
	Role          SignerRole
	TokenHash     string
	Status        SignatureStatus
	ExpiresAt     time.Time
	UsedAt        *time.Time
	UsedIP        *string
	UsedUserAgent *string
	RevokedAt     *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
