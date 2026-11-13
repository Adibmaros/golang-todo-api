# Golang Todo API

Proyek membangun REST API dengan Go menggunakan framework [Gin](https://gin-gonic.com/), MySQL, GORM, dan autentikasi JWT.

## Tech Stack

- **Go** (`1.26.0`)
- **Web Framework:** [Gin](https://gin-gonic.com/) (`github.com/gin-gonic/gin`)
- **ORM & Database:** [GORM](https://gorm.io/) (`gorm.io/gorm`) & MySQL (`gorm.io/driver/mysql`)
- **Konfigurasi:** `github.com/joho/godotenv`
- **Validasi:** `github.com/go-playground/validator/v10`
- **Autentikasi:** JWT (`github.com/golang-jwt/jwt/v5`) & Bcrypt (`golang.org/x/crypto/bcrypt`)

## Struktur Proyek

```text
todo-api/
├── cmd/
│   └── main.go                      # Entry point aplikasi, dependency wiring & routing
├── internal/
│   ├── auth/
│   │   └── jwt.go                   # Helper generate & parse token JWT
│   ├── database/
│   │   └── mysql.go                 # Koneksi database MySQL via GORM
│   ├── dto/
│   │   ├── auth.go                  # Request DTO untuk register & login
│   │   └── todo.go                  # Request/response DTO Todo & pagination
│   ├── handler/
│   │   ├── auth.go                  # Handler register & login
│   │   ├── health.go                # Handler health check
│   │   ├── todo.go                  # Handler CRUD Todo
│   │   └── validation.go            # Helper format error validasi
│   ├── middleware/
│   │   └── auth.go                  # Middleware verifikasi JWT
│   ├── model/
│   │   ├── todo.go                  # Model Todo (relasi UserID)
│   │   └── user.go                  # Model User
│   ├── repository/
│   │   ├── todo_repository.go       # Query database Todo (filter, sort, pagination)
│   │   └── user_repository.go       # Query database User
│   ├── response/
│   │   └── response.go              # Helper response JSON & error
│   └── service/
│       ├── auth_service.go          # Business logic register, login & bcrypt
│       └── todo_service.go          # Business logic Todo & perhitungan pagination
├── migrations/
│   └── init.sql                     # SQL pembuatan database awal
├── .env.example                     # Template environment variable
├── .gitignore
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
3. **Jalankan server:**
   ```powershell
   go run ./cmd/main.go
   ```
   Server akan berjalan di `http://localhost:8080` dan GORM `AutoMigrate` akan otomatis membuat tabel `users` dan `todos`.

### Menjalankan via Docker

Pastikan Docker dan Docker Compose sudah terpasang, lalu jalankan:

```powershell
copy .env.example .env
docker compose up --build
```

Perintah tersebut akan membangun image Go API (multi-stage build) dan menjalankan container `mysql` beserta `api` secara otomatis. Untuk menghentikan container:

```powershell
docker compose down
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
| `GET` | `/health` | Cek status server | - |
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
