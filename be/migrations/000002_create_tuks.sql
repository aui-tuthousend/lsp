-- Master data TUK (Tempat Uji Kompetensi).
CREATE TABLE tuks (
    id             varchar(32)  PRIMARY KEY,          -- short uuid, digenerate di Go
    name           varchar(255) NOT NULL,
    logo           text,
    postal_code    varchar(10)  NOT NULL,
    phone          varchar(20)  NOT NULL,
    email          citext       NOT NULL,
    province       varchar(120) NOT NULL,
    city           varchar(120) NOT NULL,
    address        text         NOT NULL,
    maps_url       text         NOT NULL,
    letterhead     text,
    created_at     timestamptz  NOT NULL DEFAULT now(),
    updated_at     timestamptz  NOT NULL DEFAULT now(),
    deleted_at     timestamptz,
    CONSTRAINT uq_tuks_email UNIQUE (email)
);

CREATE TRIGGER trg_tuks_updated_at
BEFORE UPDATE ON tuks
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_tuks_created_at_not_deleted ON tuks (created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_tuks_name_not_deleted ON tuks (name) WHERE deleted_at IS NULL;
