-- Skema sertifikasi (mis. 'Gada Pratama') beserta struktur unit kompetensi:
-- unit -> elemen -> kriteria. Dipakai sebagai master data penilaian.
-- Validasi (oleh role `validator`) di level skema; tiap aksi dicatat di riwayat.
-- Tiap tabel memakai trigger updated_at dan index untuk kolom relasi.
--
-- Catatan: validation_status NULL = belum divalidasi.
CREATE TYPE validation_status AS ENUM ('APPROVED', 'REJECTED');

CREATE TABLE certification_schemas (
    id                varchar(32)  PRIMARY KEY,          -- short uuid, digenerate di Go
    code              varchar(50)  NOT NULL,             -- slug internal, mis. 'gada_pratama'
    name              varchar(255) NOT NULL,             -- mis. 'Gada Pratama'
    schema_code       varchar(100),                      -- nomor resmi BNSP, mis. 'SS/OKU/04/GU/X/2019'
    sort_order        smallint     NOT NULL DEFAULT 0,
    is_active         boolean      NOT NULL DEFAULT true,
    
    validation_status validation_status,                 -- NULL = belum divalidasi
    validated_by      varchar(32),                       -- users(id) validator terakhir
    validated_at      timestamptz,
    validation_note   text,
    created_at        timestamptz  NOT NULL DEFAULT now(),
    updated_at        timestamptz  NOT NULL DEFAULT now(),
    CONSTRAINT uq_certification_schemas_code UNIQUE (code),
    CONSTRAINT uq_certification_schemas_name UNIQUE (name),
    CONSTRAINT fk_certification_schemas_validator FOREIGN KEY (validated_by) REFERENCES users (id)
);

CREATE TABLE competency_units (
    id                      varchar(32)  PRIMARY KEY,
    certification_schema_id varchar(32)  NOT NULL,
    code                    varchar(50)  NOT NULL,  -- mis. 'N.80PAM00.001.2' (persis dokumen resmi)
    title                   varchar(255) NOT NULL,
    sort_order              smallint     NOT NULL DEFAULT 0,
    created_at              timestamptz  NOT NULL DEFAULT now(),
    updated_at              timestamptz  NOT NULL DEFAULT now(),
    CONSTRAINT uq_competency_units_schema_code UNIQUE (certification_schema_id, code),
    CONSTRAINT fk_competency_units_schema FOREIGN KEY (certification_schema_id)
        REFERENCES certification_schemas (id) ON DELETE CASCADE
);

CREATE TABLE competency_unit_elements (
    id                 varchar(32) PRIMARY KEY,
    competency_unit_id varchar(32) NOT NULL,
    title              varchar(255),
    sort_order         smallint    NOT NULL DEFAULT 0,
    created_at         timestamptz NOT NULL DEFAULT now(),
    updated_at         timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT fk_competency_unit_elements_unit FOREIGN KEY (competency_unit_id)
        REFERENCES competency_units (id) ON DELETE CASCADE
);

CREATE TABLE competency_unit_criteria (
    id                         varchar(32) PRIMARY KEY,
    competency_unit_element_id varchar(32) NOT NULL,
    content                    text        NOT NULL,   -- teks kriteria unjuk kerja (KUK)
    sort_order                 smallint    NOT NULL DEFAULT 0,
    created_at                 timestamptz NOT NULL DEFAULT now(),
    updated_at                 timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT fk_competency_unit_criteria_element FOREIGN KEY (competency_unit_element_id)
        REFERENCES competency_unit_elements (id) ON DELETE CASCADE
);

CREATE TABLE certification_schema_validations (
    id           varchar(32)       PRIMARY KEY,
    schema_id    varchar(32)       NOT NULL,
    validator_id varchar(32)       NOT NULL,
    status_from  validation_status,                          -- NULL = belum divalidasi sebelumnya
    status_to    validation_status NOT NULL,
    note         text,
    created_at   timestamptz       NOT NULL DEFAULT now(),
    updated_at   timestamptz       NOT NULL DEFAULT now(),
    CONSTRAINT fk_certification_schema_validations_schema FOREIGN KEY (schema_id) REFERENCES certification_schemas (id) ON DELETE CASCADE,
    CONSTRAINT fk_certification_schema_validations_validator FOREIGN KEY (validator_id) REFERENCES users (id)
);

CREATE TRIGGER trg_certification_schemas_updated_at
BEFORE UPDATE ON certification_schemas
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_competency_units_updated_at
BEFORE UPDATE ON competency_units
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_competency_unit_elements_updated_at
BEFORE UPDATE ON competency_unit_elements
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_competency_unit_criteria_updated_at
BEFORE UPDATE ON competency_unit_criteria
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_certification_schema_validations_updated_at
BEFORE UPDATE ON certification_schema_validations
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_certification_schemas_sort_order ON certification_schemas (sort_order);
CREATE INDEX idx_certification_schemas_validation_status ON certification_schemas (validation_status);
CREATE INDEX idx_competency_units_schema ON competency_units (certification_schema_id, sort_order);
CREATE INDEX idx_competency_unit_elements_unit ON competency_unit_elements (competency_unit_id, sort_order);
CREATE INDEX idx_competency_unit_criteria_element ON competency_unit_criteria (competency_unit_element_id, sort_order);
CREATE INDEX idx_certification_schema_validations_schema ON certification_schema_validations (schema_id, created_at DESC);
