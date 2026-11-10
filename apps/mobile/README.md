# Kitab Dapur — mobile

Aplikasi Flutter. Butuh API yang sedang jalan (`make dev` di root repo).

```bash
make run    # iOS simulator / perangkat; Android emulator: make run API_BASE_URL=http://10.0.2.2:8080
make test
make lint
```

## Struktur

```
lib/
  main.dart          composition root
  domain/            Recipe, NewRecipe, ValidationFailure (Dart murni)
  application/       use case (KitabDapur) + port (KitabApi, KitabIdStore)
  infrastructure/    HttpKitabApi (package:http), PrefsKitabIdStore (shared_preferences)
  presentation/      halaman Flutter
```

Dependensi mengarah ke dalam: `presentation`/`infrastructure` → `application` → `domain`.
`HttpKitabApi` ditulis tangan mengikuti `api/openapi.yaml` di root repo.
