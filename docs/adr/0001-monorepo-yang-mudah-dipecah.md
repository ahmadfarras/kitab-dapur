# 0001 — Monorepo yang mudah dipecah jadi beberapa repo

Status: diterima · 2026-11-10

## Konteks

Kitab Dapur dimulai sebagai aplikasi keluarga, lalu akan dipakai orang lain. Satu repo memudahkan
perubahan yang menyentuh API dan mobile sekaligus. Tapi nanti tiap bagian mungkin perlu repo, tim,
atau jadwal rilisnya sendiri.

## Keputusan

Satu repo, dengan aturan supaya tiap `apps/*` bisa dipindah ke repo sendiri tanpa ditulis ulang:

1. **Tiap app berdiri sendiri.** `apps/api` punya `go.mod`, `apps/mobile` punya `pubspec.yaml`,
   masing-masing dengan `Makefile`, `Dockerfile`/build, dan job CI sendiri. `Makefile` di root
   hanya meneruskan perintah ke app (`make -C apps/api ...`).
2. **Tidak ada import antar-app** dan tidak ada folder `shared/`. Satu-satunya yang dipakai
   bersama adalah kontrak `api/openapi.yaml`.
3. **Kode hasil generate di-commit di dalam app** (`apps/api/internal/openapi/api.gen.go`), jadi app
   tetap bisa di-build tanpa folder lain. Path ke kontrak cukup diubah di satu tempat per app
   (`internal/openapi/generate.go`).
4. **Migrasi database milik API** (`apps/api/migrations`), karena hanya API yang menyentuh database.
5. **`deploy/` hanya merangkai image**, tidak berisi kode app.

## Cara memecah nanti

- `git filter-repo --subdirectory-filter apps/api` (ulangi untuk `apps/mobile`) supaya riwayat
  commit ikut terbawa.
- `api/openapi.yaml` ikut ke repo API. App lain mengambilnya dari sana: versi rilis, atau di-copy
  saat menjalankan `make gen`.
- Job CI di `.gitlab-ci.yml` sudah dipisah per app, tinggal dipindah.

## Konsekuensi

- Sedikit duplikasi (misal batas panjang judul resep ada di Go dan di `openapi.yaml`) diterima
  demi kemandirian tiap app.
- Perubahan kontrak yang mematahkan client harus dikerjakan bertahap (tambah dulu, hapus belakangan),
  karena setelah dipecah app tidak lagi dirilis bersamaan.
