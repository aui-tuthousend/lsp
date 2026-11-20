-- Profil khusus asesor (relasi 1:1 ke users). Dibuat setelah migrasi tuks
-- karena menyimpan afiliasi asesor ke sebuah TUK, disertai domisili
-- (provinsi/kota) dan batas penugasan harian. Pemakaian harian dihitung
-- dari tabel penugasan (derived), bukan disimpan di sini.
CREATE TABLE asesor_profiles (
    user_id                varchar(32) PRIMARY KEY,
    tuk_id                 varchar(32),
    province               varchar(120),
    city                   varchar(120),
    daily_assignment_limit integer     NOT NULL DEFAULT 0,   -- 0 = tanpa batas
    created_at             timestamptz NOT NULL DEFAULT now(),
    updated_at             timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT fk_asesor_profiles_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_asesor_profiles_tuk  FOREIGN KEY (tuk_id)  REFERENCES tuks (id),
    CONSTRAINT ck_asesor_profiles_limit CHECK (daily_assignment_limit >= 0)
);

CREATE TRIGGER trg_asesor_profiles_updated_at
BEFORE UPDATE ON asesor_profiles
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE INDEX idx_asesor_profiles_tuk ON asesor_profiles (tuk_id) WHERE tuk_id IS NOT NULL;
