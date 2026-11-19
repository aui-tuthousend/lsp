# Product Requirements Document (PRD): Sistem Informasi Sertifikasi Kompetensi (APL-01 & APL-02)

## 1. Ringkasan Produk
Aplikasi backend REST API berbasis Go untuk mengelola alur sertifikasi kompetensi profesi pada Lembaga Sertifikasi Profesi (LSP), mencakup pendaftaran permohonan sertifikasi (**FR-APL-01**) dan asesmen mandiri berbasis unit kompetensi & Elemen/KUK (**FR-APL-02**), lengkap dengan verifikasi oleh Admin LSP, rekomendasi oleh Asesor Kompetensi, penyimpanan dokumen portofolio di Cloudflare R2, serta notifikasi email otomatis melalui Mailgun.

## 2. Aktor & Hak Akses (RBAC)
1. **Asesi (`ASESI`):** Peserta sertifikasi yang mengajukan permohonan (APL-01), mengunggah bukti kelengkapan ke cloud storage, dan mengisi asesmen mandiri (APL-02).
2. **Admin LSP (`ADMIN`):** Memverifikasi kelengkapan dokumen dasar APL-01, menetapkan status kelayakan pendaftaran, dan memetakan Asesor ke Asesi.
3. **Asesor Kompetensi (`ASESOR`):** Meninjau jawaban asesmen mandiri (APL-02) serta bukti portofolio per Elemen/KUK, lalu memberikan keputusan rekomendasi (`LANJUT` / `TIDAK_LANJUT`).
4. **Super Admin (`SUPERADMIN`):** Mengelola master data Skema Sertifikasi, Unit Kompetensi, Elemen, Kriteria Unjuk Kerja (KUK), serta manajemen akun pengguna.

## 3. Core Functional Requirements (Bounded Contexts)

### 3.1. Master Skema & Unit Kompetensi (`schema`)
* **FR-SCH-01:** CRUD Skema Sertifikasi (Kode Skema, Judul, Jenis: KKNI/Okupasi/Klaster, Status Aktif).
* **FR-SCH-02:** Manajemen hierarki skema: `Skema` -> `Unit Kompetensi` -> `Elemen Kompetensi` -> `Kriteria Unjuk Kerja (KUK)`.
* **FR-SCH-03:** Endpoint pembacaan struktur lengkap skema (di-cache di Redis karena jarang berubah namun sering dibaca saat pengisian APL-02).

### 3.2. Pendaftaran & Permohonan Sertifikasi — FR-APL-01 (`apl01`)
* **FR-APL01-01:** Asesi membuat pengajuan APL-01 pada satu Skema Sertifikasi aktif (Data Pribadi, Data Pekerjaan/Pendidikan, Tujuan Asesmen: Sertifikasi / Sertifikasi Ulang / PKT / RPL).
* **FR-APL01-02:** Penautan bukti kelengkapan dasar (KTP, Ijazah, Transkrip, Sertifikat Pelatihan, CV) melalui URL hasil *client-side upload* beserta status kelengkapan per item (`MEMENUHI_SYARAT`, `TIDAK_MEMENUHI_SYARAT`, `BELUM_LENGKAP`).
* **FR-APL01-03:** Submit APL-01 oleh Asesi (mengubah status dari `DRAFT` menjadi `SUBMITTED`).
* **FR-APL01-04:** Verifikasi APL-01 oleh Admin LSP (`APPROVED` atau `REJECTED` beserta catatan verifikasi) serta plotting Asesor. Jika `APPROVED`, sistem otomatis men-generate kerangka form **APL-02** berdasarkan seluruh KUK pada skema terkait dan menjadwalkan notifikasi penugasan ke Asesor dalam satu transaksi atomik.

### 3.3. Asesmen Mandiri — FR-APL-02 (`apl02`)
* **FR-APL02-01:** Asesi melihat daftar seluruh Unit Kompetensi, Elemen, dan KUK yang harus dinilai mandiri berdasarkan APL-01 yang telah disetujui.
* **FR-APL02-02:** **Bulk Save Jawaban Asesmen Mandiri:** Asesi dapat menyimpan puluhan hingga ratusan jawaban KUK sekaligus dalam satu request (`kompeten`: `K` / `BK`, dan `bukti_relevan_ids`: array ID bukti portofolio). Harus bersifat atomik (Unit of Work) dan tahan terhadap *concurrent save*.
* **FR-APL02-03:** Submit final APL-02 oleh Asesi (`DRAFT` -> `SUBMITTED`). Validasi domain wajib memastikan seluruh KUK pada skema telah terisi sebelum di-submit, serta memicu pengiriman notifikasi email ke Asesor yang ditugaskan.
* **FR-APL02-04:** Penilaian & Rekomendasi Asesor (`RECOMMENDED` / `NOT_RECOMMENDED`) beserta catatan tinjauan per unit/keseluruhan dan tanda tangan digital/timestamp rekomendasi.

### 3.4. Manajemen Bukti Portofolio & Storage (`portfolio`)
* **FR-PRT-01:** Penerbitan *Presigned PUT URL* Cloudflare R2 agar client mengunggah file dokumen langsung ke cloud storage tanpa membebani server backend.
* **FR-PRT-02:** Registrasi metadata file bukti kompetensi (Nomor Bukti, Nama Dokumen, Object Key, File URL, Tipe Bukti) yang dapat direferensikan ulang oleh banyak KUK di APL-02.

### 3.5. Notifikasi Email Otomatis (Mailgun)
* **FR-NTF-01:** Mengirim email notifikasi ke Asesor saat Admin menetapkan plotting/penugasan Asesor pada permohonan yang disetujui (`ASESOR_ASSIGNED`).
* **FR-NTF-02:** Mengirim email pemberitahuan ke Asesor ketika Asesi selesai melakukan submit final formulir asesmen mandiri (`APL02_SUBMITTED`) lengkap dengan *deep link* peninjauan.

## 4. Non-Functional Requirements (NFR)
* **Konsistensi Data & Reliabilitas Email:** Seluruh operasi multi-tabel, penyimpanan massal jawaban KUK (APL-02), dan penjadwalan pengiriman email wajib menggunakan pola *Unit of Work* + *Transactional Outbox* agar tidak ada notifikasi hilang atau terkirim saat transaksi rollback.
* **Idempotensi:** Endpoint *submit* APL-01, *submit* APL-02, dan keputusan verifikasi wajib mendukung `Idempotency-Key` untuk mencegah duplikasi akibat *retry* jaringan.
* **Dokumentasi API:** OpenAPI 3.1 interaktif otomatis tersedia di `/docs` dan `/openapi.json`.
* **Performa Read:** Daftar data menggunakan *Offset Pagination* standar (`page`, `limit`) dengan response `< 100ms` dan caching Redis pada data master skema.
* **Developer Experience:** Mendukung *Auth Bypass Mode* (`DISABLE_AUTH=true`) dan *Log Mailer Mode* (`MAILER_DRIVER=log`) saat `APP_ENV=development` agar pengembangan lokal berjalan tanpa token JWT maupun kredensial Mailgun asli.