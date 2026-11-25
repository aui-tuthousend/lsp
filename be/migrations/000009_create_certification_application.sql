-- FR-APL-01: Pendaftaran & Permohonan Sertifikasi Kompetensi.
--
-- Ruang lingkup migrasi ini dibatasi sampai tahap APL-01:
--   user mendaftar -> mengisi data pribadi (asesi_profiles) & data pekerjaan
--   (asesi_occupations) -> memilih skema + tujuan asesmen -> melampirkan bukti
--   kelengkapan sesuai skema -> mengajukan jadwal asesmen.
--
-- Tabel verifikasi QR (form_signatures) dipisah ke migrasi tersendiri
-- (000010) karena dipakai lintas tahap, bukan hanya APL-01.

-- Status siklus hidup permohonan (FR-APL01-03 & FR-APL01-04).
CREATE TYPE application_status AS ENUM ('DRAFT', 'SUBMITTED', 'APPROVED', 'REJECTED');

-- Tujuan asesmen (istilah regulasi BNSP).
CREATE TYPE assessment_purpose AS ENUM (
    'SERTIFIKASI', 'SERTIFIKASI_ULANG', 'REKOGNISI_PEMBELAJARAN_LAMPAU', 'PENGAKUAN_KOMPETENSI_TERKINI'
);

-- Metode pelaksanaan asesmen.
CREATE TYPE assessment_mode AS ENUM ('ONLINE', 'OFFLINE');

-- Status kelengkapan tiap bukti dasar (FR-APL01-02).
CREATE TYPE document_status AS ENUM ('BELUM_LENGKAP', 'MEMENUHI_SYARAT', 'TIDAK_MEMENUHI_SYARAT');

-- Header permohonan APL-01 (1 baris = 1 pengajuan asesi).
CREATE TABLE certification_applications (
    id                 varchar(32)        PRIMARY KEY,          -- short uuid, digenerate di Go
    user_id            varchar(32)        NOT NULL,             -- asesi pemohon
    code               varchar(50)        NOT NULL,             -- kode permohonan, mis. 'APL01-2025-0001'
    status             application_status NOT NULL DEFAULT 'DRAFT',
    current_step       smallint           NOT NULL DEFAULT 1,   -- step terakhir untuk resume draft (1..4)
    submitted_at       timestamptz,                             -- saat asesi submit (DRAFT -> SUBMITTED)
    verified_by        varchar(32),                             -- admin LSP yang memverifikasi
    verified_at        timestamptz,
    verification_note  text,                                    -- catatan verifikasi (APPROVED/REJECTED)
    asesor_id          varchar(32),                             -- asesor hasil plotting
    asesor_assigned_at timestamptz,
    created_at         timestamptz        NOT NULL DEFAULT now(),
    updated_at         timestamptz        NOT NULL DEFAULT now(),
    deleted_at         timestamptz,
    CONSTRAINT uq_certification_applications_code     UNIQUE (code),
    CONSTRAINT fk_certification_applications_user     FOREIGN KEY (user_id)     REFERENCES users (id) ON DELETE RESTRICT,
    CONSTRAINT fk_certification_applications_verifier FOREIGN KEY (verified_by) REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT fk_certification_applications_asesor   FOREIGN KEY (asesor_id)   REFERENCES users (id) ON DELETE SET NULL,
    CONSTRAINT ck_certification_applications_step CHECK (current_step BETWEEN 1 AND 4)
);

-- Data sertifikasi: skema yang dipilih + tujuan asesmen (1:1 ke permohonan).
CREATE TABLE certification_application_data (
    id                      varchar(32)        PRIMARY KEY,
    application_id          varchar(32)        NOT NULL,
    certification_schema_id varchar(32)        NOT NULL,
    purpose                 assessment_purpose NOT NULL,
    created_at              timestamptz        NOT NULL DEFAULT now(),
    updated_at              timestamptz        NOT NULL DEFAULT now(),
    CONSTRAINT uq_certification_application_data_application UNIQUE (application_id),
    CONSTRAINT fk_certification_application_data_application FOREIGN KEY (application_id)
        REFERENCES certification_applications (id) ON DELETE CASCADE,
    CONSTRAINT fk_certification_application_data_schema FOREIGN KEY (certification_schema_id)
        REFERENCES certification_schemas (id) ON DELETE RESTRICT
);

-- Bukti kelengkapan dasar per jenis dokumen (FR-APL01-02).
-- Jenis dokumen yang relevan ditentukan oleh schema_document_requirements.
CREATE TABLE application_documents (
    id               varchar(32)     PRIMARY KEY,
    application_id   varchar(32)     NOT NULL,
    document_type_id varchar(32)     NOT NULL,
    file_url         varchar(255),                          -- URL hasil client-side upload (R2)
    status           document_status NOT NULL DEFAULT 'BELUM_LENGKAP',
    note             text,                                  -- catatan verifikasi admin
    verified_at      timestamptz,
    created_at       timestamptz     NOT NULL DEFAULT now(),
    updated_at       timestamptz     NOT NULL DEFAULT now(),
    CONSTRAINT uq_application_documents UNIQUE (application_id, document_type_id),
    CONSTRAINT fk_application_documents_application FOREIGN KEY (application_id)
        REFERENCES certification_applications (id) ON DELETE CASCADE,
    CONSTRAINT fk_application_documents_type FOREIGN KEY (document_type_id)
        REFERENCES document_types (id) ON DELETE RESTRICT
);

-- Jadwal asesmen yang diajukan asesi (1:1 ke permohonan).
CREATE TABLE assessment_schedules (
    id             varchar(32)     PRIMARY KEY,
    application_id varchar(32)     NOT NULL,
    mode           assessment_mode NOT NULL,             -- online/offline
    scheduled_date date            NOT NULL,             -- tanggal yang diajukan asesi
    tuk_id         varchar(32),                          -- TUK mandiri (opsional; kosong bila online)
    created_at     timestamptz     NOT NULL DEFAULT now(),
    updated_at     timestamptz     NOT NULL DEFAULT now(),
    CONSTRAINT uq_assessment_schedules_application UNIQUE (application_id),
    CONSTRAINT fk_assessment_schedules_application FOREIGN KEY (application_id)
        REFERENCES certification_applications (id) ON DELETE CASCADE,
    CONSTRAINT fk_assessment_schedules_tuk FOREIGN KEY (tuk_id)
        REFERENCES tuks (id) ON DELETE SET NULL
);

CREATE TRIGGER trg_certification_applications_updated_at
BEFORE UPDATE ON certification_applications
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_certification_application_data_updated_at
BEFORE UPDATE ON certification_application_data
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_application_documents_updated_at
BEFORE UPDATE ON application_documents
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_assessment_schedules_updated_at
BEFORE UPDATE ON assessment_schedules
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Daftar permohonan milik asesi (beranda / resume draft) — hanya yang belum dihapus.
CREATE INDEX idx_certification_applications_user ON certification_applications (user_id, created_at DESC) WHERE deleted_at IS NULL;
-- Daftar verifikasi admin per status.
CREATE INDEX idx_certification_applications_status ON certification_applications (status, created_at DESC) WHERE deleted_at IS NULL;
-- Daftar penugasan per asesor.
CREATE INDEX idx_certification_applications_asesor ON certification_applications (asesor_id) WHERE asesor_id IS NOT NULL;
