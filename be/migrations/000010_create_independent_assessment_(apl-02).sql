-- FR-APL-02: Asesmen Mandiri (self-assessment).
--
-- Struktur ini diadaptasi dari aplikasi lama agar sesuai alur baru:
--   * merujuk permohonan APL-01 (certification_applications), bukan tabel lama.
--   * token QR TIDAK disimpan plaintext di sini; verifikasi QR memakai tabel
--     bersama form_signatures (entity 'asesmen_mandiri', role asesi/asesor).
--   * jawaban per UNIT kompetensi (boolean kompeten/tidak) + banyak berkas
--     bukti per unit (dikirim saat asesi menyatakan kompeten).
--
-- Asesmen mandiri bersifat 1:1 dengan permohonan dan dibuat "lazy" saat asesi
-- pertama kali membuka laman APL-02 (setelah permohonan berstatus SUBMITTED).

-- Status siklus hidup asesmen mandiri.
CREATE TYPE independent_assessment_status AS ENUM ('DRAFT', 'SUBMITTED', 'REVIEWED');

-- Rekomendasi/keputusan asesor atas asesmen mandiri.
CREATE TYPE asesor_decision AS ENUM ('RECOMMENDED', 'NOT_RECOMMENDED');

CREATE TABLE independent_assessments (
    id               varchar(32)                  PRIMARY KEY,  -- short uuid, digenerate di Go
    application_id   varchar(32)                  NOT NULL,     -- permohonan APL-01 terkait
    code             varchar(50)                  NOT NULL,     -- kode asesmen, mis. 'APL02-2025-0001'
    status           independent_assessment_status NOT NULL DEFAULT 'DRAFT',
    asesor_note      text,                                       -- catatan tinjauan asesor
    asesor_decision  asesor_decision,                            -- rekomendasi asesor
    submitted_at     timestamptz,                                -- saat asesi submit APL-02
    reviewed_by      varchar(32),                                -- asesor peninjau
    reviewed_at      timestamptz,
    zoom_link        varchar(255),                               -- link pertemuan (dikirim asesor)
    zoom_date        date,                                       -- jadwal pertemuan yang disepakati
    zoom_time        varchar(5),                                 -- jam pertemuan 'HH:MM'
    created_at       timestamptz                  NOT NULL DEFAULT now(),
    updated_at       timestamptz                  NOT NULL DEFAULT now(),
    CONSTRAINT uq_independent_assessments_application UNIQUE (application_id),
    CONSTRAINT uq_independent_assessments_code        UNIQUE (code),
    CONSTRAINT fk_independent_assessments_application FOREIGN KEY (application_id)
        REFERENCES certification_applications (id) ON DELETE CASCADE,
    CONSTRAINT fk_independent_assessments_reviewer    FOREIGN KEY (reviewed_by)
        REFERENCES users (id) ON DELETE SET NULL
);

-- Jawaban asesi per unit kompetensi (kerangka digenerate dari unit milik skema).
CREATE TABLE independent_assessment_units (
    id                 varchar(32) PRIMARY KEY,
    assessment_id      varchar(32) NOT NULL,
    competency_unit_id varchar(32) NOT NULL,
    is_competent       boolean,                              -- pernyataan mandiri asesi (null = belum dijawab)
    created_at         timestamptz NOT NULL DEFAULT now(),
    updated_at         timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT uq_independent_assessment_units UNIQUE (assessment_id, competency_unit_id),
    CONSTRAINT fk_independent_assessment_units_assessment FOREIGN KEY (assessment_id)
        REFERENCES independent_assessments (id) ON DELETE CASCADE,
    CONSTRAINT fk_independent_assessment_units_unit FOREIGN KEY (competency_unit_id)
        REFERENCES competency_units (id) ON DELETE RESTRICT
);

-- Berkas bukti per jawaban unit (banyak baris per unit).
CREATE TABLE independent_assessment_unit_files (
    id                 varchar(32)  PRIMARY KEY,
    assessment_unit_id varchar(32)  NOT NULL,
    file_url           varchar(255) NOT NULL,                -- URL/object key hasil upload (R2)
    created_at         timestamptz  NOT NULL DEFAULT now(),
    CONSTRAINT fk_independent_assessment_unit_files_unit FOREIGN KEY (assessment_unit_id)
        REFERENCES independent_assessment_units (id) ON DELETE CASCADE
);

CREATE TRIGGER trg_independent_assessments_updated_at
BEFORE UPDATE ON independent_assessments
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_independent_assessment_units_updated_at
BEFORE UPDATE ON independent_assessment_units
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_independent_assessments_status ON independent_assessments (status);
CREATE INDEX idx_independent_assessment_units_assessment ON independent_assessment_units (assessment_id);
CREATE INDEX idx_independent_assessment_unit_files_unit ON independent_assessment_unit_files (assessment_unit_id);
