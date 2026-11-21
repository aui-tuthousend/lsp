```markdown
# Coding Standards & Implementation Rules for AI Agent

Ikuti aturan berikut secara ketat saat men-generate kode untuk proyek ini:

## 1. Aturan App Container & Cross-Module Isolation (`internal/bootstrap/`)
* **Dilarang** membuat dynamic service locator berbasis `map[string]any`, `reflect`, atau `interface{}`. Seluruh dependency pada `bootstrap.App` dan sub-modulnya wajib bertipe eksplisit (*compile-time safe*).
* **Dilarang** menaruh logika *wiring* repository/handler di dalam `cmd/api/main.go`. `main.go` hanya boleh memanggil `bootstrap.NewApp(cfg, pool, rdb, objStorage)`.
* **Cross-Module Isolation:** Paket `internal/features/apl01` dilarang keras meng-import paket apa pun di bawah `internal/features/apl02` (dan sebaliknya). Jika `apl01` perlu memicu pembuatan form `apl02` saat approval, definisikan port interface (misal `APL02Initializer`) di `internal/features/apl01/domain/repository.go`, implementasikan di `internal/features/apl02/adapters/postgres/`, dan sambungkan keduanya hanya di `internal/bootstrap/app.go`.
* Pisahkan konstruktor modul per *bounded context* (`module_schema.go`, `module_apl01.go`, `module_apl02.go`, `module_portfolio.go`) di dalam paket `internal/bootstrap`.
* `bootstrap.NewApp` tidak boleh memanggil `os.Exit` atau `log.Fatal`; biarkan `cmd/api/main.go` yang menangani lifecycle startup dan shutdown.

## 2. Aturan Email Service — Mailgun & Outbox (`internal/platform/mailer` & `worker`)
* **Dilarang** memanggil `mailer.Send(...)` secara langsung di dalam HTTP handler atau di dalam blok `TxManager.WithinTransaction(...)`.
* Setiap aksi bisnis yang memicu notifikasi email ke Asesor (seperti `VerifyAPL01` saat plotting Asesor atau `SubmitAPL02` saat Asesi submit asesmen mandiri) **wajib** menulis event ke tabel `outbox_events` menggunakan `txCtx` yang sama dengan mutasi database utamanya.
* Eksekusi pengiriman email melalui **Mailgun (`mailgun/mailgun-go/v4`)** hanya boleh dilakukan oleh worker asinkron (`internal/platform/worker/outbox_processor.go`) yang menarik antrean dengan `FOR UPDATE SKIP LOCKED`.
* Sediakan switch konfigurasi `MAILER_DRIVER` (`mailgun` untuk produksi/staging, `log` untuk local development yang mencetak subjek, penerima, dan isi email ke `slog`).

## 3. Aturan Short UUID & Client-Side Cloud Storage (`internal/platform/id` & `storage`)
* **Short UUID:** Selalu generate ID baru di aplikasi Go menggunakan `id.New()` dari `internal/platform/id`. Jangan menyerahkan pembuatan Primary Key ke `DEFAULT` database.
* **Client-Side Presigned Upload (Pola B):**
  * Backend Go **tidak menerima** file biner `multipart/form-data`. Seluruh upload file dilakukan langsung dari client ke Cloudflare R2 menggunakan `PresignedPutURL`.
  * Saat client mengirimkan `file_url` / `object_key` untuk disimpan ke database, backend **wajib memvalidasi** bahwa `object_key` diawali dengan prefix folder milik user yang sedang login (mencegah *IDOR / object hijacking* antar-user) dan URL sesuai dengan base URL R2 yang dikonfigurasi.
  * Selalu injeksikan interface `storage.ObjectStorage` ke Command Handler (bukan struct S3 SDK langsung).

## 4. Aturan Huma HTTP Adapter & Validasi (`adapters/http/`)
* **Pisahkan 4 file per fitur:** `dto.go` (struct input/output & tag validasi), `routes.go` (registrasi grup route), `handler.go` (method handler tipis 5–15 baris), dan `errors.go` (pemetaan sentinel error domain ke `huma.ErrorXXX`).
* **Validasi Request:** Gunakan tag Huma (`minLength`, `maxLength`, `pattern`, `format`, `enum`) pada `dto.go` untuk validasi sintaksis; gunakan tipe non-pointer untuk field wajib dan pointer (`*string`, `*int`) untuk field opsional. Untuk validasi lintas-field pada input, implementasikan `Resolve(ctx huma.Context) []error` di `dto.go`.
* **Dilarang** menggunakan *anonymous closure* di dalam `routes.go`. Gunakan method receiver pada struct handler (`huma.Post(grp, "/submit", h.Submit, httpplatform.WithDoc(...))`).
* **Dilarang** menaruh tag `Authorization` dengan `required:"true"` di dalam struct Request DTO agar `DevAuthMiddleware` dapat bekerja tanpa header Authorization saat development. Daftarkan `bearerAuth` pada level konfigurasi OpenAPI global.
* **Pagination Standar:** Semua endpoint list wajib meng-embed `httpplatform.PaginationQuery` pada struct Request dan mengembalikan `*httpplatform.PagedResponse[ItemDTO]`.

## 5. Aturan Raw SQL & Repository (`adapters/postgres/`)
* Selalu ambil koneksi melalui `db := database.GetExecutor(ctx, r.pool)` di baris pertama setiap method repository agar kompatibel dengan `TxManager`.
* Gunakan parameter binding PostgreSQL (`$1`, `$2`, ...) — dilarang keras melakukan string concatenation untuk nilai input user.
* Untuk query list dengan pagination, gunakan Window Function `COUNT(*) OVER() AS total_count` dalam satu query, disertai fallback `SELECT COUNT(*)` hanya jika `len(items) == 0 && offset > 0`.
* Konversikan error teknis database (`pgx.ErrNoRows`, `*pgconn.PgError` kode `23505` unique violation, `23503` foreign key violation) menjadi *domain sentinel errors* (`domain.ErrNotFound`, `domain.ErrConflict`, dll.) di dalam repository. Jangan pernah membocorkan `*pgconn.PgError` ke layer `app/` atau `http/`.

## 6. Aturan Domain Layer (`domain/`)
* File di dalam folder `domain/` **tidak boleh** meng-import `database/sql`, `github.com/jackc/pgx/v5`, `github.com/aws/aws-sdk-go-v2`, `github.com/mailgun/mailgun-go/v4`, `net/http`, atau `github.com/danielgtaylor/huma/v2`. Import `internal/platform/id` diperbolehkan.
* Validasi *state transition* (misal: `apl01.Submit()`, `apl01.Approve()`, `apl02.BulkUpdateAnswers()`) harus dienkapsulasi sebagai method pada struct Entity/Aggregate Root di `domain/`.

## 7. Database Schema Conventions (PostgreSQL)
* Gunakan `VARCHAR(26)` sebagai tipe kolom Primary Key dan Foreign Key untuk menampung Short UUID dari Go.
* Gunakan `TIMESTAMPTZ` (UTC) untuk seluruh kolom waktu (`created_at`, `updated_at`, `submitted_at`, `verified_at`).
* Wajib membuat *composite index* pada kolom yang sering difilter dan diurutkan bersamaan (contoh: `CREATE INDEX idx_apl01_status_created ON apl01_applications (status, created_at DESC);`).