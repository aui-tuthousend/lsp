-- Token verifikasi QR (disimpan sebagai hash, bukan plaintext).
--
-- Dipisah dari migrasi permohonan APL-01 karena dipakai lintas tahap
-- (APL.01, APL.02, AK.01, ...). Satu token aktif per (permohonan, formulir,
-- peran); riwayat/re-issue tetap tersimpan.

-- Formulir yang ditandatangani/diverifikasi lewat QR.
CREATE TYPE form_entitas AS ENUM (
    'permohonan',            -- FR.APL.01
    'asesmen_mandiri',       -- FR.APL.02
    'persetujuan_asesmen',
    'pra_asesmen',           -- FR.AK.01
    'proses_asesmen',        -- ujian
    'rekaman_asesmen',       -- FR.AK.02
    'umpan_balik'
);

-- Pihak yang melakukan verifikasi QR.
CREATE TYPE signer_role AS ENUM ('asesi', 'asesor', 'asesor_ulang');

-- Status sebuah token verifikasi.
CREATE TYPE signature_status AS ENUM ('menunggu', 'dipakai', 'kadaluarsa', 'dibatalkan');

CREATE TABLE form_signatures (
    id              varchar(32)      PRIMARY KEY,
    application_id  varchar(32)      NOT NULL,
    entity          form_entitas     NOT NULL,             -- formulir yang diverifikasi
    role            signer_role      NOT NULL,             -- pihak penandatangan
    token_hash      varchar(64)      NOT NULL,             -- sha256(token), BUKAN plaintext
    status          signature_status NOT NULL DEFAULT 'menunggu',
    expires_at      timestamptz      NOT NULL,             -- masa berlaku QR
    used_at         timestamptz,                           -- kapan di-scan (single-use)
    used_ip         inet,                                  -- audit
    used_user_agent text,                                  -- audit
    revoked_at      timestamptz,                           -- pembatalan (re-issue)
    created_at      timestamptz      NOT NULL DEFAULT now(),
    updated_at      timestamptz      NOT NULL DEFAULT now(),
    CONSTRAINT uq_form_signatures_token UNIQUE (token_hash),
    CONSTRAINT fk_form_signatures_application FOREIGN KEY (application_id)
        REFERENCES certification_applications (id) ON DELETE CASCADE,
    CONSTRAINT ck_form_signatures_used  CHECK (status <> 'dipakai'    OR used_at    IS NOT NULL),
    CONSTRAINT ck_form_signatures_batal CHECK (status <> 'dibatalkan' OR revoked_at IS NOT NULL)
);

CREATE TRIGGER trg_form_signatures_updated_at
BEFORE UPDATE ON form_signatures
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Hanya SATU token aktif per (permohonan, formulir, peran).
CREATE UNIQUE INDEX uq_form_signatures_aktif
    ON form_signatures (application_id, entity, role)
    WHERE status = 'menunggu';

-- Lookup utama saat scan/verifikasi.
CREATE INDEX idx_form_signatures_status_expires ON form_signatures (status, expires_at);
-- (token_hash sudah ter-index via UNIQUE)
