-- Master jenis dokumen (menggantikan string bebas 'ktp','ijazah',...).
-- Baris diisi oleh seeder (cmd/seed), id digenerate di Go.

CREATE TABLE document_types (
    id                 varchar(32)  PRIMARY KEY,          -- short uuid, digenerate di Go
    code               varchar(50)  NOT NULL,             -- 'ktp','ijazah','sertifikat_pelatihan',...
    label              varchar(255) NOT NULL,
    -- MIME yang diizinkan untuk diunggah (subset allowlist uploads purpose 'document').
    allowed_mime_types varchar[]    NOT NULL DEFAULT ARRAY['application/pdf', 'image/jpeg', 'image/png']::varchar[],
    created_at         timestamptz  NOT NULL DEFAULT now(),
    updated_at         timestamptz  NOT NULL DEFAULT now(),
    CONSTRAINT uq_document_types_code UNIQUE (code),
    CONSTRAINT ck_document_types_allowed_mime_types CHECK (array_length(allowed_mime_types, 1) >= 1)
);

-- Dokumen wajib per skema (inilah penentu "dokumen relevan" APL.02).
CREATE TABLE schema_document_requirements (
    id                      varchar(32) PRIMARY KEY,          -- short uuid, digenerate di Go
    certification_schema_id varchar(32) NOT NULL,
    document_type_id        varchar(32) NOT NULL,
    is_required             boolean     NOT NULL DEFAULT true,
    created_at              timestamptz NOT NULL DEFAULT now(),
    updated_at              timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT fk_schema_document_requirements_schema FOREIGN KEY (certification_schema_id)
        REFERENCES certification_schemas (id) ON DELETE CASCADE,
    CONSTRAINT fk_schema_document_requirements_type FOREIGN KEY (document_type_id)
        REFERENCES document_types (id) ON DELETE RESTRICT,
    CONSTRAINT uq_schema_document_requirements UNIQUE (certification_schema_id, document_type_id)
);

CREATE TRIGGER trg_document_types_updated_at
BEFORE UPDATE ON document_types
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_schema_document_requirements_updated_at
BEFORE UPDATE ON schema_document_requirements
FOR EACH ROW EXECUTE FUNCTION set_updated_at();
