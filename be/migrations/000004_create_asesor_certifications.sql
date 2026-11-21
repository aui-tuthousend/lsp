-- Daftar sertifikasi milik seorang asesor (relasi 1:N ke asesor_profiles).
-- Tiap baris adalah record mandiri: nama sertifikasi/skema, lembaga penerbit,
-- nomor sertifikat, serta masa berlaku (opsional). Menghapus profil asesor
-- akan ikut menghapus seluruh sertifikasinya (ON DELETE CASCADE).
CREATE TABLE asesor_certifications (
    id         varchar(32)  PRIMARY KEY,          -- short uuid, digenerate di Go
    user_id    varchar(32)  NOT NULL,
    name       varchar(255) NOT NULL,             -- nama sertifikasi/skema
    issuer     varchar(255) NOT NULL,             -- lembaga penerbit
    number     varchar(120),                      -- nomor sertifikat
    issued_at  date,                              -- tanggal terbit
    expired_at date,                              -- tanggal kadaluarsa
    created_at timestamptz  NOT NULL DEFAULT now(),
    updated_at timestamptz  NOT NULL DEFAULT now(),
    CONSTRAINT fk_asesor_certifications_profile FOREIGN KEY (user_id) REFERENCES asesor_profiles (user_id) ON DELETE CASCADE,
    CONSTRAINT ck_asesor_certifications_period  CHECK (expired_at IS NULL OR issued_at IS NULL OR expired_at >= issued_at)
);

CREATE TRIGGER trg_asesor_certifications_updated_at
BEFORE UPDATE ON asesor_certifications
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_asesor_certifications_user ON asesor_certifications (user_id);
