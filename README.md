# User Management Microservice (GoFiber MVC)

Author: **Amas Parera**

Proyek ini adalah microservice User Management modern berbasis Go, terinspirasi dari ORY Kratos namun dengan struktur MVC dan framework GoFiber. Dirancang untuk kebutuhan sistem terdistribusi (microservices) dengan dokumentasi lengkap dan kode berbahasa Indonesia.

---

## Fitur Utama
- Registrasi, login, dan manajemen user berbasis JWT
- Verifikasi email, recovery password, audit log
- Dukungan TOTP (2FA) dan custom attributes
- Dokumentasi API interaktif (Swagger UI)
- Siap dikembangkan untuk kebutuhan enterprise (OAuth2, session management, dsb)

## Struktur Project
- `cmd/` — Entry point aplikasi (opsional)
- `config/` — Konfigurasi aplikasi
- `controller/` — Logika controller
- `model/` — Struktur data/model database
- `repository/` — Abstraksi akses data
- `service/` — Bisnis logic
- `routes/` — Routing aplikasi
- `migration/` — Migrasi database
- `docs/` — Dokumentasi OpenAPI/Swagger
- `main.go` — Bootstrap aplikasi

## Instalasi
1. Pastikan Go 1.20+ sudah terinstall
2. Clone repository ini
3. Jalankan:
   ```bash
   go mod tidy
   go run main.go
   ```
4. Server berjalan di `http://localhost:3001`

## Dokumentasi API
- Dokumentasi interaktif Swagger UI tersedia di:  
  [http://localhost:3001/docs/swagger/](http://localhost:3001/docs/swagger/)
- Spesifikasi OpenAPI: `docs/openapi.yaml`
- Dokumentasi endpoint ringkas: `API_DOKUMENTASI.md`

## Roadmap Singkat
- [x] Fitur dasar user management (register, login, JWT, recovery, dsb)
- [x] Swagger UI & dokumentasi API
- [ ] Verifikasi email modular & flow-based
- [ ] TOTP production ready (setup, enable, backup code)
- [ ] Password policy & brute-force protection
- [ ] Social login (OAuth2)
- [ ] Self-service flow & UI node

## Kontribusi & Lisensi
Project ini terbuka untuk kontribusi. Silakan buat issue atau pull request untuk perbaikan/fitur baru.

---

**Amas Parera — 2025**
