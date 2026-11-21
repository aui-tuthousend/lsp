# System & Software Architecture Specification

## 1. High-Level Architectural Style
Sistem dibangun menggunakan **Domain-Oriented Hexagonal Architecture (Ports & Adapters)** yang dipadukan dengan **CQRS (Command Query Responsibility Segregation)**, **Unit of Work (UoW)**, **Transactional Outbox**, **Realtime Events (SSE)**, dan **Typed App Container (`internal/bootstrap`)**.

Repositori berbentuk **monorepo**: backend Go di `be/`, frontend di `fe/` (bun workspace), dengan tooling bersama (`Makefile`, `docker-compose.yml`) di root.

* **Prinsip Utama:**
  * **Feature/Use-Case Centric Services:** Tidak ada `Service` 1-on-1 per tabel database. Logika bisnis diorganisasikan berdasarkan *Use Case* spesifik (`app/command` dan `app/query`).
  * **Bounded Context Isolation:** `apl01` dan `apl02` dipisahkan sebagai dua modul independen karena perbedaan aktor (Admin LSP vs Asesor), pola transaksi (single row vs bulk KUK upsert), dan siklus hidupnya. Komunikasi lintas modul hanya melalui *Port Interface* yang disambungkan di `internal/bootstrap`.
  * **Dependency Rule:** Ketergantungan kode hanya mengarah ke dalam (`adapters` -> `app` -> `domain`). Layer `domain` adalah Go murni tanpa import SQL, Redis, S3 SDK, Mailgun SDK, maupun HTTP.
  * **Typed App Container (Compile-Time Safe DI):** Seluruh *wiring* dependency dirakit secara eksplisit dan *type-safe* di dalam paket `internal/bootstrap` (`App` struct), menjaga `cmd/api/main.go` tetap ramping dan memudahkan *end-to-end integration testing*.

---

## 2. Struktur Direktori Proyek

