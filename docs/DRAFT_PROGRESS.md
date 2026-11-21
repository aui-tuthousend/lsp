# Konsep: Draft Progress Pendaftaran (disimpan untuk nanti)

> Status: **konsep / belum diimplementasikan**. Catatan ini disimpan sebagai acuan
> saat alur pendaftaran multi-step mulai dikerjakan.

## 1. Masalah

Proses pendaftaran terdiri dari banyak step. User sering berhenti di tengah
(mis. mengisi form lalu menutup tab). Kita ingin:

1. User bisa **menyimpan progress** pendaftaran yang belum selesai.
2. User bisa **melanjutkan** dari step terakhir.
3. Beranda menampilkan **notif** kalau ada progress yang belum selesai.

## 2. Keputusan yang sudah diambil

| Aspek | Keputusan |
|---|---|
| Tempat draft | **Status `DRAFT` di tabel form itu sendiri** (bukan tabel draft terpisah) |
| Cara simpan | **Autosave + tombol "Simpan"** eksplisit |
| Notif beranda | **Query biasa** (bukan push/SSE) |
| Cakupan awal | **Satu alur dulu** (APL-01) sebagai pola percontohan |

## 3. Konsep inti

Draft = record pendaftaran dengan `status = 'DRAFT'`, bukan tabel terpisah.
"Simpan" = update record tersebut. Submit final = ubah status ke `SUBMITTED`.

```
User isi form
   │
   ├─(autosave / tombol simpan)──▶ record status=DRAFT  (current_step, draft_data, updated_at)
   │
   ├─ tutup / pergi ──▶ record tetap DRAFT di DB
   │
   └─ buka beranda ──▶ query "pendaftaran DRAFT milik saya" ──▶ banner "lanjutkan"
                          └─ klik ──▶ resume ke step terakhir
```

Keuntungan:
- Satu sumber kebenaran; tidak ada tabel draft terpisah yang harus disinkronkan.
- Submit final cukup ubah status (atomic, idempotent-friendly).
- Notif beranda murni derived → tidak ada state notifikasi yang perlu dibersihkan.

## 4. Bentuk data (gambaran, untuk APL-01)

```
apl01
 ├─ id
 ├─ user_id          -- pemilik draft
 ├─ schema_id        -- skema sertifikasi yang dipilih
 ├─ status           -- DRAFT | SUBMITTED | APPROVED | REJECTED
 ├─ current_step     -- step terakhir yang dikerjakan (untuk resume)
 ├─ draft_data       -- jsonb: isi form per step
 ├─ created_at
 ├─ updated_at       -- dipakai untuk "terakhir disimpan"
 └─ submitted_at
```

Batasan: satu draft aktif per `(user_id, schema_id)`:

```
unique index (user_id, schema_id) where status = 'DRAFT'
```

## 5. Endpoint (gambaran)

| Method | Path | Fungsi |
|---|---|---|
| POST | `/apl01` | buat draft baru (atau ambil yang sudah ada) |
| GET | `/apl01/mine` | daftar draft milik saya → untuk beranda |
| GET | `/apl01/{id}` | detail → resume ke `current_step` |
| PATCH | `/apl01/{id}` | simpan partial (autosave & tombol simpan) |
| POST | `/apl01/{id}/submit` | validasi lengkap → `SUBMITTED` |

Validasi dibedakan:
- **PATCH (draft)**: longgar — hanya field yang dikirim yang divalidasi.
- **submit**: ketat — semua field wajib + bukti kelengkapan harus lengkap.

## 6. Alur FE

```
Mulai isi form
  └─ POST /apl01 ──▶ dapat id

Pindah step / berkala (debounce)
  └─ PATCH /apl01/{id}  (current_step + draft_data)

Klik "Simpan"
  └─ PATCH /apl01/{id}  + toast "Tersimpan"

Buka beranda
  └─ GET /apl01/mine  ──▶ banner "Lanjutkan pendaftaran (step 2/4)"
        └─ klik ──▶ GET /apl01/{id} ──▶ lanjut dari current_step

Submit final
  └─ POST /apl01/{id}/submit ──▶ status SUBMITTED
        └─ otomatis hilang dari daftar draft beranda
```

TanStack Query: key `["apl01","mine","DRAFT"]`, `staleTime` ~15s.

## 7. Yang belum final (perlu diputuskan saat implementasi)

1. **Jumlah draft per user** — saran: 1 draft aktif per `(user_id, schema_id)`.
2. **Konkurensi multi-tab** — autosave dari 2 tab bisa saling menimpa; opsi
   optimistic concurrency pakai `updated_at`.
3. **Kedaluwarsa draft** — auto-arsip setelah sekian lama, atau biarkan dulu.
4. **Idempotency submit** — pakai `Idempotency-Key` (sudah diminta di PRD).
5. **Struktur `draft_data` jsonb** — perlu kesepakatan key per step FE & BE.
