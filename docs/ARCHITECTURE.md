# System & Software Architecture Specification

## 1. High-Level Architectural Style
Sistem dibangun menggunakan **Domain-Oriented Hexagonal Architecture (Ports & Adapters)** yang dipadukan dengan **CQRS (Command Query Responsibility Segregation)**, **Unit of Work (UoW)**, **Transactional Outbox**, dan **Typed App Container (`internal/bootstrap`)**.

* **Prinsip Utama:**
  * **Feature/Use-Case Centric Services:** Tidak ada `Service` 1-on-1 per tabel database. Logika bisnis diorganisasikan berdasarkan *Use Case* spesifik (`app/command` dan `app/query`).
  * **Bounded Context Isolation:** `apl01` dan `apl02` dipisahkan sebagai dua modul independen karena perbedaan aktor (Admin LSP vs Asesor), pola transaksi (single row vs bulk KUK upsert), dan siklus hidupnya. Komunikasi lintas modul hanya melalui *Port Interface* yang disambungkan di `internal/bootstrap`.
  * **Dependency Rule:** Ketergantungan kode hanya mengarah ke dalam (`adapters` -> `app` -> `domain`). Layer `domain` adalah Go murni tanpa import SQL, Redis, S3 SDK, Mailgun SDK, maupun HTTP.
  * **Typed App Container (Compile-Time Safe DI):** Seluruh *wiring* dependency dirakit secara eksplisit dan *type-safe* di dalam paket `internal/bootstrap` (`App` struct), menjaga `cmd/api/main.go` tetap ramping dan memudahkan *end-to-end integration testing*.

---

## 2. Struktur Direktori Proyek

