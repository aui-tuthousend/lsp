-- Data pekerjaan terkini asesi (relasi 1:1 ke asesi_profiles).
--
-- Dibuat sebagai tabel independen (bukan bagian dari permohonan sertifikasi)
-- agar bisa dipakai ulang pada setiap pengajuan: saat mengajukan, data ini
-- di-upsert (cari berdasarkan asesi, update bila ada / insert bila belum).
CREATE TABLE asesi_occupations (
    id                 varchar(32) PRIMARY KEY,           -- short uuid, digenerate di Go
    asesi_profile_id   varchar(32) NOT NULL,
    company_name       varchar(100) NOT NULL,             -- nama perusahaan/instansi
    position           varchar(100) NOT NULL,             -- jabatan
    tenure_months      integer,                           -- lama bekerja (bulan)
    subordinates_count integer,                           -- jumlah bawahan
    address            text         NOT NULL,             -- alamat tempat kerja
    province           varchar(100) NOT NULL,             -- provinsi
    city               varchar(100) NOT NULL,             -- kota/kabupaten
    district           varchar(100) NOT NULL,             -- kecamatan
    postal_code        varchar(10),                       -- kode pos
    office_email       citext,                            -- email kantor
    office_phone       varchar(20),                       -- no. telepon kantor
    created_at         timestamptz NOT NULL DEFAULT now(),
    updated_at         timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT fk_asesi_occupations_profile FOREIGN KEY (asesi_profile_id) REFERENCES asesi_profiles (id) ON DELETE CASCADE,
    CONSTRAINT uq_asesi_occupations_profile UNIQUE (asesi_profile_id),
    CONSTRAINT ck_asesi_occupations_tenure CHECK (tenure_months IS NULL OR tenure_months >= 0),
    CONSTRAINT ck_asesi_occupations_subordinates CHECK (subordinates_count IS NULL OR subordinates_count >= 0)
);

CREATE TRIGGER trg_asesi_occupations_updated_at
BEFORE UPDATE ON asesi_occupations
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_asesi_occupations_created_at ON asesi_occupations (created_at DESC);
