-- Profil data diri asesi (relasi 1:1 ke users). Menyimpan identitas lengkap
-- peserta sertifikasi: data pribadi, domisili, dan berkas pendukung (pas foto,
-- KTP, CV, NPWP). Berkas disimpan sebagai URL/path hasil upload ke object storage.
--
-- Catatan: tipe enum dibuat lebih dulu karena dipakai sebagai tipe kolom.
CREATE TYPE gender AS ENUM ('LAKI-LAKI', 'PEREMPUAN');

CREATE TYPE education_level AS ENUM (
    'SD', 'SMP', 'SMA_SMK', 'D1', 'D2', 'D3', 'D4', 'S1', 'S2', 'S3'
);

CREATE TABLE asesi_profiles (
    id                   varchar(32) PRIMARY KEY,            -- short uuid, digenerate di Go
    user_id              varchar(32) NOT NULL,
    full_name            varchar(255) NOT NULL,
    nik                  varchar(16)  NOT NULL,
    phone                varchar(20)  NOT NULL,
    gender               gender       NOT NULL,
    birth_place          varchar(255) NOT NULL,
    birth_date           date         NOT NULL,
    nationality          varchar(100) NOT NULL DEFAULT 'Indonesia',
    occupation           varchar(255) NOT NULL,
    education_level      education_level NOT NULL,

    ktp_address          text         NOT NULL,
    domicile_same_as_ktp boolean      NOT NULL DEFAULT false,
    domicile_address     text,                               -- wajib bila domisili beda dari KTP
    province             varchar(255) NOT NULL,
    city                 varchar(255) NOT NULL,
    district             varchar(255) NOT NULL,
    village              varchar(255) NOT NULL,
    postal_code          varchar(10),

    photo                varchar(255),
    ktp_file             varchar(255),
    cv_file              varchar(255),
    npwp_file            varchar(255),

    created_at           timestamptz  NOT NULL DEFAULT now(),
    updated_at           timestamptz  NOT NULL DEFAULT now(),
    CONSTRAINT fk_asesi_profiles_user     FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT uq_asesi_profiles_user     UNIQUE (user_id),
    CONSTRAINT uq_asesi_profiles_nik      UNIQUE (nik),
    CONSTRAINT ck_asesi_profiles_domicile CHECK (domicile_same_as_ktp OR domicile_address IS NOT NULL)
);

CREATE TRIGGER trg_asesi_profiles_updated_at
BEFORE UPDATE ON asesi_profiles
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_asesi_profiles_created_at ON asesi_profiles (created_at DESC);
