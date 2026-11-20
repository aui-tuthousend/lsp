CREATE EXTENSION IF NOT EXISTS citext;

-- Master data role. Baris role diisi oleh seeder (cmd/seed), id digenerate di Go.
CREATE TABLE roles (
    id          varchar(32)  PRIMARY KEY,          -- short uuid, digenerate di Go
    code        varchar(50)  NOT NULL,             -- mis. 'admin', 'manager_sertifikasi'
    name        varchar(100) NOT NULL,             -- mis. 'Manager Sertifikasi'
    description text,
    created_at  timestamptz  NOT NULL DEFAULT now(),
    updated_at  timestamptz  NOT NULL DEFAULT now(),
    CONSTRAINT uq_roles_code UNIQUE (code)
);

CREATE TABLE users (
    id                varchar(32)  PRIMARY KEY,                 -- short uuid, digenerate di Go
    role_id           varchar(32)  NOT NULL,
    isvalid           boolean      not null default true,
    name              varchar(255) NOT NULL,
    email             citext       NOT NULL,
    phone             varchar(20),
    nomor_registrasi  varchar(50),
    password_hash     text         NOT NULL,
    email_verified_at timestamptz,
    remember_token    varchar(100),
    created_at        timestamptz  NOT NULL DEFAULT now(),
    updated_at        timestamptz  NOT NULL DEFAULT now(),
    deleted_at        timestamptz,
    CONSTRAINT fk_users_role             FOREIGN KEY (role_id) REFERENCES roles (id),
    CONSTRAINT uq_users_email            UNIQUE (email),
    CONSTRAINT uq_users_nomor_registrasi UNIQUE (nomor_registrasi)
);

CREATE TABLE outbox_events (
    id             varchar(32) PRIMARY KEY,
    event_type     varchar(120) NOT NULL,
    aggregate_type varchar(80) NOT NULL,
    aggregate_id   varchar(32) NOT NULL,
    payload        jsonb NOT NULL,
    status         varchar(20) NOT NULL DEFAULT 'pending',
    attempts       integer NOT NULL DEFAULT 0,
    available_at   timestamptz NOT NULL DEFAULT now(),
    processed_at   timestamptz,
    error          text,
    created_at     timestamptz NOT NULL DEFAULT now(),
    updated_at     timestamptz NOT NULL DEFAULT now()
);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_roles_updated_at
BEFORE UPDATE ON roles
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER trg_outbox_events_updated_at
BEFORE UPDATE ON outbox_events
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_users_role_created_at_not_deleted ON users (role_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at_not_deleted ON users (created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX idx_outbox_events_pending ON outbox_events (available_at, created_at) WHERE status = 'pending';
