# Kitab Dapur

Aplikasi untuk menulis dan menyimpan resep keluarga: resep dari orang tua dan saudara,
supaya tidak hilang. Nantinya bisa dipakai orang lain juga.

| Bagian | Teknologi |
| --- | --- |
| `apps/api` | Go, `net/http`, PostgreSQL (pgx), golang-migrate |
| `apps/mobile` | Flutter |
| `api/openapi.yaml` | Kontrak API, satu-satunya sumber kebenaran |

Monorepo ini disusun supaya tiap app gampang dipindah ke repo sendiri. Aturannya ada di
[docs/adr/0001](docs/adr/0001-monorepo-yang-mudah-dipecah.md).

## Mulai

Butuh: Docker, Go 1.25, Flutter.

```bash
make .env   # buat .env dari .env.example, lalu ganti password-nya
make dev    # postgres + api di http://localhost:8080
make test   # test semua app (test API memakai postgres dari docker)
make gen    # generate ulang kode dari api/openapi.yaml
```

Port bentrok dengan project lain? Ganti `POSTGRES_PORT` atau `API_PORT` di `.env`.
Aplikasi mobile: lihat [apps/mobile/README.md](apps/mobile/README.md).

Coba API-nya:

```bash
curl -X POST localhost:8080/kitabs -H 'Content-Type: application/json' -d '{"name":"Keluarga Syafrin"}'
curl -X POST localhost:8080/kitabs/<id-kitab>/recipes -H 'Content-Type: application/json' \
  -d '{"title":"Rendang","author":"Ibu","ingredients":["1 kg daging"],"steps":["Masak santan"]}'
curl localhost:8080/kitabs/<id-kitab>/recipes
```

## Struktur

```
apps/api/
  cmd/server/main.go        composition root: satu-satunya tempat yang merangkai semua bagian
  internal/
    recipe/                 domain + use case (hanya standard library)
      handler/              adapter HTTP: menjalankan operasi dari openapi.yaml
      repository/           adapter PostgreSQL
    kitab/                  koleksi resep milik satu keluarga, dengan susunan yang sama
    auth/  media/           belum dibuat
    openapi/                hasil generate dari api/openapi.yaml (jangan diedit)
    platform/               config, koneksi db, http server, health check
  migrations/               SQL migration
apps/mobile/                aplikasi Flutter
api/openapi.yaml            kontrak API
deploy/docker-compose.yml   lingkungan dev lokal
docs/                       catatan desain dan keputusan arsitektur
```

Arah dependensi selalu ke dalam: `handler`/`repository` → `recipe` (domain + use case).
Package domain tidak boleh meng-import `net/http`, pgx, atau kode hasil generate.

## Status

- [x] Kitab dan resep: buat, lihat daftar, lihat detail
- [ ] **Auth.** Saat ini siapa pun yang tahu id kitab bisa membaca dan menulis isinya. Jangan
      dibuka ke internet sebelum `internal/auth` selesai.
- [ ] Upload foto resep ke Google Cloud Storage (`internal/media`)
- [x] Aplikasi mobile: buat kitab, tulis resep, lihat daftar dan detail (belum dicoba di HP atau emulator)
- [ ] Config server produksi (`deploy/prod/`)
