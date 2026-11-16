# Golang Todo API

Proyek membangun REST API dengan Go menggunakan framework [Gin](https://gin-gonic.com/), MySQL, GORM, autentikasi JWT, dan dokumentasi Swagger.

## Tech Stack

- **Go** (`1.26.0`)
- **Web Framework:** [Gin](https://gin-gonic.com/) (`github.com/gin-gonic/gin`)
- **ORM & Database:** [GORM](https://gorm.io/) (`gorm.io/gorm`) & MySQL (`gorm.io/driver/mysql`)
- **Konfigurasi:** `github.com/joho/godotenv`
- **Validasi:** `github.com/go-playground/validator/v10`
- **Autentikasi:** JWT (`github.com/golang-jwt/jwt/v5`) & Bcrypt (`golang.org/x/crypto/bcrypt`)
- **Dokumentasi API:** Swagger/OpenAPI (`github.com/swaggo/swag`, `github.com/swaggo/gin-swagger`, `github.com/swaggo/files`)
- **Frontend:** Vue 3 + Vite (disajikan via Nginx pada Docker)
- **Containerization:** Docker & Docker Compose

## Fitur Production Readiness

- **Graceful Shutdown** (`http.Server` + signal handling `SIGINT` / `SIGTERM`)
- **Structured Logging** menggunakan `log/slog` (JSON format)
- **Database Connection Pooling** (`SetMaxIdleConns`, `SetMaxOpenConns`, `SetConnMaxLifetime`)
- **Middleware CORS & Rate Limiting** (`60` request/menit per IP)
- **Health Check Database** (`PingContext` ke MySQL pada `GET /health`)

## Struktur Proyek

```text
todo-api/
├── cmd/
│   └── main.go                      # Entry point, wiring, routing & graceful shutdown
├── deploy/
│   └── Caddyfile.example            # Contoh konfigurasi reverse proxy Caddy
├── docs/
│   ├── docs.go                      # Swagger generated docs
│   ├── swagger.json                 # OpenAPI spec (JSON)
│   └── swagger.yaml                 # OpenAPI spec (YAML)
├── frontend/
│   ├── src/
│   │   ├── App.vue                  # Halaman Login/Register & Dashboard Todo
│   │   ├── api.js                   # Client fetch & manajemen JWT
│   │   ├── main.js
│   │   └── style.css                # Styling light mode
│   ├── Dockerfile                   # Multi-stage build Node + Nginx
│   ├── nginx.conf                   # SPA fallback & proxy /api ke backend
│   ├── package.json
│   └── vite.config.js               # Dev proxy /api ke backend lokal
├── internal/
│   ├── auth/
│   │   └── jwt.go                   # Helper generate & parse token JWT
│   ├── database/
│   │   └── mysql.go                 # Koneksi database MySQL & connection pooling
│   ├── dto/
│   │   ├── auth.go                  # Request DTO untuk register & login
│   │   └── todo.go                  # Request/response DTO Todo & pagination
│   ├── handler/
│   │   ├── auth.go                  # Handler register & login
│   │   ├── health.go                # Handler health check & DB ping
│   │   ├── todo.go                  # Handler CRUD Todo
│   │   ├── todo_test.go             # HTTP handler test (httptest)
│   │   └── validation.go            # Helper format error validasi
│   ├── middleware/
│   │   ├── auth.go                  # Middleware verifikasi JWT
│   │   ├── cors.go                  # Middleware CORS
│   │   └── ratelimit.go             # Middleware rate limiter per IP
│   ├── model/
│   │   ├── todo.go                  # Model Todo (relasi UserID)
│   │   └── user.go                  # Model User (relasi Todos)
│   ├── repository/
│   │   ├── todo_repository.go       # Query database Todo (filter, sort, pagination)
│   │   ├── todo_repository_test.go  # Integration test MySQL
│   │   └── user_repository.go       # Query database User
│   ├── response/
│   │   └── response.go              # Helper response JSON & error
│   └── service/
│       ├── auth_service.go          # Business logic register, login & bcrypt
│       ├── todo_service.go          # Business logic Todo & perhitungan pagination
│       └── todo_service_test.go     # Unit test TodoService
├── migrations/
│   └── init.sql                     # SQL pembuatan database awal
├── .dockerignore
├── .env.example                     # Template environment variable
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

## Cara Menjalankan

1. **Siapkan database MySQL** (pastikan MySQL berjalan di `127.0.0.1:3306`), lalu jalankan query di `migrations/init.sql`:
   ```sql
   CREATE DATABASE todo_db;
   ```
2. **Konfigurasi environment variable:**
   Salin `.env.example` menjadi `.env`, lalu sesuaikan nilainya (koneksi database, `APP_PORT`, dan `JWT_SECRET`):
   ```powershell
   copy .env.example .env
   ```
3. **Jalankan backend:**
   ```powershell
   go run ./cmd/main.go
   ```
   Server akan berjalan sesuai `APP_PORT` di `.env` (atau `http://localhost:8080` jika default) dan GORM `AutoMigrate` akan otomatis membuat tabel `users` dan `todos`.
4. **Jalankan frontend (opsional, untuk development lokal):**
   ```powershell
   cd frontend
   npm install
   npm run dev
   ```
   Frontend dev server berjalan di `http://localhost:3001` dengan proxy `/api` ke `http://localhost:8080`.

### Menjalankan via Docker

Pastikan Docker dan Docker Compose sudah terpasang, lalu jalankan:

```powershell
copy .env.example .env
docker compose up -d --build
```

Perintah tersebut akan menjalankan 3 container: `mysql`, `api` (port `8081`), dan `frontend` Vue.js + Nginx (port `3001`).
- Akses UI di `http://localhost:3001`
- Contoh konfigurasi Caddy untuk domain `todo.adibmaros.web.id` tersedia di `deploy/Caddyfile.example`

Untuk menghentikan container:

```powershell
docker compose down
```

## Dokumentasi API (Swagger)

Saat server berjalan, buka Swagger UI di browser melalui:

```text
http://localhost:8080/swagger/index.html
```

Untuk menguji endpoint `/todos`, login melalui `POST /auth/login` terlebih dahulu, lalu klik tombol **Authorize** di Swagger UI dan masukkan `Bearer <token>`.

Jika mengubah anotasi Swagger pada handler atau `cmd/main.go`, generate ulang dokumentasi dengan:

```powershell
go run github.com/swaggo/swag/cmd/swag init -g cmd/main.go
```

## Menjalankan Test

Jalankan seluruh unit test & handler test:

```powershell
go test -v ./...
```

Untuk ikut menjalankan **integration test** ke database MySQL khusus test:
1. Buat database test di MySQL, misalnya:
   ```sql
   CREATE DATABASE todo_db_test;
   ```
2. Set environment variable lalu jalankan `go test`:
   ```powershell
   $env:TEST_DB_NAME="todo_db_test"
   go test -v ./...
   ```
   *(Opsional: gunakan `TEST_DB_HOST`, `TEST_DB_PORT`, `TEST_DB_USER`, `TEST_DB_PASS`, atau langsung `TEST_DB_DSN` jika konfigurasi berbeda dari default lokal `root@127.0.0.1:3306`.)*

## Daftar Endpoint

### Public

| Method | Endpoint | Deskripsi | Contoh Body |
| :--- | :--- | :--- | :--- |
| `GET` | `/` | Informasi dasar API | - |
| `GET` | `/health` | Cek status server & koneksi database | - |
| `GET` | `/swagger/*any` | Dokumentasi interaktif Swagger UI | - |
| `POST` | `/auth/register` | Registrasi user baru | `{"email": "user@mail.com", "password": "secret123"}` |
| `POST` | `/auth/login` | Login & dapatkan JWT | `{"email": "user@mail.com", "password": "secret123"}` |

### Protected (`/todos`)

Semua endpoint `/todos` wajib menyertakan header `Authorization: Bearer <token>` dan hanya mengakses data milik user yang sedang login.

| Method | Endpoint | Deskripsi | Contoh Body / Query |
| :--- | :--- | :--- | :--- |
| `GET` | `/todos` | Ambil daftar todo (pagination, filter, sort) | `?page=1&limit=10&done=false&sort=created_at` |
| `GET` | `/todos/:id` | Ambil detail todo berdasarkan ID | - |
| `POST` | `/todos` | Tambah todo baru | `{"title": "Belajar Go"}` |
| `PUT` | `/todos/:id` | Update todo berdasarkan ID | `{"title": "Belajar Go", "done": true}` |
| `DELETE` | `/todos/:id` | Hapus todo berdasarkan ID | - |

**Query Parameter `GET /todos`:**
- `page` (opsional, default `1`): nomor halaman (`>= 1`)
- `limit` (opsional, default `10`, maks `100`): jumlah data per halaman
- `done` (opsional): filter status selesai (`true` / `false`)
- `sort` (opsional, default `created_at`): urutan data (`created_at`, `created_asc`, `title`, `done`)