```text
├── cmd/
│   ├── api/
│   │   └── main.go                     # Entrypoint HTTP Server tipis (< 45 baris), load config & graceful shutdown
│   └── worker/
│       └── main.go                     # Entrypoint Background Worker (Outbox Email Processor & Async Cleanup)
├── docs/
│   ├── PRD.md
│   ├── ARCHITECTURE.md
│   └── CODING_STANDARDS.md
├── migrations/                         # File migrasi SQL berurutan (Goose / Golang-Migrate)
├── internal/
│   ├── bootstrap/                      # Composition Root & Typed App Container
│   │   ├── app.go                      # Struct App & NewApp() (merakit seluruh modul & router)
│   │   ├── module_schema.go            # Wiring khusus bounded context schema
│   │   ├── module_apl01.go             # Wiring khusus bounded context apl01
│   │   ├── module_apl02.go             # Wiring khusus bounded context apl02
│   │   └── module_portfolio.go         # Wiring khusus bounded context portfolio & presigned storage
│   │
│   ├── platform/                       # Shared Infrastructure & Cross-Cutting Concerns
│   │   ├── config/
│   │   │   └── config.go               # Env loader & validation (DB, Redis, R2 S3, Mailgun, JWT)
│   │   ├── id/
│   │   │   └── short_uuid.go           # Generator & validator Short UUID (22-char base57)
│   │   ├── storage/
│   │   │   ├── storage.go              # Port interface ObjectStorage (PresignedPutURL, PresignedGetURL, Delete, Exists)
│   │   │   └── s3_storage.go           # Adapter Cloudflare R2 / S3 (aws-sdk-go-v2)
│   │   ├── mailer/
│   │   │   ├── mailer.go               # Port interface Mailer & struct EmailMessage
│   │   │   ├── mailgun_mailer.go       # Adapter produksi Mailgun API (mailgun/mailgun-go/v4)
│   │   │   ├── log_mailer.go           # Adapter development (cetak email ke slog terminal tanpa kirim ke internet)
│   │   │   └── templates.go            # Renderer HTML email template (Penugasan Asesor & Submit APL-02)
│   │   ├── database/
│   │   │   ├── postgres.go             # pgxpool setup & tuning
│   │   │   ├── tx_manager.go           # Unit of Work (WithinTransaction via context.Context)
│   │   │   ├── querier.go              # DBTX interface & GetExecutor(ctx, pool)
│   │   │   ├── outbox_repo.go          # Repository tabel outbox_events (Insert & FetchPendingForUpdate)
│   │   │   └── retry.go                # Retry wrapper untuk PostgreSQL Deadlock (40P01) & Serialization (40001)
│   │   ├── cache/
│   │   │   └── redis.go                # go-redis/v9 client wrapper
│   │   ├── auth/
│   │   │   ├── claims.go               # UserClaims struct & context injector
│   │   │   ├── jwt.go                  # golang-jwt/v5 sign & verify
│   │   │   └── huma_middleware.go      # Switchable JWTAuthMiddleware (Prod) vs DevAuthMiddleware (Dev)
│   │   ├── http/
│   │   │   ├── middleware.go           # Chain(), RequestID, Structured Logger, Panic Recoverer, Timeout, MaxBodySize
│   │   │   ├── pagination.go           # Generic PaginationQuery, PaginationMeta, PagedResponse[T]
│   │   │   └── openapi_helper.go       # Helper WithDoc() untuk Huma operations
│   │   └── worker/
│   │       └── outbox_processor.go     # Fixed Worker Pool + FOR UPDATE SKIP LOCKED untuk eksekusi event & email
│   │
│   ├── schema/                         # Bounded Context 1: Master Skema, Unit, Elemen, KUK
│   │   ├── domain/
│   │   ├── app/
│   │   │   ├── command/
│   │   │   └── query/
│   │   └── adapters/
│   │       ├── postgres/
│   │       ├── redis/
│   │       └── http/
│   │
│   ├── apl01/                          # Bounded Context 2: Pendaftaran & Verifikasi FR-APL-01
│   │   ├── domain/
│   │   │   ├── application.go          # Aggregate Root APL-01 & state machine rules
│   │   │   ├── errors.go               # Domain sentinel errors
│   │   │   └── repository.go           # Port interfaces (termasuk APL02Initializer & OutboxWriter)
│   │   ├── app/
│   │   │   ├── command/
│   │   │   │   ├── create_apl01.go
│   │   │   │   ├── submit_apl01.go
│   │   │   │   └── verify_apl01.go     # Approve APL-01 + generate draft APL-02 + outbox email Asesor dalam 1 Tx
│   │   │   └── query/
│   │   │       ├── get_apl01_detail.go
│   │   │       └── list_apl01.go
│   │   └── adapters/
│   │       ├── postgres/
│   │       │   ├── apl01_repo.go
│   │       │   └── query_model.go
│   │       └── http/
│   │           ├── dto.go
│   │           ├── routes.go
│   │           ├── handler.go
│   │           └── errors.go
│   │
│   ├── apl02/                          # Bounded Context 3: Asesmen Mandiri FR-APL-02
│   │   ├── domain/
│   │   │   ├── assessment.go           # Aggregate Root APL-02 & KUK answer validation
│   │   │   ├── errors.go
│   │   │   └── repository.go
│   │   ├── app/
│   │   │   ├── command/
│   │   │   │   ├── bulk_save_kuk.go    # Simpan puluhan jawaban KUK sekaligus + anti-deadlock
│   │   │   │   ├── submit_apl02.go     # Submit APL-02 + simpan event outbox email ke Asesor dalam 1 Tx
│   │   │   │   └── recommend_apl02.go  # Rekomendasi Asesor
│   │   │   └── query/
│   │   │       └── get_apl02_form.go
│   │   └── adapters/
│   │       ├── postgres/
│   │       ├── redis/
│   │       └── http/
│   │
│   └── portfolio/                      # Bounded Context 4: Manajemen Dokumen & Bukti Portofolio
│       ├── domain/
│       ├── app/
│       │   ├── command/
│       │   │   ├── generate_upload_url.go # Generate Presigned PUT URL + object_key ke client
│       │   │   ├── register_document.go   # Terima file_url/object_key setelah client upload & simpan DB
│       │   │   └── delete_document.go
│       │   └── query/
│       └── adapters/
│           ├── postgres/
│           └── http/
└── go.mod