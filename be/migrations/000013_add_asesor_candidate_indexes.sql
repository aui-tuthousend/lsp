-- Index pendukung pemilihan asesor (plotting): menghitung jumlah penugasan
-- asesor pada tanggal jadwal tertentu, difilter per TUK.

-- Pencarian penugasan aktif per asesor pada tanggal jadwal.
CREATE INDEX idx_assessment_schedules_scheduled_date ON assessment_schedules (scheduled_date);

-- Penugasan aktif (PENDING/ACCEPTED) per asesor (dipakai untuk hitung kapasitas).
CREATE INDEX idx_asesor_assignments_active_asesor
    ON asesor_assignments (asesor_id, status)
    WHERE status IN ('PENDING', 'ACCEPTED');