```text
├── be/                                 # Backend Go
│   ├── cmd/
│   │   ├── api/
│   │   │   └── main.go                 # Entrypoint HTTP Server tipis (< 45 baris), load config & graceful shutdown
│   │   ├── worker/
│   │   │   └── main.go                 # Entrypoint Background Worker (Outbox Email Processor & publish event SSE)
│   │   └── seed/
│   │       └── main.go                 # Seeder user awal (idempoten)
│   ├── migrations/                     # File migrasi SQL berurutan (Goose / Golang-Migrate)
│   ├── internal/
│   │   ├── bootstrap/                  # Composition Root & Typed App Container
│   │   │   ├── app.go                  # Struct App & NewApp() (merakit seluruh modul & router)
│   │   │   ├── module_schema.go        # Wiring khusus bounded context schema
│   │   │   ├── module_apl01.go         # Wiring khusus bounded context apl01
│   │   │   ├── module_apl02.go         # Wiring khusus bounded context apl02
│   │   │   ├── module_portfolio.go     # Wiring khusus bounded context portfolio & presigned storage
│   │   │   └── module_notifications.go # Wiring khusus bounded context notifikasi realtime (SSE)
│   │   │   └── module_asesor.go        # Wiring khusus bounded context profil asesor
│   │   │
│   │   ├── platform/                   # Shared Infrastructure & Cross-Cutting Concerns
│   │   │   ├── config/
│   │   │   │   └── config.go           # Env loader & validation (DB, Redis, R2 S3, Mailgun, JWT)
│   │   │   ├── id/
│   │   │   │   └── short_uuid.go       # Generator & validator Short UUID (22-char base57)
│   │   │   ├── storage/
│   │   │   │   ├── storage.go          # Port interface ObjectStorage (PresignedPutURL, PresignedGetURL, Delete, Exists)
│   │   │   │   └── s3_storage.go       # Adapter Cloudflare R2 / S3 (aws-sdk-go-v2)
│   │   │   ├── mailer/
│   │   │   │   ├── mailer.go           # Port interface Mailer & struct EmailMessage
│   │   │   │   ├── mailgun_mailer.go   # Adapter produksi Mailgun API (mailgun/mailgun-go/v4)
│   │   │   │   ├── log_mailer.go       # Adapter development (cetak email ke slog terminal tanpa kirim ke internet)
│   │   │   │   └── templates.go        # Renderer HTML email template (Penugasan Asesor & Submit APL-02)
│   │   │   ├── database/
│   │   │   │   ├── postgres.go         # pgxpool setup & tuning
│   │   │   │   ├── tx_manager.go       # Unit of Work (WithinTransaction via context.Context)
│   │   │   │   ├── querier.go          # DBTX interface & GetExecutor(ctx, pool)
│   │   │   │   ├── outbox_repo.go      # Repository tabel outbox_events (Insert & FetchPendingForUpdate)
│   │   │   │   └── retry.go            # Retry wrapper untuk PostgreSQL Deadlock (40P01) & Serialization (40001)
│   │   │   ├── cache/
│   │   │   │   └── redis.go            # go-redis/v9 client wrapper
│   │   │   ├── events/
│   │   │   │   ├── events.go           # Port Bus (Publish/Subscribe) + Event + helper topic
│   │   │   │   └── redis_bus.go        # Adapter Redis Pub/Sub (fan-out event realtime lintas instance)
│   │   │   ├── auth/
│   │   │   │   ├── claims.go           # UserClaims struct & context injector
│   │   │   │   ├── jwt.go              # golang-jwt/v5 sign & verify
│   │   │   │   └── huma_middleware.go  # JWTAuthMiddleware (Prod) / DevAuthMiddleware (Dev) + RequireRoles (RBAC)
│   │   │   ├── http/
│   │   │   │   ├── middleware.go       # Chain(), RequestID, Structured Logger, Panic Recoverer, Timeout, MaxBodySize
│   │   │   │   ├── pagination.go       # Generic PaginationQuery, PaginationMeta, PagedResponse[T]
│   │   │   │   └── openapi_helper.go   # Helper WithDoc() untuk Huma operations
│   │   │   └── worker/
│   │   │       └── outbox_processor.go # Worker Pool + FOR UPDATE SKIP LOCKED, kirim email & publish event SSE
│   │   │
│   │   ├── features/                   # Bounded Contexts (fitur bisnis)
│   │   │   ├── notifications/          # Notifikasi realtime (SSE)
│   │   │   │   ├── domain/
│   │   │   │   │   └── event.go            # NotificationEvent (payload SSE)
│   │   │   │   └── adapters/
│   │   │   │       └── http/
│   │   │   │           ├── dto.go          # EventsInput
│   │   │   │           ├── routes.go       # GET /events (Huma sse.Register)
│   │   │   │           └── handler.go      # Subscribe bus + heartbeat + graceful shutdown
│   │   │   │
│   │   │   ├── users/                  # Akun & autentikasi (login, register, RBAC, master roles)
│   │   │   │   ├── domain/                 # User, Role, repository port, sentinel errors
│   │   │   │   ├── app/                    # Service login/register/user CRUD
│   │   │   │   └── adapters/
│   │   │   │       ├── postgres/           # Tabel users & roles (+ IsAsesor/IsAsesi checker)
│   │   │   │       └── http/               # /auth/* & /users
│   │   │   │
│   │   │   ├── asesi/                  # Profil data diri asesi (relasi 1:1 ke users)
│   │   │   │   ├── domain/                 # Profile, Patch, repository port, AsesiChecker
│   │   │   │   ├── app/                    # CRUD service profil asesi
│   │   │   │   └── adapters/
│   │   │   │       ├── postgres/           # Tabel asesi_profiles
│   │   │   │       └── http/               # CRUD /asesi-profiles (admin-only)
│   │   │   │
│   │   │   ├── asesor/                 # Profil asesor (afiliasi TUK, domisili, batas penugasan) + sertifikasi
│   │   │   │   ├── domain/
│   │   │   │   │   ├── profile.go          # Profile & ListItem + sentinel errors
│   │   │   │   │   ├── certification.go    # Certification & Patch
│   │   │   │   │   └── repository.go       # Port Repository & AsesorChecker
│   │   │   │   ├── app/
│   │   │   │   │   └── service.go          # Upsert/Patch/Get/List + CRUD sertifikasi
│   │   │   │   └── adapters/
│   │   │   │       ├── postgres/           # Tabel asesor_profiles & asesor_certifications
│   │   │   │       └── http/               # /asesor-profiles & /asesor-certifications
│   │   │   │
│   │   │   ├── tuk/                    # Master data Tempat Uji Kompetensi
│   │   │   │   ├── domain/                 # TUK model, repository port, sentinel errors
│   │   │   │   ├── app/                    # CRUD service
│   │   │   │   └── adapters/
│   │   │   │       ├── postgres/           # Tabel tuks
│   │   │   │       └── http/               # CRUD /tuks (admin-only)
│   │   │   │
│   │   │   ├── schema/                 # Master Skema, Unit, Elemen, KUK
│   │   │   │   ├── domain/
│   │   │   │   ├── app/
│   │   │   │   │   ├── command/
│   │   │   │   │   └── query/
│   │   │   │   └── adapters/
│   │   │   │       ├── postgres/
│   │   │   │       ├── redis/
│   │   │   │       └── http/
│   │   │   │
│   │   │   ├── apl01/                  # Pendaftaran & Verifikasi FR-APL-01
│   │   │   │   ├── domain/
│   │   │   │   │   ├── application.go      # Aggregate Root APL-01 & state machine rules
│   │   │   │   │   ├── errors.go           # Domain sentinel errors
│   │   │   │   │   └── repository.go       # Port interfaces (termasuk APL02Initializer & OutboxWriter)
│   │   │   │   ├── app/
│   │   │   │   │   ├── command/
│   │   │   │   │   │   ├── create_apl01.go
│   │   │   │   │   │   ├── submit_apl01.go
│   │   │   │   │   │   └── verify_apl01.go # Approve APL-01 + generate draft APL-02 + outbox email Asesor dalam 1 Tx
│   │   │   │   │   └── query/
│   │   │   │   │       ├── get_apl01_detail.go
│   │   │   │   │       └── list_apl01.go
│   │   │   │   └── adapters/
│   │   │   │       ├── postgres/
│   │   │   │       │   ├── apl01_repo.go
│   │   │   │       │   └── query_model.go
│   │   │   │       └── http/
│   │   │   │           ├── dto.go
│   │   │   │           ├── routes.go
│   │   │   │           ├── handler.go
│   │   │   │           └── errors.go
│   │   │   │
│   │   │   ├── apl02/                  # Asesmen Mandiri FR-APL-02
│   │   │   │   ├── domain/
│   │   │   │   │   ├── assessment.go       # Aggregate Root APL-02 & KUK answer validation
│   │   │   │   │   ├── errors.go
│   │   │   │   │   └── repository.go
│   │   │   │   ├── app/
│   │   │   │   │   ├── command/
│   │   │   │   │   │   ├── bulk_save_kuk.go    # Simpan puluhan jawaban KUK sekaligus + anti-deadlock
│   │   │   │   │   │   ├── submit_apl02.go     # Submit APL-02 + simpan event outbox email ke Asesor dalam 1 Tx
│   │   │   │   │   │   └── recommend_apl02.go  # Rekomendasi Asesor
│   │   │   │   │   └── query/
│   │   │   │   │       └── get_apl02_form.go
│   │   │   │   └── adapters/
│   │   │   │       ├── postgres/
│   │   │   │       ├── redis/
│   │   │   │       └── http/
│   │   │   │
│   │   │   └── portfolio/              # Manajemen Dokumen & Bukti Portofolio
│   │   │       ├── domain/
│   │   │       ├── app/
│   │   │       │   ├── command/
│   │   │       │   │   ├── generate_upload_url.go # Generate Presigned PUT URL + object_key ke client
│   │   │       │   │   ├── register_document.go   # Terima file_url/object_key setelah client upload & simpan DB
│   │   │       │   │   └── delete_document.go
│   │   │       │   └── query/
│   │   │       └── adapters/
│   │   │           ├── postgres/
│   │   │           └── http/
│   ├── go.mod
│   └── .env.example                    # Contoh konfigurasi environment
├── fe/                                 # Frontend (placeholder, bun workspace)
│   └── package.json
├── docker/
│   ├── postgres.Dockerfile             # Postgres 16 + auto-apply migrasi (be/migrations)
│   └── dragonfly.Dockerfile            # DragonflyDB (kompatibel Redis)
├── docs/
│   ├── PRD.md
│   ├── ARCHITECTURE.md
│   └── CODING_STANDARDS.md
├── docker-compose.yml                  # Postgres + Dragonfly untuk development
├── Makefile                            # Target dev/up/down/seed/migrate/build/test
└── package.json                        # Root bun workspace
```

