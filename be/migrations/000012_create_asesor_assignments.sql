-- Penugasan asesor untuk sebuah permohonan (plotting oleh admin LSP / manager
-- sertifikasi). Dipisah dari certification_applications agar menyimpan status
-- dan riwayat: asesor yang ditunjuk bisa menolak, lalu manager me-plot ulang
-- (penugasan lama ditandai CANCELLED, penugasan baru dibuat PENDING).

CREATE TYPE asesor_assignment_status AS ENUM ('PENDING', 'ACCEPTED', 'REJECTED', 'CANCELLED');

CREATE TABLE asesor_assignments (
    id             varchar(32)              PRIMARY KEY,
    application_id varchar(32)              NOT NULL,   -- permohonan APL-01 terkait
    asesor_id      varchar(32)              NOT NULL,   -- asesor yang ditunjuk
    status         asesor_assignment_status NOT NULL DEFAULT 'PENDING',
    assigned_by    varchar(32),                         -- manager/admin yang menugaskan
    assigned_at    timestamptz              NOT NULL DEFAULT now(),
    responded_at   timestamptz,                         -- saat asesor menerima/menolak
    note           text,                                -- catatan / alasan penolakan
    created_at     timestamptz              NOT NULL DEFAULT now(),
    updated_at     timestamptz              NOT NULL DEFAULT now(),
    CONSTRAINT fk_asesor_assignments_application FOREIGN KEY (application_id) REFERENCES certification_applications (id) ON DELETE CASCADE,
    CONSTRAINT fk_asesor_assignments_asesor      FOREIGN KEY (asesor_id)      REFERENCES users (id) ON DELETE RESTRICT,
    CONSTRAINT fk_asesor_assignments_assigner    FOREIGN KEY (assigned_by)    REFERENCES users (id) ON DELETE SET NULL
);

-- Hanya boleh ada satu penugasan aktif (PENDING/ACCEPTED) per permohonan.
CREATE UNIQUE INDEX uq_asesor_assignments_active
    ON asesor_assignments (application_id)
    WHERE status IN ('PENDING', 'ACCEPTED');

-- Daftar penugasan per asesor.
CREATE INDEX idx_asesor_assignments_asesor ON asesor_assignments (asesor_id);

CREATE TRIGGER trg_asesor_assignments_updated_at
BEFORE UPDATE ON asesor_assignments
FOR EACH ROW EXECUTE FUNCTION set_updated_at();