---

## 3. Realtime Events (Server-Sent Events)

Selain notifikasi email asinkron, sistem menyediakan kanal **realtime satu arah** (server -> klien) melalui SSE untuk memberi umpan balik instan di UI (mis. status verifikasi APL-01, penugasan Asesor, submit APL-02).

* **Port Fan-out (`internal/platform/events`):**
  * Interface `Bus` dengan `Publish` dan `Subscribe(topics...)`, ditambah helper `UserTopic(id)` (`user:<id>`) dan `RoleTopic(role)` (`role:<role>`).
  * Adapter produksi `RedisBus` memakai Redis Pub/Sub, sehingga event bisa didistribusikan **lintas instance API** (bukan hanya in-process).
* **Producer:** `internal/platform/worker/outbox_processor.go`. Setelah event outbox diproses dan transaksi **commit**, worker mem-publish payload ke bus berdasarkan metadata `user_id` dan/atau `role` di payload outbox. Publish bersifat *best-effort* (kegagalan hanya di-log), sementara pengiriman email tetap andal lewat outbox.
* **Consumer:** `internal/features/notifications/adapters/http` mengekspos `GET /events` memakai `huma/v2/sse`. Handler berlangganan topic milik user yang login, mengirim event awal `connected`, lalu meneruskan setiap event sebagai nama event `notification` (jenis asli ada pada field `type`).
* **Autentikasi:** middleware JWT menerima token dari header `Authorization: Bearer <token>` **atau** query `?token=<token>`, karena `EventSource` browser tidak bisa mengirim header custom.
* **Keep-alive & Shutdown:**
  * Heartbeat berupa komentar SSE (`: ping`) setiap 20 detik agar koneksi tidak ditutup proxy/load balancer.
  * Handler berhenti saat request context selesai atau saat aplikasi di-shutdown (channel `ctx.Done()` diteruskan dari `main`), sehingga `server.Shutdown` tidak menggantung.
* **Middleware khusus:**
  * `httpplatform.Timeout` menerima daftar path yang dikecualikan; `/events` **wajib** dikecualikan karena `http.TimeoutHandler` membuffer respons dan tidak mendukung flush.
  * `http.Server` **tidak** boleh diberi `WriteTimeout`, karena akan memutus stream berumur panjang.
* **Format wire (Redis Pub/Sub):** `{"type":"<event_type>","data":<payload>}` pada channel = nama topic.

**Alur ringkas:**
```
Command handler ──(1 Tx)──> outbox_events ──> Worker ──> Bus.Publish ──> Redis Pub/Sub
                                                                              │
Frontend EventSource  <──  GET /events (SSE)  <── Bus.Subscribe  <────────────┘
```
